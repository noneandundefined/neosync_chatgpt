import { useTranslation } from 'react-i18next';
import { useTableContext } from './TableContext';
import GUISelect from '@/components/ui/Select/GUISelect';

type PageItem = number | '-';

const getPagination = (page: number, total: number): PageItem[] => {
	if (total <= 0) {
		return [1];
	}

	if (total <= 7) {
		return Array.from({ length: total }, (_, i) => i + 1);
	}

	if (page <= 4) {
		return [1, 2, 3, 4, 5, '-', total];
	}

	if (page >= total - 3) {
		return [1, '-', total - 4, total - 3, total - 2, total - 1, total];
	}

	return [1, '-', page - 1, page, page + 1, '-', total];
};

function PaginationTable<T>() {
	const { t } = useTranslation();

	const { page, limit, totalPages, setLimit, setPage } = useTableContext<T>();

	const lastPage = Math.max(1, totalPages);
	const canPrev = page > 1;
	const canNext = page < lastPage;

	const items = getPagination(page, lastPage);

	return (
		<div className="flex items-center w-full gap-9 justify-center md:justify-end">
			<div className="flex justify-center md:justify-end gap-6 items-center w-full text-[14px]">
				<button
					type="button"
					disabled={!canPrev}
					onClick={() => canPrev && setPage(page - 1)}
					className={`justify-self-start whitespace-nowrap transition ${canPrev ? 'text-[#49525f] hover:underline' : 'text-[#ccc] pointer-events-none'}`}
				>
					{t('label.prev')}
				</button>

				<div className="flex items-center justify-center text-[14px] min-w-0">
					{items.map((item, index) => (
						<div key={`${item}-${index}`} className="flex items-center px-1">
							{item === '-' ? (
								<span className="text-[#49525f]">....</span>
							) : item === page ? (
								<span className="font-medium text-[#000]">{item}</span>
							) : (
								<button type="button" onClick={() => setPage(item)} className="font-medium text-[#49525f] hover:underline">
									{item}
								</button>
							)}
						</div>
					))}
				</div>

				<button
					type="button"
					disabled={!canNext}
					onClick={() => canNext && setPage(page + 1)}
					className={`justify-self-end whitespace-nowrap transition ${canNext ? 'text-[#49525f] hover:underline' : 'text-[#ccc] pointer-events-none'}`}
				>
					{t('label.next')}
				</button>
			</div>

			<div className="hidden md:flex items-center gap-[1rem]">
				<p className="whitespace-nowrap text-[13px]">{t('message.table-elements-page')}:</p>

				<GUISelect className="relative min-w-[3rem]" value={limit} onChange={(e) => setLimit(Number(e.target.value))}>
					<option value={5}>5</option>
					<option value={10}>10</option>
					<option value={15}>15</option>
					<option value={25}>25</option>
				</GUISelect>
			</div>
		</div>
	);
}

export default PaginationTable;
