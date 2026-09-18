const indexKey = 'neosync';

export const CACHEKEYs = {
	GROUP_KEY: `${indexKey}:groups:sess`,
	CFG_HASHES: `${indexKey}:configuration:hashes`,
	USER_ACCESS: `${indexKey}:user:access`,
	USER_AUTH: `${indexKey}:auth:user`,
	AUTH_REMEMBER_ME: `${indexKey}:auth:remember_me`,
	TERMINAL_LLS_ADDRESS: `${indexKey}:terminal:{imei}:lls_addresses`,
	TERMINAL_LLS_MAPPING: `${indexKey}:terminal:{imei}:lls_mapping`,
	TERMINAL_TARIRATION_TABLES: `${indexKey}:terminal:{imei}:tariration_tables`,
	TERMINAL_TARIRATION_SELECTED: `${indexKey}:terminal:{imei}:tariration_selected`,
	TERMINAL_TELEMETRY: `${indexKey}:terminal:{imei}:telemetry`,
	TERMINAL_TRACK_EXPERT_SETTINGs: `${indexKey}:terminal:{imei}:track:expert_settings`,
	NEOSYNC_COMMAND_HISTORY: `${indexKey}:command_history`,
	NEOSYNC_VERSION: `${indexKey}:s:version`,
	NEOSYNC_X_REQ_ID: `${indexKey}:settings-client-req`,
	NEOSYNC_CMD_SESSION_ID: `${indexKey}:cmd:session_id`,
	ANALYTIC_WIDGETS: `${indexKey}:analytics:widgets`,
	SIDEBAR_COLLAPSED: `${indexKey}:sidebar:collapsed`,
	COOKIE_CONSENT: `${indexKey}:cookies:consent`,
};

export const ACTIONKEYs = {
	DEVICE_COMMAND_SEND: (cmd: string) => `${indexKey}:action:device:command:${cmd}:send`,
	DEVICE_COMMAND_REBOOT: (imei: string) => `${indexKey}:action:device:command:reboot:${imei}`,
	DEVICE_COMMAND_ERASE_EEPROM: (imei: string) => `${indexKey}:action:device:command:erase_eeprom:${imei}`,
	DEVICE_COMMAND_ERASE_FLASH: (imei: string) => `${indexKey}:action:device:command:erase_flash:${imei}`,
	DEVICE_COMMAND_FIND_BLE: (imei: string) => `${indexKey}:action:device:command:find_ble:${imei}`,
	DEVICE_COMMAND_FIND_OWIRE: (imei: string) => `${indexKey}:action:device:command:find_owire:${imei}`,
	DEVICE_COMMAND_DELETE: (id: number) => `${indexKey}:action:device:command:delete:${id}`,
	DEVICE_COMMAND_CANCEL: (id: number) => `${indexKey}:action:device:command:cancel:${id}`,
	DEVICE_IMPORT: `${indexKey}:action:device:import`,
	DEVICE_DELETE_MASSIVE: `${indexKey}:action:device:delete_massive`,
	CONFIGURATION_APPLY: (imei: string) => `${indexKey}:action:configuration:apply:${imei}`,
	CONFIGURATION_TELEMETRY_REBOOT: (imei: string) => `${indexKey}:action:configuration:telemetry:${imei}`,
	CONFIGURATION_IMPORT: (imei: string) => `${indexKey}:action:configuration:import:${imei}`,
	CONFIGURATION_IMPORT_BY_HISTORY: (imei: string, cfgHash: number) => `${indexKey}:action:configuration:import_history:${imei}:${cfgHash}`,
	GROUP_DELETE_MASSIVE: `${indexKey}:action:group:delete_massive`,
	USER_DELETE_MASSIVE: `${indexKey}:action:user:delete_massive`,
};

export const CACHEKEYs_TERMINAL_LLS_ADDRESS = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:lls_addresses`;
};

export const CACHEKEYs_TERMINAL_LLS_MAPPING = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:lls_mapping`;
};

export const CACHEKEYs_TERMINAL_TARIRATION_TABLES = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:tariration_tables`;
};

export const CACHEKEYs_TERMINAL_TARIRATION_SELECTED = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:tariration_selected`;
};

export const CACHEKEYs_TERMINAL_TELEMETRY = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:telemetry`;
};

export const CACHEKEYs_TERMINAL_TRACK_EXPERT_SETTINGs = (imei: string): string => {
	return `${indexKey}:terminal:${imei}:track:expert_settings`;
};

export const CACHEKEYs_GROUP_OBJECTS_KEY = (id: number): string => {
	return `${indexKey}:group:objects:sess:${id}`;
};

export const CACHEKEYs_PROTOCOL = (imei: string): string => {
	return `${indexKey}:protocol:sess:${imei}`;
};

export const CACHEKEYs_TABLE_VIEW_COLUMNS = (key: string) => {
	return `${indexKey}:table:${key}:view:columns`;
};

export const CACHEKEYs_TABLE_LIMIT_EL = (key: string) => {
	return `${indexKey}:table:${key}:limit:el`;
};

export const CACHEKEYs_TABLE_COLUMN_SORT = (key: string) => {
	return `${indexKey}:table:${key}:columnSort`;
};
