package ngsi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/matryer/is"
)

func TestGetBeachAsGeoJSON(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", createURL("/entities?type=Beach"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Beach", "temperature")

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
	b1 := fiware.NewBeach("omaha", "Omaha Beach", location).WithDescription("This is a nice beach!")
	b1.WaterTemperature = types.NewNumberProperty(7.2)

	contextSource.GetEntitiesFunc = func(q Query, callback QueryEntitiesCallback) error {
		callback(b1)
		return nil
	}

	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, 200) // failed to get geojson data
}

func TestGetLongLatForGeoPropertyPoint(t *testing.T) {
	is := is.New(t)

	lat := 65.2789
	lon := 17.2961
	location := geojson.CreateGeoJSONPropertyFromWGS84(lon, lat)

	beach := fiware.NewBeach("beach", "Beach Beach", location).WithDescription("A very beachy beach.")

	is.Equal(beach.Location.GetAsPoint().Latitude(), lat)
	is.Equal(beach.Location.GetAsPoint().Longitude(), lon)
}

func TestGetLongLatForGeoPropertyMultiPolygon(t *testing.T) {
	is := is.New(t)

	lon := 17.2961
	lat := 65.2789
	poly := [][][][]float64{}
	row1 := [][][]float64{}
	row2 := [][]float64{}
	row3 := []float64{lon, lat}
	row2 = append(row2, row3)
	row1 = append(row1, row2)
	poly = append(poly, row1)

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon(poly)

	beach := fiware.NewBeach("beach", "Beach Beach", location).WithDescription("A very beachy beach.")
	is.Equal(beach.Location.GetAsPoint().Latitude(), lat)
	is.Equal(beach.Location.GetAsPoint().Longitude(), lon)
}

func TestGetWaterQualityObservedAsGeoJSON(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", createURL("/entities?type=WaterQualityObserved&options=keyValues"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("WaterQualityObserved", "temperature")

	wqo1 := fiware.NewWaterQualityObserved("badtempsensor", 64.2789, 17.2961, "2021-04-22T17:23:41Z")

	contextSource.GetEntitiesFunc = func(q Query, callback QueryEntitiesCallback) error {
		callback(wqo1)
		return nil
	}

	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, 200) // failed to get geojson data
}

func TestGetEntitiesAsSimplifiedGeoJSON(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities?type=Beach&options=keyValues"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Beach", "temperature")

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
	b1 := fiware.NewBeach("omaha", "Omaha Beach", location).WithDescription("This is a nice beach!")
	b1.WaterTemperature = types.NewNumberProperty(7.2)

	contextSource.GetEntitiesFunc = func(q Query, callback QueryEntitiesCallback) error {
		callback(b1)
		return nil
	}

	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	is.Equal(w.Code, 200) // failed to get geojson data
}
