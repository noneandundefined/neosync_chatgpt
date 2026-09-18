export interface UserCreateRequest {
	email: string;
	language: string;
	locality: string;
	name_organization: string;
	phone: string;
	password: string;
	role_code: string;
	accesses: {
		access_treker_create: boolean;
		access_treker_edit: boolean;
		access_treker_delete: boolean;
		access_group_manage: boolean;
		access_configuration_read: boolean;
		access_configuration_apply: boolean;
		access_configuration_history: boolean;
		access_command_send: boolean;
		access_log_read: boolean;
	};
}
