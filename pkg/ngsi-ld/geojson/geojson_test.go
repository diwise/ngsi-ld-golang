package geojson

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestLineString(t *testing.T) {
	is := is.New(t)

	p := CreateGeoJSONPropertyFromLineString([][]float64{{12.0, 14.0}, {13.0, 15.0}})
	is.Equal(p.Value.GeoPropertyType(), "LineString")

	jsonBytes, err := json.Marshal(p)
	is.NoErr(err)

	p2 := CreateGeoJSONPropertyFromJSON(jsonBytes)
	is.Equal(p2.Value.GeoPropertyType(), "LineString")
}

func TestMultiPolygon(t *testing.T) {
	is := is.New(t)

	p := CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{{{{12.0, 14.0}, {13.0, 15.0}}}})
	is.Equal(p.Value.GeoPropertyType(), "MultiPolygon")

	jsonBytes, err := json.Marshal(p)
	is.NoErr(err)

	p2 := CreateGeoJSONPropertyFromJSON(jsonBytes)
	is.Equal(p2.Value.GeoPropertyType(), "MultiPolygon")
}

func TestWGS84(t *testing.T) {
	is := is.New(t)

	p := CreateGeoJSONPropertyFromWGS84(12.0, 14.0)
	is.Equal(p.Value.GeoPropertyType(), "Point")

	jsonBytes, err := json.Marshal(p)
	is.NoErr(err)

	p2 := CreateGeoJSONPropertyFromJSON(jsonBytes)
	is.Equal(p2.Value.GeoPropertyType(), "Point")
}

func TestUnpackGeoJSONFeatureCollection(t *testing.T) {
	is := is.New(t)
	var unpackedCount int32

	UnpackGeoJSONToCallback([]byte(beachFeatureCollectionJSON), func(GeoJSONFeature) error {
		unpackedCount++
		return nil
	})

	is.Equal(unpackedCount, int32(1))
}

func TestUnpackGeoJSONFeature(t *testing.T) {
	is := is.New(t)
	var unpackedCount int32

	UnpackGeoJSONToCallback([]byte(beachFeatureJSON), func(GeoJSONFeature) error {
		unpackedCount++
		return nil
	})

	is.Equal(unpackedCount, int32(1))
}

const beachFeatureJSON string = `{"id":"urn:ngsi-ld:Beach:42","type": "Feature",
	"geometry": {
		"type": "MultiPolygon",
		"coordinates": [[[
			[16.826877016818194,62.371366230256456],[16.82746858045308,62.37197792385098],
			[16.826075957396505,62.37229386059263],[16.825800236618605,62.37160561482045],
			[16.826877016818194,62.371366230256456]
			]]]
	},
	"properties": {
	  "description": "En fin liten strand.",
	  "location": {
		"type": "MultiPolygon",
		"coordinates": [[[
			  [16.826877016818194,62.371366230256456],[16.82746858045308,62.37197792385098],
			  [16.826075957396505,62.37229386059263],[16.825800236618605,62.37160561482045],
			  [16.826877016818194,62.371366230256456]
			]]]
	  },
	  "name": "Stranden",
	  "refSeeAlso": [
		"urn:ngsi-ld:Device:tempsensor-19"
	  ],
	  "type": "Beach"
	}}`

const beachFeatureCollectionJSON string = `{"type": "FeatureCollection","features": [` + beachFeatureJSON + `]}`
