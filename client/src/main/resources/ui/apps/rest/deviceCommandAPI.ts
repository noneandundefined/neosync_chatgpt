import i18n from '@/utils/i18n';
import axiosClient from './axios';
import { generateUUID } from '@/utils/UuidUtils';
import { runGuardedAction } from '@/utils/RequestControlUtils';
import { ACTIONKEYs, CACHEKEYs } from '@/constants/CacheKeys.constants';
import { showServerErrorToast, showServerSuccessToast } from '@/utils/ToastUtils';
import { DeviceSendCommandRequest } from '@/interface/device/deviceSendCommandRequest.interface';

const apiPath = '/devices';

export interface DeviceCommand {
	id: number;
	created_at: string;
	updated_at: string;
	execute_until: string;
	user_uuid: string;
	send_mode: 'instant' | 'on_connect';
	status: 'pending' | 'inprogress' | 'completed' | 'notcompleted' | 'executionerror';
	command: string;
}

export interface DeviceCommandExecution {
	id: number;
	created_at: string;
	updated_at: string;
	command_id: number;
	device_imei: string;
	status: 'pending' | 'inprogress' | 'completed' | 'notcompleted' | 'executionerror';
	response: string | null;
}

export interface DeviceCommandExecutionShort {
	imei: string;
	status: 'pending' | 'inprogress' | 'completed' | 'notcompleted' | 'executionerror';
	response: string | null;
}

export interface DeviceCommandWithExecutions {
	id: number;
	created_at: string;
	updated_at: string;
	execute_until: string;
	user_uuid: string;
	send_mode: 'instant' | 'on_connect';
	status: 'pending' | 'inprogress' | 'notcompleted' | 'completed' | 'executionerror';
	command: string;
	executions: DeviceCommandExecutionShort[];
}

export interface DeviceShort {
	imei: string;
	status: boolean;
	activated: boolean;
	group_id?: number;
}

export interface DevicesForSendCommand {
	group_id: number;
	name: string | null;
	device_count: number;
}

export interface DevicesStateDevicesResponse {
	items: DeviceShort[];
	page: number;
	limit: number;
	total: number;
	has_more: boolean;
}

/**
 * Отправка команды
 */
export const basicDeviceSendCommand = async (payload: DeviceSendCommandRequest) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_COMMAND_SEND(payload.command),
		async () => {
			let sessionId = sessionStorage.getItem(CACHEKEYs.NEOSYNC_CMD_SESSION_ID);
			if (!sessionId) {
				sessionId = generateUUID();
				sessionStorage.setItem(CACHEKEYs.NEOSYNC_CMD_SESSION_ID, sessionId);
			}
			payload.sess_id = sessionId;

			await axiosClient.post(`${apiPath}/cmd`, payload, { actionKey: ACTIONKEYs.DEVICE_COMMAND_SEND } as any);
		},
		{
			cooldownMs: 1500,
			getCooldownMessage: (seconds) => i18n.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Отправка команды с ожиданием ответа
 */
export const basicDeviceSendCommandAndWait = async (imei: string, command: string) => {
	const response = await axiosClient.post(`${apiPath}/${imei}/cmd/send-and-wait`, { command });
	return response.data.message;
};

/**
 * Получение истории команд
 */
export const basicDeviceCommandHistory = async (): Promise<DeviceCommand[]> => {
	const response = await axiosClient.get(`${apiPath}/cmd`);
	return response.data.message;
};

/**
 * Получение истории команд по сессии
 */
export const basicDeviceCommandRecentHistory = async (): Promise<DeviceCommandWithExecutions[]> => {
	let sessionId = sessionStorage.getItem(CACHEKEYs.NEOSYNC_CMD_SESSION_ID);
	if (!sessionId) {
		sessionId = generateUUID();
		sessionStorage.setItem(CACHEKEYs.NEOSYNC_CMD_SESSION_ID, sessionId);
	}

	const response = await axiosClient.get(`${apiPath}/cmd/sessions/${sessionId}`);
	return response.data.message;
};

/**
 * Удаление команды
 */
export const basicDeviceCommandDelete = async (id: number) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_COMMAND_DELETE(id),
		async () => {
			const response = await axiosClient.delete(`${apiPath}/cmd/${id}/a/delete`, { actionKey: ACTIONKEYs.DEVICE_COMMAND_DELETE(id) } as any);
			showServerSuccessToast(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18n.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Отмена команды
 */
export const basicDeviceCommandCancel = async (id: number) => {
	const response = await axiosClient.patch(`${apiPath}/cmd/${id}/a/cancel`, undefined, { actionKey: `device:command:cancel:${id}` } as any);
	showServerSuccessToast(response.data.message);
};

/**
 * Получение деталей команды
 */
export const basicDeviceCommandById = async (id: number): Promise<DeviceCommandWithExecutions> => {
	const response = await axiosClient.get(`${apiPath}/cmd/${id}`);
	return response.data.message;
};

/**
 * Отправка команды перезагрузки
 */
export const basicDeviceCommandReboot = async (imei: string) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_COMMAND_REBOOT(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/cmd/reboot`, undefined, { actionKey: ACTIONKEYs.DEVICE_COMMAND_REBOOT(imei) } as any);
			showServerSuccessToast(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18n.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Отправка команды сброса до заводских настроек
 */
export const basicDeviceCommandEraseEeprom = async (imei: string) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_COMMAND_ERASE_EEPROM(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/cmd/erase_eeprom`, undefined, { actionKey: ACTIONKEYs.DEVICE_COMMAND_ERASE_EEPROM(imei) } as any);
			showServerSuccessToast(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18n.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Отправка команды очистка памяти устройствва
 */
export const basicDeviceCommandEraseFlash = async (imei: string) => {
	await runGuardedAction(
		ACTIONKEYs.DEVICE_COMMAND_ERASE_FLASH(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/cmd/erase_flash`, undefined, { actionKey: ACTIONKEYs.DEVICE_COMMAND_ERASE_FLASH(imei) } as any);
			showServerSuccessToast(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18n.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Отправка команды поиска BLE датчиков
 */
export const basicDeviceCommandFindBleSensors = async (imei: string, command: string) => {
	const response = await axiosClient.post(`${apiPath}/${imei}/cmd/find_ble_sensors`, { command }, { actionKey: ACTIONKEYs.DEVICE_COMMAND_FIND_BLE(imei) } as any);
	showServerSuccessToast(response.data.message);
};

/**
 * Отправка команды поиска датчиков температуры
 */
export const basicDeviceCommandOWireSensors = async (imei: string) => {
	const response = await axiosClient.post(`${apiPath}/${imei}/cmd/find_owire_sensors`, undefined, { actionKey: ACTIONKEYs.DEVICE_COMMAND_FIND_OWIRE(imei) } as any);
	showServerSuccessToast(response.data.message);
};
