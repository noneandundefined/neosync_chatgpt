import Close from '../@icons/close';
import Tooltip from '../ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { useEffect, useRef, useState } from 'react';
import { useModalContext } from '@/context/useModalContext';

interface ModalProps {
	title: string;
	argv?: number[] | string[];
	width?: string;
	children: React.ReactNode;
}

/* Global component for create modal components */
const Modal: React.FC<ModalProps> = ({ title, argv, width = '600px', children }) => {
	const { t } = useTranslation();

	const { close } = useModalContext();

	const modalRef = useRef<HTMLDivElement>(null);

	const [animationClasses, setAnimationClasses] = useState('opacity-0 scale-95');

	useEffect(() => {
		const timer = setTimeout(() => {
			setAnimationClasses('opacity-100 scale-100');
		}, 1);

		return () => clearTimeout(timer);
	}, []);

	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (modalRef.current && !modalRef.current.contains(event.target as Node)) {
				close();
			}
		};

		document.addEventListener('mousedown', handleClickOutside);
		return () => {
			document.removeEventListener('mousedown', handleClickOutside);
		};
	}, [close]);

	return (
		<div className={`fixed backdrop-blur-[2px] inset-0 z-[1005] flex items-center !justify-center bg-black/30 transition-opacity duration-100 ${animationClasses.split(' ').find((cls) => cls.startsWith('opacity-'))}`}>
			<div ref={modalRef} style={{ width }} className={`bg-white rounded-lg shadow-lg p-3 sm:p-4 max-w-[96%] sm:max-w-[90%] relative transition-transform duration-100 ${animationClasses}`} role="dialog" aria-modal="true">
				<div className="flex items-center justify-between mb-4">
					<div>
						<p className="text-sm sm:text-[15px] text-left font-medium text-[#333]">
							{title} {argv}
						</p>
					</div>

					<Tooltip title={t('label.close')} position="bottom">
						<div className="text-[#666] bg-white hover:bg-[#f1f1f1] hover:bg-[#f9f9f9] rounded-[8px] p-[8px] cursor-pointer text-[1.1rem]" onClick={close}>
							<Close fill="#444" size={17} />
						</div>
					</Tooltip>
				</div>

				{children}
			</div>
		</div>
	);
};

export default Modal;
