package firmware_handler_v1

import (
	"time"
)

type Fota struct {
	CreatedAt   time.Time `json:"created_at"`
	DeviceModel string    `json:"device_model"`
	Firmware    string    `json:"firmware"`
	ReleaseDate string    `json:"release_date"`
	Type        string    `json:"type"`
}

type FotasWPResponse struct {
	Items      []Fota `json:"items"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
}
