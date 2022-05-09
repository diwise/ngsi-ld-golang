package fiware

import (
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
