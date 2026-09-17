import { useTableContext } from '../TableContext';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';

function CommonViewColumn<T>() {
	const { columns, setColumns } = useTableContext<T>();

	const toggleColumn = (key: any, value: boolean) => {
		setColumns((prev) => prev.map((c) => (c.key === key ? { ...c, visible: value } : c)));
	};

	return (
		<div className="scale-100 bg-white min-w-[15rem] flex flex-col" style={{ boxShadow: '0 9px 17px 0 rgba(0, 0, 0, 0.3)' }}>
			{columns.map((col) => (
				<label key={String(col.key)} className="hover:bg-[#f3f3f3] flex items-center px-3 p-2 cursor-pointer">
					<GUICheckbox checked={col.visible !== false} onChange={(v) => toggleColumn(col.key, v)} />
					{col.title}
				</label>
			))}
		</div>
	);
}

export default CommonViewColumn;
