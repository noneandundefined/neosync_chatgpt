import General from './General';
import Navigate from './Navigate';
import AboutDevice from './AboutDevice';
import CommandDevice from './CommandDevice';
import BatterySaving from './BatterySaving';
import { useTranslation } from 'react-i18next';
import AuthorizedPhones from './AuthorizedPhones';
import VoltageTelemetry from './VoltageTelemetry';
import IndexConfiguration from '../IndexConfiguration';
import PulseInputsTelemetry from './PulseInputsTelemetry';
import AnalogInputsTelemetry from './AnalogInputsTelemetry';
import DeepStaticModeSettings from './DeepStaticModeSettings';

import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import SwitchToBeaconMode from './SwitchToBeaconMode';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigDeviceInner = () => {
	const { t } = useTranslation();

	const { support, manager, isTemplate } = useConfigurationWrapperContext();

	const staticPowerMode = useConfigurationVal(manager, UIDs.STATIC_POWER_MODE) ?? 0;
	const deviceFunctions = useConfigurationVal(manager, UIDs.DEVICE_FUNCTIONS);
	const deviceMode = useConfigurationVal(manager, UIDs.DEVICE_MODE) ?? 0;

	const modeDeviceFuncTracker = !support.ModeHybrid && BitHelperUtils.checkBit(deviceFunctions, BITs.DEVICE_FUNCTIONS.DEVICE_MODE_TRACKER);
	const modeDeviceFuncBeacon = !support.ModeHybrid && !BitHelperUtils.checkBit(deviceFunctions, BITs.DEVICE_FUNCTIONS.DEVICE_MODE_TRACKER);

	const showDeepStatic = isTemplate || (support.ModeHybrid ? deviceMode === 1 || deviceMode === 2 : modeDeviceFuncTracker);
	/** Это маяк(старая/новая версия) или гибрид */
	const showSwitchToBeacon = isTemplate || (support.ModeHybrid ? deviceMode === 0 || deviceMode === 2 : modeDeviceFuncBeacon || staticPowerMode === 4);

	return (
		<div className="flex flex-col lg:flex-row flex-wrap gap-x-4 w-full">
			{!isTemplate && <IndexConfiguration title={t('label.about-device')} content={<AboutDevice />} />}

			<IndexConfiguration title={t('label.general')} content={<General />} />

			{!isTemplate && (
				<>
					<IndexConfiguration title={t('label.navigation')} content={<Navigate />} />
					<div className="flex flex-row md:flex-col gap-x-4">
						<IndexConfiguration title={t('label.mode-analog-inputs')} content={<AnalogInputsTelemetry />} />
						<IndexConfiguration title={t('label.mode-pulse-inputs')} content={<PulseInputsTelemetry />} />
					</div>
					<IndexConfiguration title={t('label.voltage')} content={<VoltageTelemetry />} />
				</>
			)}

			{showDeepStatic && <IndexConfiguration title={t('message.deep-static-mode-settings')} content={<DeepStaticModeSettings staticPowerMode={staticPowerMode} />} />}

			{showSwitchToBeacon && <SwitchToBeaconMode />}

			{support.BatterySaving && <IndexConfiguration title={t('message.battery-saving')} content={<BatterySaving />} />}

			<IndexConfiguration title={t('message.authorized-numbers')} content={<AuthorizedPhones />} />

			{!isTemplate && <IndexConfiguration title={t('message.sending-command')} content={<CommandDevice />} className="w-full lg:basis-full" maxWidth="100%" />}
		</div>
	);
};

const ConfigDevice = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigDeviceInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigDevice;
