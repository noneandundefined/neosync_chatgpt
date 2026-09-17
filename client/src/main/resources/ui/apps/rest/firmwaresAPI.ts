import { toast } from 'react-toastify';
import axiosClient from './axios';

const apiPath = '/firmwares';

export interface DeviceModelSource {
	id: number;
	created_at: string;
	device_model: string;
	firmware_url: string;
	firmware_cache_key: string;
}

export interface Firmwares_DeviceModelSource {
	created_at: string;
	device_model: string;
	firmware: string;
	release_date: string;
	type: string;
}

export interface FirmwaresWPResponse {
	items: Firmwares_DeviceModelSource[];
	page: number;
	limit: number;
	total: number;
	total_pages: number;
}

/**
 * Вывод моделей для трекеров
 */
export const basicFirmwareDeviceModelsGet = async (signal?: AbortSignal): Promise<string[]> => {
	const response = await axiosClient.get(`${apiPath}/models`, { signal });
	return response.data.message;
};

/**
 * Вывод sources терминалов
 */
export const basicFirmwareSources = async (page: number, limit: number, search: string = '', signal?: AbortSignal): Promise<FirmwaresWPResponse> => {
	const response = await axiosClient.get(`${apiPath}/sources?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}`, { signal });
	return response.data.message;
};

/**
 * Обновление версии прошивки трекера
 */
export const basicFirmwareUpdate = async (imei: string) => {
	const response = await axiosClient.post(`${apiPath}/devices/${imei}/firmware`);
	toast.success(response.data.message);
};
