package fiware

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/matryer/is"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestUnmarshalJSONToAirQualityObserved(t *testing.T) {
	is := is.New(t)
	aqo := AirQualityObserved{}

	err := json.Unmarshal([]byte(airQualityObservedJson), &aqo)

	is.NoErr(err)                         // Could not unmarshal json string into AirQualityObserved
	is.True(aqo.CO2 != nil)               // CO2 should be present
	is.Equal(aqo.CO2.Value, float64(449)) // CO2 value should be 449
}

func TestUnmarshalJSONToWaterConsumptionObserved(t *testing.T) {
	is := is.New(t)
	wco := WaterConsumptionObserved{}

	err := json.Unmarshal([]byte(wcoJson), &wco)

	is.NoErr(err) // Could not unmarshal json string into WaterConsumptionObserved
}

func TestNewWaterConsumptionObserved(t *testing.T) {
	is := is.New(t)
	id := "waterConsumption01"
	observedAt := time.Now().UTC()

	wco := NewWaterConsumptionObserved(id)

	wco.WithConsumption(id, 806040.0, observedAt)

	is.Equal(wco.ID, WaterConsumptionObservedIDPrefix+id)         // Expected wco.ID to match
	is.Equal(wco.WaterConsumption.NumberProperty.Value, 806040.0) // Value of water consumption does not match

}

func TestRoadSegment(t *testing.T) {
	location := [][2]float64{{1.0, 1.0}, {2.0, 2.0}}
	ts := time.Now()
	rs := NewRoadSegment("segid", "segname", "roadid", location, &ts)

	if rs.ID != "urn:ngsi-ld:RoadSegment:segid" {
		t.Errorf("Expectation failed. Road id %s != %s", rs.ID, "segid")
	}

	probability := 0.8
	surfaceType := "snow"
	rs = rs.WithSurfaceType(surfaceType, probability)

	if rs.SurfaceType.Probability != probability {
		t.Errorf("Expectation failed. Surface type probability %f != %f", rs.SurfaceType.Probability, probability)
	}

	if rs.SurfaceType.Value != surfaceType {
		t.Errorf("Expectation failed. Surface type %s != %s", rs.SurfaceType.Value, surfaceType)
	}
}

func TestDeviceModel(t *testing.T) {
	categories := []string{"temperature"}

	deviceModel := NewDeviceModel("id", categories)

	if deviceModel.Category.Value[0] != categories[0] {
		t.Errorf("Expectation failed. Category does not match %s", categories)
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

func TestRoadAccident(t *testing.T) {

	ra := NewRoadAccident(RoadAccidentIDPrefix + "roadAccidentTest")

	ra.AccidentDate = *types.CreateDateTimeProperty("2021-05-23T23:14:16.000Z")

	if ra == nil {
		t.Error("Expectation failed. RoadAccident is empty")
	}

	if ra.AccidentDate.Value.Value != "2021-05-23T23:14:16.000Z" {
		t.Error("Expected dates to be identical.")
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

const wcoJson string = `{
    "id": "urn:ngsi-ld:WaterConsumptionObserved:Consumer01",
    "type": "WaterConsumptionObserved",
    "acquisitionStageFailure": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 0,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmFlowPersistence": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": "Nothing to report",
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmInProgress": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 1,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmMetrology": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 1,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmStopsLeaks": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 0,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmSystem": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 1,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmTamper": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 0,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "alarmWaterQuality": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 0,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "location": {
        "type": "GeoProperty",
        "value": {
            "type": "Point",
            "coordinates": [
                -4.128871,
                50.95822
            ]
        }
    },
    "maxFlow": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 620,
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "E32"
    },
    "minFlow": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 1,
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "E32"
    },
    "moduleTampered": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 1,
        "observedAt": "2021-05-23T23:14:16.000Z"
    },
    "persistenceFlowDuration": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": "3h < 6h",
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "HUR"
    },
    "waterConsumption": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 191051,
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "LTR"
    },
    "@context": [
        "https://raw.githubusercontent.com/easy-global-market/ngsild-api-data-models/master/WaterSmartMeter/jsonld-contexts/waterSmartMeter-compound.jsonld"
    ]
}`

const airQualityObservedJson string = `{
    "id": "urn:ngsi-ld:AirQualityObserved:Madrid-AmbientObserved-28079004-2016-11-30T07:00:00",
    "type": "AirQualityObserved",
    "dateObserved": {
        "type": "Property",
        "value": {
            "@type": "DateTime",
            "@value": "2016-11-30T07:00:00.00Z"
        }
    },
    "airQualityLevel": {
        "type": "Property",
        "value": "moderate"
    },
    "CO2": {
        "type": "Property",
        "value": 449,
        "unitCode": "59"
    },
    "temperature": {
        "type": "Property",
        "value": 12.2
    },
    "NO": {
        "type": "Property",
        "value": 45,
        "unitCode": "GQ"
    },
    "refPointOfInterest": {
        "type": "Relationship",
        "object": "urn:ngsi-ld:PointOfInterest:28079004-Pza.deEspanya"
    },
    "windDirection": {
        "type": "Property",
        "value": 186
    },
    "source": {
        "type": "Property",
        "value": "http://datos.madrid.es"
    },
    "windSpeed": {
        "type": "Property",
        "value": 0.64
    },
    "SO2": {
        "type": "Property",
        "value": 11,
        "unitCode": "GQ"
    },
    "NOx": {
        "type": "Property",
        "value": 139,
        "unitCode": "GQ"
    },
    "location": {
        "type": "GeoProperty",
        "value": {
            "type": "Point",
            "coordinates": [-3.712247222222222, 40.423852777777775]
        }
    },
    "airQualityIndex": {
        "type": "Property",
        "value": 65
    },
    "address": {
        "type": "Property",
        "value": {
            "addressCountry": "ES",
            "addressLocality": "Madrid",
            "streetAddress": "Plaza de Espa\u00f1a",
            "type": "PostalAddress"
        }
    },
    "reliability": {
        "type": "Property",
        "value": 0.7
    },
    "relativeHumidity": {
        "type": "Property",
        "value": 0.54
    },
    "precipitation": {
        "type": "Property",
        "value": 0
    },
    "NO2": {
        "type": "Property",
        "value": 69,
        "unitCode": "GQ"
    },
    "CO_Level": {
        "type": "Property",
        "value": "moderate"
    },
    "@context": [
        "https://schema.lab.fiware.org/ld/context",
        "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"
    ]
}`
