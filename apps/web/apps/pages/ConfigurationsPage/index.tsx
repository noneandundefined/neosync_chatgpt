import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import ConfigurationsPage from './ConfigurationsPage';

const IConfigurationsPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-configurations" ogTitleKey="message.og-title-configurations" ogDescriptionKey="message.og-description-configurations" path="/configurations" />

			<ConfigurationsPage />
		</React.Fragment>
	);
};

export default IConfigurationsPage;
