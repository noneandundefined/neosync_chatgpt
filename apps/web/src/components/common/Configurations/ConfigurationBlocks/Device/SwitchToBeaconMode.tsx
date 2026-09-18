import React from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import IndexConfiguration from '../IndexConfiguration';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { DETERMING_COORDINATES } from '@/constants/SwitchToBeaconMode.constant';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const SwitchToBeaconMode = () => {
	const { t } = useTranslation();

	return (
		<React.Fragment>
			<IndexConfiguration title={t('message.setting-activity-mode')} content={<SettingActivityMode />} />
			<IndexConfiguration title={t('message.setting-sleep-mode')} content={<SettingSleepMode />} />
		</React.Fragment>
	);
};

const SettingActivityMode = () => {
	const { t } = useTranslation();
	const { support } = useConfigurationWrapperContext();

	const modes = [
		{
			label: 'message.min-time-active',
			description: t('message.tooltip-minimum-activity-time'),
			index: 0,
		},
		{ label: 'message.max-time-active', description: t('message.tooltip-maximum-activity-time'), index: 1 },
	];

	return (
		<div className="flex min-w-auto sm:min-w-[17rem] flex-col space-y-5">
			{modes.map((mode) => (
				<LabeledField label={mode.label} key={mode.index}>
					<Tooltip title={mode.description} position="bottom">
						<div className="flex items-center gap-2">
							<CFGArrayItemInput uid={UIDs.BEACON_TIME_SEC} index={mode.index} />
							<label>{t('label.s')}</label>
						</div>
					</Tooltip>
				</LabeledField>
			))}

			{support.Gsm && (
				<LabeledField label="message.determining-coordinates">
					<Tooltip title={t('message.tooltip-coordinates-detection-source')}>
						<CFGSelect uid={UIDs.BEACON_MODE_GSM_GNSS} index={BITs.BEACON_MODE_GNSS_ACTIVE.BEACON_MODE_GNSS_ACTIVE} array={DETERMING_COORDINATES} />
					</Tooltip>
				</LabeledField>
			)}
		</div>
	);
};

const SettingSleepMode = () => {
	const { t } = useTranslation();
	const { support } = useConfigurationWrapperContext();

	const modes = [
		{ label: 'message.sleep-time', description: t('message.tooltip-sleep-time'), index: 2 },
		{ label: 'message.sleep-time-after-error', description: t('message.tooltip-sleep-time-after-error'), index: 3 },
	];

	return (
		<div className="flex min-w-auto sm:min-w-[17rem] flex-col space-y-5">
			{modes.map((mode) => (
				<LabeledField label={mode.label} key={mode.index}>
					<Tooltip title={mode.description} position="bottom">
						<div className="flex items-center gap-2">
							<CFGArrayItemInput uid={UIDs.BEACON_TIME_SEC} index={mode.index} />
							<label>{t('label.s')}</label>
						</div>
					</Tooltip>
				</LabeledField>
			))}

			{support.Gsm && (
				<LabeledField label="message.gsm-module" className="flex flex-row items-center justify-between">
					<CFGSwitch uid={UIDs.BEACON_MODE_GSM_GNSS} index={BITs.BEACON_MODE_GNSS_ACTIVE.BEACON_MODE_GSM_SLEEP} vOff={0} vOn={1} />
				</LabeledField>
			)}
		</div>
	);
};

export default SwitchToBeaconMode;
