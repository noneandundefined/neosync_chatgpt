import { Sims } from './Sim';
import SimPriority from './SimPriority';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import IndexConfiguration from '../IndexConfiguration';
import MultiApn, { GetModeMultiApn } from './MultiApn';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { DEVICE_MODELS, getStaticSupport } from '@/constants/DeviceModels.constant';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigSimInner = () => {
	const { t } = useTranslation();
	const { model, manager, isTemplate } = useConfigurationWrapperContext();

	const staticSupport = getStaticSupport((model ?? DEVICE_MODELS.ADM333V2) as keyof typeof DEVICE_MODELS, isTemplate);

	const deviceFunction3 = useConfigurationVal(manager, UIDs.DEVICE_FUNCTION_3) ?? 0;
	const enabledApn = GetModeMultiApn(deviceFunction3) !== 0;

	return (
		<div className="flex flex-col sm:flex-row flex-wrap gap-x-4">
			{/* Сим карта */}
			<Sims enabledApn={enabledApn} />

			{/* Multi apn (settings) */}
			{staticSupport.MultiApn && <IndexConfiguration title={t('label.multi-apn').toUpperCase()} content={<MultiApn />} />}

			{/* Sim priority (if not enabled apn) (if sims count > 1) */}
			{!enabledApn && staticSupport.SimCount > 1 && <IndexConfiguration title={t('label.priority')} content={<SimPriority />} />}
		</div>
	);
};

const ConfigSim = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigSimInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigSim;
