import { TFunction } from 'i18next';
import { CompanyDetailsResponse, CompanyTaskResponse } from '@/rest/companyAPI';
import { encodeBulkCsv } from './bulkConfigurationCsv';

export type BulkTaskDisplayStatus = 'completed' | 'failed' | 'pending' | 'running' | 'cancelled' | 'expired' | 'queued' | 'sending' | 'awaiting_confirmation';

export const BULK_TASK_STATUS_COLORS: Record<BulkTaskDisplayStatus, string> = {
	completed: '#1e772e',
	failed: '#e2483d',
	pending: '#0969da',
	running: '#cd9a00',
	queued: '#cd9a00',
	sending: '#cd9a00',
	awaiting_confirmation: '#cd9a00',
	cancelled: '#8c959f',
	expired: '#49525f',
};

export const resolveTaskDisplayStatus = (task: CompanyTaskResponse): BulkTaskDisplayStatus => task.status;

export const getTaskStatusLabel = (status: BulkTaskDisplayStatus, t: TFunction) => t(`message.bulk-task-status-${status}`);

export const getCompanyStatusLabel = (status: CompanyDetailsResponse['status'], t: TFunction) => t(`message.bulk-company-status-${status}`);

export const getCompanyStatusClass = (status: CompanyDetailsResponse['status']) => {
	switch (status) {
		case 'running':
			return 'bg-[#dafbe1] text-[#1e772e] border-[#1e772e]';
		case 'completed':
			return 'bg-[#dafbe1] text-[#116329] border-[#aceebb]';
		case 'failed':
			return 'bg-[#ffecef] text-[#d1242f] border-[#d1242f]';
		case 'cancelled':
			return 'bg-[#f6f8fa] text-[#656d76] border-[#d0d7de]';
		default:
			return 'bg-[#fffcec] text-[#cd9a00] border-[#cd9a00]';
	}
};

export const getConfigurationSourceLabel = (source: CompanyDetailsResponse['configuration_source'], t: TFunction) => {
	switch (source) {
		case 'device_sources':
			return t('message.provisioning-configuration-source-device');
		case 'template_sources':
			return t('message.provisioning-configuration-source-template');
		case 'file_sources':
			return t('message.provisioning-configuration-source-file');
		default:
			return source;
	}
};

export const getExistingTaskActionLabel = (action: CompanyDetailsResponse['existing_task_action'], t: TFunction) => {
	if (action === 'skip') return t('message.provisioning-action-skip');
	return t('message.provisioning-action-replace');
};

export const getLaunchModeLabel = (mode: CompanyDetailsResponse['launch_mode'], t: TFunction) => {
	if (mode === 'all') return t('message.provisioning-launch-mode-all');
	return t('message.provisioning-launch-mode-selected');
};

export const getTtlLabel = (ttl: number, t: TFunction) => {
	if (ttl === 1) return t('message.provisioning-ttl-1-day');
	if (ttl === 7) return t('message.provisioning-ttl-7-days');
	if (ttl === 30) return t('message.provisioning-ttl-30-days');
	return t('message.ttl-days', { count: ttl });
};

export const exportCompanyTasksCsv = (company: CompanyDetailsResponse, t: TFunction) => {
	const metadata = company.source_metadata;
	const common = {
		csv_version: '1',
		company_id: company.id,
		name: company.name,
		priority_model: metadata?.priority_model ?? '',
		ttl: company.ttl,
		existing_task_action: company.existing_task_action,
		launch_mode: company.launch_mode,
		incompatible_action: company.incompatible_action,
		configuration_cource: company.configuration_source,
		configuration_source_label: getConfigurationSourceLabel(company.configuration_source, t),
		configuration_template_id: metadata?.configuration_template_id ?? '',
		configuration_device_cource: metadata?.configuration_device_cource ?? '',
		configuration_file_name: metadata?.configuration_file_name ?? '',
		configuration_file_relative_path: metadata?.configuration_file_name ?? '',
		configuration_file_base64: company.configuration_source === 'file_sources' ? company.configuration_snapshot_base64 ?? '' : '',
		configuration_snapshot_base64: company.configuration_snapshot_base64 ?? '',
		selected_imeis_json: JSON.stringify(metadata?.selected_imeis ?? company.tasks.flatMap((task) => task.device_imei ? [task.device_imei] : [])),
		company_status: company.status,
		created_at: company.created_at,
		updated_at: company.updated_at,
		launched_at: company.launched_at ?? '',
		expires_at: company.expires_at,
		metadata_available: metadata?.version === 1,
	};
	const rows = (company.tasks.length ? company.tasks : [null]).map((task) => ({
		...common,
		task_id: task?.id ?? '',
		device_id: task?.device_id ?? '',
		imei: task?.device_imei ?? '',
		device_model: task?.device_model ?? '',
		task_status: task?.status ?? '',
		task_status_label: task ? getTaskStatusLabel(resolveTaskDisplayStatus(task), t) : '',
		attempts: task?.attempts ?? '',
		error_message: task?.error_message ?? '',
		task_created_at: task?.created_at ?? '',
		task_updated_at: task?.updated_at ?? '',
	}));
	const blob = new Blob([encodeBulkCsv(rows)], { type: 'text/csv;charset=utf-8;' });
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = `neosync_company_${company.name.replace(/[^\w\-]+/g, '_')}.csv`;
	link.click();
	URL.revokeObjectURL(url);
};

export const buildStatusChartSegments = (stats: CompanyDetailsResponse['status_stats']) => {
	const items: Array<{ key: BulkTaskDisplayStatus; value: number; color: string }> = [
		{ key: 'completed', value: stats.completed, color: BULK_TASK_STATUS_COLORS.completed },
		{ key: 'failed', value: stats.failed, color: BULK_TASK_STATUS_COLORS.failed },
		{ key: 'pending', value: stats.pending, color: BULK_TASK_STATUS_COLORS.pending },
		{ key: 'running', value: stats.running, color: BULK_TASK_STATUS_COLORS.running },
		{ key: 'expired', value: stats.expired, color: BULK_TASK_STATUS_COLORS.expired },
		{ key: 'cancelled', value: stats.cancelled, color: BULK_TASK_STATUS_COLORS.cancelled },
	];

	return items.filter((item) => item.value > 0);
};
