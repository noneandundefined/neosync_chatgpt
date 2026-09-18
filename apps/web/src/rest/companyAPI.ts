import axiosClient from './axios';
import { toast } from 'react-toastify';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { runSpamGuardedAction } from '@/utils/SpamGuardUtils';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';

const apiPath = '/companies';

export interface CompanyResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	name: string;
	status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
	ttl: number;
	existing_task_action: 'skip' | 'overwrite' | 'replace';
	launch_mode: 'all' | 'selected';
	incompatible_action: 'exclude' | 'fail';
	configuration_source: 'device_sources' | 'template_sources' | 'file_sources';
	tasks_count?: number;
}

export interface CompanyTaskResponse {
	id: number;
	created_at: string;
	updated_at: string;
	company_id: number;
	device_id: number;
	device_imei?: string;
	device_model?: string;
	status: 'pending' | 'queued' | 'sending' | 'awaiting_confirmation' | 'completed' | 'failed' | 'cancelled' | 'expired';
	expires_at: string;
	last_sent_at?: string | null;
	confirmation_deadline?: string | null;
	next_attempt_at?: string | null;
	error_message?: string | null;
	attempts: number;
}

export interface CompanyTaskStatusStats {
	total: number;
	completed: number;
	failed: number;
	pending: number;
	running: number;
	cancelled: number;
	expired: number;
}

export interface CompanyDetailsResponse extends CompanyResponse {
	source_metadata?: {
		version?: number;
		priority_model?: string;
		selected_imeis?: string[];
		configuration_template_id?: number | null;
		configuration_device_cource?: string | null;
		configuration_file_name?: string | null;
	};
	configuration_snapshot_base64?: string;
	tasks: CompanyTaskResponse[];
	status_stats: CompanyTaskStatusStats;
	expires_at: string;
	launched_at?: string | null;
}

export interface CompaniesWPResponse {
	items: CompanyResponse[];
	page: number;
	limit: number;
	total: number;
	total_all: number;
	total_pages: number;
}

export interface CompanyCreatePayload extends CompanyCreateRequest {
	imeis: string[];
}

/**
 * Создание массовой настройки
 */
export const basicCompanyCreate = async (payload: CompanyCreatePayload): Promise<void> => {
	return await runSpamGuardedAction(
		`company:create:${payload.name}:${payload.ttl}`,
		async () => {
			const response = await axiosClient.post(`${apiPath}`, payload);
			toast.success(response.data.message);
		},
		{
			maxAttempts: 6,
			windowMs: 10_000,
			blockMs: 60_000,
			onBlocked: ({ message }) => {
				showServerErrorToast(message);
			},
		}
	);
};

/**
 * Получение списка массовых настроек
 */
export const basicCompaniesGet = async (signal?: AbortSignal): Promise<CompaniesWPResponse> => {
	const response = await axiosClient.get(`${apiPath}?page=1&limit=100`, { signal });
	return response.data.message;
};

/**
 * Получение списка массовых настроек с пагинацией
 */
export const basicCompaniesGetWithParams = async (page: number, limit: number, search: string = '', columnSort: ColumnSortState, signal?: AbortSignal): Promise<CompaniesWPResponse> => {
	const response = await axiosClient.get(`${apiPath}?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}&columnSortKey=${columnSort.key}&columnSortDir=${columnSort.direction}`, { signal });
	return response.data.message;
};

/**
 * Получение массовой настройки по ID
 */
export const basicCompanyGetById = async (id: number, signal?: AbortSignal): Promise<CompanyDetailsResponse> => {
	const response = await axiosClient.get(`${apiPath}/${id}`, { signal });
	return response.data.message;
};

export const basicCompanyRetryFailed = async (id: number): Promise<string> => {
	const response = await axiosClient.post(`${apiPath}/${id}/retry-failed`);
	return response.data.message.message;
};

export const basicCompanyCancelPending = async (id: number): Promise<string> => {
	const response = await axiosClient.post(`${apiPath}/${id}/cancel-pending`);
	return response.data.message.message;
};
