package fiware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Beach is a Fiware entity
type Beach struct {
	ngsi.BaseEntity
	Name             *ngsi.TextProperty            `json:"name,omitempty"`
	Description      *ngsi.TextProperty            `json:"description"`
	Location         geojson.GeoJSONGeometry       `json:"location,omitempty"`
	RefSeeAlso       *ngsi.MultiObjectRelationship `json:"refSeeAlso,omitempty"`
	SameAs           *ngsi.TextProperty            `json:"sameAs,omitempty"`
	WaterTemperature *ngsi.NumberProperty          `json:"waterTemperature,omitempty"`
	DateCreated      *ngsi.DateTimeProperty        `json:"dateCreated,omitempty"`
	DateModified     *ngsi.DateTimeProperty        `json:"dateModified,omitempty"`
}

type beachDTO struct {
	ngsi.BaseEntity
	Name             *ngsi.TextProperty            `json:"name,omitempty"`
	Description      *ngsi.TextProperty            `json:"description"`
	Location         json.RawMessage               `json:"location,omitempty"`
	RefSeeAlso       *ngsi.MultiObjectRelationship `json:"refSeeAlso,omitempty"`
	SameAs           *ngsi.TextProperty            `json:"sameAs,omitempty"`
	WaterTemperature *ngsi.NumberProperty          `json:"waterTemperature,omitempty"`
	DateCreated      *ngsi.DateTimeProperty        `json:"dateCreated,omitempty"`
	DateModified     *ngsi.DateTimeProperty        `json:"dateModified,omitempty"`
}

func (b Beach) ToGeoJSONFeature(propertyName string, simplified bool) (geojson.GeoJSONFeature, error) {
	g := geojson.NewGeoJSONFeature(b.ID, b.Type, b.Location.GeoPropertyValue())

	if simplified {
		g.SetProperty("name", b.Name.Value)
		g.SetProperty(propertyName, b.Location.GeoPropertyValue())

		if b.Description != nil {
			g.SetProperty("description", b.Description.Value)
		}

		if b.WaterTemperature != nil {
			g.SetProperty("waterTemperature", b.WaterTemperature.Value)
		}

		if b.RefSeeAlso != nil {
			g.SetProperty("refSeeAlso", b.RefSeeAlso.Object)
		}

		if b.SameAs != nil {
			g.SetProperty("sameAs", b.SameAs.Value)
		}
	} else {
		g.SetProperty("name", b.Name)
		g.SetProperty(propertyName, b.Location)

		g.SetProperty("description", b.Description)
		g.SetProperty("waterTemperature", b.WaterTemperature)
		g.SetProperty("refSeeAlso", b.RefSeeAlso)
		g.SetProperty("sameAs", b.SameAs)
	}

	return g, nil
}

//NewBeach creates a new Beach from given ID and name
func NewBeach(id, name string, location geojson.GeoJSONGeometry) *Beach {
	if !strings.HasPrefix(id, BeachIDPrefix) {
		id = BeachIDPrefix + id
	}

	return &Beach{
		Name:     ngsi.NewTextProperty(name),
		Location: location,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: BeachTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

//WithDescription adds a text property named Deescription to this Beach instance
func (b *Beach) WithDescription(description string) *Beach {
	b.Description = ngsi.NewTextProperty(description)
	return b
}

func (b *Beach) UnmarshalJSON(data []byte) error {
	dto := &beachDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		b.ID = dto.ID
		b.Type = dto.Type
		b.Name = dto.Name
		b.Description = dto.Description

		b.DateCreated = dto.DateCreated
		b.DateModified = dto.DateModified
		b.RefSeeAlso = dto.RefSeeAlso
		b.SameAs = dto.SameAs
		b.WaterTemperature = dto.WaterTemperature

		b.Context = dto.Context

		b.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
	}

	return err
}
