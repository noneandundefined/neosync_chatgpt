import Server from './Server';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import IndexConfiguration from '../IndexConfiguration';
import TransmittedDataBlocks from './TransmittedDataBlocks';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { getServerCount } from '@/constants/DeviceModels.constant';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigServerInner = () => {
	const { t } = useTranslation();
	const { support, manager, isTemplate } = useConfigurationWrapperContext();

	const serverCount = getServerCount(support, isTemplate);

	const hosts = useConfigurationVal(manager, UIDs.SERVER_HOST);
	const hostArray = Array.isArray(hosts) ? hosts : [];

	const activeCount = hostArray.filter((host) => host && host !== '').length;

	return (
		<div>
			<div className="flex flex-col sm:flex-row flex-wrap gap-x-3">
				{Array.from({ length: serverCount }, (_, index) => {
					const hostValue = hostArray[index] ?? '';

					return <IndexConfiguration key={index} title={`${t('label.server')} ${index + 1}`} content={<Server index={index} hostValue={hostValue} isLastActive={activeCount === 1 && hostValue !== ''} />} />;
				})}
			</div>

			<IndexConfiguration title={t('message.data-blocks-transmitted')} content={<TransmittedDataBlocks />} />
		</div>
	);
};

const ConfigServer = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigServerInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigServer;
