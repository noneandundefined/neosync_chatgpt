import { ReactNode, useEffect, useRef } from 'react';

interface DropdownProps {
	children: ReactNode;
	open: boolean;
	close: () => void;
	className?: string;
	stopPropagation?: boolean;
}

const Dropdown: React.FC<DropdownProps> = ({ children, open, close, className, stopPropagation }) => {
	const ref = useRef<HTMLDivElement>(null);

	useEffect(() => {
		if (!open) return;

		const handleClickOutside = (e: MouseEvent) => {
			if (ref.current && !ref.current.contains(e.target as Node)) {
				close();
			}
		};

		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	}, [open, close]);

	if (!open) return null;

	return (
		<div ref={ref} className={className} onClick={stopPropagation ? (e) => e.stopPropagation() : undefined}>
			{children}
		</div>
	);
};

export default Dropdown;
