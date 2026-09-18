import { useCallback, useMemo, useState } from 'react';

export function useTableSelection<T>(data: T[], getId: (item: T) => string | number) {
	const [selectedSet, setSelectedSet] = useState<Set<string | number>>(new Set());
	const [selectionMode, setSelectionMode] = useState<'page' | 'allPages'>('page');
	const ids = data.map(getId);

	const checkedItems = useMemo(() => {
		const map: Record<string | number, boolean> = {};
		selectedSet.forEach((id) => {
			map[id] = true;
		});
		return map;
	}, [selectedSet]);

	const allChecked = ids.length > 0 && ids.every((id) => selectedSet.has(id));
	const selectedIds = Array.from(selectedSet);
	const selectedCount = selectedIds.length;
	const hasSelected = selectedCount > 0;

	const toggleAll = useCallback(() => {
		if (selectionMode === 'allPages' && allChecked) {
			setSelectedSet(new Set());
			setSelectionMode('page');
			return;
		}

		const next = !allChecked;
		setSelectedSet((prev) => {
			const updated = new Set(prev);

			ids.forEach((id) => {
				if (next) {
					updated.add(id);
				} else {
					updated.delete(id);
				}
			});

			return updated;
		});
		setSelectionMode('page');
	}, [allChecked, ids, selectionMode]);

	const toggleOne = useCallback((id: string | number) => {
		setSelectedSet((prev) => {
			const updated = new Set(prev);

			if (updated.has(id)) {
				updated.delete(id);
			} else {
				updated.add(id);
			}

			return updated;
		});
		setSelectionMode('page');
	}, []);

	const reset = useCallback(() => {
		setSelectedSet(new Set());
		setSelectionMode('page');
	}, []);

	const replaceSelection = useCallback((nextIds: (string | number)[], mode: 'page' | 'allPages' = 'page') => {
		setSelectedSet(new Set(nextIds));
		setSelectionMode(mode);
	}, []);

	return {
		checkedItems,
		allChecked,
		selectedIds,
		selectedCount,
		hasSelected,
		toggleAll,
		toggleOne,
		resetSelection: reset,
		replaceSelection,
		selectionMode,
	};
}
