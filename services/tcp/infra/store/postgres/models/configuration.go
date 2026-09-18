package models

import "time"

type Configuration struct {
	ID            uint64    `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeviceID      uint64    `json:"device_id"`
	CfgHash       uint32    `json:"cfg_hash"`
	CfgData       []byte    `json:"cfg_data"`
	CfgSyncStatus string    `json:"cfg_sync_status" db:"cfg_sync_status"`
	CfgSyncError  *string   `json:"cfg_sync_error" db:"cfg_sync_error"`
	CfgPushedHash uint32    `json:"cfg_pushed_hash" db:"cfg_pushed_hash"`
}
