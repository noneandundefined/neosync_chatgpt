export interface GroupCreateRequest {
	name: string;
	description?: string;
	turnstile_token: string;
	can_edit_group: boolean;
	can_manage_devices: boolean;
	can_read_config: boolean;
	can_edit_config: boolean;
	can_send_commands: boolean;
}
