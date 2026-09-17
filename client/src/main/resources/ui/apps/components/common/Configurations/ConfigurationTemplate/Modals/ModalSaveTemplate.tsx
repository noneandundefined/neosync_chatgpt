import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { useQueryClient } from '@tanstack/react-query';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUIButton from '@/components/ui/Button/GUIButton';
import { Dispatch, SetStateAction, useState } from 'react';
import LabeledField from '@/components/ui/Form/LabeledField';
import { basicConfigurationTemplateCreate } from '@/rest/configurationTemplateAPI';
import { compactTableQueryKey } from '@/components/common/Table/CompactTable/types';
import { clearTemplateDraft, getTemplateDraftCfgHash, getTemplateDraftChanges } from '@/utils/TemplateDraftUtils';

interface ModalSaveTemplateProps {
	saveLoading: boolean;
	setSaveLoading: Dispatch<SetStateAction<boolean>>;
	onClose: () => void;
}

const ModalSaveTemplate: React.FC<ModalSaveTemplateProps> = ({ saveLoading, setSaveLoading, onClose }) => {
	const { t } = useTranslation();

	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const [name, setName] = useState('');
	const [typeOfSaving, setTypeOfSaving] = useState<'modified_ones' | 'full'>('modified_ones');

	const handleSave = async () => {
		if (!name.trim()) return;

		const changes = getTemplateDraftChanges();
		if (Object.keys(changes).length === 0) return;

		try {
			setSaveLoading(true);

			await basicConfigurationTemplateCreate({
				name: name.trim(),
				type_of_saving: typeOfSaving,
				cfg_hash: getTemplateDraftCfgHash(),
				changes,
			});

			clearTemplateDraft();
			onClose();

			await queryClient.invalidateQueries({ queryKey: compactTableQueryKey('configuration-templates') });
			navigate(`${ROUTES.CONFIGURATIONS}?tab=templates`);
		} finally {
			setSaveLoading(false);
		}
	};

	return (
		<div className="space-y-5">
			<LabeledField label="label.title">
				<GUInput value={name} onChange={(e) => setName(e.target.value)} placeholder={t('message.enter-title')} />
			</LabeledField>

			<LabeledField label="label.save-type">
				<GUISelect value={typeOfSaving} onChange={(e) => setTypeOfSaving(e.target.value as 'modified_ones' | 'full')}>
					<option value="modified_ones">{t('message.configuration-template-type-modified')}</option>
				</GUISelect>
			</LabeledField>

			<div className="flex justify-end gap-2 pt-2">
				<GUIButton onClick={onClose}>{t('label.cancel')}</GUIButton>
				<GUIButton onClick={handleSave} disabled={saveLoading}>
					{t('label.save')}
				</GUIButton>
			</div>
		</div>
	);
};

export default ModalSaveTemplate;
