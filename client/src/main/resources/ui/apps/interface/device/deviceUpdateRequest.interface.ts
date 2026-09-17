export interface DeviceUpdateRequest {
	imei?: string;
	model?: string;
	name?: string | null;
	phone?: string | null;
	name_organization?: string | null;
	request_configuration_on_connect?: boolean;
}
