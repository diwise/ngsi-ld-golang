package fiware

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
)

type WCOTextProperty struct {
	ngsi.TextProperty
	ObservedBy ngsi.SingleObjectRelationship `json:"observedBy,omitempty"`
	ObservedAt time.Time                     `json:"observedAt,omitempty"`
	UnitCode   string                        `json:"unitCode,omitempty"`
}

type WCONumberProperty struct {
	ngsi.NumberProperty
	ObservedBy ngsi.SingleObjectRelationship `json:"observedBy,omitempty"`
	ObservedAt time.Time                     `json:"observedAt,omitempty"`
	UnitCode   string                        `json:"unitCode,omitempty"`
}

//WaterConsumptionObserved is a fiware entity
type WaterConsumptionObserved struct {
	ngsi.BaseEntity
	AlarmFlowPersistence    *WCOTextProperty         `json:"alarmFlowPersistence,omitempty"`
	AlarmInProgress         *WCONumberProperty       `json:"alarmInProgress,omitempty"`
	AlarmMetrology          *WCONumberProperty       `json:"alarmMetrology,omitempty"`
	AlarmStopsLeaks         *WCONumberProperty       `json:"alarmStopsLeaks,omitempty"`
	AlarmSystem             *WCONumberProperty       `json:"alarmSystem,omitempty"`
	AlarmTamper             *WCONumberProperty       `json:"alarmTamper,omitempty"`
	AlarmWaterQuality       *WCONumberProperty       `json:"alarmWaterQuality,omitempty"`
	MaxFlow                 *WCONumberProperty       `json:"maxFlow,omitempty"`
	MinFlow                 *WCONumberProperty       `json:"minFlow,omitempty"`
	ModuleTampered          *WCONumberProperty       `json:"moduleTampered,omitempty"`
	PersistenceFlowDuration *WCOTextProperty         `json:"persistenceFlowDuration,omitempty"`
	WaterConsumption        *WCONumberProperty       `json:"waterConsumption,omitempty"`
	Location                *geojson.GeoJSONProperty `json:"location,omitempty"`
}

type waterConsumptionObservedDTO struct {
	ngsi.BaseEntity
	AlarmFlowPersistence    *WCOTextProperty   `json:"alarmFlowPersistence,omitempty"`
	AlarmInProgress         *WCONumberProperty `json:"alarmInProgress,omitempty"`
	AlarmMetrology          *WCONumberProperty `json:"alarmMetrology,omitempty"`
	AlarmStopsLeaks         *WCONumberProperty `json:"alarmStopsLeaks,omitempty"`
	AlarmSystem             *WCONumberProperty `json:"alarmSystem,omitempty"`
	AlarmTamper             *WCONumberProperty `json:"alarmTamper,omitempty"`
	AlarmWaterQuality       *WCONumberProperty `json:"alarmWaterQuality,omitempty"`
	MaxFlow                 *WCONumberProperty `json:"maxFlow,omitempty"`
	MinFlow                 *WCONumberProperty `json:"minFlow,omitempty"`
	ModuleTampered          *WCONumberProperty `json:"moduleTampered,omitempty"`
	PersistenceFlowDuration *WCOTextProperty   `json:"persistenceFlowDuration,omitempty"`
	WaterConsumption        *WCONumberProperty `json:"waterConsumption,omitempty"`
	Location                json.RawMessage    `json:"location,omitempty"`
}

//NewWaterConsumptionObserved creates a new instance of WaterConsumptionObserved
func NewWaterConsumptionObserved(id string) *WaterConsumptionObserved {
	if !strings.HasPrefix(id, WaterConsumptionObservedIDPrefix) {
		id = WaterConsumptionObservedIDPrefix + id
	}

	wco := &WaterConsumptionObserved{
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "WaterConsumptionObserved",
			Context: []string{
				"https://raw.githubusercontent.com/easy-global-market/ngsild-api-data-models/master/WaterSmartMeter/jsonld-contexts/waterSmartMeter-compound.jsonld",
			},
		},
	}

	return wco
}

func (wco *WaterConsumptionObserved) UnmarshalJSON(data []byte) error {
	dto := &waterConsumptionObservedDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		wco.ID = dto.ID
		wco.Type = dto.Type
		wco.Context = dto.Context

		if dto.AlarmFlowPersistence != nil {
			wco.AlarmFlowPersistence = dto.AlarmFlowPersistence
		}

		if dto.AlarmInProgress != nil {
			wco.AlarmInProgress = dto.AlarmInProgress
		}

		if dto.AlarmMetrology != nil {
			wco.AlarmMetrology = dto.AlarmMetrology
		}

		if dto.AlarmStopsLeaks != nil {
			wco.AlarmStopsLeaks = dto.AlarmStopsLeaks
		}

		if dto.AlarmSystem != nil {
			wco.AlarmSystem = dto.AlarmSystem
		}

		if dto.AlarmTamper != nil {
			wco.AlarmTamper = dto.AlarmTamper
		}

		if dto.AlarmWaterQuality != nil {
			wco.AlarmWaterQuality = dto.AlarmWaterQuality
		}

		if dto.MaxFlow != nil {
			wco.MaxFlow = dto.MaxFlow
		}

		if dto.MinFlow != nil {
			wco.MinFlow = dto.MinFlow
		}

		if dto.ModuleTampered != nil {
			wco.ModuleTampered = dto.ModuleTampered
		}

		if dto.PersistenceFlowDuration != nil {
			wco.PersistenceFlowDuration = dto.PersistenceFlowDuration
		}

		if dto.WaterConsumption != nil {
			wco.WaterConsumption = dto.WaterConsumption
		}

		if dto.Location != nil {
			wco.Location = geojson.CreateGeoJSONPropertyFromJSON(dto.Location)
		}
	}

	return err
}
