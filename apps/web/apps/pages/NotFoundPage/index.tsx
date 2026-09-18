import React from 'react';
import NotFoundPage from './NotFoundPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const INotFoundPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-not-found" ogTitleKey="message.og-title-not-found" ogDescriptionKey="message.og-description-not-found" noIndex />

			<NotFoundPage />
		</React.Fragment>
	);
};

export default INotFoundPage;
