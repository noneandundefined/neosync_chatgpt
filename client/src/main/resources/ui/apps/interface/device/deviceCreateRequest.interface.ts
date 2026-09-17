export interface DeviceCreateRequest {
	imei: string;
	model: string;
	name: string | null;
	phone: string | null;
	name_organization: string | null;
	request_configuration_on_connect: boolean;
	turnstile_token: string;
}

export interface DeviceImportRequest {
	imei: string;
	model?: string | null;
	owner_email?: string | null;
	group_name?: string | null;
}
