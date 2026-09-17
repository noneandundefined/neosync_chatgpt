import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import BulkConfigurationDetailsPage from './BulkConfigurationDetailsPage';

const IBulkConfigurationDetailsPage = () => {
	return (
		<React.Fragment>
			<PageMeta
				descriptionKey="message.meta-description-bulk-configuration-details"
				ogTitleKey="message.og-title-bulk-configuration-details"
				ogDescriptionKey="message.og-description-bulk-configuration-details"
				path="/configurations/bulk/:id"
				noIndex
			/>

			<BulkConfigurationDetailsPage />
		</React.Fragment>
	);
};

export default IBulkConfigurationDetailsPage;
