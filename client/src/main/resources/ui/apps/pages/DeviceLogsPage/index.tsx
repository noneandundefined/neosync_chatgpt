import React from 'react';
import DeviceLogsPage from './DeviceLogsPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IDeviceLogsPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-device-logs" ogTitleKey="message.og-title-device-logs" ogDescriptionKey="message.og-description-device-logs" noIndex />

			<DeviceLogsPage />
		</React.Fragment>
	);
};

export default IDeviceLogsPage;
