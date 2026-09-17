export interface PtrAccessFields {
	access_treker_create?: boolean;
	access_treker_edit?: boolean;
	access_treker_delete?: boolean;
	access_group_manage?: boolean;
	access_configuration_read?: boolean;
	access_configuration_apply?: boolean;
	access_configuration_history?: boolean;
	access_command_send?: boolean;
	access_log_read?: boolean;
}

export interface UserUpdateRequest {
	email?: string;
	password?: string;
	phone?: string;
	name_organization?: string;
	locality?: string;
	language?: string;
	accesses?: PtrAccessFields;
	turnstile_token?: string;
}
