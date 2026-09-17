type EventView = {
	text: string;
	color: string;
};

export const EVENT_VIEW_MAP: Record<string, EventView> = {
	event_device_conn: {
		text: 'label.log-terminal-conn',
		color: '#039709',
	},
	event_device_disconn: {
		text: 'label.log-terminal-disconn',
		color: '#000',
	},
	event_device_sync: {
		text: 'label.log-sync',
		color: '#000',
	},
	event_device_cfgrecev: {
		text: 'label.log-cfg-recev',
		color: '#000',
	},
	event_device_cfgrequest: {
		text: 'label.log-cfg-request',
		color: '#000',
	},
	event_device_cfgsend: {
		text: 'label.log-cfg-send',
		color: '#000',
	},
	event_device_sendcommand: {
		text: 'label.log-command-send',
		color: '#000',
	},
	event_device_recevcommand: {
		text: 'label.log-command-recev',
		color: '#000',
	},
};
