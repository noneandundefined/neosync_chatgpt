import i18next from 'i18next';
import axiosClient from './axios';
import { toast } from 'react-toastify';
import { createSSE, SseCallbacks } from './sse';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { FieldType } from '@/constants/FieldType.constant';
import { ACTIONKEYs } from '@/constants/CacheKeys.constants';
import { runSpamGuardedAction } from '@/utils/SpamGuardUtils';
import { EXPORT_CFG_FILE } from '@/constants/Files.constants';
import { EVENT_CONFIGURATION_DRAFT } from '@/hooks/useDraftCheck';
import { ClientGuardError, runGuardedAction } from '@/utils/RequestControlUtils';
import { CACHEKEYs_TERMINAL_LLS_ADDRESS } from '@/constants/CacheKeys.constants';

import { TelemetryResponse } from '@/interface/configuration/configurationTelemetryResponse.interface';
import { ConfigurationChangePassRequest } from '@/interface/configuration/configurationChangePassRequest.interface';
import { ConfigurationDraftInsertRequest } from '@/interface/configuration/configurationDraftInsertRequest.interface';
import { ConfigurationExportLLSTariration } from '@/interface/configuration/configurationExportLLSTariration.interface';

const apiPath = '/devices';

const TELEMETRY_REBOOT_COOLDOWN_MS = 15_000;

export interface FieldConfigurationParsed {
	uid: string;
	size: number;
	value: any;
}

export interface FieldSchema {
	uid: number;
	resolve_uid: number;
	name: string;
	type: FieldType;
	default?: any;
	min?: number;
	max?: number;
	min_length?: number;
	max_length?: number;
	max_items?: number;
	is_digits?: boolean;
	section: string;
}

export interface ConfigurationParsedResponse {
	device_imei: string;
	cfg_hash: number;
	config_parsed: FieldConfigurationParsed[];
	schema: FieldSchema[];
	timestamp: number;
}

export interface ConfigurationSectionParsedResponse {
	device_imei: string;
	cfg_hash: number;
	section: string;
	config_parsed: FieldConfigurationParsed[];
	schema: FieldSchema[];
	timestamp: number;
}

export interface ConfigurationRawResponse {
	device_imei: string;
	cfg_hash: number;
	config_hex: string;
	timestamp: number;
}

export interface ConfigurationDraftResponse {
	device_imei: string;
	cfg_hash: number;
	section: string;
	changes: Record<string, any>;
	schema: FieldSchema[];
	timestamp: number;
}

export interface ConfigurationHistory {
	id: number;
	apply_at: string;
	device_id: number;
	cfg_hash: number;
	cfg_data: Uint8Array;
}

/**
 * Запрос всей конфигурации терминала
 */
export const basicConfigurationParsed = async (imei: string, signal?: AbortSignal): Promise<ConfigurationParsedResponse> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/parsed`, { signal });
	return response.data.message;
};

/**
 * Запрос всего hex пакета конфигурации терминала
 */
export const basicConfigurationRaw = async (imei: string, signal?: AbortSignal): Promise<ConfigurationRawResponse> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/raw`, { signal });
	return response.data.message;
};

/**
 * Запрос одного блока конфигурации
 */
export const basicConfigurationSectionParsedSSEGet = (imei: string, section: string, callbacks: SseCallbacks): EventSource => {
	const eventSource = createSSE(`${apiPath}/${imei}/configuration/parsed/${section}`);

	eventSource.addEventListener('open', (event: MessageEvent) => {
		if (event.data) {
			try {
				const data = JSON.parse(event.data);
				callbacks.onOpen?.(data.message);
			} catch {}
		} else {
			callbacks.onOpen?.('connected');
		}
	});

	eventSource.addEventListener('progress', (event: MessageEvent) => {
		if (event.data) {
			try {
				const data = JSON.parse(event.data);
				callbacks.onProgress?.(data.message);
			} catch {}
		}
	});

	eventSource.addEventListener('configuration', (event: MessageEvent) => {
		if (event.data) {
			try {
				const data: ConfigurationSectionParsedResponse = JSON.parse(event.data);
				callbacks.onConfiguration?.(data);
			} catch {}
		}
	});

	eventSource.addEventListener('error', (event: MessageEvent) => {
		if (event.data) {
			try {
				const data = JSON.parse(event.data);
				callbacks.onError?.(data.message);
			} catch {}
		}
	});

	eventSource.addEventListener('done', (event: MessageEvent) => {
		if (event.data) {
			try {
				const data = JSON.parse(event.data);
				callbacks.onDone?.(data.message);
			} catch {}
		}

		callbacks.onClose?.();
		eventSource.close();
	});

	eventSource.onerror = () => {
		callbacks.onError?.(i18next.t('message.sse-connection-error'));
		callbacks.onClose?.();
		eventSource.close();
	};

	return eventSource;
};

/**
 * Запрос черновика измененной конфигурации
 */
export const basicConfigurationGetDraft = async (imei: string, section: string, signal?: AbortSignal): Promise<ConfigurationDraftResponse> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/draft/${section}`, { signal });
	return response.data.message;
};

/**
 * Запрос истории конфигурации устройства
 */
export const basicConfigurationGetHistories = async (imei: string, signal?: AbortSignal): Promise<ConfigurationHistory[]> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/histories`, { signal });
	return response.data.message;
};

/**
 * Импорт истории конфигурации устройства из истории конфигурации
 */
// export const basicConfigurationImportByHistory = async (imei: string, cfgHash: number): Promise<void> => {
// 	await runGuardedAction(
// 		ACTIONKEYs.CONFIGURATION_IMPORT_BY_HISTORY(imei, cfgHash),
// 		async () => {
// 			const payload: ConfigurationImportByHistoryRequest = {
// 				cfg_hash: cfgHash,
// 			};

// 			const response = await axiosClient.post(`${apiPath}/${imei}/configuration/histories/import`, payload, {
// 				actionKey: ACTIONKEYs.CONFIGURATION_IMPORT_BY_HISTORY(imei, cfgHash),
// 			} as any);
// 			toast.success(response.data.message);
// 		},
// 		{
// 			cooldownMs: 2000,
// 			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
// 			onBlocked: (msg) => showServerErrorToast(msg),
// 		}
// 	);
// };

/**
 * Сброс черновика измененной конфигурации
 */
export const basicConfigurationResetDraft = async (imei: string) => {
	await axiosClient.delete(`${apiPath}/${imei}/configuration/draft/reset`);
	sessionStorage.removeItem(CACHEKEYs_TERMINAL_LLS_ADDRESS(imei));

	/** Clear GUI event */
	window.dispatchEvent(new CustomEvent(EVENT_CONFIGURATION_DRAFT, { detail: false }));
};

/**
 * Запись черновика изменения конфигурации
 */
export const basicConfigurationInsertDraft = async (imei: string, section: string, cfg_hash: number, changes: Record<string, any>) => {
	const payload: ConfigurationDraftInsertRequest = {
		device_imei: imei,
		cfg_hash: cfg_hash,
		changes: changes,
		timestamp: Math.floor(Date.now() / 1000),
	};

	await runSpamGuardedAction(
		`configuration-draft:${imei}:${section}`,
		async () => {
			await axiosClient.post(`${apiPath}/${imei}/configuration/draft/${section}`, payload);
		},
		{
			maxAttempts: 20,
			windowMs: 10_000,
			blockMs: 30_000,
			onBlocked: ({ message }) => {
				showServerErrorToast(message);
			},
		}
	);
};

/**
 * Приминение конфигурации
 */
export const basicConfigurationApply = async (imei: string) => {
	await runGuardedAction(
		ACTIONKEYs.CONFIGURATION_APPLY(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/configuration/apply`, undefined, {
				actionKey: ACTIONKEYs.CONFIGURATION_APPLY(imei),
			} as any);

			sessionStorage.removeItem(CACHEKEYs_TERMINAL_LLS_ADDRESS(imei));

			/** Clear GUI event */
			window.dispatchEvent(new CustomEvent(EVENT_CONFIGURATION_DRAFT, { detail: false }));

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
 * Экспорт конфигурации в файл
 */
export const basicConfigurationExport = async (imei: string) => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/export`, {
		responseType: 'blob',
	});

	const blob = new Blob([response.data], {
		type: 'application/octet-stream',
	});

	const url = window.URL.createObjectURL(blob);
	const link = document.createElement('a');

	link.href = url;
	link.setAttribute('download', EXPORT_CFG_FILE);
	document.body.appendChild(link);

	link.click();
	link.remove();

	toast.success(i18next.t('message.file-export-success'));
};

/**
 * Импорт конфигурации из файла
 */
export const basicConfigurationImport = async (imei: string, file: File) => {
	await runGuardedAction(
		ACTIONKEYs.CONFIGURATION_IMPORT(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/configuration/import`, file, {
				headers: {
					'Content-Type': 'application/octet-stream',
				},
				actionKey: ACTIONKEYs.CONFIGURATION_IMPORT(imei),
			} as any);

			sessionStorage.removeItem(CACHEKEYs_TERMINAL_LLS_ADDRESS(imei));

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
 * Экспорт тарирование датчиков терминала в файл
 */
export const basicConfigurationExportLLSTariration = async (imei: string, payload: ConfigurationExportLLSTariration) => {
	const response = await axiosClient.post(`${apiPath}/${imei}/configuration/export/llstariration`, payload, { responseType: 'blob' });

	let filename = 'LLS Tariration.csv';

	const disposition = response.headers['content-disposition'];
	if (disposition) {
		const match = disposition.match(/filename="?(.+)"?/);
		if (match) filename = decodeURIComponent(match[1]);
	}

	const url = URL.createObjectURL(response.data);
	const link = document.createElement('a');

	link.href = url;
	link.download = filename;
	link.click();

	URL.revokeObjectURL(url);

	toast.success(i18next.t('message.file-export-llstariration-success'));
};

/**
 * Получение конфигурации телеметрии данных устройства
 */
export const basicConfigurationTelemetry = async (imei: string, signal?: AbortSignal): Promise<TelemetryResponse> => {
	const response = await axiosClient.get(`${apiPath}/${imei}/configuration/telemetry`, { signal });
	return response.data.message;
};

/**
 * Считывание конфигурации телеметрии данных устройства
 */
export const basicRebootConfigurationTelemetry = async (imei: string) => {
	try {
		await runGuardedAction(
			ACTIONKEYs.CONFIGURATION_TELEMETRY_REBOOT(imei),
			async () => {
				const response = await axiosClient.post(`${apiPath}/${imei}/configuration/telemetry`, undefined, {
					actionKey: ACTIONKEYs.CONFIGURATION_TELEMETRY_REBOOT(imei),
				} as any);
				toast.success(response.data.message);
			},
			{ cooldownMs: TELEMETRY_REBOOT_COOLDOWN_MS }
		);
	} catch (err) {
		if (err instanceof ClientGuardError) {
			toast.success(i18next.t('message.sensor-data-updated'));
			return;
		}
		throw err;
	}
};

/**
 * Считывание конфигурации телеметрии данных устройства
 */
export const basicConfigurationChangePass = async (imei: string, payload: ConfigurationChangePassRequest) => {
	const response = await axiosClient.post(`${apiPath}/${imei}/configuration/pass`, payload);
	toast.success(response.data.message);
};
