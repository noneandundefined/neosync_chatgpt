import React from 'react';
import NeosyncHelpPage from './NeosyncHelpPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const INeosyncHelpPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.help-meta-description" ogTitleKey="message.help-og-title" ogDescriptionKey="message.help-og-description" path="/neosync/help" />

			<NeosyncHelpPage />
		</React.Fragment>
	);
};

export default INeosyncHelpPage;
