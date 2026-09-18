package models

import (
	"encoding/json"
	"time"
)

type Company struct {
	ID                      uint64          `json:"id" db:"id"`
	CreatedAt               time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at" db:"updated_at"`
	UserUUID                string          `json:"user_uuid" db:"user_uuid"`
	Name                    string          `json:"name" db:"name"`
	Status                  string          `json:"status" db:"status"`
	TTL                     int             `json:"ttl" db:"ttl"`
	ExistingTaskAction      string          `json:"existing_task_action" db:"existing_task_action"`
	LaunchMode              string          `json:"launch_mode" db:"launch_mode"`
	IncompatibleAction      string          `json:"incompatible_action" db:"incompatible_action"`
	ConfigurationSource     string          `json:"configuration_source" db:"configuration_source"`
	ConfigurationSourceData []byte          `json:"-" db:"configuration_source_data"`
	SourceMetadata          json.RawMessage `json:"source_metadata,omitempty" db:"source_metadata"`
	TasksCount              int             `json:"tasks_count" db:"tasks_count"`
}

type CompanyTask struct {
	ExpiresAt            time.Time  `json:"expires_at" db:"expires_at"`
	LastSentAt           *time.Time `json:"last_sent_at,omitempty" db:"last_sent_at"`
	ConfirmationDeadline *time.Time `json:"confirmation_deadline,omitempty" db:"confirmation_deadline"`
	NextAttemptAt        *time.Time `json:"next_attempt_at,omitempty" db:"next_attempt_at"`
	ID                   uint64     `json:"id" db:"id"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	CompanyID            uint64     `json:"company_id" db:"company_id"`
	DeviceID             uint64     `json:"device_id" db:"device_id"`
	DeviceImei           *string    `json:"device_imei,omitempty" db:"device_imei"`
	DeviceModel          *string    `json:"device_model,omitempty" db:"device_model"`
	Status               string     `json:"status" db:"status"`
	ErrorMessage         *string    `json:"error_message" db:"error_message"`
	Attempts             int        `json:"attempts" db:"attempts"`
	CfgData              []byte     `json:"-" db:"cfg_data"`
}

type CompanyWithTasks struct {
	Company
	ConfigurationSnapshot []byte        `json:"configuration_snapshot_base64,omitempty"`
	Tasks                 []CompanyTask `json:"tasks"`
}
