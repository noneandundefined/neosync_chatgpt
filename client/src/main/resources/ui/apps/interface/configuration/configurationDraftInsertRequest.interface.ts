export interface ConfigurationDraftInsertRequest {
	device_imei: string;
	cfg_hash: number;
	changes: Record<string, any>;
	timestamp: number;
}
