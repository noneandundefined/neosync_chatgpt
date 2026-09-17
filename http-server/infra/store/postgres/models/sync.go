package models

import "time"

type Sync struct {
	ID                   uint64     `json:"id" db:"id"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeviceID             uint64     `json:"device_id" db:"device_id"`
	FirmwareVersion      uint16     `json:"firmware_version" db:"firmware_version"`
	FirmwareVersionUpd   uint16     `json:"firmware_version_upd" db:"firmware_version_upd"`
	FirmwareUpdatedAt    *time.Time `json:"firmware_updated_at" db:"firmware_updated_at"`
	FirmwareUpdateStatus string     `json:"firmware_update_status" db:"firmware_update_status"`
	CfgVersion           uint8      `json:"cfg_version" db:"cfg_version"`
	LastModTime          uint32     `json:"last_mod_time" db:"last_mod_time"`
	CfgHash              uint32     `json:"cfg_hash" db:"cfg_hash"`
}

type SyncData struct {
	FirmwareVersion uint16 `json:"firmware_version"`
	CfgVersion      uint8  `json:"cfg_version"`
	LastModTime     uint32 `json:"last_mod_time"`
	CfgHash         uint32 `json:"cfg_hash"`
}
