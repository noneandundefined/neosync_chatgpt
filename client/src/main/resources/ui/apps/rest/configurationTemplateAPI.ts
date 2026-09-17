import i18next from 'i18next';
import axiosClient from './axios';
import { toast } from 'react-toastify';
import { createSSE, SseCallbacks } from './sse';
import { runGuardedAction } from '@/utils/RequestControlUtils';
import { ACTIONKEYs } from '@/constants/CacheKeys.constants';
import { ConfigurationSectionParsedResponse } from './configurationAPI';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';

const apiPath = '/devices';

export interface ConfigurationTemplateResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	name: string;
	model?: string | null;
	type_of_saving: 'modified_ones' | 'full';
}

export interface ConfigurationTemplatesWPResponse {
	items: ConfigurationTemplateResponse[];
	page: number;
	limit: number;
	total: number;
	total_all: number;
	total_pages: number;
}

export interface ConfigurationTemplateSectionParsedResponse extends Omit<ConfigurationSectionParsedResponse, 'device_imei'> {}

export interface TemplateSseCallbacks extends Omit<SseCallbacks, 'onConfiguration'> {
	onConfiguration?: (data: ConfigurationTemplateSectionParsedResponse) => void;
}

export interface ConfigurationTemplateCreateRequest {
	name: string;
	type_of_saving: 'modified_ones' | 'full';
	cfg_hash: number;
	changes: Record<string, any>;
}

export interface ConfigurationTemplateUpdateRequest {
	name?: string;
	type_of_saving?: 'modified_ones' | 'full';
	cfg_hash: number;
	changes: Record<string, any>;
}

export interface ConfigurationTemplateCreateResponse {
	id: number;
	message: string;
}

const bindTemplateSectionSse = (eventSource: EventSource, callbacks: TemplateSseCallbacks) => {
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
				const data: ConfigurationTemplateSectionParsedResponse = JSON.parse(event.data);
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
 * Получение шаблона конфигурации по id
 */
export const basicConfigurationTemplateGetById = async (id: number, signal?: AbortSignal): Promise<ConfigurationTemplateResponse> => {
	const response = await axiosClient.get(`${apiPath}/configuration-templates/${id}`, { signal });
	return response.data.message;
};

/**
 * Обновление шаблона конфигурации
 */
export const basicConfigurationTemplateUpdate = async (id: number, payload: ConfigurationTemplateUpdateRequest): Promise<void> => {
	return await runGuardedAction(
		`configuration-template:update:${id}`,
		async () => {
			const response = await axiosClient.put(`${apiPath}/configuration-templates/${id}`, payload);
			const message = typeof response.data.message === 'string' ? response.data.message : response.data.message?.message;

			toast.success(message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
		}
	);
};

/**
 * Получение списка шаблонов конфигурации
 */
export const basicConfigurationTemplatesGet = async (signal?: AbortSignal): Promise<ConfigurationTemplatesWPResponse> => {
	const response = await axiosClient.get(`${apiPath}/configuration-templates?page=1&limit=100`, { signal });
	return response.data.message;
};

/**
 * Получение списка шаблонов конфигурации с пагинацией
 */
export const basicConfigurationTemplatesGetWithParams = async (page: number, limit: number, search: string = '', columnSort: ColumnSortState, signal?: AbortSignal): Promise<ConfigurationTemplatesWPResponse> => {
	const response = await axiosClient.get(`${apiPath}/configuration-templates?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}&columnSortKey=${columnSort.key}&columnSortDir=${columnSort.direction}`, { signal });
	return response.data.message;
};

/**
 * Создание шаблона конфигурации
 */
export const basicConfigurationTemplateCreate = async (payload: ConfigurationTemplateCreateRequest): Promise<ConfigurationTemplateCreateResponse> => {
	return await runGuardedAction(
		`configuration-template:create:${payload.name}`,
		async () => {
			const response = await axiosClient.post(`${apiPath}/configuration-templates`, payload);
			const message = typeof response.data.message === 'string' ? response.data.message : response.data.message?.message;

			toast.success(message);
			return response.data.message;
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
		}
	);
};

/**
 * Удаление шаблона конфигурации
 */
export const basicConfigurationTemplateDelete = async (id: number): Promise<void> => {
	const response = await axiosClient.delete(`${apiPath}/configuration-templates/${id}`);
	toast.success(response.data.message);
};

/**
 * Запрос блока конфигурации шаблона по секции (значения Default из схемы)
 */
export const basicConfigurationTemplateSectionParsedSSEGet = (section: string, callbacks: TemplateSseCallbacks): EventSource => {
	const eventSource = createSSE(`${apiPath}/templates/configuration/parsed/${section}`);
	return bindTemplateSectionSse(eventSource, callbacks);
};

/**
 * Запрос блока сохранённого шаблона конфигурации по секции
 */
export const basicConfigurationTemplateSectionParsedSSEGetById = (templateId: number, section: string, callbacks: TemplateSseCallbacks): EventSource => {
	const eventSource = createSSE(`${apiPath}/configuration-templates/${templateId}/configuration/parsed/${section}`);
	return bindTemplateSectionSse(eventSource, callbacks);
};

/**
 * Применение шаблона конфигурации к устройству
 */
export const basicConfigurationTemplateApply = async (imei: string, templateId: number) => {
	await runGuardedAction(
		ACTIONKEYs.CONFIGURATION_APPLY(imei),
		async () => {
			const response = await axiosClient.post(`${apiPath}/${imei}/configuration/apply-template/${templateId}`, undefined, {
				actionKey: ACTIONKEYs.CONFIGURATION_APPLY(imei),
			} as any);

			toast.success(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
		}
	);
};
