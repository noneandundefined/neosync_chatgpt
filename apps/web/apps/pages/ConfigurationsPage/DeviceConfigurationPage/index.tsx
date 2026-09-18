import React from 'react';
import { useTranslation } from 'react-i18next';
import PageMeta from '@/components/PageMeta/PageMeta';
import { useDecryptedImei } from '@/hooks/useDecryptedImei';
import DeviceConfigurationPage from './DeviceConfigurationPage';

const IDeviceConfigurationPage = () => {
	const { t } = useTranslation();
	const imei = useDecryptedImei();

	return (
		<React.Fragment>
			<PageMeta
				descriptionKey="message.meta-description-terminal-config"
				ogTitleKey="message.og-title-terminal-config"
				ogDescriptionKey="message.og-description-terminal-config"
				title={`*${imei?.slice(-4)} - ${t('label.page-terminal-config')}`}
				noIndex
			/>

			<DeviceConfigurationPage />
		</React.Fragment>
	);
};

export default IDeviceConfigurationPage;
