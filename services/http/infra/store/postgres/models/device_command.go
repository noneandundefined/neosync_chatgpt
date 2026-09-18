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
	Status       string     `json:"status" db:"status"`       /* pending | inprogress | completed | notcompleted | executionerror */
	Command      string     `json:"command" db:"command"`
}

type DeviceCommandExecution struct {
	ID         uint64    `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	CommandID  uint64    `json:"command_id" db:"command_id"`
	DeviceIMEI string    `json:"device_imei" db:"device_imei"`
	Status     string    `json:"status" db:"status"` /* pending | inprogress | completed | notcompleted | executionerror */
	Response   *string   `json:"response" db:"response"`
}

type DeviceCommandWithExecutions struct {
	ID           uint64                        `json:"id" db:"id"`
	CreatedAt    time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at" db:"updated_at"`
	ExecuteUntil *time.Time                    `json:"execute_until" db:"execute_until"`
	SessionID    string                        `json:"session_id" db:"session_id"`
	UserUUID     string                        `json:"user_uuid" db:"user_uuid"`
	SendMode     string                        `json:"send_mode" db:"send_mode"` /* instant | on_connect */
	Status       string                        `json:"status" db:"status"`       /* pending | inprogress | completed | notcompleted | executionerror */
	Command      string                        `json:"command" db:"command"`
	Executions   []DeviceCommandExecutionShort `json:"executions" db:"executions"`
}

type DeviceCommandExecutionShort struct {
	IMEI     string  `json:"imei" db:"imei"`
	Status   string  `json:"status" db:"status"`
	Response *string `json:"response" db:"response"`
}

type DevicesForSendCommand struct {
	GroupID     uint64  `json:"group_id" db:"group_id"`
	GroupName   *string `json:"name" db:"name"`
	DeviceCount int     `json:"device_count" db:"device_count"`
}

type DeviceShort struct {
	IMEI      string `json:"imei" db:"imei"`
	Status    bool   `json:"status" db:"status"`
	Activated *bool  `json:"activated" db:"activated"`
}
