export const ACCESS_RULES = [
	{
		key: 'access_treker_create',
		label: 'message.access-create-treker',
		group: 'main',
	},
	{
		key: 'access_treker_edit',
		label: 'message.access-edit-treker',
		group: 'main',
	},
	{
		key: 'access_treker_delete',
		label: 'message.access-delete-treker',
		group: 'main',
	},
	{
		key: 'access_group_manage',
		label: 'message.access-manage-group',
		group: 'main',
	},
	{
		key: 'access_configuration_read',
		label: 'message.access-read-configuration',
		group: 'configuration',
	},
	{
		key: 'access_configuration_apply',
		label: 'message.access-edit-configuration',
		group: 'configuration',
	},
	{
		key: 'access_configuration_history',
		label: 'message.access-history-configuration',
		group: 'configuration',
	},
	{
		key: 'access_command_send',
		label: 'message.access-send-commands',
		group: 'command',
	},
	{
		key: 'access_log_read',
		label: 'message.access-log-read',
		group: 'log',
	},
];

export const GROUP_RIGHTs_RULES = [
	{
		key: 'can_edit_group',
		label: 'message.editing-group',
		group: 'main',
	},
	{
		key: 'can_manage_devices',
		label: 'message.actions-objects-in-group',
		group: 'main',
	},
	{
		key: 'can_read_config',
		label: 'message.access-read-configuration',
		group: 'device',
	},
	{
		key: 'can_edit_config',
		label: 'message.access-edit-configuration',
		group: 'device',
	},
	{
		key: 'can_send_commands',
		label: 'message.access-send-commands',
		group: 'device',
	},
];

export const BASE_ACCESS_MASK = {
	access_treker_create: true,
	access_treker_edit: true,
	access_treker_delete: true,
	access_group_manage: true,
	access_configuration_read: true,
	access_configuration_apply: true,
	access_configuration_history: true,
	access_command_send: true,
	access_log_read: true,
};

export const BASE_GROUP_RIGHTs_MASK = {
	can_edit_group: true,
	can_manage_devices: true,
	can_read_config: true,
	can_edit_config: true,
	can_send_commands: true,
};
