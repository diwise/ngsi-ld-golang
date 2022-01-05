package fiware

import (
	"bytes"
	"encoding/json"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//WeatherObserved is an observation of weather conditions at a certain place and time.
type WeatherObserved struct {
	ngsi.BaseEntity
	DateCreated  *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved ngsi.DateTimeProperty          `json:"dateObserved"`
	Location     geojson.GeoJSONProperty        `json:"location"`
	RefDevice    *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	SnowHeight   *ngsi.NumberProperty           `json:"snowHeight,omitempty"`
	Temperature  *ngsi.NumberProperty           `json:"temperature,omitempty"`
}

type weatherObservedDTO struct {
	ngsi.BaseEntity
	DateCreated  *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved ngsi.DateTimeProperty          `json:"dateObserved"`
	Location     json.RawMessage                `json:"location"`
	RefDevice    *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	SnowHeight   *ngsi.NumberProperty           `json:"snowHeight,omitempty"`
	Temperature  *ngsi.NumberProperty           `json:"temperature,omitempty"`
}

//NewWeatherObserved creates a new instance of WeatherObserved
func NewWeatherObserved(device string, latitude float64, longitude float64, observedAt string) *WeatherObserved {
	dateTimeValue := ngsi.CreateDateTimeProperty(observedAt)
	refDevice := CreateDeviceRelationshipFromDevice(device)

	if refDevice == nil {
		device = "manual"
	}

	id := WeatherObservedIDPrefix + device + ":" + observedAt

	return &WeatherObserved{
		DateObserved: *dateTimeValue,
		Location:     *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		RefDevice:    refDevice,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: WeatherObservedTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

func (wo WeatherObserved) ToGeoJSONFeature(propertyName string, simplified bool) (geojson.GeoJSONFeature, error) {
	g := geojson.NewGeoJSONFeature(wo.ID, wo.Type, wo.Location.GeoPropertyValue())

	if simplified {
		g.SetProperty(propertyName, wo.Location.GeoPropertyValue())
		g.SetProperty("dateObserved", wo.DateObserved.Value.Value)

		if wo.SnowHeight != nil {
			g.SetProperty("snowHeight", wo.SnowHeight.Value)
		}

		if wo.Temperature != nil {
			g.SetProperty("temperature", wo.Temperature.Value)
		}

		if wo.RefDevice != nil {
			g.SetProperty("refDevice", wo.RefDevice.Object)
		}
	} else {
		g.SetProperty(propertyName, wo.Location)
		g.SetProperty("dateObserved", wo.DateObserved)

		g.SetProperty("snowHeight", wo.SnowHeight)
		g.SetProperty("temperature", wo.Temperature)
		g.SetProperty("refDevice", wo.RefDevice)
	}

	return g, nil
}

func (wo *WeatherObserved) UnmarshalJSON(data []byte) error {
	dto := &weatherObservedDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		wo.ID = dto.ID
		wo.Type = dto.Type

		wo.DateCreated = dto.DateCreated
		wo.DateModified = dto.DateModified
		wo.DateObserved = dto.DateObserved
		wo.RefDevice = dto.RefDevice
		wo.SnowHeight = dto.SnowHeight
		wo.Temperature = dto.Temperature

		wo.Context = dto.Context

		wo.Location = *geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
	}

	return err
}
