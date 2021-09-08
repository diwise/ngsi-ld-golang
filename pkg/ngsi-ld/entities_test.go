package ngsi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	ngsierrors "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/errors"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/matryer/is"
)

func createURL(path string, params ...string) string {
	url := "http://localhost:8080/ngsi-ld/v1" + path

	if len(params) > 0 {
		url = url + "?"

		for _, p := range params {
			url = url + p + "&"
		}

		url = strings.TrimSuffix(url, "&")
	}

	return url
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNewRoadSegment(t *testing.T) {
	id := "urn:ngsi-ld:RoadSegment:road1"
	roadID := "road"
	name := "segName"

	coords := [][2]float64{
		{0.0, 0.0},
		{1.1, 6.5},
	}

	roadSegment := fiware.NewRoadSegment(id, name, roadID, coords, nil)

	if len(roadSegment.Location.Value.Coordinates) != 2 {
		t.Error("Number of coords not as expected.")
	}

	for index, point := range roadSegment.Location.Value.Coordinates {
		if point != coords[index] {
			t.Error("Coords do not match.")
		}
	}
}

const tfoStr string = `{
    "id": "urn:ngsi-ld:TrafficFlowObserved:TrafficFlowObserved",
    "type": "TrafficFlowObserved",
    "dateObserved": {
        "type": "Property",
        "value": "2016-12-07T11:10:00Z"
    },
    "dateObservedFrom": {
        "type": "Property",
        "value": {
            "@type": "DateTime",
            "@value": "2016-12-07T11:10:00Z"
        }
    },
    "dateObservedTo": {
        "type": "Property",
        "value": {
            "@type": "DateTime",
            "@value": "2016-12-07T11:15:00Z"
        }
    },
    "intensity": {
        "type": "Property",
        "value": 197
    },
    "laneId": {
        "type": "Property",
        "value": 1
    },
    "location": {
      	"type": "GeoProperty",
      	"value": {
        		"coordinates": [
          		17.414662,
          		62.421033
        		],
        		"type": "Point"
      	}
    },
    "averageVehicleSpeed": {
        "type": "Property",
        "value": 52.6
    },
    "@context": [
        "https://schema.lab.fiware.org/ld/context",
        "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"
    ]
}`

func TestCreateEntityWorksForTrafficFlowObserved(t *testing.T) {
	is := is.New(t)
	entityID := fiware.TrafficFlowObservedIDPrefix + "TrafficFlowObserved"
	tfo := fiware.TrafficFlowObserved{}

	err := json.Unmarshal([]byte(tfoStr), &tfo)
	is.NoErr(err)

	jsonBytes, _ := json.Marshal(tfo)
	byteReader := bytes.NewBuffer(jsonBytes)
	typeName := tfo.Type

	req, _ := http.NewRequest("POST", createURL("/entities"), byteReader)
	w := httptest.NewRecorder()

	ctxReg, ctxSrc := newContextRegistryWithSourceForType(typeName)

	NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusCreated)         // failed to create entity
	is.Equal(ctxSrc.createdEntityType, typeName) // create entity called with wrong type name
	is.Equal(ctxSrc.createdEntity, entityID)     // create entity called with wrong entity id
}

func TestCreateEntityUsesCorrectTypeAndID(t *testing.T) {
	is := is.New(t)
	entityID := fiware.DeviceIDPrefix + "livboj"
	byteReader, typeName := newEntityAsByteBuffer(entityID)
	req, _ := http.NewRequest("POST", createURL("/entities"), byteReader)
	w := httptest.NewRecorder()

	ctxReg, ctxSrc := newContextRegistryWithSourceForType(typeName)

	NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusCreated)         // failed to create entity
	is.Equal(ctxSrc.createdEntityType, typeName) // create entity called with wrong type name
	is.Equal(ctxSrc.createdEntity, entityID)     // create entity called with wrong entity id

}

func TestThatProblemsAreReportedWithCorrectContentType(t *testing.T) {
	is := is.New(t)
	byteBuffer, _ := newEntityAsByteBuffer("id")
	req, _ := http.NewRequest("POST", createURL("/entities"), byteBuffer)
	w := httptest.NewRecorder()

	NewCreateEntityHandler(NewContextRegistry()).ServeHTTP(w, req)

	response := w.Result()
	contentTypes := response.Header["Content-Type"]

	if len(contentTypes) != 1 {
		t.Errorf("Response returned %d content types. Expected 1!", len(contentTypes))
	} else {
		contentType := contentTypes[0]
		is.Equal(contentType, ngsierrors.ProblemReportContentType) // problem reported with wrong content type
	}
}

func TestCreateEntityFailsWithNoContextSources(t *testing.T) {
	is := is.New(t)
	byteBuffer, _ := newEntityAsByteBuffer("id")
	req, _ := http.NewRequest("POST", createURL("/entities"), byteBuffer)
	w := httptest.NewRecorder()

	NewCreateEntityHandler(NewContextRegistry()).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusBadRequest) // wrong response code when posting device with no context sources.
}

func TestCreateEntityHandlesFailureFromContextSource(t *testing.T) {
	is := is.New(t)
	byteBuffer, typeName := newEntityAsByteBuffer("id")
	req, _ := http.NewRequest("POST", createURL("/entities"), byteBuffer)
	w := httptest.NewRecorder()

	ctxReg, ctxSrc := newContextRegistryWithSourceForType(typeName)
	ctxSrc.createEntityShouldFailWithError = errors.New("failure")

	NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusBadRequest) // wrong response code when create entity fails.
}

func TestGetEntitiesWithoutAttributesOrTypesFails(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", createURL("/entitites"), nil)
	w := httptest.NewRecorder()

	NewQueryEntitiesHandler(NewContextRegistry()).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusBadRequest) // wrong response code when type and attrs are missing.
}

func TestGetEntitiesWithAttribute(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entitites", "attrs=snowHeight"), nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextRegistry.Register(newMockedContextSource(
		"", "snowHeight",
		e(""),
	))

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusOK) // failed to get entity by attribute
}

func TestGetEntitiesForDevice(t *testing.T) {
	is := is.New(t)

	deviceID := fiware.DeviceIDPrefix + "mydevice"
	req, _ := http.NewRequest("GET", createURL("/entitites", "attrs=snowHeight", "q=refDevice==\""+deviceID+"\""), nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("", "snowHeight")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(contextSource.queriedDevice, deviceID) // queried device did not match expectations
	is.Equal(w.Code, http.StatusOK)                 // unexpected response code
}

func TestGetEntitiesWithGeoQueryNearPoint(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL(
		"/entitites",
		"type=RoadSegment",
		"georel=near;maxDistance==2000",
		"geometry=Point",
		"coordinates=[8,40]"),
		nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("RoadSegment", "")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusOK) // unexpected response code

	query := contextSource.generatedQuery
	is.True(query.IsGeoQuery()) // expected a geo query from QueryEntitiesHandler

	geo := query.Geo()
	is.Equal(geo.GeoRel, "near") // geospatial relation not correctly saved in geo query

	distance, _ := geo.Distance()
	is.Equal(distance, uint32(2000)) // unexpected near distance parsed from geo query

	x, y, _ := geo.Point()
	is.Equal(fmt.Sprintf("(%.1f,%.1f)", x, y), "(8.0,40.0)") // mismatching point
}

func TestGetEntitiesWithGeoQueryWithinRect(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", createURL(
		"/entitites",
		"type=RoadSegment",
		"georel=within",
		"geometry=Polygon",
		"coordinates=[[8,40],[9,41],[10,42]]"),
		nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("RoadSegment", "")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusOK) // unexpected response code

	query := contextSource.generatedQuery
	is.True(query.IsGeoQuery()) // expected a GeoQuery from the QueryEntititesHandler

	geo := query.Geo()
	is.Equal(geo.GeoRel, "within") // geospatial relation not correctly saved in geo query

	lon0, lat0, lon1, lat1, _ := geo.Rectangle()
	is.Equal(fmt.Sprintf("[(%.1f,%.1f),(%.1f,%.1f)]", lat0, lon0, lat1, lon1), "[(40.0,8.0),(42.0,10.0)]") // bad coordinates in GeoQuery rect
}

func TestRetrieveEntity(t *testing.T) {
	is := is.New(t)
	deviceID := fiware.DeviceIDPrefix + "mydevice"
	req, _ := http.NewRequest("GET", createURL("/entities/"+deviceID), nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Device", "")
	contextRegistry.Register(contextSource)

	NewRetrieveEntityHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(contextSource.retrievedEntity, deviceID) // DeviceID does not match retrievedEntity
	is.Equal(w.Code, http.StatusOK)                   // unexpected response code
}

func TestUpdateEntitityAttributes(t *testing.T) {
	is := is.New(t)

	deviceID := fiware.DeviceIDPrefix + "mydevice"
	jsonBytes, _ := json.Marshal(e("testvalue"))

	req, _ := http.NewRequest("PATCH", createURL("/entities/"+deviceID+"/attrs/"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("", "value")
	contextRegistry.Register(contextSource)

	NewUpdateEntityAttributesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(contextSource.patchedEntity, deviceID) // patched entity did not match expectations.
}

type mockEntity struct {
	Value string
}

func e(val string) mockEntity {
	return mockEntity{Value: val}
}

func newContextRegistryWithSourceForType(typeName string) (ContextRegistry, *mockCtxSource) {
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource(typeName, "")
	contextRegistry.Register(contextSource)
	return contextRegistry, contextSource
}

func newEntityAsByteBuffer(entityID string) (io.Reader, string) {

	device := fiware.NewDevice(entityID, "")
	jsonBytes, _ := json.Marshal(device)
	return bytes.NewBuffer(jsonBytes), device.Type
}

func newMockedContextSource(typeName string, attributeName string, e ...mockEntity) *mockCtxSource {
	source := &mockCtxSource{typeName: typeName, attributeName: attributeName}
	for _, entity := range e {
		source.entities = append(source.entities, entity)
	}
	return source
}

type mockCtxSource struct {
	typeName      string
	attributeName string
	entities      []Entity

	createEntityShouldFailWithError error

	queriedDevice     string
	createdEntity     string
	createdEntityType string
	patchedEntity     string
	retrievedEntity   string

	generatedQuery Query
}

func (s *mockCtxSource) CreateEntity(typeName, entityID string, r Request) error {
	if s.createEntityShouldFailWithError == nil {
		s.createdEntity = entityID
		s.createdEntityType = typeName

		entity := &types.BaseEntity{}
		return r.DecodeBodyInto(entity)
	}

	return s.createEntityShouldFailWithError
}

func (s *mockCtxSource) GetEntities(q Query, cb QueryEntitiesCallback) error {

	s.generatedQuery = q

	if q.HasDeviceReference() {
		s.queriedDevice = q.Device()
	}

	for _, e := range s.entities {
		cb(e)
	}
	return nil
}

func (s *mockCtxSource) UpdateEntityAttributes(entityID string, req Request) error {
	s.patchedEntity = entityID
	e := &mockEntity{}
	return req.DecodeBodyInto(e)
}

func (s *mockCtxSource) ProvidesAttribute(attributeName string) bool {
	return s.attributeName == attributeName
}

func (s *mockCtxSource) ProvidesEntitiesWithMatchingID(entityID string) bool {
	return true
}

func (s *mockCtxSource) GetProvidedTypeFromID(entityID string) (string, error) {
	return "", errors.New("not implemented")
}

func (s *mockCtxSource) ProvidesType(typeName string) bool {
	return s.typeName == typeName
}

func (s *mockCtxSource) RetrieveEntity(entityID string, req Request) (Entity, error) {
	s.retrievedEntity = entityID
	e := &mockEntity{}
	return e, nil
}
