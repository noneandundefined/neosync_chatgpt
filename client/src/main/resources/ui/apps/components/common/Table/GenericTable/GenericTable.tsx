import { keepPreviousData } from '@tanstack/react-query';
import React, { useCallback, useEffect, useState } from 'react';

import TableBody from './TableBody';
import TableHead from './TableHeader';
import TableToolbar from './TableToolbar';
import PaginationTable from './PaginationTable';

import { Column } from './types';
import TableContext from './TableContext';
import useQueryState from '@/hooks/useQueryState';
import Fallback from '@/components/Fallback/Fallback';
import { useTableLimit } from '@/hooks/Table/useTableLimit';
import { useTableColumns } from '@/hooks/Table/useTableColumns';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useTableSelection } from '@/hooks/Table/useTableSelection';
import { CACHEKEYs_TABLE_COLUMN_SORT } from '@/constants/CacheKeys.constants';

const defaultLimit: number = 10;

export interface ColumnSortState {
	key: string;
	direction: 'asc' | 'desc';
}

export interface FetchParams {
	page: number;
	limit: number;
	search: string;
	columnSort: ColumnSortState;
}

export interface FetchResult<T> {
	items: T[];
	total: number;
	totalPages: number;
}

export type FetchFn<T> = (params: FetchParams, signal?: AbortSignal) => Promise<FetchResult<T>>;

interface GenericTableProps<T> {
	/** Идентификатор таблиц */
	tableKey: string;
	/** Запрос к серверу с пагинацией */
	fetchFn: FetchFn<T>;
	/** Описание колонок таблицы (заголовки, ключи, кастомные рендеры) */
	columns: Column<T>[];
	/** Токен для обновления данных в таблице */
	sysToken: number;
	/** Уникальное значение каждой строки в таблицы */
	getRowId: (item: T) => string | number;
	/**
	 * Функция для рендера блока действий в строке
	 * @param item - объект данных текущей строки
	 * @param index - индекс строки (с учётом пагинации)
	 * @param isHovered - true, если мышь наведена на строку
	 */
	rowActions?: (item: T, index: number, isHovered: boolean) => React.ReactNode;
	/**
	 * Обработчик клика по строке (например, для перехода на детальную страницу)
	 * @param item - объект данных текущей строки
	 * @param index - индекс строки
	 */
	onRowClick?: (item: T, index: number) => void;
	/** Индекс раскрытой строки (для expandable row), либо null если ничего не раскрыто */
	expandedRow?: number | null;
	/** Функция для изменения раскрытой строки */
	setExpandedRow?: (idx: number | null) => void;
	/**
	 * Функция для рендера раскрытой строки (expandable row)
	 * @param item - объект данных текущей строки
	 * @param index - индекс строки
	 * @returns ReactNode, который будет отрисован под основной строкой
	 */
	renderExpandedRow?: (item: T, index: number) => React.ReactNode;
	/**
	 * Функция для удаления обьектов
	 * @param id - Уникальное значение каждой строки в таблицы (getRowId)
	 */
	onDeleteFn: (ids: (string | number)[]) => Promise<void>;
	/** Компонент для пустого состояния таблицы */
	emptyComponent?: React.ReactNode;
	/** Включить action массового трансфера в toolbar */
	enableTransferSelectedAction?: boolean;
	/** Обработчик массового трансфера выбранных строк */
	onTransferSelected?: (ids: (string | number)[]) => void | Promise<void>;
	/** Интервал автообновления таблицы, мс */
	refetchInterval?: number;
}

function GenericTable<T>({ tableKey, columns, fetchFn, sysToken, getRowId, rowActions, onRowClick, renderExpandedRow, onDeleteFn, emptyComponent, enableTransferSelectedAction = false, onTransferSelected }: GenericTableProps<T>) {
	const { columns: columnsState, setColumns } = useTableColumns(tableKey, columns);
	const { limit, setLimit } = useTableLimit(tableKey, defaultLimit);

	const [hoveredRow, setHoveredRow] = useState<number | null>(null);

	const [pageParam, setPageParam] = useQueryState('page', '1');
	const [search, setSearch] = useQueryState('search', '');

	const page = Math.max(1, parseInt(pageParam, 10) || 1);
	const setPage = useCallback(
		(value: React.SetStateAction<number>) => {
			setPageParam(String(typeof value === 'function' ? value(page) : value));
		},
		[page, setPageParam]
	);

	/** Column sort */
	const [columnSort, setColumnSort] = useState<ColumnSortState>(() => {
		const saved = localStorage.getItem(CACHEKEYs_TABLE_COLUMN_SORT(tableKey));
		return saved ? JSON.parse(saved) : { key: 'created_at', direction: 'desc' };
	});

	useEffect(() => {
		localStorage.setItem(CACHEKEYs_TABLE_COLUMN_SORT(tableKey), JSON.stringify(columnSort));
	}, [columnSort, tableKey]);

	/** If change limit go to page (1) */
	useEffect(() => {
		setPage(1);
	}, [limit]);

	/** Get data for table */
	const fetcher = useCallback((signal?: AbortSignal) => fetchFn({ page, limit, search, columnSort }, signal), [fetchFn, page, limit, search, columnSort]);
	const { data, loading } = useHandleServer([tableKey, page, limit, search, columnSort, sysToken], fetcher, { placeholderData: keepPreviousData });

	const items = data?.items ?? [];
	const totalAll = data?.total ?? 0;
	const totalPages = data?.totalPages ?? 1;

	useEffect(() => {
		if (!data) return;
		if (page > data.totalPages) {
			setPage(Math.max(1, data.totalPages));
		}
	}, [data, page]);

	const { checkedItems, allChecked, selectedIds, selectedCount, hasSelected, toggleAll, toggleOne, resetSelection, replaceSelection, selectionMode } = useTableSelection(items, getRowId);

	const selectAllAcrossPages = useCallback(async () => {
		if (totalAll === 0) {
			replaceSelection([], 'page');
			return;
		}

		const tryFetchAllOnce = async () => {
			const result = await fetchFn({ page: 1, limit: Math.max(totalAll, 1), search, columnSort });
			return result.items.map(getRowId);
		};

		const fetchInChunks = async (chunkSize: number) => {
			const pageSize = Math.max(chunkSize, 1);
			const pages = Math.max(1, Math.ceil(totalAll / pageSize));
			const allIds: (string | number)[] = [];

			for (let nextPage = 1; nextPage <= pages; nextPage += 1) {
				const result = await fetchFn({ page: nextPage, limit: pageSize, search, columnSort });
				result.items.forEach((item) => allIds.push(getRowId(item)));
			}

			return allIds;
		};

		let allIds: (string | number)[] = [];

		try {
			allIds = await tryFetchAllOnce();

			if (allIds.length < totalAll) {
				allIds = await fetchInChunks(1000);
			}
		} catch {
			allIds = await fetchInChunks(1000);
		}

		replaceSelection(allIds, 'allPages');
	}, [totalAll, replaceSelection, limit, fetchFn, search, columnSort, getRowId]);

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
		setColumnSort,
		checkedItems,
		allChecked,
		toggleAll,
		toggleOne,
		selectedIds,
		selectedCount,
		hasSelected,
		resetSelection,
		replaceSelection,
		selectAllAcrossPages,
		selectionMode,
		enableTransferSelectedAction,
		onTransferSelected,
		hoveredRow,
		setHoveredRow,
		columns: columnsState,
		setColumns,
		getRowId,
		rowActions,
		onRowClick,
		onDeleteFn,
		renderExpandedRow,
	};

	if (loading) {
		return <Fallback />;
	}

	if (items.length === 0 && !search && emptyComponent) {
		return <>{emptyComponent}</>;
	}

	return (
		<TableContext.Provider value={ctxValue}>
			<div className="flex flex-col flex-1 min-h-0">
				<TableToolbar />

				<div className="min-h-0 overflow-y-auto border border-[#dedede] rounded-[6px] mt-[15px] mb-[15px]">
					<table id="table" className="w-full">
						<TableHead />
						<TableBody />
					</table>
				</div>

				<div className="relative shrink-0">
					<PaginationTable />
				</div>
			</div>
		</TableContext.Provider>
	);
}

export default GenericTable;
