export interface GroupUpdateRequest {
	name: string;
	description: string;
	devices_add?: number[];
	devices_remove?: number[];
	can_edit_group: boolean;
	can_manage_devices: boolean;
	can_read_config: boolean;
	can_edit_config: boolean;
	can_send_commands: boolean;
}
