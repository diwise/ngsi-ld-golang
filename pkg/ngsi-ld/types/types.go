package types

import "strconv"

//BaseEntity contains the required base properties an Entity must have
type BaseEntity struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Context []string `json:"@context"`
}

//Property contains the mandatory Type property
type Property struct {
	Type string `json:"type"`
}

//DateTimeProperty stores date and time values (surprise, surprise ...)
type DateTimeProperty struct {
	Property
	Value struct {
		Type  string `json:"@type"`
		Value string `json:"@value"`
	} `json:"value"`
}

//CreateDateTimeProperty creates a property from a UTC time stamp
func CreateDateTimeProperty(value string) *DateTimeProperty {
	dtp := &DateTimeProperty{
		Property: Property{Type: "Property"},
	}

	dtp.Value.Type = "DateTime"
	dtp.Value.Value = value

	return dtp
}

//GeoJSONProperty is used to store lat/lon coordinates
type GeoJSONProperty struct {
	Property
	Value struct {
		Type        string     `json:"type"`
		Coordinates [2]float64 `json:"coordinates"`
	} `json:"value"`
}

//CreateGeoJSONPropertyFromWGS84 creates a GeoJSONProperty from a WGS84 coordinate
func CreateGeoJSONPropertyFromWGS84(latitude, longitude float64) GeoJSONProperty {
	p := GeoJSONProperty{
		Property: Property{Type: "GeoProperty"},
	}

	p.Value.Type = "Point"
	p.Value.Coordinates[0] = longitude
	p.Value.Coordinates[1] = latitude

	return p
}

//NumberProperty holds a float64 Value
type NumberProperty struct {
	Property
	Value float64 `json:"value"`
}

//NewNumberProperty is a convenience function for creating NumberProperty instances
func NewNumberProperty(value float64) *NumberProperty {
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    value,
	}
}

func NewNumberPropertyFromInt(value int) *NumberProperty {
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    float64(value),
	}
}

//TextProperty stores values of type text
type TextProperty struct {
	Property
	Value string `json:"value"`
}

func NewNumberPropertyFromString(value string) *NumberProperty {
	number, _ := strconv.ParseFloat(value, 64)
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    number,
	}
}

func NewTextProperty(value string) *TextProperty {
	return &TextProperty{
		Property: Property{Type: "Property"},
		Value:    value,
	}
}

//Relationship is a base type for all types of relationships
type Relationship struct {
	Type string `json:"type"`
}

//DeviceRelationship stores information about an entity's relation to a certain Device
type DeviceRelationship struct {
	Relationship
	Object string `json:"object"`
}

//CreateDeviceRelationshipFromDevice create a DeviceRelationship from a Device
func CreateDeviceRelationshipFromDevice(device string) *DeviceRelationship {
	if len(device) == 0 {
		return nil
	}

	return &DeviceRelationship{
		Relationship: Relationship{Type: "Relationship"},
		Object:       "urn:ngsi-ld:Device:" + device,
	}
}
