import { UIDs } from './UID.constant';

export const PRESET_MAPPING: Record<string, Record<string, string>> = {
	operator: {
		apn: UIDs.APN_NAME,
		login: UIDs.APN_USER,
		password: UIDs.APN_PASS,
	},
	server: {
		address: UIDs.SERVER_HOST,
		port: UIDs.SERVER_PORT,
		protocol: UIDs.TRAFFIC_PROTOCOL,
	},
};
