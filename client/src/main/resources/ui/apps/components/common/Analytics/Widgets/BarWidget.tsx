import Widget from './Widget';
import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicAnalyticQuery, WidgetResponse } from '@/rest/analyticAPI';
import { useAnalyticRequest } from '@/hooks/Analytics/useAnalyticRequest';

interface BarWidgetProps {
	widget: WidgetResponse;
}

type BucketRow = {
	bucket: string;
	value: number;
	percent: number;
};

const DEFAULT_CATEGORY = 'chart';

const withPercents = (rows: BucketRow[]): BucketRow[] => {
	const total = rows.reduce((sum, row) => sum + row.value, 0);
	if (total <= 0) return rows;

	return rows.map((row) => ({
		...row,
		percent: Math.round((row.value * 100) / total),
	}));
};

const BarWidget: React.FC<BarWidgetProps> = ({ widget }) => {
	const { t } = useTranslation();

	const request = useAnalyticRequest(widget.type, widget.query);

	const fetchAnalyticQuery = useCallback(() => {
		if (!request) return Promise.resolve(null);
		return basicAnalyticQuery(request);
	}, [request]);

	const { data: respAnalyticQuery, loading: analyticLoading } = useHandleServer(['respAnalyticPopularQuery', widget.id, request], fetchAnalyticQuery, { enabled: !!request });

	const sections = useMemo(() => {
		const columns = respAnalyticQuery?.columns ?? [];
		const rows = respAnalyticQuery?.rows ?? [];

		const categoryIdx = columns.indexOf('category');
		const bucketIdx = columns.indexOf('bucket');
		const valueIdx = columns.indexOf('value');
		const percentIdx = columns.indexOf('percent');

		if (bucketIdx < 0 || valueIdx < 0) {
			return [] as { key: string; rows: BucketRow[] }[];
		}

		const grouped = new Map<string, BucketRow[]>();

		rows.forEach((row) => {
			const category = categoryIdx >= 0 ? String(row[categoryIdx] ?? '') : DEFAULT_CATEGORY;
			const bucket = String(row[bucketIdx] ?? '');
			const value = Number(row[valueIdx] ?? 0);
			const percentRaw = percentIdx >= 0 ? Number(row[percentIdx] ?? NaN) : NaN;

			if (!bucket || !Number.isFinite(value) || value <= 0) return;

			const key = category || DEFAULT_CATEGORY;
			if (!grouped.has(key)) grouped.set(key, []);

			grouped.get(key)!.push({
				bucket,
				value,
				percent: Number.isFinite(percentRaw) ? percentRaw : 0,
			});
		});

		return Array.from(grouped.entries()).map(([key, sectionRows]) => {
			const needsPercent = sectionRows.some((row) => row.percent <= 0);
			const rows = needsPercent ? withPercents(sectionRows) : sectionRows;

			return {
				key,
				rows: rows.sort((a, b) => b.value - a.value).slice(0, 8),
			};
		});
	}, [respAnalyticQuery]);

	const hasData = sections.some((section) => section.rows.length > 0);

	const sectionTitle = (key: string) => {
		if (key === DEFAULT_CATEGORY) return widget.title;
		const translated = t(`label.analytics-popular-category-${key}`, { defaultValue: '' });
		return translated || key;
	};

	return (
		<Widget title={widget.title} loading={analyticLoading} className="relative" editable={false}>
			{!hasData ? (
				<section className="text-sm text-[#49525f]">{t('label.no-data')}</section>
			) : (
				<div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
					{sections.map((section) => (
						<section key={section.key} className="min-w-0">
							{sections.length > 1 && <div className="text-sm font-medium text-[#49525f] my-3">{sectionTitle(section.key)}</div>}
							<div className="space-y-2.5">
								{section.rows.map((row) => {
									const width = Math.max(8, Math.min(100, row.percent));

									return (
										<div key={`${section.key}-${row.bucket}`} className="min-w-0">
											<div className="flex items-center justify-between gap-2 text-sm">
												<span className="text-[#49525f] truncate">{t(`label.analytics-popular-bucket-${row.bucket}`, { defaultValue: row.bucket })}</span>
												<span className="text-xs text-[#9AA4B2] shrink-0">{row.value}</span>
											</div>
											<div className="flex items-center gap-2">
												<div className="flex-1 my-2 h-2.5 bg-[#E8EDF5] overflow-hidden">
													<div className="h-full bg-[#395d95]" style={{ width: `${width}%` }} />
												</div>
												<span className="text-xs text-[#49525f] font-medium w-10 text-right">{row.percent}%</span>
											</div>
										</div>
									);
								})}
							</div>
						</section>
					))}
				</div>
			)}
		</Widget>
	);
};

export default BarWidget;
