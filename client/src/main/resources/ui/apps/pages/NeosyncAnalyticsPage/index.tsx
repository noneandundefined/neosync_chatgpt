import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import NeosyncAnalyticsPage from './NeosyncAnalyticsPage';

const INeosyncAnalyticsPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.analytics-meta-description" ogTitleKey="message.analytics-og-title" ogDescriptionKey="message.analytics-og-description" path="/neosync/analytics" />

			<NeosyncAnalyticsPage />
		</React.Fragment>
	);
};

export default INeosyncAnalyticsPage;
