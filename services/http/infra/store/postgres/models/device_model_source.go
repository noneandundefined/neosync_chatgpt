package models

import "time"

type DeviceModelSource struct {
	ID               uint64    `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	DeviceModel      string    `json:"device_model" db:"device_model"`
	FirmwareUrl      string    `json:"firmware_url" db:"firmware_url"`
	FirmwareCacheKey string    `json:"firmware_cache_key" db:"firmware_cache_key"`
}

type DeviceModelResult struct {
	DeviceModel string `db:"device_model"`
}

type FirmwareUrlResult struct {
	FirmwareUrl string `db:"firmware_url"`
}
