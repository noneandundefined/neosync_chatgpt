package models

import "time"

type AnalyticsConfiguration struct {
	ID           uint64    `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UserUUID     string    `db:"user_uuid" json:"user_uuid"`
	DeviceID     uint64    `db:"device_id" json:"device_id"`
	CfgHash      uint32    `db:"cfg_hash" json:"cfg_hash"`
	RMQPublished bool      `db:"rmq_published" json:"rmq_published"`
	DeviceOnline bool      `db:"device_online" json:"device_online"`
	ErrorReason  *string   `db:"error_reason" json:"error_reason,omitempty"`
}
