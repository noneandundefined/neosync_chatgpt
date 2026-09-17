package models

import "time"

type DeviceCommand struct {
	ID           uint64     `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	ExecuteUntil *time.Time `json:"execute_until" db:"execute_until"`
	SessionID    string     `json:"session_id" db:"session_id"`
	UserUUID     string     `json:"user_uuid" db:"user_uuid"`
	SendMode     string     `json:"send_mode" db:"send_mode"` /* instant | on_connect */
	Status       string     `json:"status" db:"status"`       /* pending | inprogress | completed | executionerror */
	Command      string     `json:"command" db:"command"`
}

type DeviceCommandExecution struct {
	ID         uint64    `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	CommandID  uint64    `json:"command_id" db:"command_id"`
	DeviceIMEI string    `json:"device_imei" db:"device_imei"`
	Status     string    `json:"status" db:"status"` /* pending | inprogress | completed | executionerror */
	Response   *string   `json:"response" db:"response"`
}

type DeviceCommandWithExecutions struct {
	ID          uint64 `json:"id" db:"id"`
	ExecutionID uint64 `json:"execution_id" db:"execution_id"`
	SendMode    string `json:"send_mode" db:"send_mode"`
	Command     string `json:"command" db:"command"`
	IMEI        string `json:"imei" db:"device_imei"`
	Status      string `json:"status" db:"status"`
}
