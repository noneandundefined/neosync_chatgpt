import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import BulkConfigurationPage from './BulkConfigurationPage';

const IBulkConfigurationPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-provisioning-wizard" ogTitleKey="message.og-title-provisioning-wizard" ogDescriptionKey="message.og-description-provisioning-wizard" path="/configurations/bulk/new" />

			<BulkConfigurationPage />
		</React.Fragment>
	);
};

export default IBulkConfigurationPage;
