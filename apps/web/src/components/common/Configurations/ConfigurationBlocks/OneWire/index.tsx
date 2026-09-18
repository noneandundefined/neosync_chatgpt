import IButton from './IButton';
import OWTempSn from './OWTempSn';
import { useTranslation } from 'react-i18next';
import AutoSettingTemp from './AutoSettingTemp';
import IndexConfiguration from '../IndexConfiguration';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigWireInner = () => {
	const { t } = useTranslation();
	const { isTemplate } = useConfigurationWrapperContext();

	return (
		<>
			<IndexConfiguration title="IButton" content={<IButton />} />
			{!isTemplate && <IndexConfiguration title={t('message.temperature-sensor-auto-tuning')} content={<AutoSettingTemp />} />}

			<div>
				<p className="font-medium my-3">{t('message.temperature-sensors')}</p>
				<OWTempSn />
			</div>
		</>
	);
};

const ConfigWire = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigWireInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigWire;
