import React from 'react';
import FirmwarePage from './FirmwarePage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IFirmwarePage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-firmwares-devices" ogTitleKey="message.og-title-firmwares-devices" ogDescriptionKey="message.og-description-firmwares-devices" path="/firmwares" />

			<FirmwarePage />
		</React.Fragment>
	);
};

export default IFirmwarePage;
