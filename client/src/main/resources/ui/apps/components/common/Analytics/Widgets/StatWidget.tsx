import Widget from './Widget';
import { useMemo } from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useCallback, useState } from 'react';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { useAnalyticsFilters } from '@/context/useAnalyticsFiltersContext';
import { useTranslation } from 'react-i18next';

import Poll from '@/components/@icons/poll';

interface StatWidgetProps {
	widget: WidgetResponse;
}

const StatWidget: React.FC<StatWidgetProps> = ({ widget }) => {
	const { i18n } = useTranslation();
	const { dateFilter } = useAnalyticsFilters();

	const [title, setTitle] = useState(widget.title);
	const [query, setQuery] = useState(widget.query ?? '');

	const dateRangeText = useMemo(() => {
		if (!dateFilter?.from || !dateFilter?.to) return null;

		const from = new Date(dateFilter.from);
		const to = new Date(dateFilter.to);
		if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) return null;

		const fmt = new Intl.DateTimeFormat(i18n.language, { day: '2-digit', month: '2-digit' });
		return `${fmt.format(from)} - ${fmt.format(to)}`;
	}, [dateFilter, i18n.language]);

	const request = useAnalyticRequest(widget.type, query);

	const fetchAnalyticQuery = useCallback(() => {
		if (!request) return Promise.resolve(null);
		return basicAnalyticQuery(request);
	}, [request]);

	const { data: respAnalyticQuery, loading: analyticLoading } = useHandleServer(['respAnalyticQuery', 'stat', widget.id, query, request], fetchAnalyticQuery, { enabled: !!request });

	const firstRow = respAnalyticQuery?.rows?.[0] ?? [];
	const columns = respAnalyticQuery?.columns ?? [];
	const valueIdx = columns.indexOf('value');

	const valueRaw = valueIdx >= 0 ? firstRow[valueIdx] : 0;
	const value = Number(valueRaw ?? 0);
	const valueText = Number.isFinite(value) ? value.toLocaleString(i18n.language) : '0';

	return (
		<Widget loading={analyticLoading} widgetId={widget.id} editableTitle={title} editableQuery={query} onApplyQuery={setQuery} onApplyTitle={setTitle}>
			<div className="flex flex-col sm:flex-row gap-3 min-w-0">
				<div className="w-10 h-10 rounded-full bg-[#E8EDF5] flex items-center justify-center shrink-0">
					<Poll fill="#395d95" size={21} />
				</div>

				<div className="flex flex-col gap-2.5 min-w-0 flex-1">
					<Tooltip title={title}>
						<div className="flex items-start justify-between gap-3 min-w-0">
							<div className="flex flex-col min-w-0">
								<span className="text-sm text-[#49525f] break-words">{title}</span>
							</div>
						</div>
					</Tooltip>

					<span className="text-2xl sm:text-3xl font-semibold text-gray-900 break-words">{valueText}</span>

					{widget.type !== 'stat_live' && !widget.subtitle && dateRangeText && <span className="text-[13px] text-gray-400">{dateRangeText}</span>}
				</div>
			</div>
		</Widget>
	);
};

export default StatWidget;
