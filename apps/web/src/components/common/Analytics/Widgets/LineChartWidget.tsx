import Widget from './Widget';
import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { WidgetResponse, basicAnalyticQuery } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';

interface LineChartWidgetProps {
	widget: WidgetResponse;
}

type ChartPoint = {
	date: string;
	created: number;
	sent: number;
};

const chartWidth = 760;
const chartHeight = 280;
const padding = { top: 20, right: 16, bottom: 34, left: 44 };

const formatDayLabel = (value: string) => {
	const date = new Date(value);
	if (!Number.isNaN(date.getTime())) {
		const dd = String(date.getDate()).padStart(2, '0');
		const mm = String(date.getMonth() + 1).padStart(2, '0');
		return `${dd}.${mm}`;
	}

	if (value.length >= 10) {
		return `${value.slice(8, 10)}.${value.slice(5, 7)}`;
	}

	return value;
};

const LineChartWidget: React.FC<LineChartWidgetProps> = ({ widget }) => {
	const { t } = useTranslation();

	const request = useAnalyticRequest(widget.type, widget.query);

	const fetchAnalyticQuery = useCallback(() => {
		if (!request) return Promise.resolve(null);
		return basicAnalyticQuery(request);
	}, [request]);

	const { data: respAnalyticQuery, loading: analyticLoading } = useHandleServer(['respAnalyticQuery', widget.type, widget.id, request], fetchAnalyticQuery, { enabled: !!request });

	const points = useMemo<ChartPoint[]>(() => {
		if (!respAnalyticQuery?.rows?.length || !respAnalyticQuery.columns?.length) return [];

		const dateIdx = respAnalyticQuery.columns.indexOf('date');
		const createdIdx = respAnalyticQuery.columns.indexOf('created');

		const sentIdx = respAnalyticQuery.columns.indexOf('sent');
		if (dateIdx < 0 || createdIdx < 0 || sentIdx < 0) return [];

		return respAnalyticQuery.rows.map((row: any) => ({
			date: String(row[dateIdx] ?? ''),
			created: Number(row[createdIdx] ?? 0),
			sent: Number(row[sentIdx] ?? 0),
		}));
	}, [respAnalyticQuery]);

	const maxY = useMemo(() => {
		const max = Math.max(0, ...points.flatMap((p) => [p.created, p.sent]));
		return max > 0 ? Math.ceil(max / 10) * 10 : 10;
	}, [points]);

	const toX = (index: number) => {
		const innerWidth = chartWidth - padding.left - padding.right;
		if (points.length <= 1) return padding.left;
		return padding.left + (index / (points.length - 1)) * innerWidth;
	};

	const toY = (value: number) => {
		const innerHeight = chartHeight - padding.top - padding.bottom;
		return padding.top + innerHeight - (value / maxY) * innerHeight;
	};

	const buildPath = (key: 'created' | 'sent') => points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${toX(i)} ${toY(p[key])}`).join(' ');

	const yTicks = 5;
	const yLabels = Array.from({ length: yTicks + 1 }, (_, i) => Math.round((maxY / yTicks) * (yTicks - i)));

	return (
		<Widget title={widget.title} loading={analyticLoading} editable={false}>
			<div className="flex items-center gap-6 mb-3 text-sm text-[#49525f]">
				<div className="flex items-center gap-2">
					<span className="w-2.5 h-2.5 rounded-full bg-[#395d95]" />
					<span>{t('label.analytics-created-configs')}</span>
				</div>
				<div className="flex items-center gap-2">
					<span className="w-2.5 h-2.5 rounded-full bg-[#24b778]" />
					<span>{t('label.analytics-sent-configs')}</span>
				</div>
			</div>

			<div className="w-full overflow-x-auto">
				<svg viewBox={`0 0 ${chartWidth} ${chartHeight}`} className="w-full h-[300px]">
					{yLabels.map((label) => {
						const y = toY(label);
						return (
							<g key={label}>
								<line x1={padding.left} y1={y} x2={chartWidth - padding.right} y2={y} stroke="#e5e7eb" strokeWidth="1" />
								<text x={padding.left - 10} y={y + 4} textAnchor="end" className="fill-gray-400 text-[11px]">
									{label}
								</text>
							</g>
						);
					})}

					{points.length > 1 && (
						<>
							<path d={buildPath('created')} fill="none" stroke="#395d95" strokeWidth="3" strokeLinecap="round" />
							<path d={buildPath('sent')} fill="none" stroke="#24b778" strokeWidth="3" strokeLinecap="round" />
						</>
					)}

					{points.map((point, idx) => (
						<g key={`${point.date}-${idx}`}>
							<circle cx={toX(idx)} cy={toY(point.created)} r="3.5" fill="#395d95" />
							<circle cx={toX(idx)} cy={toY(point.sent)} r="3.5" fill="#24b778" />
						</g>
					))}

					{points.map((point, idx) => (
						<text key={`x-${point.date}-${idx}`} x={toX(idx)} y={chartHeight - 10} textAnchor="middle" className="fill-gray-400 text-[11px]">
							{formatDayLabel(point.date)}
						</text>
					))}
				</svg>
			</div>
		</Widget>
	);
};

export default LineChartWidget;
