package models

import "time"

type Device struct {
	ID                  uint64    `json:"id" db:"id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	UserUUID            *string   `json:"user_uuid" db:"user_uuid"`
	Name                *string   `json:"name" db:"name"`
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

type RequestConfigurationOnConnect struct {
	IMEI                          string `json:"imei"`
	RequestConfigurationOnConnect bool   `json:"request_configuration_on_connect"`
}

type Device_DeviceConf_Sync struct {
	/* From devices (devices) */
	Device

	/* From device_confs (device_confs) */
	AutoUpdate                    bool `json:"auto_update"`
	RequestConfigurationOnConnect bool `json:"request_configuration_on_connect" db:"request_configuration_on_connect"`

	/* From syncs (syncs) */
	FirmwareVersion uint16 `json:"firmware_version" db:"firmware_version"`
	CfgVersion      uint8  `json:"cfg_version" db:"cfg_version"`
	LastModTime     uint32 `json:"last_mod_time" db:"last_mod_time"`
	SyncCfgHash     uint32 `json:"sync_cfg_hash" db:"sync_cfg_hash"`

	/* Configuration (configurations) */
	CfgUpdatedAt  time.Time `json:"cfg_updated_at" db:"cfg_updated_at"`
	CfgHash       uint32    `json:"cfg_hash" db:"cfg_hash"`
	CfgData       []byte    `json:"cfg_data" db:"cfg_data"`
	CfgSyncStatus string    `json:"cfg_sync_status" db:"cfg_sync_status"`
	CfgSyncError  *string   `json:"cfg_sync_error" db:"cfg_sync_error"`
	CfgPushedHash uint32    `json:"cfg_pushed_hash" db:"cfg_pushed_hash"`

	/* From user_cores (user_cores) */
	UserEmail                         *string `json:"user_email" db:"user_email"`
	ForceNeosyncConfigurationPriority bool    `json:"force_neosync_configuration_priority" db:"force_neosync_configuration_priority"`

	/* From groups (groups) and group_user_devices (gud) */
	GroupName *string `json:"group_name" db:"group_name"`
	GroupID   *uint64 `json:"group_id" db:"group_id"`
}
