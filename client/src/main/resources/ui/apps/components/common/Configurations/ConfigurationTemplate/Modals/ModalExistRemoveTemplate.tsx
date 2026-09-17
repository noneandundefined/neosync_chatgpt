import { useTranslation } from 'react-i18next';
import GUIButton from '@/components/ui/Button/GUIButton';

interface ModalExistRemoveTemplateProps {
	onClose: () => void;
	onConfirm: () => void;
}

const ModalExistRemoveTemplate: React.FC<ModalExistRemoveTemplateProps> = ({ onClose, onConfirm }) => {
	const { t } = useTranslation();

	return (
		<div className="space-y-5">
			<div className="mt-4 sm:px-0 space-y-2">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">
					{t('message.exit-configuration-template-confirmation')} {t('message.exit-configuration-template-description')}
				</p>
			</div>

			<div className="flex justify-end gap-2 pt-2">
				<GUIButton onClick={onClose}>{t('label.cancel')}</GUIButton>
				<GUIButton onClick={onConfirm}>{t('label.confirm')}</GUIButton>
			</div>
		</div>
	);
};

export default ModalExistRemoveTemplate;
