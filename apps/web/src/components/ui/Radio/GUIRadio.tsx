import { UseFormRegisterReturn } from 'react-hook-form';

interface Props {
	label: string;
	value: string;
	size?: string;
	register?: UseFormRegisterReturn;
	checked?: boolean;
	onChange?: (checked: boolean) => void;
}

const GUIRadio = ({ label, value, size = '18px', register, checked, onChange }: Props) => {
	return (
		<label className="inline-flex items-center cursor-pointer gap-3">
			<input
				type="radio"
				value={value}
				{...register}
				className="w-auto"
				checked={checked}
				onChange={(e) => {
					register?.onChange(e);
					onChange?.(e.target.checked);
				}}
				style={{ width: size, height: size }}
			/>
			<span>{label}</span>
		</label>
	);
};

export default GUIRadio;
