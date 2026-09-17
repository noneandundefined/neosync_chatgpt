export interface DeviceSendCommandRequest {
	sess_id?: string;
	send_mode: string;
	command: string;
	imeis: string[];
}
