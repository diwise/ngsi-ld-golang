package ngsi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/errors"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//Entity is an informational representative of something that is supposed to exist in the real world, physically or conceptually
type Entity interface {
}

//QueryEntitiesCallback is used when queried context sources should pass back any
//entities matching the query that has been passed in
type QueryEntitiesCallback func(entity Entity) error

func getEntityConverterFromRequest(r *http.Request) (string, func(interface{}) interface{}, *geojson.GeoJSONFeatureCollection) {
	// Default entity converter doesn't actually convert anything
	entityConverter := func(e interface{}) interface{} { return e }

	responseContentType := "application/ld+json;charset=utf-8"
	var geoJSONFeatureCollection *geojson.GeoJSONFeatureCollection

	// Check Accept to find out what kind of data the client wants
	for _, acceptableType := range r.Header["Accept"] {
		if strings.HasPrefix(acceptableType, geojson.ContentType) {
			options := r.URL.Query().Get("options")
			geoJSONFeatureCollection = geojson.NewGeoJSONFeatureCollection([]geojson.GeoJSONFeature{}, true)
			entityConverter = geojson.NewEntityConverter("location", options == "keyValues", geoJSONFeatureCollection)
			responseContentType = acceptableType
		}
	}

	return responseContentType, entityConverter, geoJSONFeatureCollection
}

//NewQueryEntitiesHandler handles GET requests for NGSI entities
func NewQueryEntitiesHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseContentType, entityConverter, geoJSONFeatureCollection := getEntityConverterFromRequest(r)

		entityTypeNames := r.URL.Query().Get("type")
		attributeNames := r.URL.Query().Get("attrs")

		if entityTypeNames == "" && attributeNames == "" {
			errors.ReportNewBadRequestData(
				w,
				"A request for entities MUST specify at least one of type or attrs.",
			)
			return
		}

		entityTypes := strings.Split(entityTypeNames, ",")
		attributes := strings.Split(attributeNames, ",")

		q := r.URL.Query().Get("q")
		query, err := newQueryFromParameters(r, entityTypes, attributes, q)
		if err != nil {
			errors.ReportNewBadRequestData(
				w, err.Error(),
			)
			return
		}

		contextSources := ctxReg.GetContextSourcesForQuery(query)

		var entities = []Entity{}
		var entityCount = uint64(0)
		var entityMaxCount = uint64(18446744073709551615) // uint64 max

		if query.PaginationLimit() > 0 {
			entityMaxCount = query.PaginationLimit()
		}

		for _, source := range contextSources {
			err = source.GetEntities(query, func(entity Entity) error {
				if entityCount < entityMaxCount {
					entities = append(entities, entityConverter(entity))
					entityCount++
				}
				return nil
			})
			if err != nil {
				break
			}
		}

		if err != nil {
			errors.ReportNewInternalError(
				w,
				"An internal error was encountered when trying to get entities from the context source: "+err.Error(),
			)
			return
		}

		var bytes []byte

		if geoJSONFeatureCollection != nil {
			bytes, err = json.MarshalIndent(geoJSONFeatureCollection, "", "  ")
		} else {
			bytes, err = json.MarshalIndent(entities, "", "  ")
		}

		if err != nil {
			errors.ReportNewInternalError(w, "Failed to encode response.")
			return
		}

		w.Header().Add("Content-Type", responseContentType)
		// TODO: Add a RFC 8288 Link header with information about previous and/or next page if they exist
		w.Write(bytes)
	})
}

type UpdateEntityAttributesCompletionCallback func(entityType, entityID string, request Request, logger zerolog.Logger)

//NewUpdateEntityAttributesHandler handles PATCH requests for NGSI entitity attributes
func NewUpdateEntityAttributesHandler(ctxReg ContextRegistry) http.HandlerFunc {
	noop := func(string, string, Request, zerolog.Logger) {}
	return NewUpdateEntityAttributesHandlerWithCallback(
		ctxReg, log.With().Logger(), noop,
	)
}

//NewUpdateEntityAttributesHandlerWithCallback handles PATCH requests for NGSI entitity
//attributes and calls a callback on successful completion
func NewUpdateEntityAttributesHandlerWithCallback(
	ctxReg ContextRegistry,
	logger zerolog.Logger,
	onsuccess UpdateEntityAttributesCompletionCallback) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sublogger := decorateLogger(r, logger)

		// TODO: Replace this string manipulation with a callback that can use the http router's
		//		 functionality to extract URL params ...
		entitiesIdx := strings.Index(r.URL.Path, "/entities/")
		attrsIdx := strings.LastIndex(r.URL.Path, "/attrs/")

		if entitiesIdx == -1 || attrsIdx == -1 || attrsIdx < entitiesIdx {
			errors.ReportNewBadRequestData(
				w,
				"The supplied URL is invalid.",
			)
			return
		}

		entityID := r.URL.Path[entitiesIdx+10 : attrsIdx]

		request := newRequestWrapper(r)
		contextSources := ctxReg.GetContextSourcesForEntity(entityID)

		if len(contextSources) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err := contextSources[0].UpdateEntityAttributes(entityID, request)
		if err != nil {
			errors.ReportNewInvalidRequest(w, "Unable to update entity attributes: "+err.Error())
			return
		}

		entityType, err := contextSources[0].GetProvidedTypeFromID(entityID)
		if err == nil {
			// Call the success callback with the type and ID of the updated entity and the request instance
			onsuccess(entityType, entityID, request, sublogger)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type CreateEntityCompletionCallback func(entityType, entityID string, request Request, logger zerolog.Logger)

//NewCreateEntityHandler handles incoming POST requests for NGSI entities
func NewCreateEntityHandler(ctxReg ContextRegistry) http.HandlerFunc {
	noop := func(string, string, Request, zerolog.Logger) {}
	return NewCreateEntityHandlerWithCallback(ctxReg, log.With().Logger(), noop)
}

func NewCreateEntityHandlerWithCallback(
	ctxReg ContextRegistry,
	logger zerolog.Logger,
	onsuccess CreateEntityCompletionCallback) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sublogger := decorateLogger(r, logger)

		request := newRequestWrapper(r)

		entity := &types.BaseEntity{}
		err := request.DecodeBodyInto(entity)
		if err != nil {
			errors.ReportNewInvalidRequest(
				w,
				fmt.Sprintf("Unable to decode request payload: %s", err.Error()),
			)
		}

		contextSources := ctxReg.GetContextSourcesForEntityType(entity.Type)

		if len(contextSources) == 0 {
			errors.ReportNewInvalidRequest(
				w,
				fmt.Sprintf("No context sources found matching the provided type %s", entity.Type),
			)
			return
		}

		for _, source := range contextSources {
			err := source.CreateEntity(entity.Type, entity.ID, request)
			if err != nil {
				errors.ReportNewInvalidRequest(w, "Failed to create entity: "+err.Error())
				return
			}
		}

		onsuccess(entity.Type, entity.ID, request, sublogger)

		w.WriteHeader(http.StatusCreated)
	})
}

//NewRetrieveEntityHandler retrieves entity by ID.
func NewRetrieveEntityHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: A more elegant way to select the response content type ...
		responseContentType, entityConverter, _ := getEntityConverterFromRequest(r)

		entitiesIdx := strings.Index(r.URL.Path, "/entities/")

		if entitiesIdx == -1 {
			errors.ReportNewBadRequestData(
				w,
				"The supplied URL is invalid.",
			)
			return
		}

		entityID := r.URL.Path[entitiesIdx+10 : len(r.URL.Path)]

		contextSources := ctxReg.GetContextSourcesForEntity(entityID)

		if len(contextSources) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		request := newRequestWrapper(r)

		var entity Entity
		var err error

		for _, source := range contextSources {
			entity, err = source.RetrieveEntity(entityID, request)
			if err != nil {
				errors.ReportNewInvalidRequest(w, "Failed to find entity: "+err.Error())
				return
			}
			break
		}

		if entity == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		bytes, _ := json.Marshal(entityConverter(entity))

		w.Header().Add("Content-Type", responseContentType)
		w.Write(bytes)
	})
}

//decorateLogger looks for b3 trace headers and adds them to the logger if found
func decorateLogger(r *http.Request, logger zerolog.Logger) zerolog.Logger {
	traceHeaders := []string{
		http.CanonicalHeaderKey("x-b3-traceid"),
		http.CanonicalHeaderKey("x-b3-parentspanid"),
		http.CanonicalHeaderKey("x-b3-spanid"),
	}
	_, ok := r.Header[traceHeaders[0]]

	if ok {
		ctx := logger.With().Str("traceID", r.Header[traceHeaders[0]][0])

		for _, hdr := range traceHeaders {
			if len(r.Header[hdr]) > 0 {
				ctx = ctx.Str(strings.ToLower(hdr), r.Header[hdr][0])
			}
		}

		logger = ctx.Logger()
	}

	return logger
}
