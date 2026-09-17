package models

import "neomatica/neosync/infra/constants"

type Access struct {
	ID      uint   `json:"id" db:"id"`
	Code    string `json:"code" db:"code"`
	NameKey string `json:"name_key" db:"name_key"`
}

type AccessFields struct {
	AccessTrekerCreate         bool `json:"access_treker_create" db:"access_treker_create"`
	AccessTrekerEdit           bool `json:"access_treker_edit" db:"access_treker_edit"`
	AccessTrekerDelete         bool `json:"access_treker_delete" db:"access_treker_delete"`
	AccessGroupManage          bool `json:"access_group_manage" db:"access_group_manage"`
	AccessConfigurationRead    bool `json:"access_configuration_read" db:"access_configuration_read"`
	AccessConfigurationApply   bool `json:"access_configuration_apply" db:"access_configuration_apply"`
	AccessConfigurationHistory bool `json:"access_configuration_history" db:"access_configuration_history"`
	AccessCommandSend          bool `json:"access_command_send" db:"access_command_send"`
	AccessLogRead              bool `json:"access_log_read" db:"access_log_read"`
}

type PtrAccessFields struct {
	AccessTrekerCreate         *bool `json:"access_treker_create" db:"access_treker_create"`
	AccessTrekerEdit           *bool `json:"access_treker_edit" db:"access_treker_edit"`
	AccessTrekerDelete         *bool `json:"access_treker_delete" db:"access_treker_delete"`
	AccessGroupManage          *bool `json:"access_group_manage" db:"access_group_manage"`
	AccessConfigurationRead    *bool `json:"access_configuration_read" db:"access_configuration_read"`
	AccessConfigurationApply   *bool `json:"access_configuration_apply" db:"access_configuration_apply"`
	AccessConfigurationHistory *bool `json:"access_configuration_history" db:"access_configuration_history"`
	AccessCommandSend          *bool `json:"access_command_send" db:"access_command_send"`
	AccessLogRead              *bool `json:"access_log_read" db:"access_log_read"`
}

func (a *AccessFields) Has(access constants.Access) bool {
	switch access {
	case constants.AccessTrekerCreate:
		return a.AccessTrekerCreate

	case constants.AccessTrekerEdit:
		return a.AccessTrekerEdit

	case constants.AccessTrekerDelete:
		return a.AccessTrekerDelete

	case constants.AccessGroupManage:
		return a.AccessGroupManage

	case constants.AccessConfigurationRead:
		return a.AccessConfigurationRead

	case constants.AccessConfigurationApply:
		return a.AccessConfigurationApply

	case constants.AccessConfigurationHistory:
		return a.AccessConfigurationHistory

	case constants.AccessCommandSend:
		return a.AccessCommandSend

	case constants.AccessLogRead:
		return a.AccessLogRead

	default:
		return false
	}
}
