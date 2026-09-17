// import BarWidget from './Widgets/BarWidget';
import BarWidget from './Widgets/BarWidget';
import SQLWidget from './Widgets/SQLWidget';
import StatWidget from './Widgets/StatWidget';
import TableWidget from './Widgets/TableWidget';
import FunnelWidget from './Widgets/FunnelWidget';
import { useTranslation } from 'react-i18next';
import { getAnalyticsWidgetTitle } from '@/utils/AnalyticsFormatUtils';
import { WidgetResponse } from '@/rest/analyticAPI';
import LineChartWidget from './Widgets/LineChartWidget';
import GaugeErrorWidget from './Widgets/GaugeErrorWidget';
import ReportTableWidget from './Widgets/ReportTableWidget';
import ReportLineWidget from './Widgets/ReportLineWidget';

interface WidgetRendererProps {
	widget: WidgetResponse;
	allWidgets?: WidgetResponse[];
}

const getWidgetLayoutClass = (widget: WidgetResponse) => {
	switch (widget.type) {
		case 'stat':
		case 'stat_live':
			return 'col-span-12 sm:col-span-6 xl:col-span-3 min-w-0';
		case 'line_chart':
			return 'col-span-12 xl:col-span-5 min-w-0';
		case 'bar':
			return 'col-span-12 xl:col-span-7 min-w-0';
		case 'gauge_error':
			return 'col-span-12 lg:col-span-5 xl:col-span-6 min-w-0';
		case 'table':
			return 'col-span-12 lg:col-span-7 xl:col-span-6 min-w-0';
		case 'sql':
			return 'col-span-12 min-w-0';
		case 'report_line':
			return 'col-span-12 min-w-0';
		case 'report_table':
		case 'funnel':
			return 'col-span-12 xl:col-span-6 min-w-0';
		default:
			return 'col-span-12 min-w-0';
	}
};

const WidgetRenderer: React.FC<WidgetRendererProps> = ({ widget, allWidgets }) => {
	const { t } = useTranslation();
	const localizedWidget = { ...widget, title: getAnalyticsWidgetTitle(t, widget) };

	if (widget.type === 'gauge_top_reasons') {
		return null;
	}

	const layoutClass = getWidgetLayoutClass(widget);
	const topReasonsQuery = allWidgets?.find((w) => w.id === 11)?.query;

	switch (widget.type) {
		case 'stat':
		case 'stat_live':
			return (
				<div className={layoutClass}>
					<StatWidget widget={localizedWidget} />
				</div>
			);

		case 'line_chart':
			return (
				<div className={layoutClass}>
					<LineChartWidget widget={localizedWidget} />
				</div>
			);

		case 'bar':
			return (
				<div className={layoutClass}>
					<BarWidget widget={localizedWidget} />
				</div>
			);

		case 'table':
			return (
				<div className={layoutClass}>
					<TableWidget widget={localizedWidget} />
				</div>
			);

		case 'gauge_error':
			return (
				<div className={layoutClass}>
					<GaugeErrorWidget widget={localizedWidget} topReasonsQuery={topReasonsQuery} />
				</div>
			);

		case 'sql':
			return (
				<div className={layoutClass}>
					<SQLWidget />
				</div>
			);
		case 'report_table':
			return (
				<div className={layoutClass}>
					<ReportTableWidget widget={localizedWidget} />
				</div>
			);
		case 'report_line':
			return (
				<div className={layoutClass}>
					<ReportLineWidget widget={localizedWidget} />
				</div>
			);
		case 'funnel':
			return (
				<div className={layoutClass}>
					<FunnelWidget widget={localizedWidget} />
				</div>
			);

		default:
			return null;
	}
};

export default WidgetRenderer;
