import React from 'react';
import { useEffect, useRef, useState } from 'react';

interface GUISelectProps extends Omit<React.HTMLAttributes<HTMLDivElement>, 'onChange'> {
	children: React.ReactNode;
	value?: string | number;
	onChange?: React.ChangeEventHandler<HTMLSelectElement>;
	placeholder?: string;
}

type Option = {
	value: string | number;
	label: React.ReactNode;
	disabled?: boolean;
};

const arrowIcon = `url("data:image/svg+xml;utf8,<svg fill='gray' height='21' viewBox='0 0 24 24' width='21' xmlns='http://www.w3.org/2000/svg'><path d='M7 10l5 5 5-5z'/></svg>")`;

const GUISelect: React.FC<GUISelectProps> = ({ children, value, onChange, placeholder = '', ...rest }) => {
	const [internalValue, setInternalValue] = useState<string | number>();

	const [openUp, setOpenUp] = useState(false);
	const [open, setOpen] = useState<boolean>(false);

	const ref = useRef<HTMLDivElement>(null);

	const options: Option[] = React.Children.toArray(children)
		.filter((child): child is React.ReactElement => React.isValidElement(child) && child.type === 'option')
		.map((child) => ({
			value: child.props.value,
			label: child.props.children,
			disabled: child.props.disabled,
		}));

	const selectedValue = value ?? internalValue ?? options[0]?.value;
	const selected = options.find((o) => o.value === selectedValue);

	useEffect(() => {
		const handleClickOutside = (e: MouseEvent) => {
			if (ref.current && !ref.current.contains(e.target as Node)) {
				setOpen(false);
			}
		};

		document.addEventListener('mousedown', handleClickOutside);
		return () => document.removeEventListener('mousedown', handleClickOutside);
	}, []);

	const triggerChange = (val: string | number) => {
		if (value === undefined) {
			setInternalValue(val);
		}

		if (!onChange) return;

		const event = {
			target: { value: val },
		} as React.ChangeEvent<HTMLSelectElement>;

		onChange(event);
	};

	const toggleOpen = () => {
		if (!ref.current) return;

		const rect = ref.current.getBoundingClientRect();
		const viewportHeight = window.innerHeight;

		const DROPDOWN_MAX_HEIGHT = viewportHeight * 0.35;

		const spaceBelow = viewportHeight - rect.bottom;
		const spaceAbove = rect.top;

		setOpenUp(spaceBelow < DROPDOWN_MAX_HEIGHT && spaceAbove > spaceBelow);
		setOpen((prev) => !prev);
	};

	return (
		<div ref={ref} className={`relative min-w-[16vw] w-full ${rest.className}`} style={rest.style} {...rest}>
			<div
				onClick={toggleOpen}
				className="flex items-center pr-7 pb-3 pt-2 bg-transparent text-[13px] cursor-pointer border-b border-[#d1d5db] hover:border-black"
				style={{
					backgroundImage: arrowIcon,
					backgroundRepeat: 'no-repeat',
					backgroundPosition: 'right center',
				}}
			>
				<span className="min-w-0 truncate">{selected?.label || placeholder}</span>
			</div>

			<ul
				className={`bg-white absolute ${openUp ? 'bottom-0' : 'top-0'} left-0 w-full origin-top transition-all duration-100 ease-out z-[9999] max-h-[35vh] overflow-y-auto ${open ? 'scale-100 opacity-100' : 'scale-50 opacity-0 pointer-events-none'}`}
				style={{ boxShadow: '0 9px 17px 0 rgba(0, 0, 0, 0.3)' }}
			>
				{options.map((option) => {
					const active = option.value === selectedValue;

					return (
						<li
							key={option.value}
							className={`cursor-pointer text-[13px] ${active ? 'text-black' : 'text-[#7b838d]'} p-3 hover:bg-[#f6f6f6]`}
							onClick={() => {
								if (option.disabled) return;

								triggerChange(option.value);
								setOpen(false);
							}}
						>
							{option.label}
						</li>
					);
				})}
			</ul>
		</div>
	);
};

export default GUISelect;
