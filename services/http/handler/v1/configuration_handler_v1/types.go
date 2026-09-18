package configuration_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/adm"
)

type ConfigurationParsed struct {
	DeviceImei   string                         `json:"device_imei" validate:"imei"`
	CfgHash      uint32                         `json:"cfg_hash"`
	ConfigParsed []adm.FieldConfigurationParsed `json:"config_parsed"`
	Schema       []constants.FieldSchema        `json:"schema"`
	Timestamp    int64                          `json:"timestamp"`
}

type ConfigurationSectionParsed struct {
	DeviceImei   string                         `json:"device_imei" validate:"imei"`
	CfgHash      uint32                         `json:"cfg_hash"`
	Section      string                         `json:"section"`
	ConfigParsed []adm.FieldConfigurationParsed `json:"config_parsed"`
	Schema       []constants.FieldSchema        `json:"schema"`
	Timestamp    int64                          `json:"timestamp"`
}

type ConfigurationRaw struct {
	DeviceImei string `json:"device_imei" validate:"imei"`
	CfgHash    uint32 `json:"cfg_hash"`
	ConfigHex  string `json:"config_hex"`
	Timestamp  int64  `json:"timestamp"`
}

type ConfigurationDraftGet struct {
	DeviceImei string                  `json:"device_imei" validate:"imei"`
	CfgHash    uint32                  `json:"cfg_hash"`
	Section    string                  `json:"section"`
	Changes    map[string]any          `json:"changes"`
	Schema     []constants.FieldSchema `json:"schema"`
	Timestamp  int64                   `json:"timestamp"`
}

type ConfigurationDraftInsert struct {
	DeviceImei string         `json:"device_imei" validate:"required,imei"`
	CfgHash    uint32         `json:"cfg_hash"`
	Changes    map[string]any `json:"changes"`
	Timestamp  int64          `json:"timestamp" validate:"required"`
}

type ConfigurationExportLLSTarirationRaw struct {
	Id    int     `json:"id"`
	Value float64 `json:"value"`
	Level float64 `json:"level"`
}

type ConfigurationExportLLSTarirationPayload struct {
	Tariration map[int][]ConfigurationExportLLSTarirationRaw `json:"tariration" validate:"required"`
}

type CsvLLSTarirationLabels struct {
	SensorHeader string
	ColIndex     string
	ColValue     string
	ColLevel     string
}

type ConfigurationImportByHistoryPayload struct {
	CfgHash uint32 `json:"cfg_hash" validate:"required"`
}

type ConfigurationChangePassPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ConfigurationTemplateChangesData struct {
	CfgHash uint32         `json:"cfg_hash"`
	Changes map[string]any `json:"changes"`
}

type ConfigurationTemplateCreate struct {
	Name         string         `json:"name" validate:"required,max=255"`
	Model        string         `json:"model,omitempty"`
	TypeOfSaving string         `json:"type_of_saving" validate:"required,oneof=modified_ones full"`
	CfgHash      uint32         `json:"cfg_hash"`
	Changes      map[string]any `json:"changes"`
}

type ConfigurationTemplateUpdate struct {
	Name         string         `json:"name" validate:"omitempty,max=255"`
	TypeOfSaving string         `json:"type_of_saving" validate:"omitempty,oneof=modified_ones full"`
	CfgHash      uint32         `json:"cfg_hash"`
	Changes      map[string]any `json:"changes"`
}

type ConfigurationTemplateSectionParsed struct {
	CfgHash      uint32                         `json:"cfg_hash"`
	Section      string                         `json:"section"`
	ConfigParsed []adm.FieldConfigurationParsed `json:"config_parsed"`
	Schema       []constants.FieldSchema        `json:"schema"`
	Timestamp    int64                          `json:"timestamp"`
}

type ConfigurationTemplatesWPResponse struct {
	Items      []models.ConfigurationTemplate `json:"items"`
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	Total      int                            `json:"total"`
	TotalAll   int                            `json:"total_all"`
	TotalPages int                            `json:"total_pages"`
}

type ConfigurationTemplateCreateResponse struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}
