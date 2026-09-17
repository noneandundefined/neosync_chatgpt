import { ConfigurationParsedResponse, ConfigurationSectionParsedResponse, FieldConfigurationParsed } from '@/rest/configurationAPI';
import { UIDs } from './UID.constant';

export const DEVICE_MODELS = {
	ADMP50: 'ADMP50',
	ADMP50LTE: 'ADMP50LTE',
	ADM007: 'ADM007',
	ADM007BLE: 'ADM007BLE',
	ADM333: 'ADM333',
	ADM333V2: 'ADM333V2',
	ADM333BLE: 'ADM333BLE',
	ADM500: 'ADM500',
};

/** Модель по умолчанию для шаблонов конфигурации (без привязки к устройству) */
export const TEMPLATE_STATIC_MODEL: keyof typeof DEVICE_MODELS = 'ADM333V2';

export const SUPPORT_CONFIG = (cfg: ConfigurationSectionParsedResponse | ConfigurationParsedResponse) => {
	const fields = fieldsToMap(cfg.config_parsed);

	const result: Record<string, any> = {};

	for (const [feature, rule] of Object.entries(FEATURE_DETECTORS)) {
		if ('detect' in rule) {
			result[feature] = rule.detect(cfg);
			continue;
		}

		const exist = rule.required.every((uid) => fields[uid] !== undefined);
		result[feature] = exist;
	}

	return result;
};

export const LLS_COUNT = 3;

export const STATIC_SUPPORT_CONFIG = {
	ADM500: {
		SimCount: 2, // Количество SIM
		MultiApn: false, // Поддержка apn
		OutputCount: 2, // Количество выходов
		StaticFreezingCoordinates: false,
		AlternativeSourcesPositioning: true,
	},
	ADM007BLE: {
		SimCount: 1, // Количество SIM
		MultiApn: true, // Поддержка apn
		OutputCount: 1, // Количество выходов
		StaticFreezingCoordinates: true,
		AlternativeSourcesPositioning: true,
	},
	ADM007: {
		SimCount: 2, // Количество SIM
		MultiApn: true, // Поддержка apn
		OutputCount: 2, // Количество выходов
		StaticFreezingCoordinates: true,
		AlternativeSourcesPositioning: true,
	},
	ADM333: {
		SimCount: 2, // Количество SIM
		MultiApn: true, // Поддержка apn
		OutputCount: 1, // Количество выходов
		StaticFreezingCoordinates: true,
		AlternativeSourcesPositioning: true,
	},
	ADM333BLE: {
		SimCount: 2, // Количество SIM
		MultiApn: true, // Поддержка apn
		OutputCount: 1, // Количество выходов
		StaticFreezingCoordinates: true,
		AlternativeSourcesPositioning: true,
	},
	ADM333V2: {
		SimCount: 2, // Количество SIM
		MultiApn: true, // Поддержка apn
		OutputCount: 1, // Количество выходов
		StaticFreezingCoordinates: true,
		AlternativeSourcesPositioning: true,
	},
	ADMP50: {
		SimCount: 2, // Количество SIM
		MultiApn: false, // Поддержка apn
		OutputCount: 0, // Количество выходов
		StaticFreezingCoordinates: false,
		AlternativeSourcesPositioning: false,
	},
	ADMP50LTE: {
		SimCount: 2, // Количество SIM
		MultiApn: false, // Поддержка apn
		OutputCount: 0, // Количество выходов
		StaticFreezingCoordinates: false,
		AlternativeSourcesPositioning: true,
	},
};

/** Лимиты UI для шаблонов — максимум по всем моделям, без привязки к устройству */
export const TEMPLATE_STATIC_SUPPORT = {
	SimCount: 2,
	MultiApn: true,
	OutputCount: 2,
	StaticFreezingCoordinates: true,
	AlternativeSourcesPositioning: true,
} as const;

export const getStaticSupport = (model: keyof typeof DEVICE_MODELS, isTemplate?: boolean) => (isTemplate ? TEMPLATE_STATIC_SUPPORT : STATIC_SUPPORT_CONFIG[model]);

export const getServerCount = (support: Record<string, unknown>, isTemplate?: boolean) => (isTemplate ? 2 : support.ModeHybrid ? 1 : 2);

export const FEATURE_DETECTORS = {
	RS485_Tab: {
		//Вкладка
		required: [UIDs.MODBUS_DEVICES_ADDR, UIDs.MODBUS_BYTES_ORDER, UIDs.MODBUS_BYTES_COUNT, UIDs.MODBUS_DEVICES_BD_RATE, UIDs.MODBUS_REGISTERS_TYPE, UIDs.MODBUS_REGISTERS_ADDR_2],
	},
	Sim_Tab: {
		//Вкладка
		required: [UIDs.APN_USER, UIDs.APN_NAME, UIDs.APN_PASS, UIDs.SIM_PIN, UIDs.OPS_WHITE_LIST_2, UIDs.OPS_BLACK_LIST],
	},
	Server_Tab: {
		//Вкладка
		required: [UIDs.SERVER_HOST, UIDs.SERVER_PORT, UIDs.TRAFFIC_PROTOCOL, UIDs.BLACK_BOX_STRUCTURE],
	},
	Event_Tab: {
		//Вкладка
		required: [UIDs.ALARM_PHONE],
	},
	Output_Tab: {
		//Вкладка
		required: [UIDs.OUTPUT],
	},
	Bluetooth_Tab: {
		// Вкладка Bluetooth
		required: [UIDs.BLE_SENSOR_TYPE],
	},
	OneWire_Tab: {
		//Вкладка
		required: [UIDs.IBUTTON_ENABLED, UIDs.OW_TEMP_SN],
	},
	Track_Tab: {
		//Вкладка
		required: [UIDs.POINTS_PERIOD_MOVE, UIDs.POINTS_PERIOD_HOLD, UIDs.POINTS_TRACK_COURSE, UIDs.POINTS_TRACK_CROSSTRACK],
	},
	Input_Tab: {
		//Вкладка
		required: [UIDs.AIN_FALSE_HIGH, UIDs.AIN_TRUE_LOW],
	},
	Modbus: {
		required: [UIDs.MODBUS_DEVICES_ADDR, UIDs.MODBUS_BYTES_ORDER, UIDs.MODBUS_BYTES_COUNT, UIDs.MODBUS_DEVICES_BD_RATE, UIDs.MODBUS_REGISTERS_TYPE, UIDs.MODBUS_REGISTERS_ADDR_2],
	},
	ADM20: {
		required: [UIDs.ADM20_ADDR, UIDs.ADM20_ID_MODE, UIDs.ADM20_MODE, UIDs.ADM20_ADDR, UIDs.ADM20_OUTPUT_STATE_MASK, UIDs.ADM20_ALARM_STATE],
	},
	ADM40: {
		required: [UIDs.ADM40_ADDR],
	},
	LLS: {
		required: [UIDs.FUEL_SENSOR_ADDRESS_LIST],
	},
	ModeHybrid: {
		required: [UIDs.DEVICE_MODE],
	},
	SimNet: {
		required: [UIDs.NET_MODE],
	},
	PowerBtn: {
		required: [UIDs.CUSTOM_MASK_DEFAULT_FALSE_1],
	},
	BatterySaving: {
		required: [UIDs.BATTERY_SAFE_MODE],
	},
	Gsm: {
		required: [UIDs.BEACON_MODE_GSM_GNSS],
	},
	LowBatteryAlarm: {
		required: [UIDs.LOW_BATTERY_LEVEL, UIDs.CUSTOM_MASK_DEFAULT_TRUE_1],
	},
	Sos: {
		required: [UIDs.CUSTOM_MASK_DEFAULT_FALSE_1],
	},
	InDoor: {
		required: [UIDs.CUSTOM_MASK_DEFAULT_FALSE_1],
	},
	StaticFreezingCoordinates: {
		required: [UIDs.STATIC_MODE, UIDs.STATIC_AIN_NUMBER, UIDs.STATIC_PROGRAM_SPEED, UIDs.DEVICE_FUNCTION_2],
	},
	InputPulse: {
		required: [UIDs.IMPULSE_INPUTS],
	},
	BleRele: {
		required: [UIDs.BLE_OUTPUT, UIDs.BLE_BROADCAST_ACTIVE_PERIOD_SEC, UIDs.BLE_BROADCAST_SILENT_PERIOD_SEC, UIDs.BLE_OUTPUT_ADDR, UIDs.BLE_OUTPUT_KEY],
	},
	MultiApn: {
		required: [UIDs.SIMTIME_MIN_CONF, UIDs.DEVICE_FUNCTION_3],
	},
	NavigationFilter: {
		required: [UIDs.NAV_FILTER_SATS_COUNT, UIDs.NAV_FILTER_HDOP_TENTHS, UIDs.NAV_FILTER_HDOP_MAX, UIDs.NAV_FILTER_VALID_HEIGHT_DECA_HECTO_METERS],
	},
	AlternativeSourcesPositioning: {
		required: [
			UIDs.ALTERNATIVE_SOURCES_POSITIONING_OPERATING_MODE,
			UIDs.ALTERNATIVE_SOURCES_POSITIONING_QUERY_PERIOD,
			UIDs.ALTERNATIVE_SOURCES_POSITIONING_MIN_TIME_REQUESTS,
			UIDs.ALTERNATIVE_SOURCES_POSITIONING_ACCURACY_COORDINATES,
			UIDs.ALTERNATIVE_SOURCES_POSITIONING_TOKEN,
		],
	},
	FilteringEmissions: {
		required: [UIDs.DISCARDED_POINTS_DURING_EJECTION, UIDs.DISTANCE_BETWEEN_POINTS],
	},
	SubstitutionCoordinatesAlternativeSource: {
		required: [UIDs.COORDINATE_SUBSTITUTION_ALTERNATIVE, UIDs.MULTIPLIER_COORDINATE_VALIDITY_AREA],
	},
	Navtimesync: {
		required: [UIDs.NAVTIMESYNC, UIDs.NAVTIMESYNC_MAX_TIME_DIFF_SEC],
	},
	InputCount: {
		detect: (cfg: ConfigurationParsedResponse | ConfigurationSectionParsedResponse) => {
			const entry = cfg.config_parsed.find((c) => [UIDs.AIN_FALSE_HIGH, UIDs.AIN_TRUE_LOW].includes(c.uid as UIDs));

			if (!entry || !Array.isArray(entry.value)) return 0;

			return entry.value.length;
		},
	},
};

const fieldsToMap = (configuration: FieldConfigurationParsed[]) => {
	const map: Record<string, FieldConfigurationParsed> = {};
	for (const field of configuration) {
		map[field.uid] = field;
	}

	return map;
};
