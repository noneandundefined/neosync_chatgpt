import { createContext, Dispatch, SetStateAction, useContext } from 'react';
import { Column } from './types';
import { ColumnSortState } from './GenericTable';

interface TableContextType<T> {
	// search
	search: string;
	setSearch: (s: string) => void;

	// pagination limit
	limit: number;
	setLimit: React.Dispatch<React.SetStateAction<number>>;

	// pagination
	page: number;
	setPage: Dispatch<SetStateAction<number>>;
	totalPages: number;
	totalAll: number;
	pageData: T[];
	loading: boolean;

	// sort
	columnSort: ColumnSortState;
	setColumnSort: Dispatch<React.SetStateAction<ColumnSortState>>;

	// checked
	checkedItems: Record<string | number, boolean>;
	allChecked: boolean;
	toggleAll: () => void;
	toggleOne: (id: string | number) => void;
	selectedIds: (string | number)[];
	selectedCount: number;
	hasSelected: boolean;
	resetSelection: () => void;
	replaceSelection: (ids: (string | number)[], mode?: 'page' | 'allPages') => void;
	selectAllAcrossPages: () => Promise<void>;
	selectionMode: 'page' | 'allPages';
	enableTransferSelectedAction: boolean;
	onTransferSelected?: (ids: (string | number)[]) => void | Promise<void>;

	// hover
	hoveredRow: number | null;
	setHoveredRow: (idx: number | null) => void;

	// columns
	columns: Column<T>[];
	setColumns: React.Dispatch<React.SetStateAction<Column<T>[]>>;

	getRowId: (item: T) => string | number;

	//action/clicks
	onRowClick?: (item: T, idx: number) => void;
	rowActions?: (item: T, idx: number, hovered: boolean) => React.ReactNode;
	onDeleteFn: (ids: (string | number)[]) => Promise<void>;

	// expandable
	renderExpandedRow?: (item: T, idx: number) => React.ReactNode;
}

const TableContext = createContext<TableContextType<any> | null>(null);

export function useTableContext<T>(): TableContextType<T> {
	const ctx = useContext<TableContextType<T> | null>(TableContext);
	if (!ctx) throw new Error('TableContext used outside of provider');

	return ctx as TableContextType<T>;
}

export default TableContext;
