package diwise

import (
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

type Lifebuoy struct {
	ngsi.BaseEntity
	Location     geojson.GeoJSONProperty `json:"location"`
	Status       *ngsi.TextProperty      `json:"status"`
	DateObserved *ngsi.DateTimeProperty  `json:"dateObserved,omitempty"`
}

const LifebuoyTypeName string = "Lifebuoy"
const LifebuoyIDPrefix string = urnPrefix + LifebuoyTypeName + ":"

func NewLifebuoy(id string, status string) *Lifebuoy {
	if !strings.HasPrefix(id, LifebuoyIDPrefix) {
		id = LifebuoyIDPrefix + id
	}

	return &Lifebuoy{
		Status: ngsi.NewTextProperty(status),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: LifebuoyTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
