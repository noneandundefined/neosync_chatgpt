import React from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { STATIC_POWER_MODE_3, STATIC_POWER_MODE_BASE } from '@/constants/StaticPowerMode.constant';

interface DeepStaticModeSettingsProps {
	staticPowerMode: number | null;
}

const DeepStaticModeSettings: React.FC<DeepStaticModeSettingsProps> = ({ staticPowerMode }) => {
	const { t } = useTranslation();
	const { support, isTemplate } = useConfigurationWrapperContext();

	if (staticPowerMode == null) return;

	return (
		<div className="space-y-4">
			<div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-5">
				<p className="whitespace-nowrap text-[14px]">{t('message.deep-static-mode')}</p>
				<CFGSelect uid={UIDs.STATIC_POWER_MODE} array={isTemplate || support.ModeHybrid ? STATIC_POWER_MODE_3 : STATIC_POWER_MODE_BASE} />
			</div>

			<div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-2 lg:gap-5">
				<p className="whitespace-nowrap text-[14px]">{t('message.deep-static-mode-enter-timeout')}</p>
				<Tooltip title={t('message.tooltip-deep-sleep-time')}>
					<div className="flex items-center gap-2">
						<CFGInput uid={UIDs.STATIC_POWER_DELAY_MIN} disabled={staticPowerMode === 0} />
						<label>{t('label.min')}</label>
					</div>
				</Tooltip>
			</div>
		</div>
	);
};

export default DeepStaticModeSettings;
