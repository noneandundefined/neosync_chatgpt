package models

import "time"

type CompanyTaskWithCompany struct {
	DeviceImei           string     `db:"device_imei"`
	DatabaseNow          time.Time  `db:"database_now"`
	ID                   uint64     `db:"id"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
	CompanyID            uint64     `db:"company_id"`
	DeviceID             uint64     `db:"device_id"`
	Status               string     `db:"status"`
	ErrorMessage         *string    `db:"error_message"`
	Attempts             int        `db:"attempts"`
	ExpiresAt            time.Time  `db:"expires_at"`
	LastSentAt           *time.Time `db:"last_sent_at"`
	ConfirmationDeadline *time.Time `db:"confirmation_deadline"`
	NextAttemptAt        *time.Time `db:"next_attempt_at"`
	AttemptToken         *string    `db:"attempt_token"`
	LeaseUntil           *time.Time `db:"lease_until"`
	CfgData              []byte     `db:"cfg_data"`
	CompanyCreatedAt     time.Time  `db:"company_created_at"`
	CompanyStatus        string     `db:"company_status"`
	CompanyTTL           int        `db:"company_ttl"`
}

type CompanyTaskDelivery struct {
	Status               string
	ErrorMessage         *string
	AttemptToken         *string
	ConfirmationDeadline *time.Time
	NextAttemptAt        *time.Time
	LeaseUntil           *time.Time
	Sent                 bool
	Claim                bool
}
