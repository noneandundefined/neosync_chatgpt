import Check from '@/components/@icons/check';

interface Props {
	size?: string;
	checked?: boolean;
	label?: React.ReactNode;
	onChange?: (checked: boolean) => void;
}

const GUICheckbox = ({ size = '14px', label, checked = false, onChange }: Props) => {
	return (
		<label className="inline-flex items-center cursor-pointer mr-2">
			<input type="checkbox" checked={checked} className="peer hidden" onChange={(e) => onChange?.(e.target.checked)} />
			<div className="border border-gray-400 rounded-sm flex items-center justify-center peer-checked:bg-blue-600 peer-checked:border-blue-600 transition-colors" style={{ width: size, height: size }}>
				<div>
					<Check
						fill="#fff"
						size={16}
						style={{
							pointerEvents: 'none',
						}}
					/>
				</div>
			</div>

			{label && <div className="ml-2">{label}</div>}
		</label>
	);
};

export default GUICheckbox;
