import Widget from './Widget';
import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import GUIDonut from '@/components/ui/GUIDonut';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { useAnalyticsFilters } from '@/context/useAnalyticsFiltersContext';

interface GaugeErrorWidgetProps {
	widget: WidgetResponse;
	topReasonsQuery?: string;
}

type ErrorReasonRow = {
	reason: string;
	count: number;
};

type ReasonBucketKey = 'tracker_offline' | 'delivery_timeout' | 'unsupported_firmware' | 'config_error';

const getReasonBucket = (reason: string): ReasonBucketKey => {
	const lowered = reason.toLowerCase();

	if (lowered.includes('tracker offline')) return 'tracker_offline';
	if (lowered.includes('timeout доставки') || lowered.includes('timeout delivery') || lowered.includes('delivery timeout')) return 'delivery_timeout';
	if (lowered.includes('неподдерживаемая прошивка') || lowered.includes('unsupported firmware')) return 'unsupported_firmware';

	return 'config_error';
};

const GaugeErrorWidget: React.FC<GaugeErrorWidgetProps> = ({ widget, topReasonsQuery }) => {
	const { t, i18n } = useTranslation();

	const { dateFilter } = useAnalyticsFilters();

	const dateRangeText = useMemo(() => {
		if (!dateFilter?.from || !dateFilter?.to) return null;

		const from = new Date(dateFilter.from);
		const to = new Date(dateFilter.to);
		if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) return null;

		const fmt = new Intl.DateTimeFormat(i18n.language, { day: '2-digit', month: '2-digit' });
		return `${fmt.format(from)} - ${fmt.format(to)}`;
	}, [dateFilter, i18n.language]);

	const gaugeRequest = useAnalyticRequest(widget.type, widget.query);
	const topReasonsRequest = useAnalyticRequest('gauge_top_reasons', topReasonsQuery);

	const fetchAnalyticQuery = useCallback(() => {
		if (!gaugeRequest) return Promise.resolve(null);
		return basicAnalyticQuery(gaugeRequest);
	}, [gaugeRequest]);

	const fetchTopReasonsQuery = useCallback(() => {
		if (!topReasonsRequest) return Promise.resolve(null);
		return basicAnalyticQuery(topReasonsRequest);
	}, [topReasonsRequest]);

	const { data: respAnalyticQuery, loading: analyticLoading } = useHandleServer(['respAnalyticQuery', widget.id, gaugeRequest], fetchAnalyticQuery, { enabled: !!gaugeRequest });

	const { data: respTopReasonsQuery, loading: topReasonsLoading } = useHandleServer(['respTopReasonsQuery', topReasonsRequest], fetchTopReasonsQuery, { enabled: !!topReasonsRequest });

	const { success, failed } = useMemo(() => {
		if (!respAnalyticQuery?.rows?.length || !respAnalyticQuery.columns?.length) {
			return { success: 0, failed: 0 };
		}

		const statusIdx = respAnalyticQuery.columns.indexOf('status');
		const valueIdx = respAnalyticQuery.columns.indexOf('value');

		if (statusIdx < 0 || valueIdx < 0) {
			return { success: 0, failed: 0 };
		}

		let success = 0;
		let failed = 0;

		respAnalyticQuery.rows.forEach((row) => {
			const status = String(row[statusIdx]);
			const value = Number(row[valueIdx] ?? 0);

			if (status === 'success') success = value;
			if (status === 'failed') failed = value;
		});

		return { success, failed };
	}, [respAnalyticQuery]);

	const total = success + failed;
	const successRate = total > 0 ? Math.round((success / total) * 100) : 0;

	const topReasonRows = useMemo<ErrorReasonRow[]>(() => {
		const columns = respTopReasonsQuery?.columns ?? [];
		const rows = respTopReasonsQuery?.rows ?? [];
		const reasonIdx = columns.indexOf('reason');
		const countIdx = columns.indexOf('reason_count');

		if (reasonIdx < 0 || countIdx < 0) return [];

		const grouped: Record<ReasonBucketKey, number> = {
			tracker_offline: 0,
			delivery_timeout: 0,
			unsupported_firmware: 0,
			config_error: 0,
		};

		rows.forEach((row) => {
			const reason = String(row[reasonIdx] ?? '');
			const count = Number(row[countIdx] ?? 0);
			if (!reason || !Number.isFinite(count)) return;

			const bucket = getReasonBucket(reason);
			grouped[bucket] += count;
		});

		return (Object.keys(grouped) as ReasonBucketKey[])
			.filter((bucket) => grouped[bucket] > 0)
			.map((bucket) => ({ reason: t(`label.analytics-reason-${bucket}`), count: grouped[bucket] }))
			.sort((a, b) => b.count - a.count);
	}, [respTopReasonsQuery, t]);

	const maxReasonCount = Math.max(1, ...topReasonRows.map((row) => row.count));

	return (
		<Widget title={widget.title} loading={analyticLoading || topReasonsLoading} editable={false}>
			<div className="flex flex-col xl:flex-row gap-6 xl:gap-8 min-w-0">
				<div className="flex justify-center">
					<GUIDonut
						size={192}
						stroke={16}
						trackColor="#e5e7eb"
						segments={[
							{ key: 'failed', value: failed, color: '#e2483d' },
							{ key: 'success', value: success, color: '#1e772e' },
						]}
						center={
							<div className="flex flex-col space-y-1 items-center">
								<div className="text-2xl font-bold text-black">{successRate}%</div>
								<div className="text-sm text-[#49525f]">{t('label.analytics-success-rate')}</div>
								<div className="text-xs text-gray-400">
									{success.toLocaleString(i18n.language)} / {total.toLocaleString(i18n.language)}
								</div>
								{dateRangeText && <span className="text-[13px] text-gray-400">{dateRangeText}</span>}
							</div>
						}
					/>
				</div>

				<div className="flex-1 min-w-0">
					<h4 className="text-base font-medium text-[#49525f] mb-3">{t('label.analytics-top-error-reasons')}</h4>
					<div className="space-y-2.5">
						{topReasonRows.length === 0 && <div className="text-sm text-[#49525f]">{t('label.no-data')}</div>}
						{topReasonRows.map((row) => {
							const barWidth = Math.max(6, Math.round((row.count / maxReasonCount) * 100));

							return (
								<div key={row.reason} className="grid grid-cols-[1fr_auto] gap-3 items-center">
									<div className="min-w-0">
										<div className="text-sm text-[#49525f] truncate">{row.reason}</div>
										<div className="mt-1 h-1.5 bg-[#EEF2F7] rounded-full overflow-hidden">
											<div className="h-full bg-[#e2483d] rounded-full" style={{ width: `${barWidth}%` }} />
										</div>
									</div>
									<div className="text-sm text-[#49525f] font-medium min-w-[56px] text-right">{row.count}</div>
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</Widget>
	);
};

export default GaugeErrorWidget;
