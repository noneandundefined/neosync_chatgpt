package group_handler_v1

import "neomatica/neosync/infra/store/postgres/models"

type CreateGroupPayload struct {
	Name           string `json:"name" validate:"required,max=30"`
	Description    string `json:"description" validate:"omitempty,max=300"`
	TurnstileToken string `json:"turnstile_token" validate:"required"`

	CanEditGroup     bool `json:"can_edit_group"`
	CanManageDevices bool `json:"can_manage_devices"`
	CanReadConfig    bool `json:"can_read_config"`
	CanEditConfig    bool `json:"can_edit_config"`
	CanSendCommands  bool `json:"can_send_commands"`
}

type UpdateGroupPayload struct {
	Name        string `json:"name" validate:"required,max=30"`
	Description string `json:"description" validate:"omitempty,max=300"`

	DevicesAdd    []uint64 `json:"devices_add" validate:"omitempty,dive,gt=0"`
	DevicesRemove []uint64 `json:"devices_remove" validate:"omitempty,dive,gt=0"`

	CanEditGroup     bool `json:"can_edit_group"`
	CanManageDevices bool `json:"can_manage_devices"`
	CanReadConfig    bool `json:"can_read_config"`
	CanEditConfig    bool `json:"can_edit_config"`
	CanSendCommands  bool `json:"can_send_commands"`
}

type DevicesGroupResponse struct {
	DevicesInGroup    []models.GDevice `json:"devices_in_group"`
	DevicesNotInGroup []models.GDevice `json:"devices_not_in_group"`
}

type ActionDevicesGroupPayload struct {
	DevicesId []uint64 `json:"devices_id" validate:"required,min=1,dive,gt=0"`
}

type CommandGroupPayload struct {
	Command string `json:"command" validate:"required,min=3"`
}

type CommandGroupResponse struct {
	Imei   string `json:"imei" validate:"imei"`
	Status int    `json:"status"`
	Result string `json:"result"`
}

type GroupsWPResponse struct {
	Items      []models.Group `json:"items"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalAll   int            `json:"total_all"`
	TotalPages int            `json:"total_pages"`
}

type DeleteGroupsPayload struct {
	IDs []uint64 `json:"ids" validate:"required"`
}

type GroupMembersBatchPayload struct {
	MemberUUIDs    []string `json:"member_uuids" validate:"required,min=1,dive,required"`
	TurnstileToken string   `json:"turnstile_token" validate:"required"`
}
