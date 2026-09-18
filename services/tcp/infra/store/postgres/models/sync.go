package models

import "time"

type Sync struct {
	ID                 uint64     `json:"id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeviceID           uint64     `json:"device_id"`
	FirmwareVersion    uint16     `json:"firmware_version"`
	FirmwareVersionUpd uint16     `json:"firmware_version_upd"`
	FirmwareUpdatedAt  *time.Time `json:"firmware_updated_at" db:"firmware_updated_at"`
	CfgVersion         uint8      `json:"cfg_version"`
	LastModTime        uint32     `json:"last_mod_time"`
	CfgHash            uint32     `json:"cfg_hash"`
}

type SyncData struct {
	FirmwareVersion uint16 `json:"firmwareVersion"`
	CfgVersion      uint8  `json:"cfgVersion"`
	LastModTime     uint32 `json:"lastModTime"`
	CfgHash         uint32 `json:"cfgHash"`
}
