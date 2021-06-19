package fiware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Device is a Fiware entity
type Device struct {
	ngsi.BaseEntity
	Value                 *ngsi.TextProperty             `json:"value"`
	DateLastValueReported *ngsi.DateTimeProperty         `json:"dateLastValueReported,omitempty"`
	DateCreated           *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified          *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	Location              *geojson.GeoJSONProperty       `json:"location,omitempty"`
	RefDeviceModel        *ngsi.SingleObjectRelationship `json:"refDeviceModel,omitempty"`
}

type deviceDTO struct {
	ngsi.BaseEntity
	Value                 *ngsi.TextProperty             `json:"value"`
	DateLastValueReported *ngsi.DateTimeProperty         `json:"dateLastValueReported,omitempty"`
	DateCreated           *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified          *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	Location              json.RawMessage                `json:"location,omitempty"`
	RefDeviceModel        *ngsi.SingleObjectRelationship `json:"refDeviceModel,omitempty"`
}

//NewDevice creates a new Device from given ID and Value
func NewDevice(id string, value string) *Device {
	if !strings.HasPrefix(id, DeviceIDPrefix) {
		id = DeviceIDPrefix + id
	}

	return &Device{
		Value: ngsi.NewTextProperty(value),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "Device",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

//CreateDeviceRelationshipFromDevice create a DeviceRelationship from a Device
func CreateDeviceRelationshipFromDevice(device string) *ngsi.SingleObjectRelationship {
	if len(device) == 0 {
		return nil
	}

	if !strings.HasPrefix(device, DeviceIDPrefix) {
		device = DeviceIDPrefix + device
	}

	return ngsi.NewSingleObjectRelationship(device)
}

func (d *Device) UnmarshalJSON(data []byte) error {
	dto := &deviceDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		d.ID = dto.ID
		d.Type = dto.Type

		d.Value = dto.Value
		d.DateLastValueReported = dto.DateLastValueReported

		d.DateCreated = dto.DateCreated
		d.DateModified = dto.DateModified

		if d.Location != nil {
			d.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
		}

		d.RefDeviceModel = dto.RefDeviceModel

		d.Context = dto.Context
	}

	return err
}
