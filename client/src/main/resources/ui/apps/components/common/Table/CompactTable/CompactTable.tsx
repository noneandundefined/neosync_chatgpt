import CompactTableBody from './CompactTableBody';
import useQueryState from '@/hooks/useQueryState';
import Fallback from '@/components/Fallback/Fallback';
import CompactTableSearch from './CompactTableSearch';
import CompactTableContext from './CompactTableContext';
import CompactTableToolbar from './CompactTableToolbar';
import { keepPreviousData } from '@tanstack/react-query';
import { useCallback, useEffect, useState } from 'react';
import { useTableLimit } from '@/hooks/Table/useTableLimit';
import CompactTablePagination from './CompactTablePagination';
import { ColumnSortState } from '../GenericTable/GenericTable';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { CompactTableProps, compactTableQueryKey } from './types';
import { useTableSelection } from '@/hooks/Table/useTableSelection';
import { CACHEKEYs_TABLE_COLUMN_SORT } from '@/constants/CacheKeys.constants';

const defaultLimit = 25;

function CompactTable<T>({
	tableKey,
	fetchFn,
	sysToken,
	getRowId,
	renderTitle,
	renderSubtitle,
	renderLabel,
	renderMeta,
	renderTime,
	isUnread,
	onRowClick,
	onTitleClick,
	onDeleteFn,
	bulkActions,
	emptyComponent,
	footer,
	enableSearch = true,
	selectionMode = 'multiple',
	controlledSelectedIds,
	onSelectionChange,
	persistSearchInUrl = true,
}: CompactTableProps<T>) {
	const { limit, setLimit } = useTableLimit(`compact-${tableKey}`, defaultLimit);

	const [hoveredRow, setHoveredRow] = useState<number | null>(null);
	const [pageParam, setPageParam] = useQueryState(`compact-page-${tableKey}`, '1', persistSearchInUrl);
	const [urlSearch, setUrlSearch] = useQueryState(`compact-search-${tableKey}`, '', persistSearchInUrl);

	const [localSearchApplied, setLocalSearchApplied] = useState('');
	const [localSearch, setLocalSearch] = useState(() => (persistSearchInUrl ? urlSearch : ''));

	const [localPage, setLocalPage] = useState(1);

	const page = persistSearchInUrl ? Math.max(1, parseInt(pageParam, 10) || 1) : localPage;
	const setPage = useCallback(
		(value: React.SetStateAction<number>) => {
			const next = typeof value === 'function' ? value(page) : value;
			if (persistSearchInUrl) {
				setPageParam(String(next));
				return;
			}
			setLocalPage(next);
		},
		[page, persistSearchInUrl, setPageParam]
	);

	const search = persistSearchInUrl ? urlSearch : localSearchApplied;
	const setSearch = persistSearchInUrl ? setUrlSearch : setLocalSearchApplied;

	const [columnSort] = useState<ColumnSortState>(() => {
		const fallback: ColumnSortState = { key: 'created_at', direction: 'desc' };
		const cacheKey = CACHEKEYs_TABLE_COLUMN_SORT(tableKey);
		const saved = localStorage.getItem(cacheKey);

		if (!saved) {
			return fallback;
		}

		try {
			const parsed = JSON.parse(saved) as Partial<ColumnSortState>;
			if (typeof parsed.key === 'string' && (parsed.direction === 'asc' || parsed.direction === 'desc')) {
				return parsed as ColumnSortState;
			}
		} catch {
			localStorage.removeItem(cacheKey);
		}

		return fallback;
	});

	useEffect(() => {
		setPage(1);
	}, [limit]);

	const fetcher = useCallback((signal?: AbortSignal) => fetchFn({ page, limit, search, columnSort }, signal), [fetchFn, page, limit, search, columnSort]);
	const { data, loading } = useHandleServer([...compactTableQueryKey(tableKey), page, limit, search, columnSort, sysToken], fetcher, { placeholderData: keepPreviousData });

	const items = data?.items ?? [];
	const totalAll = data?.total ?? 0;
	const totalPages = data?.totalPages ?? 1;

	useEffect(() => {
		if (!data) return;
		if (page > data.totalPages) {
			setPage(Math.max(1, data.totalPages));
		}
	}, [data, page]);

	const { checkedItems, allChecked, selectedIds, selectedCount, hasSelected, toggleAll, toggleOne, replaceSelection } = useTableSelection(items, getRowId);

	useEffect(() => {
		if (controlledSelectedIds === undefined) return;

		replaceSelection(controlledSelectedIds);
	}, [controlledSelectedIds, replaceSelection]);

	const handleToggleOne = useCallback(
		(id: string | number) => {
			if (selectionMode === 'single') {
				const nextIds = checkedItems[id] ? [] : [id];
				replaceSelection(nextIds);
				onSelectionChange?.(nextIds);
				return;
			}

			toggleOne(id);
		},
		[selectionMode, checkedItems, replaceSelection, onSelectionChange, toggleOne]
	);

	useEffect(() => {
		const timeout = setTimeout(() => {
			if (localSearch === search) {
				return;
			}

			setSearch(localSearch);
			setPage(1);
		}, 300);

		return () => clearTimeout(timeout);
	}, [localSearch, search, setSearch]);

	const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setLocalSearch(e.target.value);
	};

	const ctxValue = {
		search,
		setSearch,
		limit,
		setLimit,
		page,
		setPage,
		totalPages,
		totalAll,
		pageData: items,
		loading,
		columnSort,
		checkedItems,
		allChecked,
		toggleAll,
		toggleOne: handleToggleOne,
		selectedIds,
		selectedCount,
		hasSelected,
		hoveredRow,
		setHoveredRow,
		getRowId,
		onRowClick,
		onTitleClick,
		onDeleteFn,
		bulkActions,
		renderTitle,
		renderSubtitle,
		renderLabel,
		renderMeta,
		renderTime,
		isUnread,
		footer,
		selectionMode,
	};

	if (loading) {
		return <Fallback />;
	}

	if (items.length === 0 && !search && emptyComponent) {
		return <>{emptyComponent}</>;
	}

	return (
		<CompactTableContext.Provider value={ctxValue}>
			{enableSearch && <CompactTableSearch value={localSearch} onChange={handleSearchChange} />}

			<div className="border border-[#e9e9e9] rounded-[6px] bg-white">
				<CompactTableToolbar />
				<div className="overflow-hidden rounded-b-[6px]">
					<CompactTableBody />
					<CompactTablePagination />
				</div>
			</div>
		</CompactTableContext.Provider>
	);
}

export default CompactTable;
