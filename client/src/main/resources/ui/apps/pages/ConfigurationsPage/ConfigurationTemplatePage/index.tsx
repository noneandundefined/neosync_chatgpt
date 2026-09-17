import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import ConfigurationTemplateCreatePage from './ConfigurationTemplateCreatePage';

const ConfigurationTemplatePage = () => {
	return (
		<React.Fragment>
			<PageMeta
				descriptionKey="message.meta-description-configuration-template-create"
				ogTitleKey="message.og-title-configuration-template-create"
				ogDescriptionKey="message.og-description-configuration-template-create"
				path="/configurations/devices/template"
			/>

			<ConfigurationTemplateCreatePage />
		</React.Fragment>
	);
};

export default ConfigurationTemplatePage;
