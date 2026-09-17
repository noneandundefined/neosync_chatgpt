package models

import "time"

func (d *Device) GetUserUUID() *string {
	return d.UserUUID
}

/* Device основная таблица терминалов */
type Device struct {
	ID                  uint64    `json:"id" db:"id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	UserUUID            *string   `json:"user_uuid" db:"user_uuid"`
	OwnerUUID           *string   `json:"owner_uuid" db:"owner_uuid"`
	Name                *string   `json:"name" db:"name"`
	Phone               *string   `json:"phone" db:"phone"`
	NameOrganization    *string   `json:"name_organization" db:"name_organization"`
	IMEI                string    `json:"imei" db:"imei"`
	Status              bool      `json:"status" db:"status"`
	GroupName           *string   `json:"group_name" db:"group_name"`
	DeviceModel         *string   `json:"device_model" db:"device_model"`
	DeviceExtendedModel *string   `json:"device_extended_model" db:"device_extended_model"`
	Activated           *bool     `json:"activated" db:"activated"`
}

/* DeviceConf основная таблица доп. инфы терминалов */
type DeviceConf struct {
	ID                            uint64    `json:"id" db:"id"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
	DeviceID                      uint64    `json:"device_id" db:"device_id"`
	Password                      string    `json:"password" db:"password"`
	RequestConfigurationOnConnect bool      `json:"request_configuration_on_connect" db:"request_configuration_on_connect"`
}

/* DeviceAndDeviceConf таблица из двух моделей терминала */
type DeviceWithDeviceConf struct {
	Device
	Password                      string `json:"password" db:"password"`
	RequestConfigurationOnConnect bool   `json:"request_configuration_on_connect" db:"request_configuration_on_connect"`
}

type Device_DeviceConf_Sync struct {
	/* From devices (devices) */
	Device

	/* From device_confs (device_confs) */
	RequestConfigurationOnConnect bool `json:"request_configuration_on_connect" db:"request_configuration_on_connect"`

	/* From configurations */
	ConfigurationUpdatedAt  *time.Time `json:"configuration_updated_at" db:"configuration_updated_at"`
	ConfigurationSyncStatus *string    `json:"configuration_sync_status" db:"configuration_sync_status"`
	ConfigurationSyncError  *string    `json:"configuration_sync_error" db:"configuration_sync_error"`

	/* From syncs (syncs) */
	FirmwareVersion      uint16     `json:"firmware_version" db:"firmware_version"`
	FirmwareVersionUpd   uint16     `json:"firmware_version_upd" db:"firmware_version_upd"`
	FirmwareUpdatedAt    *time.Time `json:"firmware_updated_at" db:"firmware_updated_at"`
	FirmwareUpdateStatus string     `json:"firmware_update_status" db:"firmware_update_status"`
	CfgVersion           uint8      `json:"cfg_version" db:"cfg_version"`
	LastModTime          uint32     `json:"last_mod_time" db:"last_mod_time"`
	CfgHash              uint32     `json:"cfg_hash" db:"cfg_hash"`

	/* From user_cores (user_cores) */
	UserEmail *string `json:"user_email" db:"user_email"`

	/* From groups (groups) and group_user_devices (gud) */
	GroupName *string `json:"group_name" db:"group_name"`
	GroupID   *uint64 `json:"group_id" db:"group_id"`
}

/* DeviceStatus статус терминала */
type DeviceStatus struct {
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	Status              bool      `json:"status" db:"status"`
	DeviceExtendedModel *string   `json:"device_extended_model" db:"device_extended_model"`
}

type DeviceDetailedStatus struct {
	ID                      uint64  `json:"id" db:"id"`
	Imei                    string  `json:"imei" db:"imei"`
	Status                  bool    `json:"status" db:"status"`
	DeviceModel             *string `json:"device_model" db:"device_model"`
	ConfigurationSyncStatus *string `json:"configuration_sync_status" db:"configuration_sync_status"`
	FirmwareVersion         *uint16 `json:"firmware_version" db:"firmware_version"`
}

type UpdateDevice struct {
	Imei                          *string `json:"imei" db:"imei"`
	DeviceModel                   *string `json:"model" db:"model"`
	Name                          *string `json:"name" db:"name"`
	Phone                         *string `json:"phone" db:"phone"`
	NameOrganization              *string `json:"name_organization" db:"name_organization"`
	RequestConfigurationOnConnect *bool   `json:"request_configuration_on_connect" db:"request_configuration_on_connect"`
}
