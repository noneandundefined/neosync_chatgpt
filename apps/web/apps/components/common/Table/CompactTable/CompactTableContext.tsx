import { ColumnSortState } from '../GenericTable/GenericTable';
import { createContext, Dispatch, SetStateAction, useContext } from 'react';

interface CompactTableContextType<T> {
	search: string;
	setSearch: (s: string) => void;
	limit: number;
	setLimit: Dispatch<SetStateAction<number>>;
	page: number;
	setPage: Dispatch<SetStateAction<number>>;
	totalPages: number;
	totalAll: number;
	pageData: T[];
	loading: boolean;
	columnSort: ColumnSortState;
	checkedItems: Record<string | number, boolean>;
	allChecked: boolean;
	toggleAll: () => void;
	toggleOne: (id: string | number) => void;
	selectedIds: (string | number)[];
	selectedCount: number;
	hasSelected: boolean;
	hoveredRow: number | null;
	setHoveredRow: (idx: number | null) => void;
	getRowId: (item: T) => string | number;
	onRowClick?: (item: T, idx: number) => void;
	onTitleClick?: (item: T, idx: number) => void;
	onDeleteFn?: (ids: (string | number)[]) => Promise<void>;
	bulkActions?: React.ReactNode;
	renderTitle: (item: T, index: number) => React.ReactNode;
	renderSubtitle?: (item: T, index: number) => React.ReactNode;
	renderLabel?: (item: T, index: number) => React.ReactNode;
	renderMeta?: (item: T, index: number) => React.ReactNode;
	renderTime?: (item: T, index: number) => React.ReactNode;
	isUnread?: (item: T) => boolean;
	footer?: React.ReactNode;
	selectionMode: 'multiple' | 'single';
}

const CompactTableContext = createContext<CompactTableContextType<any> | null>(null);

export function useCompactTableContext<T>(): CompactTableContextType<T> {
	const ctx = useContext<CompactTableContextType<T> | null>(CompactTableContext);
	if (!ctx) throw new Error('CompactTableContext used outside of provider');

	return ctx as CompactTableContextType<T>;
}

export default CompactTableContext;
