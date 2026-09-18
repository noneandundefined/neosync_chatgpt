package device_handler_v1

import (
	"context"
	"neomatica/neosync/infra/store/postgres/models"
)

type CreateDevicePayload struct {
	Imei                          string  `json:"imei" validate:"required,imei"`
	Model                         string  `json:"model" validate:"required"`
	Name                          *string `json:"name" validate:"omitempty"`
	Phone                         *string `json:"phone" validate:"omitempty"`
	NameOrganization              *string `json:"name_organization" validate:"omitempty,max=100"`
	RequestConfigurationOnConnect bool    `json:"request_configuration_on_connect" validate:"omitempty"`
	TurnstileToken                string  `json:"turnstile_token" validate:"required"`
}

type ImportDevicePayload struct {
	Imei       string  `json:"imei" validate:"required,imei"`
	Model      string  `json:"model" validate:"omitempty,max=10"`
	OwnerEmail *string `json:"owner_email" validate:"omitempty,email"`
	GroupName  *string `json:"group_name" validate:"omitempty,max=255"`
}

type ImportError struct {
	Imei  string `json:"imei"`
	Error string `json:"error"`
}

type ImportCheckPayload struct {
	Imeis []string `json:"imeis" validate:"required"`
}

type EditDevicePayload struct {
	Imei                          *string `json:"imei" validate:"omitempty,imei"`
	Model                         *string `json:"model" validate:"omitempty,max=10"`
	Name                          *string `json:"name" validate:"omitempty"`
	Phone                         *string `json:"phone" validate:"omitempty"`
	NameOrganization              *string `json:"name_organization" validate:"omitempty,max=100"`
	RequestConfigurationOnConnect *bool   `json:"request_configuration_on_connect" validate:"omitempty"`
}

type DevicesStateDevicesResponse struct {
	Items   []models.DeviceShort `json:"items"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	Total   int                  `json:"total"`
	HasMore bool                 `json:"has_more"`
}

type DevicesWPResponse struct {
	Items      []models.Device_DeviceConf_Sync `json:"items"`
	Page       int                             `json:"page"`
	Limit      int                             `json:"limit"`
	Total      int                             `json:"total"`
	TotalAll   int                             `json:"total_all"`
	TotalPages int                             `json:"total_pages"`
}

type CommandPayload struct {
	Imei    string `json:"imei" validate:"required,imei"`
	Command string `json:"command" validate:"required,min=2"`
}

type TransferAccountPayload struct {
	IMEIs          []string `json:"imeis" validate:"required"`
	UUID           *string  `json:"uuid" validate:"optional_uuid"`
	TurnstileToken string   `json:"turnstile_token" validate:"required"`
}

type BleTask struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type DeleteDevicesPayload struct {
	IMEIs []string `json:"imeis" validate:"required"`
}

type DeviceDetailedStatusPayload struct {
	IMEIs []string `json:"imeis" validate:"required"`
}
