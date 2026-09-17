import { Device } from '@/rest/deviceAPI';
import { UserResponse } from '@/rest/userAPI';

export interface TableColumnDevice extends Device {
	firmware_version: number;
	cfg_hash: number;
	cfg_version: number;
	last_mod_time: number;
	task?: any;
}

export interface TableColumnUsers extends UserResponse {}
