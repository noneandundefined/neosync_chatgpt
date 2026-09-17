export interface ADM20Info {
	emptiness_mask: number | null;
	card_uid: number | null;
	tag_id: number | null;
	tag_sn: number | null;
	tag_rssi: number | null;
	tag_adc: number | null;
	tag_period: number | null;
}

export interface BLESensor {
	rssi?: number;
	voltage?: number;
	lmt?: number;
	fuel_level?: number;
	temp?: number;
}

export interface TelemetryResponse {
	gps_satellites: number;
	glonass_satellites: number;
	lat?: number;
	lon?: number;
	date_time?: number;
	v_power?: number;
	v_battery?: number;
	in_a0?: number;
	in_a1?: number;
	in_d0?: number;
	in_d1?: number;
	adm20info?: ADM20Info;
	blesensorinfo?: Record<string, BLESensor>;
	fuelinfo?: Record<string, BLESensor>;
}
