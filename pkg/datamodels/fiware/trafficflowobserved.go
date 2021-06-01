package fiware

import (
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//TrafficFlowObserved is a Fiware entity
type TrafficFlowObserved struct {
	ngsi.BaseEntity
	DateObserved        *ngsi.DateTimeProperty         `json:"dateObserved"`
	DateCreated         *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified        *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObservedTo      *ngsi.DateTimeProperty         `json:"dateObservedTo,omitempty"`
	DateObservedFrom    *ngsi.DateTimeProperty         `json:"dateObservedFrom,omitempty"`
	Location            geojson.GeoJSONProperty        `json:"location,omitempty"`
	LaneID              *ngsi.NumberProperty           `json:"laneID"`
	AverageVehicleSpeed *ngsi.NumberProperty           `json:"averageVehicleSpeed,omitempty"`
	Intensity           *ngsi.NumberProperty           `json:"intensity,omitempty"`
	RefRoadSegment      *ngsi.SingleObjectRelationship `json:"refRoadSegment,omitempty"`
}

//NewTrafficFlowObserved creates a new TrafficFlowObserved from given ID
func NewTrafficFlowObserved(id string, latitude float64, longitude float64, observedAt string, laneID int, intensity int) *TrafficFlowObserved {
	if !strings.HasPrefix(id, TrafficFlowObservedIDPrefix) {
		id = TrafficFlowObservedIDPrefix + id
	}

	dateTimeValue := ngsi.CreateDateTimeProperty(observedAt)
	lane := ngsi.NewNumberPropertyFromInt(laneID)
	intense := ngsi.NewNumberPropertyFromInt(intensity)

	return &TrafficFlowObserved{
		DateObserved: dateTimeValue,
		Location:     *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		LaneID:       lane,
		Intensity:    intense,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "TrafficFlowObserved",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
