package constants

type Access string

const (
	AccessTrekerCreate         Access = "access_treker_create"
	AccessTrekerEdit           Access = "access_treker_edit"
	AccessTrekerDelete         Access = "access_treker_delete"
	AccessGroupManage          Access = "access_group_manage"
	AccessConfigurationRead    Access = "access_configuration_read"
	AccessConfigurationApply   Access = "access_configuration_apply"
	AccessConfigurationHistory Access = "access_configuration_history"
	AccessCommandSend          Access = "access_command_send"
	AccessLogRead              Access = "access_log_read"
)
