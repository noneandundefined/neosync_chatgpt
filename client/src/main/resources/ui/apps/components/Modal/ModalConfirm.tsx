import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';

interface ModalConfirmProps {
	message: string;
	onConfirm: () => void;
	onCancel: () => void;
}

const ModalConfirm: React.FC<ModalConfirmProps> = ({ message, onConfirm, onCancel }) => {
	const { t } = useTranslation();

	const handleConfirm = () => {
		onConfirm();
	};

	useEffect(() => {
		const handler = (e: KeyboardEvent) => {
			if (e.key === 'Enter') {
				e.preventDefault();
				handleConfirm();
			}
		};

		window.addEventListener('keydown', handler);
		return () => window.removeEventListener('keydown', handler);
	}, []);

	return (
		<React.Fragment>
			<div className="mt-4 sm:px-0 space-y-2">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">{message}</p>
			</div>

			<div className="flex flex-row items-center justify-end gap-3 mt-5">
				<GUIButton onClick={onCancel}>{t('label.cancel')}</GUIButton>
				<GUIButton onClick={handleConfirm}>{t('label.confirm')}</GUIButton>
			</div>
		</React.Fragment>
	);
};

export default ModalConfirm;
