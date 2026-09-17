import Widget from './Widget';
import { useCallback } from 'react';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useTranslation } from 'react-i18next';
import { formatAnalyticsValue } from '@/utils/AnalyticsFormatUtils';

const FunnelWidget: React.FC<{ widget: WidgetResponse }> = ({ widget }) => {
	const request = useAnalyticRequest(widget.type, widget.query);
	const { t, i18n } = useTranslation();
	const fetchReport = useCallback(() => (request ? basicAnalyticQuery(request) : Promise.resolve(null)), [request]);
	const { data, loading } = useHandleServer(['analyticsFunnel', widget.id, request], fetchReport, { enabled: !!request });
	const first = Number(data?.rows?.[0]?.[1] || 0);

	return (
		<Widget title={widget.title} loading={loading} editable={false}>
			<div className="space-y-4 p-1 sm:p-3 min-h-[220px]">
				{data?.rows?.map((row, index) => {
					const value = Number(row[1] || 0);
					const percent = first ? Math.round((value / first) * 100) : 0;
					return (
						<div key={index} className="grid grid-cols-[minmax(0,1fr)_52px] sm:grid-cols-[minmax(150px,1fr)_3fr_70px] gap-2 sm:gap-3 items-center">
							<span className="text-xs sm:text-sm text-[#49525f] break-words">{formatAnalyticsValue(t, i18n.language, 'step', row[0])}</span>
							<div className="hidden sm:block h-7 bg-[#eef1f5] rounded overflow-hidden">
								<div className="h-full bg-[#395d95] rounded transition-all" style={{ width: `${Math.max(percent, value ? 4 : 0)}%` }} />
							</div>
							<span className="text-right">
								<b>{value}</b>
								<small className="block text-[#8a96a8]">{percent}%</small>
							</span>
						</div>
					);
				})}
				{!data?.rows?.length && <div className="h-40 flex items-center justify-center px-4 text-center text-sm text-[#7b838d]">{t('message.analytics-no-data')}</div>}
			</div>
		</Widget>
	);
};

export default FunnelWidget;
