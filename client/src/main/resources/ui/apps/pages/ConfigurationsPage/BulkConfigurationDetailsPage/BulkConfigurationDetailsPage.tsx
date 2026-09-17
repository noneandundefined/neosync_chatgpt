import { useCallback } from 'react';
import PageLayout from '@/pages/PageLayout';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { formatLocalDate } from '@/utils/TimeUtils';
import Fallback from '@/components/Fallback/Fallback';
import { basicCompanyGetById } from '@/rest/companyAPI';
import { useNavigate, useParams } from 'react-router-dom';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useBulkConfigurationDetailsActions } from '@/hooks/Configurations/useBulkConfigurationDetailsActions';
import BulkConfigurationTasksTable from '@/components/common/Configurations/BulkConfiguration/BulkConfigurationTasksTable';
import BulkConfigurationStatsDonut from '@/components/common/Configurations/BulkConfiguration/BulkConfigurationStatsDonut';
import { exportCompanyTasksCsv, getCompanyStatusClass, getCompanyStatusLabel, getConfigurationSourceLabel, getLaunchModeLabel, getTtlLabel } from '@/utils/bulkConfigurationDetailsUtils';

import Reload from '@/components/@icons/reload';
import ArrowLeft from '@/components/@icons/arrow-left';
import ContentCopy from '@/components/@icons/content-copy';
import ExportVariant from '@/components/@icons/export-variant';
import CloseCircleOutline from '@/components/@icons/close-circle-outline';

const BulkConfigurationDetailsPage = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();

	const { id } = useParams();

	const companyId = Number(id);
	if (!id || !Number.isFinite(companyId) || companyId <= 0) {
		return <Fallback />;
	}

	const fetchCompany = useCallback(() => basicCompanyGetById(companyId), [companyId]);
	const { data: respCompanyGetById, loading: respCompanyGetByIdLoading } = useHandleServer(['respCompanyGetById', companyId], fetchCompany, {
		refetchInterval: (query) => {
			const status = query.state.data?.status;
			return status === 'pending' || status === 'running' ? 5000 : false;
		},
	});

	const { actionLoading, handleRetryFailed, handleCancelPending, handleCreateBasedOnThis } = useBulkConfigurationDetailsActions(companyId, respCompanyGetById ?? undefined);

	if (!respCompanyGetByIdLoading && !respCompanyGetById) {
		return <Fallback />;
	}

	const paramRows = respCompanyGetById
		? [
				{ label: t('message.provisioning-task-ttl'), value: getTtlLabel(respCompanyGetById.ttl, t) },
				{ label: t('message.bulk-details-skip-active-task'), value: respCompanyGetById.existing_task_action === 'skip' ? t('label.yes') : t('label.no') },
				{ label: t('message.provisioning-launch-mode'), value: getLaunchModeLabel(respCompanyGetById.launch_mode, t) },
				{ label: t('message.provisioning-configuration-source'), value: getConfigurationSourceLabel(respCompanyGetById.configuration_source, t) },
				{ label: t('message.bulk-details-created'), value: formatLocalDate(respCompanyGetById.created_at) },
				{ label: t('message.bulk-details-launched'), value: respCompanyGetById.launched_at ? formatLocalDate(respCompanyGetById.launched_at) : t('message.no-information-available') },
				{ label: t('message.bulk-details-expires'), value: formatLocalDate(respCompanyGetById.expires_at) },
			]
		: [];

	return (
		<PageLayout>
			<div className="space-y-5">
				<div className="flex items-start gap-3">
					<button type="button" className="flex items-center gap-2 -ml-2 p-1 px-2 rounded-md hover:bg-[#f3f3f3]" onClick={() => navigate(`${ROUTES.CONFIGURATIONS}?tab=companies`)}>
						<ArrowLeft fill="#999" size={15} />
						<span className="text-[#999] text-[13px]">{t('label.bulk-configurations')}</span>
					</button>
				</div>

				{respCompanyGetByIdLoading || !respCompanyGetById ? (
					<Loading />
				) : (
					<>
						<div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
							<div className="flex items-center gap-5">
								<h1 className="text-[22px] font-semibold text-[#49525f] leading-snug">{respCompanyGetById.name}</h1>
								<span className={`inline-flex items-center rounded-[6px] border px-2 text-[13px] ${getCompanyStatusClass(respCompanyGetById.status)}`}>{getCompanyStatusLabel(respCompanyGetById.status, t)}</span>
							</div>
						</div>

						<div className="grid grid-cols-1 xl:grid-cols-3 gap-4">
							<div className="border border-[#e4e4e4] rounded bg-white p-5 xl:col-span-1">
								<p className="font-medium text-[#49525f] mb-4">{t('message.bulk-details-status-stats')}</p>
								<BulkConfigurationStatsDonut stats={respCompanyGetById.status_stats} />
							</div>

							<div className="border border-[#e4e4e4] rounded bg-white p-5 xl:col-span-1">
								<p className="font-medium text-[#49525f] mb-4">{t('message.bulk-details-campaign-params')}</p>
								<div className="space-y-3">
									{paramRows.map((row) => (
										<div key={row.label} className="flex items-start justify-between gap-4 text-[13px] border-b border-[#eef1f4] pb-3 last:border-b-0 last:pb-0">
											<span className="text-[#656d76]">{row.label}</span>
											<span className="text-[#24292f] text-right">{row.value}</span>
										</div>
									))}
								</div>
							</div>

							<div className="border border-[#e4e4e4] rounded bg-white p-5 xl:col-span-1">
								<p className="font-medium text-[#49525f] mb-4">{t('message.bulk-details-actions')}</p>
								<div className="space-y-3">
									<button
										type="button"
										className="w-full flex items-center justify-center gap-2 border border-[#395d95] text-[#395d95] rounded px-4 py-3 text-[14px] hover:bg-[#E8EDF5] disabled:opacity-40"
										disabled={respCompanyGetById.status_stats.failed === 0 || actionLoading !== null}
										onClick={handleRetryFailed}
									>
										<Reload size={18} fill="#395d95" />
										{t('message.bulk-details-retry-errors')}
									</button>

									<button
										type="button"
										className="w-full flex items-center justify-center gap-2 border border-[#d1242f] text-[#d1242f] rounded px-4 py-3 text-[14px] hover:bg-[#ffecef] disabled:opacity-40"
										disabled={respCompanyGetById.status_stats.pending + respCompanyGetById.status_stats.running === 0 || actionLoading !== null}
										onClick={handleCancelPending}
									>
										<CloseCircleOutline size={18} fill="#d1242f" />
										{t('message.bulk-details-cancel-pending')}
									</button>

									<button
										type="button"
										className="w-full flex items-center justify-center gap-2 border rounded px-4 py-3 text-[14px] hover:bg-[#E8EDF5] hover:border-[#7d92b4]"
										onClick={() => exportCompanyTasksCsv(respCompanyGetById, t)}
									>
										<ExportVariant size={18} fill="#49525f" />
										{t('message.bulk-details-download-csv')}
									</button>

									<button type="button" className="w-full flex items-center justify-center gap-2 border rounded px-4 py-3 text-[14px] hover:bg-[#E8EDF5] hover:border-[#7d92b4]" onClick={handleCreateBasedOnThis}>
										<ContentCopy size={18} fill="#49525f" />
										{t('message.bulk-details-create-based-on-this')}
									</button>
								</div>
							</div>
						</div>

						<BulkConfigurationTasksTable company={respCompanyGetById} />
					</>
				)}
			</div>
		</PageLayout>
	);
};

export default BulkConfigurationDetailsPage;
