package models

import "time"

type Configuration struct {
	ID            uint64    `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	DeviceID      uint64    `json:"device_id" db:"device_id"`
	CfgHash       uint32    `json:"cfg_hash" db:"cfg_hash"`
	CfgData       []byte    `json:"cfg_data" db:"cfg_data"`
	CfgSyncStatus string    `json:"cfg_sync_status" db:"cfg_sync_status"`
	CfgSyncError  *string   `json:"cfg_sync_error" db:"cfg_sync_error"`
	CfgPushedHash uint32    `json:"cfg_pushed_hash" db:"cfg_pushed_hash"`
}

type ConfigurationHistory struct {
	ID            uint64    `json:"id" db:"id"`
	ApplyAt       time.Time `json:"apply_at" db:"apply_at"`
	DeviceID      uint64    `json:"device_id" db:"device_id"`
	CfgHash       uint32    `json:"cfg_hash" db:"cfg_hash"`
	CfgData       []byte    `json:"-" db:"cfg_data"`
	CfgSyncStatus *string   `json:"cfg_sync_status,omitempty" db:"cfg_sync_status"`
	CfgSyncError  *string   `json:"cfg_sync_error,omitempty" db:"cfg_sync_error"`
	Origin        string    `json:"origin" db:"origin"`
}
