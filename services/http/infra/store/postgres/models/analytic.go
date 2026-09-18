package models

import "time"

type AnalyticRecord struct {
	ID           uint64    `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UserUUID     string    `db:"user_uuid" json:"user_uuid"`
	DeviceID     uint64    `db:"device_id" json:"device_id"`
	CfgHash      uint32    `db:"cfg_hash" json:"cfg_hash"`
	RMQPublished bool      `db:"rmq_published" json:"rmq_published"`
	DeviceOnline bool      `db:"device_online" json:"device_online"`
	ErrorReason  *string   `db:"error_reason" json:"error_reason,omitempty"`
}

type DatabaseSchemaColumn struct {
	ColumnName    string  `json:"column_name" db:"column_name"`
	DataType      string  `json:"data_type" db:"data_type"`
	IsNullable    string  `json:"is_nullable" db:"is_nullable"`
	ColumnDefault *string `json:"column_default" db:"column_default"`
}

// type AnalyticsConfigProfile struct {
// 	ID                 uint64    `db:"id" json:"id"`
// 	CreatedAt          time.Time `db:"created_at" json:"created_at"`
// 	UserUUID           string    `db:"user_uuid" json:"user_uuid"`
// 	DeviceID           uint64    `db:"device_id" json:"device_id"`
// 	CfgHash            uint32    `db:"cfg_hash" json:"cfg_hash"`
// 	ServerSchemaBucket string    `db:"server_schema_bucket" json:"server_schema_bucket"`
// 	SensorsBucket      string    `db:"sensors_bucket" json:"sensors_bucket"`
// 	ProtocolBucket     string    `db:"protocol_bucket" json:"protocol_bucket"`
// 	PeriodBucket       string    `db:"period_bucket" json:"period_bucket"`

// 	OnewireTempSensorCount int16  `db:"onewire_temp_sensor_count" json:"onewire_temp_sensor_count"`
// 	IbuttonEnabled         *int16 `db:"ibutton_enabled" json:"ibutton_enabled,omitempty"`
// 	/* Заполняется БД (GENERATED); при INSERT не передаётся */
// 	HasOnewirePeripheral *bool `db:"has_onewire_peripheral" json:"has_onewire_peripheral,omitempty"`

// 	DeviceMode       *int16 `db:"device_mode" json:"device_mode,omitempty"`
// 	TrafficProtocol  *int16 `db:"traffic_protocol" json:"traffic_protocol,omitempty"`
// 	PointsPeriodMove *int16 `db:"points_period_move" json:"points_period_move,omitempty"`
// 	PointsPeriodHold *int16 `db:"points_period_hold" json:"points_period_hold,omitempty"`
// 	ApnNonEmptyCount int16  `db:"apn_non_empty_count" json:"apn_non_empty_count"`
// 	AuthPhoneCount   int16  `db:"auth_phone_count" json:"auth_phone_count"`
// }

type AnalyticsUser struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserUUID  string    `db:"user_uuid" json:"user_uuid"`
	IsOnline  bool      `db:"is_online" json:"is_online"`
}

type AnalyticQueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}
