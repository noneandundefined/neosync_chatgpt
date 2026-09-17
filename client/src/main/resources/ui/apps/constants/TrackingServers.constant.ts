import { DEVICE_MODELS } from './DeviceModels.constant';

interface TrackingServerExtra {
	address: string;
	port: number;
	protocol: number;
}

export interface TrackingServer {
	id: string;
	value: string;
	extra: TrackingServerExtra;
}

export interface TrackingServerBase {
	id: string;
	value: string;
	address: string;
	protocol: number;
}

export const TRACKING_SERVERS: TrackingServerBase[] = [
	{
		id: 'default',
		value: 'label.default',
		address: 'www.test.neomatica.ru',
		protocol: 0,
	},
	{
		id: 'custom',
		value: 'label.custom',
		address: '',
		protocol: 0,
	},
	{
		id: 'neomatica',
		value: 'label.neomatica',
		address: 'www.test.neomatica.ru',
		protocol: 0,
	},
	{
		id: 'glonass-soft',
		value: 'label.glonass-soft',
		address: '77.223.106.70',
		protocol: 0,
	},
	{
		id: 'wialon',
		value: 'label.wialon',
		address: '193.193.165.165',
		protocol: 0,
	},
	{
		id: 'omnicomm',
		value: 'label.omnicomm',
		address: 'convert.omnicomm.ru',
		protocol: 0,
	},
	{
		id: 'globars',
		value: 'label.globars',
		address: '95.167.243.15',
		protocol: 0,
	},
	{
		id: 'ruhavik',
		value: 'label.ruhavik',
		address: '193.193.165.37',
		protocol: 0,
	},
	{
		id: 'rnis',
		value: 'label.rnis',
		address: 'data.rnis.mos.ru',
		protocol: 1,
	},
	{
		id: 'era-glonass',
		value: 'label.era-glonass',
		address: '10.77.60.254',
		protocol: 1,
	},
	{
		id: 'live-gps',
		value: 'label.live-gps',
		address: 'vega.livegpstracks.ru',
		protocol: 0,
	},
	{
		id: 'axenta',
		value: 'label.axenta',
		address: 'hw.axenta.cloud',
		protocol: 0,
	},
	{
		id: 'pilot-gps',
		value: 'label.pilot-gps',
		address: 'blade.pilot-gps.com',
		protocol: 0,
	},
	{
		id: 'navixy',
		value: 'label.navixy',
		address: 'tracker.navixy.com',
		protocol: 0,
	},
];

export const TRACKING_SERVER_PORTS: Record<string, number | { default: number; admp50?: number; adm007ble?: number; adm333?: number }> = {
	default: 255,
	custom: 0,
	neomatica: 12301,
	'glonass-soft': 12000,
	wialon: {
		default: 22480,
		admp50: 21282,
		adm007ble: 21404,
		adm333: 21848,
	},
	omnicomm: 15306,
	globars: 7046,
	ruhavik: 28436,
	rnis: 4444,
	'era-glonass': 30197,
	'live-gps': 3369,
	axenta: 20663,
	'pilot-gps': 20060,
	navixy: 47758,
};

export const buildTrackingServers = (model?: string): TrackingServer[] => {
	const mode: 'admp50' | 'adm007ble' | 'adm333' | undefined =
		model === DEVICE_MODELS.ADM007BLE
			? 'adm007ble'
			: model === DEVICE_MODELS.ADMP50 || model === DEVICE_MODELS.ADMP50LTE
				? 'admp50'
				: model === DEVICE_MODELS.ADM333 || model === DEVICE_MODELS.ADM333BLE || model === DEVICE_MODELS.ADM333V2
					? 'adm333'
					: undefined;

	return TRACKING_SERVERS.map((server) => {
		const config = TRACKING_SERVER_PORTS[server.id];

		let port = 0;

		if (typeof config === 'number') {
			port = config;
		} else if (config) {
			port = mode ? (config[mode] ?? config.default) : config.default;
		}

		const id = server.id === 'custom' ? 'custom' : `${server.address}:${port}`;

		return {
			id: id,
			value: server.value,
			extra: {
				address: server.address,
				port,
				protocol: server.protocol,
			},
		};
	});
};

export const TRACKING_SERVERS_PROTO = [
	{
		id: 0,
		value: 'ADM',
	},
	{
		id: 1,
		value: 'EGTS',
	},
	{
		id: 2,
		value: 'WIPS',
	},
];
