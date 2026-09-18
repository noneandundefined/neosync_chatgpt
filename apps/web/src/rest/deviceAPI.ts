import i18next from 'i18next';
import axiosClient from './axios';
import { toast } from 'react-toastify';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { ACTIONKEYs } from '@/constants/CacheKeys.constants';
import { runGuardedAction } from '@/utils/RequestControlUtils';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';
import { DeviceUpdateRequest } from '@/interface/device/deviceUpdateRequest.interface';
import { DeviceTransferRequest } from '@/interface/device/deviceTransferRequest.interface';
import { DevicesForSendCommand, DevicesStateDevicesResponse, DeviceShort } from './deviceCommandAPI';
import { DeviceCreateRequest, DeviceImportRequest } from '@/interface/device/deviceCreateRequest.interface';

const apiPath = '/devices';

export interface Device {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	owner_uuid: string | null;
	name: string | null;
	phone: string | null;
	name_organization: string | null;
	imei: string;
	status: boolean;
	group_name?: string | null;
	device_model?: string | null;
	device_extended_model?: string | null;
	activated?: boolean | null;
}

export interface AdmLogItem {
	id: string;
	ts: number;
	event: string;
	payload: any;
}

export interface DeviceConf {
	id: number;
	created_at: string;
	updated_at: string;
	device_id: number;
	password: string;
	request_configuration_on_connect: boolean;
}

export interface DeviceWithDeviceConf extends Device {
	password: string;
	request_configuration_on_connect: boolean;
}

export interface Device_DeviceConf_Sync extends Device {
	// From device_confs
	request_configuration_on_connect: boolean;

	// From configurations
	configuration_updated_at: string | null;
	configuration_sync_status: string | null;
	configuration_sync_error: string | null;

	// From syncs
	firmware_version: number;
	firmware_version_upd: number;
	firmware_updated_at: string | null;
	firmware_update_status: string | null;
	cfg_version: number;
	last_mod_time: number;
	cfg_hash: number;

	// From user_cores
	user_email: string | null;

	// From groups and group_user_devices
	group_name: string | null;
	group_id: number | null;
}

export interface DevicesWPResponse {
	items: Device_DeviceConf_Sync[];
	page: number;
	limit: number;
	total: number;
	total_all: number;
	total_pages: number;
}

export interface DeviceImportError {
	imei: string;
	error: string;
}

export interface DeviceDetailedStatusResponse {
	id: number;
	imei: string;
	status: boolean;
	device_model: string | null;
	configuration_sync_status: string | null;
	firmware_version: number | null;
}

const importDevicesRequest = async (payload: DeviceImportRequest[], showSuccessToast: boolean): Promise<{ errors?: DeviceImportError[] }> => {
	try {
		const response = await axiosClient.post(`${apiPath}/import`, payload, {
			skipErrorHandler: true,
			actionKey: ACTIONKEYs.DEVICE_IMPORT,
		} as any);

		if (showSuccessToast) {
			toast.success(response.data.message);
		}

		return {};
	} catch (err: any) {
		if (err.response?.status === 409 && err.response.data) {
			const errors: DeviceImportError[] = err.response.data.message;
			if (showSuccessToast) {
				toast.error(i18next.t('message.import-error'));
			}
			return { errors };
		}

		throw err;
	}
};

/**
 * Создания терминала
 */
export const basicDeviceCreate = async (payload: DeviceCreateRequest) => {
	const response = await axiosClient.post(`${apiPath}`, payload);
	toast.success(response.data.message);
};

/**
 * Импорт терминалов
 */
export const basicDevicesImport = async (payload: DeviceImportRequest[]): Promise<{ errors?: DeviceImportError[] }> => {
	return runGuardedAction(ACTIONKEYs.DEVICE_IMPORT, async () => importDevicesRequest(payload, true), {
		cooldownMs: 2000,
		getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
		onBlocked: (msg) => showServerErrorToast(msg),
	});
};

export const basicDevicesImportRaw = async (payload: DeviceImportRequest[]): Promise<{ errors?: DeviceImportError[] }> => {
	return importDevicesRequest(payload, false);
};

export const basicDevicesImportCheck = async (imeis: string[]): Promise<DeviceImportError[]> => {
	if (imeis.length === 0) {
		return [];
	}

	const response = await axiosClient.post(`${apiPath}/import/check`, { imeis });
	return response.data.message ?? [];
};

/**
 * Получение групп устройств для отправки команд
 */
export const basicDevicesState = async (search = ''): Promise<DevicesForSendCommand[]> => {
	const response = await axiosClient.get(`${apiPath}/state`, { params: { search } });
	return response.data.message;
};

/**
 * Получение IMEI группы для отправки команд (постранично)
 */
export const basicDevicesStateDevices = async (params: { groupId: number; search?: string; page?: number; limit?: number; all?: boolean }): Promise<DevicesStateDevicesResponse> => {
	const response = await axiosClient.get(`${apiPath}/state/devices`, {
		params: {
			group_id: params.groupId,
			search: params.search || undefined,
			page: params.all ? undefined : params.page,
			limit: params.all ? undefined : params.limit,
			all: params.all ? true : undefined,
		},
	});
	return response.data.message;
};

/**
 * Вывод терминала
 */
export const basicDeviceGetByImei = async (imei: string, signal?: AbortSignal): Promise<Device_DeviceConf_Sync> => {
	const response = await axiosClient.get(`${apiPath}/${imei}`, { signal });
	return response.data.message;
};

/**
 * Авторизация терминала
 */
export const basicDeviceAuthorization = async (imei: string, password: string): Promise<boolean> => {
	const payload = {
		password: password,
	};

	const response = await axiosClient.post(`${apiPath}/${imei}/authorization`, payload);
	return response.data.message;
};

/**
 * Вывод логов терминала
 */
export const basicDeviceLogs = async (imei: string, signal?: AbortSignal): Promise<AdmLogItem[]> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/logs`, {
		signal,
	});
	return response.data.message;
};

/**
 * Вывод статус подключения терминала
 */
export const basicDeviceStatus = async (imei: string, signal?: AbortSignal): Promise<{ updated_at: string; status: boolean; device_extended_model: string | null }> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/status`, {
		signal,
	});
	return response.data.message;
};

/**
 * Вывод подробного статуса терминалов
 */
export const basicDeviceStatusDetailed = async (imeis: string[]): Promise<DeviceDetailedStatusResponse[]> => {
	const response = await axiosClient.post(`${apiPath}/status/detailed`, { imeis });
	return response.data.message;
};

/**
 * Вывод терминалов
 */
export const basicDevicesGet = async (signal?: AbortSignal): Promise<Device_DeviceConf_Sync[]> => {
	const response = await axiosClient.get(`${apiPath}`, { signal });
	return response.data.message;
};

/**
 * Вывод терминалов с пагинацией
 */
export const basicDevicesGetWithParams = async (page: number, limit: number, search: string = '', columnSort: ColumnSortState, signal?: AbortSignal): Promise<DevicesWPResponse> => {
	const response = await axiosClient.get(`${apiPath}?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}&columnSortKey=${columnSort.key}&columnSortDir=${columnSort.direction}`, { signal });
	return response.data.message;
};

/**
 * Получение ВСЕХ терминалов
 */
export const basicDevicesGetFull = async (): Promise<DeviceShort[]> => {
	const response = await axiosClient.get(`${apiPath}/full-list`);
	return response.data.message;
};

/**
 * Удаление терминала
 */
export const basicDeviceDelete = async (imei: string) => {
	const response = await axiosClient.delete(`${apiPath}/${imei}`);
	toast.success(response.data.message);
};

export const basicDeviceMassiveDelete = async (imeis: string[]) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_DELETE_MASSIVE,
		async () => {
			const response = await axiosClient.post(`${apiPath}/delete`, { imeis }, { actionKey: ACTIONKEYs.DEVICE_DELETE_MASSIVE } as any);
			toast.success(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Трансфер устройств в аккаунты
 */
export const basicDeviceTransfer = async (payload: DeviceTransferRequest) => {
	const response = await axiosClient.patch(`${apiPath}/transfer`, payload);
	toast.success(response.data.message);
};

/**
 * Обновление данных устройства
 */
export const basicDeviceUpdate = async (imei: string, payload: DeviceUpdateRequest) => {
	const response = await axiosClient.patch(`${apiPath}/${imei}`, payload);
	toast.success(response.data.message);
};
