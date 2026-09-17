export interface DeviceConfResponse {
	id: number;
	created_at: string;
	updated_at: string;
	device_id: number;
	password: string;
	auto_update: boolean;
	request_configuration_on_connect: boolean;
}
