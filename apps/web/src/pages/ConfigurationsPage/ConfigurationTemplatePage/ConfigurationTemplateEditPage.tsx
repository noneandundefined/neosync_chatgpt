import React from 'react';
import { useParams } from 'react-router-dom';
import PageMeta from '@/components/PageMeta/PageMeta';
import ConfigurationTemplateEditorPage from './ConfigurationTemplateEditorPage';

const ConfigurationTemplateEditPage = () => {
	const { id } = useParams();
	const templateId = Number(id);

	return (
		<React.Fragment>
			<PageMeta
				descriptionKey="message.meta-description-configuration-template-edit"
				ogTitleKey="message.og-title-configuration-template-edit"
				ogDescriptionKey="message.og-description-configuration-template-edit"
				path={`/configurations/templates/${id}`}
			/>

			<ConfigurationTemplateEditorPage templateId={Number.isFinite(templateId) && templateId > 0 ? templateId : undefined} />
		</React.Fragment>
	);
};

export default ConfigurationTemplateEditPage;
