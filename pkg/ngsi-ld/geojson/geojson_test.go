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
