package telemetry

type ADM20Info struct {
	EmptinessMask *int `json:"emptiness_mask"`
	CardUId       *int `json:"card_uid"`
	TagId         *int `json:"tag_id"`
	TagSn         *int `json:"tag_sn"`
	TagRssi       *int `json:"tag_rssi"`
	TagAdc        *int `json:"tag_adc"`
	TagPeriod     *int `json:"tag_period"`
}

type BLESensor struct {
	RSSI      int     `json:"rssi,omitempty"`
	Voltage   float64 `json:"voltage,omitempty"`
	LMT       int     `json:"lmt,omitempty"`
	FuelLevel int     `json:"fuel_level,omitempty"`
	Temp      int     `json:"temp,omitempty"`
}

type Telemetry struct {
	GPSSatellites     int                  `json:"gps_satellites"`
	GlonassSatellites int                  `json:"glonass_satellites"`
	Lat               *float64             `json:"lat,omitempty"`
	Lon               *float64             `json:"lon,omitempty"`
	DateTime          uint32               `json:"date_time"`
	VPower            int                  `json:"v_power"`
	VBattery          int                  `json:"v_battery"`
	InA0              int                  `json:"in_a0"`
	InA1              int                  `json:"in_a1"`
	InD0              uint32               `json:"in_d0"`
	InD1              uint32               `json:"in_d1"`
	ADM20INFO         *ADM20Info           `json:"adm20info,omitempty"`
	BLESENSORINFO     map[string]BLESensor `json:"blesensorinfo,omitempty"`
	FUELINFO          map[string]BLESensor `json:"fuelinfo,omitempty"`
}
