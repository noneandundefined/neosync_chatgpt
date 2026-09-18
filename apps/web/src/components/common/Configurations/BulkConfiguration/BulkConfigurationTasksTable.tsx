import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { FetchFn } from '@/components/common/Table/GenericTable/GenericTable';
import CompactTable from '@/components/common/Table/CompactTable/CompactTable';
import { CompanyDetailsResponse, CompanyTaskResponse } from '@/rest/companyAPI';
import { BULK_TASK_STATUS_COLORS, getTaskStatusLabel, resolveTaskDisplayStatus } from '@/utils/bulkConfigurationDetailsUtils';

interface BulkConfigurationTasksTableProps {
	company: CompanyDetailsResponse;
}

const BulkConfigurationTasksTable: React.FC<BulkConfigurationTasksTableProps> = ({ company }) => {
	const { t } = useTranslation();

	const sysToken = useMemo(() => company.tasks.length + new Date(company.updated_at).getTime(), [company.tasks.length, company.updated_at]);

	const fetchTasks: FetchFn<CompanyTaskResponse> = useCallback(
		async ({ page, limit, search }) => {
			const query = search.trim().toLowerCase();

			const filtered = company.tasks.filter((task) => {
				if (!query) return true;

				return (task.device_imei ?? '').toLowerCase().includes(query) || (task.device_model ?? '').toLowerCase().includes(query);
			});

			const start = (page - 1) * limit;

			return {
				items: filtered.slice(start, start + limit),
				total: filtered.length,
				totalPages: Math.max(1, Math.ceil(filtered.length / limit)),
			};
		},
		[company.tasks]
	);

	const renderStatus = (task: CompanyTaskResponse) => {
		const displayStatus = resolveTaskDisplayStatus(task);

		return (
			<div className="flex items-center gap-2">
				<span className="h-2.5 w-2.5 rounded-full shrink-0" style={{ backgroundColor: BULK_TASK_STATUS_COLORS[displayStatus] }} />
				<span className="whitespace-nowrap">{getTaskStatusLabel(displayStatus, t)}</span>
			</div>
		);
	};

	return (
		<div className="space-y-4">
			<div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
				<p className="font-medium text-[#49525f]">{t('label.devices')}</p>
			</div>

			<CompactTable<CompanyTaskResponse>
				tableKey={`bulk-configuration-tasks-${company.id}`}
				fetchFn={fetchTasks}
				sysToken={sysToken}
				persistSearchInUrl={false}
				getRowId={(task) => task.id}
				renderTitle={(task) => task.device_imei ?? t('message.no-information-available')}
				renderSubtitle={(task) => task.device_model ?? t('message.no-information-available')}
				renderLabel={renderStatus}
				renderTime={(task) => (
					<span className="text-[#656d76] whitespace-nowrap">
						{t('label.attempts')}: {task.attempts}
					</span>
				)}
				emptyComponent={<div className="px-4 py-8 text-center text-sm text-[#656d76]">{t('message.data-available')}</div>}
			/>
		</div>
	);
};

export default BulkConfigurationTasksTable;
