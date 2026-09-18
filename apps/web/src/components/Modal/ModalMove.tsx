import { useEffect, useRef, useState } from 'react';

import Close from '../@icons/close';
import Tooltip from '../ui/Tooltip';
import { useTranslation } from 'react-i18next';

interface ModalMoveProps {
	open: boolean;
	children: React.ReactNode;
	width: string;
	zIndex?: number;
	onClick?: () => void;
	onClose: (event?: any) => void;
}

const ModalMove: React.FC<ModalMoveProps> = ({ open, onClose, onClick, children, width, zIndex }) => {
	const { t } = useTranslation();

	const modalRef = useRef<HTMLDivElement>(null);
	const [isVisible, setIsVisible] = useState(false);
	const [animationClasses, setAnimationClasses] = useState('opacity-0 scale-95');

	const [position, setPosition] = useState({ top: 100, left: 0 });
	const [dragging, setDragging] = useState(false);
	const dragOffset = useRef({ x: 0, y: 0 });

	useEffect(() => {
		if (open) {
			setIsVisible(true);

			const timer = setTimeout(() => {
				if (modalRef.current) {
					const rect = modalRef.current.getBoundingClientRect();
					setPosition({
						top: window.innerHeight / 2 - rect.height / 2,
						left: window.innerWidth / 2 - rect.width / 2,
					});
				}
				setAnimationClasses('opacity-100 scale-100');
			}, 5);

			return () => clearTimeout(timer);
		} else {
			setAnimationClasses('opacity-0 scale-95');

			const timer = setTimeout(() => {
				setIsVisible(false);
				setPosition({ top: window.innerHeight / 2 - 150, left: window.innerWidth / 2 - parseInt(width) / 2 });
			}, 200);

			return () => clearTimeout(timer);
		}
	}, [open, width]);

	const handleMouseDown = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		e.preventDefault();
		setDragging(true);
		const rect = modalRef.current!.getBoundingClientRect();
		dragOffset.current = {
			x: e.clientX - rect.left,
			y: e.clientY - rect.top,
		};
	};

	useEffect(() => {
		const handleMouseMove = (e: MouseEvent) => {
			if (!dragging) return;

			setPosition({
				top: e.clientY - dragOffset.current.y,
				left: e.clientX - dragOffset.current.x,
			});
		};

		const handleMouseUp = () => setDragging(false);
		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);

		return () => {
			document.removeEventListener('mousemove', handleMouseMove);
			document.removeEventListener('mouseup', handleMouseUp);
		};
	}, [dragging]);

	if (!isVisible) return null;

	return (
		<div className={`fixed transition-opacity duration-200 ${animationClasses.split(' ').find((cls) => cls.startsWith('opacity-'))}`} style={{ zIndex }}>
			<div
				ref={modalRef}
				style={{
					width,
					top: position.top,
					left: position.left,
					position: 'fixed',
					transform: 'translate(0,0)',
				}}
				className={`bg-white shadow-lg border border-[#e4e4e4] p-3 sm:p-4 max-w-[96%] sm:max-w-[90%] relative transition-transform duration-200 ${animationClasses}`}
				role="dialog"
				aria-modal="true"
				onClick={(e) => {
					e.stopPropagation();
					onClick?.();
				}}
			>
				<div className="absolute top-0 left-0 h-5 w-full cursor-move" onMouseDown={handleMouseDown} />

				<div className="flex items-center justify-between mb-4">
					<div>
						<p className="text-sm sm:text-[15px] text-left font-medium text-[#333]">{t('label.page-details-command')}</p>
					</div>

					<Tooltip title={t('label.close')} position="bottom">
						<div className="text-[#666] bg-white hover:bg-[#f1f1f1] hover:bg-[#f9f9f9] rounded-[8px] p-[8px] cursor-pointer text-[1.1rem]" onClick={onClose}>
							<Close fill="#444" size={17} />
						</div>
					</Tooltip>
				</div>

				{children}
			</div>
		</div>
	);
};

export default ModalMove;
