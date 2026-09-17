export interface DeviceFileParser {
	imei: string;
	model?: string | null;
	error?: string;
}
