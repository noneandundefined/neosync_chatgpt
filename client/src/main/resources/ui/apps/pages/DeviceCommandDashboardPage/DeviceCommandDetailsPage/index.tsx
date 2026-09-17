import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import DeviceCommandDetailsPage from './DeviceCommandDetailsPage';

const IDeviceCommandDetailsPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-details-command" ogTitleKey="message.og-title-details-command" ogDescriptionKey="message.og-description-details-command" noIndex />

			<DeviceCommandDetailsPage />
		</React.Fragment>
	);
};

export default IDeviceCommandDetailsPage;
