import useIsMobile from '@/hooks/useIsMobile';
import { useEffect, useRef, useState } from 'react';

export type PositionType = 'right' | 'left' | 'top' | 'bottom';

interface TooltipProps {
	title: string;
	children: React.ReactNode;
	position?: PositionType;
	className?: string;
}

const Tooltip: React.FC<TooltipProps> = ({ title, children, className = '', position = 'right' }) => {
	const isMobile = useIsMobile();

	const wrapperRef = useRef<HTMLSpanElement>(null);
	const tooltipRef = useRef<HTMLDivElement>(null);

	const [visible, setVisible] = useState(false);
	const [positionState, setPositionState] = useState<'right' | 'left' | 'top' | 'bottom'>(position);

	useEffect(() => {
		if (!visible) return;

		const wrapper = wrapperRef.current;
		const tooltip = tooltipRef.current;
		if (!wrapper || !tooltip) return;

		const wrapperRect = wrapper.getBoundingClientRect();
		const tooltipRect = tooltip.getBoundingClientRect();

		const spaceRight = window.innerWidth - wrapperRect.right;
		const spaceLeft = wrapperRect.left;
		const spaceTop = wrapperRect.top;
		const spaceBottom = window.innerHeight - wrapperRect.bottom;

		let newPosition = position;

		if (position === 'right' && spaceRight < tooltipRect.width + 12) {
			if (spaceLeft > tooltipRect.width) newPosition = 'left';
		}

		if (position === 'left' && spaceLeft < tooltipRect.width + 12) {
			if (spaceRight > tooltipRect.width) newPosition = 'right';
		}

		if (position === 'top' && spaceTop < tooltipRect.height + 12) {
			if (spaceBottom > tooltipRect.height) newPosition = 'bottom';
		}

		if (position === 'bottom' && spaceBottom < tooltipRect.height + 12) {
			if (spaceTop > tooltipRect.height) newPosition = 'top';
		}

		setPositionState(newPosition);
	}, [visible, position]);

	const positionClasses = {
		right: 'left-full ml-2 top-1/2 -translate-y-1/2',
		left: 'right-full mr-2 top-1/2 -translate-y-1/2',
		top: 'bottom-full mb-2 left-1/2 -translate-x-1/2',
		bottom: 'top-full mt-2 left-0',
	};

	const animationClasses = {
		right: visible ? 'opacity-100 translate-x-0' : 'opacity-0 translate-x-2',
		left: visible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-2',
		top: visible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-2',
		bottom: visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2',
	};

	if (isMobile) return children;

	return (
		<span ref={wrapperRef} className={`relative ${className}`} onMouseEnter={() => setVisible(true)} onMouseLeave={() => setVisible(false)}>
			{children}

			{title && (
				<div
					ref={tooltipRef}
					className={`w-max max-w-[450px] absolute z-[9999] break-words whitespace-normal rounded-xs bg-[#444] border border-[#333] rounded-[7px] px-[0.5rem] py-[0.3rem] text-[0.85rem] text-white pointer-events-none transition-all duration-200 ease-out
                ${positionClasses[positionState]} ${animationClasses[positionState]}`}
				>
					{title}
				</div>
			)}
		</span>
	);
};

export default Tooltip;
