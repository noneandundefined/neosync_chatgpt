package models

import "time"

type Group struct {
	ID               uint64    `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	UserUuid         string    `json:"user_uuid" db:"user_uuid"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	CanEditGroup     bool      `json:"can_edit_group" db:"can_edit_group"`
	CanManageDevices bool      `json:"can_manage_devices" db:"can_manage_devices"`
	CanReadConfig    bool      `json:"can_read_config" db:"can_read_config"`
	CanEditConfig    bool      `json:"can_edit_config" db:"can_edit_config"`
	CanSendCommands  bool      `json:"can_send_commands" db:"can_send_commands"`
	Objects          uint64    `json:"objects" db:"objects"`
}

type GroupWithDevices struct {
	Group *Group `json:"group"`

	Devices struct {
		Assigned  []GDevice `json:"assigned"`
		Available []GDevice `json:"available"`
	} `json:"devices"`
}

type GroupDevices struct {
	ID        uint64   `json:"id" db:"id"`
	UserUuid  string   `json:"user_uuid" db:"user_uuid"`
	GroupId   uint64   `json:"group_id" db:"group_id"`
	DevicesId []uint64 `json:"devices_id" db:"devices_id"`
}

type GDevice struct {
	ID   uint64 `json:"id" db:"id"`
	Imei string `json:"imei" db:"imei"`
}

type GroupName struct {
	ID   uint64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type GroupMember struct {
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	GroupID        uint64    `json:"group_id" db:"group_id"`
	MemberUserUUID string    `json:"member_user_uuid" db:"member_user_uuid"`
}
