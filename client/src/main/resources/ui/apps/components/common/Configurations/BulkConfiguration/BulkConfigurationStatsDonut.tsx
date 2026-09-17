import { useTranslation } from 'react-i18next';
import GUIDonut from '@/components/ui/GUIDonut';
import { CompanyTaskStatusStats } from '@/rest/companyAPI';
import { buildStatusChartSegments, BULK_TASK_STATUS_COLORS } from '@/utils/bulkConfigurationDetailsUtils';

interface BulkConfigurationStatsDonutProps {
	stats: CompanyTaskStatusStats;
}

const BulkConfigurationStatsDonut: React.FC<BulkConfigurationStatsDonutProps> = ({ stats }) => {
	const { t } = useTranslation();
	const segments = buildStatusChartSegments(stats);
	const total = stats.total || 0;

	return (
		<div className="flex items-center gap-9 flex-wrap">
			<GUIDonut
				size={160}
				stroke={20}
				segments={segments.map((segment) => ({
					key: segment.key,
					value: segment.value,
					color: segment.color,
				}))}
				center={
					<>
						<p className="text-[28px] font-semibold text-[#24292f] leading-none">{total}</p>
						<p className="mt-1 text-[12px] text-[#656d76] max-w-[90px]">{t('message.bulk-details-total-devices')}</p>
					</>
				}
			/>

			<div className="flex-1 min-w-[220px] space-y-2">
				{(['completed', 'failed', 'pending', 'running', 'expired', 'cancelled'] as const).map((key) => {
					const value = stats[key];
					const percent = total > 0 ? Math.round((value / total) * 100) : 0;

					return (
						<div key={key} className="flex items-center justify-between gap-3 text-[13px]">
							<div className="flex items-center gap-2 min-w-0">
								<span className="h-2.5 w-2.5 rounded-full shrink-0" style={{ backgroundColor: BULK_TASK_STATUS_COLORS[key] }} />
								<span className="text-[#24292f] truncate">{t(`message.bulk-task-status-${key}`)}</span>
							</div>
							<span className="text-[#656d76] whitespace-nowrap">
								{value} ({percent}%)
							</span>
						</div>
					);
				})}
			</div>
		</div>
	);
};

export default BulkConfigurationStatsDonut;
