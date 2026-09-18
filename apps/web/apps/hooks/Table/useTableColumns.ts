import { useEffect, useState } from 'react';
import { Column } from '@/components/common/Table/GenericTable/types';
import { CACHEKEYs_TABLE_VIEW_COLUMNS } from '@/constants/CacheKeys.constants';

export function useTableColumns<T>(tableKey: string, columns: Column<T>[]) {
	const cacheKey = CACHEKEYs_TABLE_VIEW_COLUMNS(tableKey);

	const [columnsState, setColumnsState] = useState<Column<T>[]>(() => {
		try {
			const saved = localStorage.getItem(cacheKey);
			if (!saved) return columns;

			const visibleMap: Record<string, boolean> = JSON.parse(saved);

			return columns.map((col) => ({
				...col,
				visible: visibleMap[String(col.key)] ?? col.visible ?? true,
			}));
		} catch {
			return columns;
		}
	});

	useEffect(() => {
		setColumnsState((prev) =>
			columns.map((col) => {
				const old = prev.find((c) => c.key === col.key);

				return {
					...col,
					visible: old?.visible ?? col.visible ?? true,
				};
			})
		);
	}, [columns]);

	useEffect(() => {
		const map: Record<string, boolean> = {};

		columnsState.forEach((c) => {
			map[String(c.key)] = c.visible !== false;
		});

		localStorage.setItem(cacheKey, JSON.stringify(map));
	}, [columnsState, cacheKey]);

	return {
		columns: columnsState,
		setColumns: setColumnsState,
	};
}
