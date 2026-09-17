import { useState } from 'react';

function useTableCheckbox<T>(data: T[], getRowId: (item: T) => string | number) {
	const [checkedItems, setCheckedItems] = useState<{
		[key: string | number]: boolean;
	}>({});

	const allChecked = data.length > 0 && data.every((item) => checkedItems[getRowId(item)]);
	const selectedCount = Object.values(checkedItems).filter(Boolean).length;
	const selectedIds = Object.keys(checkedItems).filter((id) => checkedItems[id]);
	const hasSelected = selectedIds.length > 0;

	const handleCheckboxChange = (index: string | number) => {
		setCheckedItems((prev) => ({
			...prev,
			[index]: !prev[index],
		}));
	};

	const handleCheckAll = () => {
		const newChecked: { [key: string]: boolean } = {};

		if (!allChecked) {
			data.forEach((item) => {
				newChecked[getRowId(item)] = true;
			});
		} else {
			data.forEach((item) => {
				newChecked[getRowId(item)] = false;
			});
		}

		setCheckedItems(newChecked);
	};

	const resetChecked = () => setCheckedItems({});

	return {
		checkedItems,
		allChecked,
		selectedCount,
		selectedIds,
		hasSelected,
		handleCheckboxChange,
		handleCheckAll,
		resetChecked,
		setCheckedItems,
	};
}

export default useTableCheckbox;
