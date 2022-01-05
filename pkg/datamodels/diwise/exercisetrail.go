package diwise

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//ExerciseTrail is a Diwise entity
type ExerciseTrail struct {
	ngsi.BaseEntity
	Name                *ngsi.TextProperty            `json:"name"`
	Description         *ngsi.TextProperty            `json:"description"`
	RefTrailSegment     *ngsi.MultiObjectRelationship `json:"refTrailSegment,omitempty"`
	Location            geojson.GeoJSONGeometry       `json:"location,omitempty"`
	Length              *ngsi.NumberProperty          `json:"length,omitempty"`
	Difficulty          *ngsi.TextProperty            `json:"difficulty,omitempty"`
	RefSeeAlso          *ngsi.MultiObjectRelationship `json:"refSeeAlso,omitempty"`
	DateCreated         *ngsi.DateTimeProperty        `json:"dateCreated,omitempty"`
	DateModified        *ngsi.DateTimeProperty        `json:"dateModified,omitempty"`
	DateLastPreparation *ngsi.DateTimeProperty        `json:"dateLastPreparation,omitempty"`
	Responsible         *ngsi.TextProperty            `json:"responsible,omitempty"`
	Source              *ngsi.TextProperty            `json:"source,omitempty"`
	Status              *ngsi.TextProperty            `json:"status,omitempty"`
}

type exerciseTrailDTO struct {
	ngsi.BaseEntity
	Name                *ngsi.TextProperty            `json:"name"`
	Description         *ngsi.TextProperty            `json:"description"`
	RefTrailSegment     *ngsi.MultiObjectRelationship `json:"refTrailSegment,omitempty"`
	Location            json.RawMessage               `json:"location,omitempty"`
	Length              *ngsi.NumberProperty          `json:"length,omitempty"`
	Difficulty          *ngsi.TextProperty            `json:"difficulty,omitempty"`
	RefSeeAlso          *ngsi.MultiObjectRelationship `json:"refSeeAlso,omitempty"`
	DateCreated         *ngsi.DateTimeProperty        `json:"dateCreated,omitempty"`
	DateModified        *ngsi.DateTimeProperty        `json:"dateModified,omitempty"`
	DateLastPreparation *ngsi.DateTimeProperty        `json:"dateLastPreparation,omitempty"`
	Responsible         *ngsi.TextProperty            `json:"responsible,omitempty"`
	Source              *ngsi.TextProperty            `json:"source,omitempty"`
	Status              *ngsi.TextProperty            `json:"status,omitempty"`
}

//NewRoad creates a new instance of Road
func NewExerciseTrail(id string, trailName string, length float64, description string, location geojson.GeoJSONGeometry) *ExerciseTrail {
	if !strings.HasPrefix(id, ExerciseTrailIDPrefix) {
		id = ExerciseTrailIDPrefix + id
	}

	trail := &ExerciseTrail{
		Name:        ngsi.NewTextProperty(trailName),
		Description: ngsi.NewTextProperty(description),
		Location:    location,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: ExerciseTrailTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}

	if length > 0.1 {
		trail.Length = ngsi.NewNumberProperty(length)
	}

	return trail
}

func (t ExerciseTrail) ToGeoJSONFeature(propertyName string, simplified bool) (geojson.GeoJSONFeature, error) {
	g := geojson.NewGeoJSONFeature(t.ID, t.Type, t.Location.GeoPropertyValue())

	if simplified {
		g.SetProperty("name", t.Name.Value)
		g.SetProperty(propertyName, t.Location.GeoPropertyValue())

		if t.Description != nil {
			g.SetProperty("description", t.Description.Value)
		}

		if t.RefSeeAlso != nil {
			g.SetProperty("refSeeAlso", t.RefSeeAlso.Object)
		}

		if t.RefTrailSegment != nil {
			g.SetProperty("refTrailSegment", t.RefTrailSegment.Object)
		}

		if t.Length != nil {
			g.SetProperty("length", t.Length.Value)
		}

		if t.Difficulty != nil {
			g.SetProperty("difficulty", t.Difficulty.Value)
		}

		if t.DateCreated != nil {
			g.SetProperty("dateCreated", t.DateCreated.Value)
		}

		if t.DateModified != nil {
			g.SetProperty("dateModified", t.DateModified.Value)
		}

		if t.DateLastPreparation != nil {
			g.SetProperty("dateLastPreparation", t.DateLastPreparation.Value)
		}

		if t.Responsible != nil {
			g.SetProperty("responsible", t.Responsible.Value)
		}

		if t.Source != nil {
			g.SetProperty("source", t.Source.Value)
		}

		if t.Status != nil {
			g.SetProperty("status", t.Status.Value)
		}
	} else {
		g.SetProperty("name", t.Name)
		g.SetProperty(propertyName, t.Location)

		g.SetProperty("description", t.Description)
		g.SetProperty("refSeeAlso", t.RefSeeAlso)

		g.SetProperty("refTrailSegment", t.RefTrailSegment)
		g.SetProperty("length", t.Length)
		g.SetProperty("difficulty", t.Difficulty)
		g.SetProperty("dateCreated", t.DateCreated)
		g.SetProperty("dateModified", t.DateModified)
		g.SetProperty("dateLastPreparation", t.DateLastPreparation)
		g.SetProperty("responsible", t.Responsible)
		g.SetProperty("source", t.Source)
		g.SetProperty("status", t.Status)
	}

	return g, nil
}

func (t *ExerciseTrail) UnmarshalJSON(data []byte) error {
	dto := &exerciseTrailDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		t.ID = dto.ID
		t.Type = dto.Type
		t.Name = dto.Name
		t.Description = dto.Description

		t.DateCreated = dto.DateCreated
		t.DateModified = dto.DateModified
		t.RefSeeAlso = dto.RefSeeAlso
		t.RefTrailSegment = dto.RefTrailSegment
		t.Length = dto.Length
		t.Difficulty = dto.Difficulty
		t.DateLastPreparation = dto.DateLastPreparation
		t.Responsible = dto.Responsible
		t.Source = dto.Source
		t.Status = dto.Status

		t.Context = dto.Context

		t.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
	}

	return err
}
