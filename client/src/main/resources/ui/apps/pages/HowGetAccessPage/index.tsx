import React from 'react';
import HowGetAccessPage from './HowGetAccessPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IHowGetAccessPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.how-get-access-meta-description" ogTitleKey="message.how-get-access-home-og-title" ogDescriptionKey="message.how-get-access-og-description" path="/how_get_access" />

			<HowGetAccessPage />
		</React.Fragment>
	);
};

export default IHowGetAccessPage;
