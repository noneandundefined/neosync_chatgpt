import Widget from './Widget';
import { useCallback, useMemo } from 'react';
import { formatDate } from '@/utils/TimeUtils';
import { useTranslation } from 'react-i18next';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { getAnalyticsErrorReasonLabel } from '@/utils/AnalyticsFormatUtils';

interface TableWidgetProps {
	widget: WidgetResponse;
}

type RecentOperationRow = {
	operationTime: string;
	accountName: string;
	deviceImei: string;
	operationReason: string;
};

const TableWidget: React.FC<TableWidgetProps> = ({ widget }) => {
	const { t } = useTranslation();

	const request = useAnalyticRequest(widget.type, widget.query);

	const fetchAnalyticQuery = useCallback(() => {
		if (!request) return Promise.resolve(null);
		return basicAnalyticQuery(request);
	}, [request]);

	const { data: respAnalyticQuery, loading: analyticLoading } = useHandleServer(['respAnalyticQuery', widget.id, request], fetchAnalyticQuery, { enabled: !!request });

	const recentRows = useMemo<RecentOperationRow[]>(() => {
		const columns = respAnalyticQuery?.columns ?? [];
		const rows = respAnalyticQuery?.rows ?? [];
		const rowTypeIdx = columns.indexOf('row_type');
		const dateIdx = columns.indexOf('operation_time');
		const accountIdx = columns.indexOf('account_name');
		const deviceIdx = columns.indexOf('device_imei');
		const reasonIdx = columns.indexOf('operation_reason');

		if (rowTypeIdx < 0 || dateIdx < 0 || accountIdx < 0 || deviceIdx < 0 || reasonIdx < 0) return [];

		return rows
			.filter((row) => String(row[rowTypeIdx] ?? '') === 'recent')
			.map((row) => ({
				operationTime: formatDate(String(row[dateIdx] ?? '')),
				accountName: String(row[accountIdx] ?? '—'),
				deviceImei: String(row[deviceIdx] ?? '—'),
				operationReason: getAnalyticsErrorReasonLabel(t, String(row[reasonIdx] ?? '')),
			}));
	}, [respAnalyticQuery, t]);

	return (
		<Widget title={widget.title} loading={analyticLoading} editable={false}>
			<div className="max-w-full overflow-x-auto">
				<table className="w-full min-w-[560px] text-sm">
					<thead>
						<tr className="text-left text-[#49525f] border-b border-gray-100">
							<th className="py-2 pr-3 font-medium">{t('label.analytics-date-time')}</th>
							<th className="py-2 pr-3 font-medium">{t('label.analytics-account')}</th>
							<th className="py-2 pr-3 font-medium">{t('label.analytics-device')}</th>
							<th className="py-2 pr-3 font-medium">{t('label.analytics-error-reason')}</th>
						</tr>
					</thead>
					<tbody>
						{recentRows.length === 0 && (
							<tr>
								<td colSpan={4} className="py-4 text-gray-400">
									{t('label.no-data')}
								</td>
							</tr>
						)}
						{recentRows.map((row, index) => (
							<tr key={`${row.operationTime}-${row.deviceImei}-${index}`} className="border-b border-gray-50">
								<td className="py-2 pr-3 whitespace-nowrap">{row.operationTime}</td>
								<td className="py-2 pr-3 whitespace-nowrap">{row.accountName}</td>
								<td className="py-2 pr-3 whitespace-nowrap">{row.deviceImei}</td>
								<td className="py-2 pr-3">{row.operationReason}</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</Widget>
	);
};

export default TableWidget;
