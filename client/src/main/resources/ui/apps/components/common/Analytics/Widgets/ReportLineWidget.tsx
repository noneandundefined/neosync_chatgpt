import Widget from './Widget';
import { useCallback, useMemo } from 'react';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useTranslation } from 'react-i18next';
import { getAnalyticsColumnLabel } from '@/utils/AnalyticsFormatUtils';

const colors = ['#395d95', '#24b778', '#e59a2f', '#e2483d', '#6f7f95'];
const width = 900;
const height = 280;
const padding = { top: 20, right: 20, bottom: 35, left: 45 };

const ReportLineWidget: React.FC<{ widget: WidgetResponse }> = ({ widget }) => {
	const request = useAnalyticRequest(widget.type, widget.query);
	const { t } = useTranslation();
	const fetchReport = useCallback(() => (request ? basicAnalyticQuery(request) : Promise.resolve(null)), [request]);
	const { data, loading } = useHandleServer(['analyticsReportLine', widget.id, request], fetchReport, { enabled: !!request });

	const series = useMemo(() => (data?.columns || []).slice(1), [data?.columns]);
	const values = useMemo(() => data?.rows || [], [data?.rows]);
	const max = Math.max(1, ...values.flatMap((row) => row.slice(1).map((value) => Number(value) || 0)));
	const x = (index: number) => padding.left + (values.length <= 1 ? 0 : (index / (values.length - 1)) * (width - padding.left - padding.right));
	const y = (value: number) => height - padding.bottom - (value / max) * (height - padding.top - padding.bottom);

	return (
		<Widget title={widget.title} loading={loading} editable={false}>
			<div className="flex flex-wrap gap-3 px-1 pt-1 text-xs sm:text-sm text-[#49525f]">
				{series.map((name, index) => (
					<span key={name} className="flex items-center gap-2">
						<i className="w-2.5 h-2.5 rounded-full" style={{ background: colors[index % colors.length] }} />
						{getAnalyticsColumnLabel(t, name)}
					</span>
				))}
			</div>
			<div className="max-w-full overflow-x-auto">
				<svg viewBox={`0 0 ${width} ${height}`} className="w-full min-w-[560px] h-[240px] sm:h-[290px]">
					{[0, 0.25, 0.5, 0.75, 1].map((ratio) => {
						const yy = y(max * ratio);
						return (
							<g key={ratio}>
								<line x1={padding.left} y1={yy} x2={width - padding.right} y2={yy} stroke="#edf0f5" />
								<text x={padding.left - 8} y={yy + 4} textAnchor="end" className="fill-[#98a2b3] text-[11px]">
									{Math.round(max * ratio)}
								</text>
							</g>
						);
					})}
					{series.map((_, seriesIndex) => (
						<path
							key={seriesIndex}
							d={values.map((row, index) => `${index ? 'L' : 'M'} ${x(index)} ${y(Number(row[seriesIndex + 1]) || 0)}`).join(' ')}
							fill="none"
							stroke={colors[seriesIndex % colors.length]}
							strokeWidth="2.5"
							strokeLinecap="round"
						/>
					))}
					{values.map((row, index) => (
						<text key={index} x={x(index)} y={height - 10} textAnchor="middle" className="fill-[#98a2b3] text-[10px]">
							{String(row[0]).slice(5, 10)}
						</text>
					))}
				</svg>
			</div>
		</Widget>
	);
};

export default ReportLineWidget;
