import React, { useMemo } from 'react';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import EmptyStateTable from './EmptyStateTable';
import { useTableContext } from './TableContext';
import { Column } from './types';
import useIsMobile from '@/hooks/useIsMobile';

interface TableRowProps<T> {
	item: T;
	globalIdx: number;
	columns: Column<T>[];
	getRowId: (item: T) => string | number;
	checkedItems: Record<string | number, boolean>;
	toggleOne: (id: string | number) => void;
	rowActions?: (item: T, index: number, isHovered: boolean) => React.ReactNode;
	hoveredRow: number | null;
	setHoveredRow: (idx: number | null) => void;
	onRowClick?: (item: T, index: number) => void;
	renderExpandedRow?: (item: T, index: number) => React.ReactNode;
}

function TableBody<T>() {
	const { pageData, columns, getRowId, checkedItems, toggleOne, rowActions, onRowClick, hoveredRow, setHoveredRow, renderExpandedRow, page, totalPages } = useTableContext<T>();

	const visibleColumns = useMemo(() => columns.filter((c) => c.visible !== false), [columns]);

	return (
		<tbody>
			{pageData.length > 0 ? (
				pageData.map((item, idx) => {
					const globalIdx = idx + page * totalPages;

					return (
						<TableRow
							key={getRowId(item)}
							item={item}
							globalIdx={globalIdx}
							columns={visibleColumns}
							getRowId={getRowId}
							checkedItems={checkedItems}
							toggleOne={toggleOne}
							rowActions={rowActions}
							hoveredRow={hoveredRow}
							setHoveredRow={setHoveredRow}
							onRowClick={onRowClick}
							renderExpandedRow={renderExpandedRow}
						/>
					);
				})
			) : (
				<EmptyStateTable colSpan={columns.length + (rowActions ? 1 : 0) + 1} />
			)}
		</tbody>
	);
}

function TableRowComponent<T>({ item, globalIdx, columns, getRowId, checkedItems, toggleOne, rowActions, hoveredRow, setHoveredRow, onRowClick, renderExpandedRow }: TableRowProps<T>) {
	const id = getRowId(item);
	const isMobile = useIsMobile();

	const rowCells = useMemo(() => {
		return columns.map((col) => (
			<td key={col.key as string} className="text-center max-w-[15rem]" data-label={col.title}>
				{col.render ? col.render(item, globalIdx) : (item as any)[col.key]}
			</td>
		));
	}, [columns, item, globalIdx]);

	const showActions = isMobile || hoveredRow === globalIdx;

	return (
		<>
			<tr
				onClick={() => onRowClick?.(item, globalIdx)}
				onMouseEnter={!isMobile ? () => setHoveredRow(globalIdx) : undefined}
				onMouseLeave={!isMobile ? () => setHoveredRow(null) : undefined}
				className="hover:bg-[#f6f6f6] text-center"
			>
				<td className="!pt-4">
					<GUICheckbox checked={!!checkedItems[id]} onChange={() => toggleOne(id)} />
				</td>
				{rowCells}
				{rowActions && <td>{rowActions(item, globalIdx, showActions)}</td>}
			</tr>
			{renderExpandedRow && renderExpandedRow(item, globalIdx)}
		</>
	);
}

const TableRow = React.memo(TableRowComponent) as typeof TableRowComponent;

export default TableBody;
