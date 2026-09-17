import { toast } from 'react-toastify';
import i18next from 'i18next';
import { QueryClient } from '@tanstack/react-query';
import { compactTableQueryKey } from '@/components/common/Table/CompactTable/types';
import { basicConfigurationTemplateUpdate, ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';
import { clearTemplateDraft, getTemplateDraftCfgHash, getTemplateDraftChanges } from '@/utils/TemplateDraftUtils';

export const saveExistingTemplateDraft = async (templateId: number, template: ConfigurationTemplateResponse, queryClient: QueryClient) => {
	const changes = getTemplateDraftChanges(templateId);

	if (Object.keys(changes).length === 0) {
		toast.error(i18next.t('message.no-changes-provided'));
		return false;
	}

	await basicConfigurationTemplateUpdate(templateId, {
		cfg_hash: getTemplateDraftCfgHash(templateId),
		changes,
		type_of_saving: template.type_of_saving,
	});

	clearTemplateDraft(templateId);
	await queryClient.invalidateQueries({ queryKey: compactTableQueryKey('configuration-templates') });
	await queryClient.invalidateQueries({ queryKey: ['configurationTemplate', templateId] });

	return true;
};
