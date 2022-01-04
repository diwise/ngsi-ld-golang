package fiware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//TrafficFlowObserved is a Fiware entity
type TrafficFlowObserved struct {
	ngsi.BaseEntity
	DateObserved        ngsi.TextProperty              `json:"dateObserved"`
	DateCreated         *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified        *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObservedTo      *ngsi.DateTimeProperty         `json:"dateObservedTo,omitempty"`
	DateObservedFrom    *ngsi.DateTimeProperty         `json:"dateObservedFrom,omitempty"`
	Location            *geojson.GeoJSONProperty       `json:"location,omitempty"`
	LaneID              *ngsi.NumberProperty           `json:"laneID"`
	AverageVehicleSpeed *ngsi.NumberProperty           `json:"averageVehicleSpeed,omitempty"`
	Intensity           *ngsi.NumberProperty           `json:"intensity,omitempty"`
	RefRoadSegment      *ngsi.SingleObjectRelationship `json:"refRoadSegment,omitempty"`
}

type trafficFlowObservedDTO struct {
	ngsi.BaseEntity
	DateObserved        ngsi.TextProperty              `json:"dateObserved"`
	DateCreated         *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified        *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObservedTo      *ngsi.DateTimeProperty         `json:"dateObservedTo,omitempty"`
	DateObservedFrom    *ngsi.DateTimeProperty         `json:"dateObservedFrom,omitempty"`
	Location            json.RawMessage                `json:"location,omitempty"`
	LaneID              *ngsi.NumberProperty           `json:"laneID"`
	AverageVehicleSpeed *ngsi.NumberProperty           `json:"averageVehicleSpeed,omitempty"`
	Intensity           *ngsi.NumberProperty           `json:"intensity,omitempty"`
	RefRoadSegment      *ngsi.SingleObjectRelationship `json:"refRoadSegment,omitempty"`
}

//NewTrafficFlowObserved creates a new TrafficFlowObserved from given ID
func NewTrafficFlowObserved(id string, observedAt string, laneID int, intensity int) *TrafficFlowObserved {
	if !strings.HasPrefix(id, TrafficFlowObservedIDPrefix) {
		id = TrafficFlowObservedIDPrefix + id
	}

	dateTimeValue := ngsi.NewTextProperty(observedAt)
	lane := ngsi.NewNumberPropertyFromInt(laneID)
	intense := ngsi.NewNumberPropertyFromInt(intensity)

	return &TrafficFlowObserved{
		DateObserved: *dateTimeValue,
		LaneID:       lane,
		Intensity:    intense,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: TrafficFlowObservedTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

func (tfo *TrafficFlowObserved) UnmarshalJSON(data []byte) error {
	dto := &trafficFlowObservedDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		tfo.ID = dto.ID
		tfo.Type = dto.Type

		tfo.DateObserved = dto.DateObserved
		tfo.DateCreated = dto.DateCreated
		tfo.DateModified = dto.DateModified
		tfo.DateObservedFrom = dto.DateObservedFrom
		tfo.DateObservedTo = dto.DateObservedTo

		tfo.LaneID = dto.LaneID
		tfo.AverageVehicleSpeed = dto.AverageVehicleSpeed
		tfo.Intensity = dto.Intensity
		tfo.RefRoadSegment = dto.RefRoadSegment

		tfo.Context = dto.Context

		if dto.Location != nil {
			tfo.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
		}
	}

	return err
}
