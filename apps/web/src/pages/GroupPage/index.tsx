import React from 'react';
import GroupPage from './GroupPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IGroupPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-group-management" ogTitleKey="message.og-title-group-management" ogDescriptionKey="message.og-description-group-management" path="/groups" />

			<GroupPage />
		</React.Fragment>
	);
};

export default IGroupPage;
