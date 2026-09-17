import { useTranslation } from 'react-i18next';
import LowBatteryAlarm from './LowBatteryAlarm';
import SosButtonSettings from './SosButtonSettings';
import DisturbingNumbers from './DisturbingNumbers';
import IndexConfiguration from '../IndexConfiguration';
import EventSetupDescription from './EventSetupDescription';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigEventsInner = () => {
	const { t } = useTranslation();
	const { support } = useConfigurationWrapperContext();

	return (
		<div className="w-full">
			{(support.LowBatteryAlarm || support.Sos) && <IndexConfiguration title={t('message.event-settings-description')} content={<EventSetupDescription />} />}

			<div className="flex flex-col md:flex-row flex-wrap gap-x-4 w-full">
				<IndexConfiguration title={t('message.disturbing-numbers')} content={<DisturbingNumbers />} maxWidth="30rem" />

				<div className="flex-1">
					{support.LowBatteryAlarm && <IndexConfiguration title={t('message.alarm-low-battery')} content={<LowBatteryAlarm />} />}
					{support.Sos && <IndexConfiguration title={t('message.sos-button-settings')} content={<SosButtonSettings />} />}
				</div>
			</div>
		</div>
	);
};

const ConfigEvent = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigEventsInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigEvent;
