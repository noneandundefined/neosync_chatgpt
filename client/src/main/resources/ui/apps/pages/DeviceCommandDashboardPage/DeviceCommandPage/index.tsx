import React from 'react';
import DeviceCommandPage from './DeviceCommandPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IDeviceCommandPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-send-commands" ogTitleKey="message.og-title-send-commands" ogDescriptionKey="message.og-description-send-commands" path="/devices/command" />

			<DeviceCommandPage />
		</React.Fragment>
	);
};

export default IDeviceCommandPage;
