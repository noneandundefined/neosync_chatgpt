import Widget from './Widget';
import { useCallback } from 'react';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useTranslation } from 'react-i18next';
import { formatAnalyticsValue, getAnalyticsColumnLabel } from '@/utils/AnalyticsFormatUtils';

const ReportTableWidget: React.FC<{ widget: WidgetResponse }> = ({ widget }) => {
	const { t, i18n } = useTranslation();
	const request = useAnalyticRequest(widget.type, widget.query);
	const fetchReport = useCallback(() => (request ? basicAnalyticQuery(request) : Promise.resolve(null)), [request]);
	const { data, loading } = useHandleServer(['analyticsReportTable', widget.id, request], fetchReport, { enabled: !!request });

	return (
		<Widget title={widget.title} loading={loading} editable={false}>
			<div className="min-h-[180px] max-w-full overflow-x-auto">
				<table className="w-full min-w-[520px] text-xs sm:text-sm">
					<thead>
						<tr className="border-b border-[#e8ebf0] text-[#7b8798]">
							{data?.columns?.map((column) => (
								<th key={column} className="text-left font-normal px-3 py-3 whitespace-nowrap">
									{getAnalyticsColumnLabel(t, column)}
								</th>
							))}
						</tr>
					</thead>
					<tbody>
						{data?.rows?.map((row, rowIndex) => (
							<tr key={rowIndex} className="border-b border-[#f0f2f5] last:border-0 hover:bg-[#f7f9fc]">
								{row.map((value, cellIndex) => (
									<td key={cellIndex} className="px-3 py-3 text-[#49525f] max-w-[260px] whitespace-normal break-words">
										{formatAnalyticsValue(t, i18n.language, data?.columns?.[cellIndex] || '', value)}
									</td>
								))}
							</tr>
						))}
					</tbody>
				</table>
				{!data?.rows?.length && <div className="h-40 flex items-center justify-center px-4 text-center text-sm text-[#7b838d]">{t('message.analytics-no-data')}</div>}
			</div>
		</Widget>
	);
};

export default ReportTableWidget;
