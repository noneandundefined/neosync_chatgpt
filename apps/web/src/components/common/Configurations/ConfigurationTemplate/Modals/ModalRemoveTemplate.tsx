import { useTranslation } from 'react-i18next';
import { Dispatch, SetStateAction } from 'react';
import GUIButton from '@/components/ui/Button/GUIButton';
import { clearTemplateDraft } from '@/utils/TemplateDraftUtils';
import { useTemplateConfigurationContext } from '@/context/useTemplateConfigurationContext';

interface ModalRemoveTemplateProps {
	saveLoading: boolean;
	setSaveLoading: Dispatch<SetStateAction<boolean>>;
	onClose: () => void;
}

const ModalRemoveTemplate: React.FC<ModalRemoveTemplateProps> = ({ saveLoading, setSaveLoading, onClose }) => {
	const { t } = useTranslation();
	const templateCtx = useTemplateConfigurationContext();

	const handleConfirm = () => {
		try {
			setSaveLoading(true);
			clearTemplateDraft(templateCtx?.templateId);

			window.location.reload();

			onClose();
		} finally {
			setSaveLoading(false);
		}
	};

	return (
		<div className="space-y-5">
			<div className="mt-4 sm:px-0 space-y-2">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">
					{t('message.remove-configuration-template-confirmation')} {t('message.remove-configuration-template-description')}
				</p>
			</div>

			<div className="flex justify-end gap-2 pt-2">
				<GUIButton onClick={onClose} disabled={saveLoading}>
					{t('label.cancel')}
				</GUIButton>

				<GUIButton onClick={handleConfirm} disabled={saveLoading}>
					{t('label.confirm')}
				</GUIButton>
			</div>
		</div>
	);
};

export default ModalRemoveTemplate;
