package fiware

import (
	"bytes"
	"encoding/json"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//AirQualityObserved is intended to represent an observation of air quality conditions at a certain place and time.
type AirQualityObserved struct {
	ngsi.BaseEntity
	DateCreated        *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified       *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved       ngsi.TextProperty              `json:"dateObserved"`
	Location           geojson.GeoJSONProperty        `json:"location"`
	RefDevice          *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	RefPointOfInterest *ngsi.SingleObjectRelationship `json:"refPointOfInterest,omitempty"`
	AreaServed         *ngsi.TextProperty             `json:"areaServed,omitempty"`
	Temperature        *ngsi.NumberProperty           `json:"temperature,omitempty"`
	CO2                *ngsi.NumberProperty           `json:"CO2,omitempty"`
	RelativeHumidity   *ngsi.NumberProperty           `json:"relativeHumidity,omitempty"`
}

type airQualityDTO struct {
	ngsi.BaseEntity
	DateCreated        *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified       *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved       ngsi.TextProperty              `json:"dateObserved"`
	Location           json.RawMessage                `json:"location"`
	RefDevice          *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	RefPointOfInterest *ngsi.SingleObjectRelationship `json:"refPointOfInterest,omitempty"`
	AreaServed         *ngsi.TextProperty             `json:"areaServed,omitempty"`
	Temperature        *ngsi.NumberProperty           `json:"temperature,omitempty"`
	CO2                *ngsi.NumberProperty           `json:"CO2,omitempty"`
	RelativeHumidity   *ngsi.NumberProperty           `json:"relativeHumidity,omitempty"`
}

//NewAirQualityObserved creates a new instance of AirQualityObserved
func NewAirQualityObserved(device string, latitude float64, longitude float64, observedAt string) *AirQualityObserved {
	dateTimeValue := ngsi.NewTextProperty(observedAt)
	refDevice := CreateDeviceRelationshipFromDevice(device)

	id := AirQualityObservedIDPrefix + device + ":" + observedAt

	return &AirQualityObserved{
		DateObserved: *dateTimeValue,
		Location:     *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		RefDevice:    refDevice,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: AirQualityObservedTypeName,
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

func (aqo AirQualityObserved) ToGeoJSONFeature(propertyName string, simplified bool) (geojson.GeoJSONFeature, error) {
	g := geojson.NewGeoJSONFeature(aqo.ID, aqo.Type, aqo.Location.GeoPropertyValue())

	if simplified {
		g.SetProperty(propertyName, aqo.Location.GeoPropertyValue())
		g.SetProperty("dateObserved", aqo.DateObserved.Value)

		if aqo.AreaServed != nil {
			g.SetProperty("areaServed", aqo.AreaServed.Value)
		}

		if aqo.CO2 != nil {
			g.SetProperty("CO2", aqo.CO2.Value)
		}

		if aqo.RelativeHumidity != nil {
			g.SetProperty("relativeHumidity", aqo.RelativeHumidity.Value)
		}

		if aqo.Temperature != nil {
			g.SetProperty("temperature", aqo.Temperature.Value)
		}

		if aqo.RefDevice != nil {
			g.SetProperty("refDevice", aqo.RefDevice.Object)
		}

		if aqo.RefPointOfInterest != nil {
			g.SetProperty("refPointOfInterest", aqo.RefPointOfInterest.Object)
		}
	} else {
		g.SetProperty(propertyName, aqo.Location)
		g.SetProperty("dateObserved", aqo.DateObserved)

		g.SetProperty("CO2", aqo.CO2)
		g.SetProperty("relativeHumidity", aqo.RelativeHumidity)
		g.SetProperty("temperature", aqo.Temperature)
		g.SetProperty("refDevice", aqo.RefDevice)
		g.SetProperty("refPointOfInterest", aqo.RefPointOfInterest)
		g.SetProperty("areaServed", aqo.AreaServed)
	}

	return g, nil
}

func (aqo *AirQualityObserved) UnmarshalJSON(data []byte) error {
	dto := &airQualityDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		aqo.ID = dto.ID
		aqo.Type = dto.Type

		aqo.DateCreated = dto.DateCreated
		aqo.DateModified = dto.DateModified
		aqo.DateObserved = dto.DateObserved
		aqo.RefDevice = dto.RefDevice
		aqo.RefPointOfInterest = dto.RefPointOfInterest
		aqo.AreaServed = dto.AreaServed
		aqo.CO2 = dto.CO2
		aqo.RelativeHumidity = dto.RelativeHumidity
		aqo.Temperature = dto.Temperature

		aqo.Context = dto.Context

		location := geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
		if location != nil {
			aqo.Location = *location
		}
	}

	return err
}

func (aqo AirQualityObserved) WithCO2(value float64) *AirQualityObserved {
	aqo.CO2 = ngsi.NewNumberPropertyWithUnitCode(value, "59")
	return &aqo
}

func (aqo AirQualityObserved) WithRelativeHumidity(value float64) *AirQualityObserved {
	aqo.RelativeHumidity = ngsi.NewNumberProperty(value)
	return &aqo
}

func (aqo AirQualityObserved) WithTemperature(value float64) *AirQualityObserved {
	aqo.Temperature = ngsi.NewNumberPropertyWithUnitCode(value, "CEL")
	return &aqo
}
