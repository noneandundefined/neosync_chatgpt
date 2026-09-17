import ApplicationPage from './ApplicationPage';
import PageMeta from '@/components/PageMeta/PageMeta';
import React from 'react';

const IApplicationPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.home-meta-description" ogTitleKey="message.home-og-title" ogDescriptionKey="message.home-og-description" path="/" />

			<ApplicationPage />
		</React.Fragment>
	);
};

export default IApplicationPage;
