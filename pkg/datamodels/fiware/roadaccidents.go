package fiware

import (
	"bytes"
	"encoding/json"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsitypes "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

type RoadAccident struct {
	ngsitypes.BaseEntity
	AccidentDate ngsitypes.DateTimeProperty  `json:"accidentDate"`
	Location     *geojson.GeoJSONProperty    `json:"location,omitempty"`
	Description  ngsitypes.TextProperty      `json:"description,omitempty"`
	DateCreated  ngsitypes.DateTimeProperty  `json:"dateCreated"`
	DateModified *ngsitypes.DateTimeProperty `json:"dateModified,omitempty"`
	Status       ngsitypes.TextProperty      `json:"status"`
}

type roadaccidentDTO struct {
	ngsitypes.BaseEntity
	AccidentDate ngsitypes.DateTimeProperty  `json:"accidentDate"`
	Location     json.RawMessage             `json:"location,omitempty"`
	Description  ngsitypes.TextProperty      `json:"description,omitempty"`
	DateCreated  ngsitypes.DateTimeProperty  `json:"dateCreated"`
	DateModified *ngsitypes.DateTimeProperty `json:"dateModified,omitempty"`
	Status       ngsitypes.TextProperty      `json:"status"`
}

func NewRoadAccident(entityID string) *RoadAccident {
	ra := &RoadAccident{
		BaseEntity: ngsitypes.BaseEntity{
			ID:   RoadAccidentIDPrefix + entityID,
			Type: RoadAccidentTypeName,
			Context: []string{
				"https://raw.githubusercontent.com/smart-data-models/dataModel.Transportation/master/context.jsonld",
			},
		},
	}

	return ra
}

func (r *RoadAccident) UnmarshalJSON(data []byte) error {
	dto := &roadaccidentDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		r.ID = dto.ID
		r.Type = dto.Type
		r.AccidentDate = dto.AccidentDate
		r.Description = dto.Description
		r.Status = dto.Status

		r.DateCreated = dto.DateCreated
		r.DateModified = dto.DateModified

		r.Context = dto.Context

		r.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
	}

	return err
}
