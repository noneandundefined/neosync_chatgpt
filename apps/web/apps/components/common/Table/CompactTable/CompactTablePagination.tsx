import { useTranslation } from 'react-i18next';
import { useCompactTableContext } from './CompactTableContext';

const CompactTablePagination = () => {
	const { t } = useTranslation();

	const { page, limit, totalPages, totalAll, setPage, footer } = useCompactTableContext<any>();

	if (totalAll === 0) return null;

	const from = (page - 1) * limit + 1;
	const to = Math.min(page * limit, totalAll);
	const canPrev = page > 1;
	const canNext = page < totalPages;

	return (
		<div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 px-4 py-3 border-t border-[#e9e9e9] min-h-[3.7rem] bg-[#f6f6f6]">
			{footer && <div className="text-sm text-[#656d76]">{footer}</div>}

			<div className={`flex items-center gap-3 ${footer ? '' : 'ml-auto'}`}>
				<p className="text-sm text-[#656d76] whitespace-nowrap">{t('message.compact-table-page-range', { from, to, total: totalAll })}</p>

				<div className="flex items-center gap-2">
					<button
						type="button"
						disabled={!canPrev}
						onClick={() => canPrev && setPage(page - 1)}
						className="px-3 py-1.5 text-sm border border-[#d0d7de] rounded-[6px] bg-white hover:bg-[#f6f8fa] text-[#49525f] disabled:opacity-50 disabled:pointer-events-none"
					>
						{t('label.prev')}
					</button>

					<button
						type="button"
						disabled={!canNext}
						onClick={() => canNext && setPage(page + 1)}
						className="px-3 py-1.5 text-sm border border-[#d0d7de] rounded-[6px] bg-white hover:bg-[#f6f8fa] text-[#49525f] disabled:opacity-50 disabled:pointer-events-none"
					>
						{t('label.next')}
					</button>
				</div>
			</div>
		</div>
	);
};

export default CompactTablePagination;
