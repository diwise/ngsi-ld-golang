package fiware

import (
	"bytes"
	"encoding/json"

	geojson "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsitypes "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

type CityWork struct {
	ngsitypes.BaseEntity
	Location     *geojson.GeoJSONProperty    `json:"location,omitempty"`
	Description  *ngsitypes.TextProperty     `json:"description,omitempty"`
	DateCreated  ngsitypes.DateTimeProperty  `json:"dateCreated"`
	DateModified *ngsitypes.DateTimeProperty `json:"dateModified,omitempty"`
	EndDate      ngsitypes.DateTimeProperty  `json:"endDate"`
	StartDate    ngsitypes.DateTimeProperty  `json:"startDate"`
}

type cityworkDTO struct {
	ngsitypes.BaseEntity
	Location     json.RawMessage             `json:"location,omitempty"`
	Description  *ngsitypes.TextProperty     `json:"description,omitempty"`
	DateCreated  ngsitypes.DateTimeProperty  `json:"dateCreated"`
	DateModified *ngsitypes.DateTimeProperty `json:"dateModified,omitempty"`
	EndDate      ngsitypes.DateTimeProperty  `json:"endDate"`
	StartDate    ngsitypes.DateTimeProperty  `json:"startDate"`
}

func NewCityWork(entityID string) CityWork {
	cw := CityWork{
		BaseEntity: ngsitypes.BaseEntity{
			ID:   CityWorkIDPrefix + entityID,
			Type: CityWorkTypeName,
			Context: []string{
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
				"https://raw.githubusercontent.com/smart-data-models/dataModel.Transportation/master/context.jsonld",
			},
		},
	}

	return cw
}

func (c *CityWork) UnmarshalJSON(data []byte) error {
	dto := &cityworkDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		c.ID = dto.ID
		c.Type = dto.Type
		c.Description = dto.Description

		c.DateCreated = dto.DateCreated
		c.DateModified = dto.DateModified
		c.StartDate = dto.StartDate
		c.EndDate = dto.EndDate

		c.Context = dto.Context

		c.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
	}

	return err
}
