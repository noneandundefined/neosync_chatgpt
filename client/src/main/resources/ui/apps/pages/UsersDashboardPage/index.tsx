import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import UsersDashboardPage from './UsersDashboardPage';

const IUsersDashboardPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-user-management" ogTitleKey="message.og-title-user-management" ogDescriptionKey="message.og-description-user-management" path="/dashboard/users" noIndex />

			<UsersDashboardPage />
		</React.Fragment>
	);
};

export default IUsersDashboardPage;
