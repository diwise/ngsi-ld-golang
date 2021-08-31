package fiware

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRoadSegment(t *testing.T) {
	location := [][2]float64{{1.0, 1.0}, {2.0, 2.0}}
	ts := time.Now()
	rs := NewRoadSegment("segid", "segname", "roadid", location, &ts)

	if rs.ID != "urn:ngsi-ld:RoadSegment:segid" {
		t.Error(fmt.Sprintf("Expectation failed. Road id %s != %s", rs.ID, "segid"))
	}

	probability := 0.8
	surfaceType := "snow"
	rs = rs.WithSurfaceType(surfaceType, probability)

	if rs.SurfaceType.Probability != probability {
		t.Error(fmt.Sprintf("Expectation failed. Surface type probability %f != %f", rs.SurfaceType.Probability, probability))
	}

	if rs.SurfaceType.Value != surfaceType {
		t.Error(fmt.Sprintf("Expectation failed. Surface type %s != %s", rs.SurfaceType.Value, surfaceType))
	}
}

func TestDeviceModel(t *testing.T) {
	categories := []string{"temperature"}

	deviceModel := NewDeviceModel("id", categories)

	if deviceModel.Category.Value[0] != categories[0] {
		t.Error(fmt.Sprintf("Expectation failed. Category does not match %s", categories))
	}

	// test devicemodel categories are as expected
}

func TestUnmarshalJSONToDeviceHasCorrectLocation(t *testing.T) {
	device := Device{}

	err := json.Unmarshal([]byte(jsonStr), &device)
	if err != nil {
		t.Error("Expectation failed. Could not unmarshal json string into device")
	}

	if device.Location.GetAsPoint().Coordinates != [2]float64{17.3069, 62.3908} {
		t.Error("Expectation failed. Value of location does not match expected value")
	}
}

func TestUnmarshalJSONToDeviceWithEmptyLocation(t *testing.T) {
	device := Device{}

	err := json.Unmarshal([]byte(jsonStrNoLocation), &device)
	if err != nil {
		t.Error("Expectation failed. Could not unmarshal json string into device")
	}

	if device.Location != nil {
		t.Error("Expectation failed. Location should be empty")
	}
}

func TestTrafficFlowObserved(t *testing.T) {
	id := TrafficFlowObservedIDPrefix + "trafficFlowObservedID"
	ts := time.Now().String()

	tfo := NewTrafficFlowObserved(id, ts, 1, 109)

	if tfo == nil {
		t.Error("Expectation failed. TrafficFlowObserved is empty")
	}
}

const jsonStr string = `{
	"@context": [
		  "https://schema.lab.fiware.org/ld/context",
		  "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"
	],
	"id": "urn:ngsi-ld:Device:device",
	"refDeviceModel": {
		  "object": "urn:ngsi-ld:DeviceModel:myDevice",
		  "type": "Relationship"
	},
	"type": "Device",
	"value": {
		  "type": "Property",
		  "value": "38"
	},
	"location": {
		"type": "GeoProperty",
		"value": {
			"coordinates": [
				17.3069,
				62.3908
			],
			"type": "Point"
		}
	}
}`

const jsonStrNoLocation string = `{
	"@context": [
		  "https://schema.lab.fiware.org/ld/context",
		  "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"
	],
	"id": "urn:ngsi-ld:Device:device",
	"refDeviceModel": {
		  "object": "urn:ngsi-ld:DeviceModel:myDevice",
		  "type": "Relationship"
	},
	"type": "Device",
	"value": {
		  "type": "Property",
		  "value": "38"
	}
}`
