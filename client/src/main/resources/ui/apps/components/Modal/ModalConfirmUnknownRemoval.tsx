import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';

interface ModalConfirmUnknownRemovalProps {
	onDelete: () => Promise<void>;
}

/* Modal component for comfirm delete element */
/* Use for typing word: DELETE */
const ModalConfirmUnknownRemoval: React.FC<ModalConfirmUnknownRemovalProps> = ({ onDelete }) => {
	const { t } = useTranslation();

	const [textDelete, setTextDelete] = useState<string>('');

	const REQUIRED_WORD = t('label.delete').toUpperCase();
	const isConfirm = textDelete.trim().toLowerCase() === REQUIRED_WORD.toLowerCase();

	const handleDelete = async () => {
		await onDelete();
	};

	return (
		<React.Fragment>
			<div className="mt-4 sm:px-0 space-y-2">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">{t('message.delete-item-title')}</p>

				<ul className="list-disc list-inside text-sm sm:text-base space-y-1">
					<li>{t('message.delete-item-consequence-1')}</li>
					<li>{t('message.delete-item-consequence-2')}</li>
					<li>{t('message.delete-item-consequence-3')}</li>
				</ul>
			</div>

			<div className="space-y-2 mt-5">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">
					{t('message.confirm-type')} <span className="text-red-700">"{t('label.delete').toUpperCase()}"</span> {t('message.box-below')}
				</p>

				<input type="text" value={textDelete} onChange={(e) => setTextDelete(e.target.value)} />
			</div>

			<GUIButton
				className={`!rounded-none mt-3`}
				onClick={() => {
					isConfirm && handleDelete();
				}}
				disabled={!isConfirm}
			>
				{t('label.delete')}
			</GUIButton>
		</React.Fragment>
	);
};

export default ModalConfirmUnknownRemoval;
