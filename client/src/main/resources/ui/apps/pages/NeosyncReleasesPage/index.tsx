import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import NeosyncReleasesPage from './NeosyncReleasesPage';

const INeosyncReleasesPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.releases-meta-description" ogTitleKey="message.releases-og-title" ogDescriptionKey="message.releases-og-description" path="/neosync/releases" />

			<NeosyncReleasesPage />
		</React.Fragment>
	);
};

export default INeosyncReleasesPage;
