import { useMemo, useState } from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { Loading } from '../../Loader/Loading';
import { useTableContext } from './TableContext';
import Dropdown from '@/components/ui/Dropdown/Dropdown';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';

import MenuUp from '@/components/@icons/menu-up';
import MenuDown from '@/components/@icons/menu-down';
import UnfoldMoreHorizontal from '@/components/@icons/unfold-more-horizontal';

function TableHead<T>() {
	const { t } = useTranslation();
	const { allChecked, toggleAll, selectAllAcrossPages, columns, rowActions, loading, columnSort, setColumnSort } = useTableContext<T>();

	const [menuOpen, setMenuOpen] = useState<boolean>(false);
	const [isSelectingAll, setIsSelectingAll] = useState<boolean>(false);

	const visibleColumns = useMemo(() => columns.filter((c) => c.visible !== false), [columns]);

	const handleSort = (key: string) => {
		setColumnSort((prev) => {
			if (!prev || prev.key !== key) {
				return { key, direction: 'asc' };
			}
			return {
				key,
				direction: prev.direction === 'asc' ? 'desc' : 'asc',
			};
		});
	};

	return (
		<thead className="relative sticky top-0 z-10 bg-white">
			{isSelectingAll && (
				<>
					<div className="fixed inset-0 z-[1200] cursor-wait" />
					<div className="fixed top-[6rem] left-1/2 -translate-x-1/2 z-[1201] bg-[#fff4bf] border border-[#e5d27a] rounded px-3 py-1 text-[12px] shadow">
						<Loading />
					</div>
				</>
			)}

			<tr>
				<th id="column-left">
					<div className="flex items-center justify-center">
						<div className="relative">
							<GUICheckbox checked={allChecked} onChange={toggleAll} />

							<Dropdown open={menuOpen} close={() => setMenuOpen(false)} stopPropagation={true} className="absolute top-7 left-0 bg-white py-1 border border-[#e4e4e4] min-w-[11rem] z-[100]">
								<p
									className={`text-sm px-3 py-[9px] md:py-[5px] ${isSelectingAll ? 'cursor-wait opacity-60' : 'hover:bg-[#f7f7f7] cursor-pointer'}`}
									onClick={async () => {
										if (isSelectingAll) return;

										try {
											setIsSelectingAll(true);

											await selectAllAcrossPages();
										} finally {
											setIsSelectingAll(false);
											setMenuOpen(false);
										}
									}}
								>
									{t('message.select-all')}
								</p>
							</Dropdown>
						</div>

						<div
							className="cursor-pointer -mt-1 -ml-1"
							onClick={(e) => {
								e.stopPropagation();
								setMenuOpen((prev) => !prev);
							}}
						>
							{menuOpen ? <MenuUp fill="#49525f" size={21} /> : <MenuDown fill="#49525f" size={21} />}
						</div>
					</div>
				</th>

				{visibleColumns.map((col) =>
					columnSort?.key === col.key ? (
						<Tooltip title={t('message.current-column-is-sorted')} key={col.key as string}>
							<th key={col.key as string} onClick={() => col.sortable && handleSort(col.key as string)} className={col.sortable ? 'cursor-pointer select-none' : ''} style={{ width: col.width, maxWidth: '10rem' }}>
								<div className="flex items-center justify-center h-full">
									{col.title}

									{col.sortable && (
										<div className="mt-[1px]">
											<UnfoldMoreHorizontal size={15} fill="#49525f" />
										</div>
									)}
								</div>
							</th>
						</Tooltip>
					) : (
						<th key={col.key as string} onClick={() => col.sortable && handleSort(col.key as string)} className={col.sortable ? 'cursor-pointer select-none' : ''} style={{ width: col.width, maxWidth: '10rem' }}>
							<div className="flex items-center justify-center">
								{col.title}

								{col.sortable && (
									<div className="mt-[1px]">
										<UnfoldMoreHorizontal size={15} fill="#49525f" />
									</div>
								)}
							</div>
						</th>
					)
				)}
				{rowActions && <th>{t('label.actions')}</th>}
			</tr>

			{loading && (
				<div className="absolute top-[3.7rem] w-full">
					<div id="loader-line" className="!w-full" />
				</div>
			)}
		</thead>
	);
}

export default TableHead;
