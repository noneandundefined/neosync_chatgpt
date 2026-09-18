import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { DEVICE_MODELS } from '@/constants/DeviceModels.constant';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { TRACK_OPERATING_MODE_ADVENCED, TRACK_OPERATING_MODE_MINIMAL } from '@/constants/Track.constants';

const AlternativeSourcesPositioning = () => {
	const { t } = useTranslation();
	const { model, isTemplate } = useConfigurationWrapperContext();

	const operatingModeArray = isTemplate || model === DEVICE_MODELS.ADM500 ? TRACK_OPERATING_MODE_ADVENCED : TRACK_OPERATING_MODE_MINIMAL;

	return (
		<div className="space-y-3">
			<LabeledField label="message.multiapn-operating-mode">
				<CFGSelect uid={UIDs.ALTERNATIVE_SOURCES_POSITIONING_OPERATING_MODE} array={operatingModeArray} />
			</LabeledField>

			<LabeledField label="message.alternative-query-period">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.ALTERNATIVE_SOURCES_POSITIONING_QUERY_PERIOD} />
					<label>{t('label.s')}</label>
				</div>
			</LabeledField>

			<LabeledField label="message.alternative-min-time-between-requests">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.ALTERNATIVE_SOURCES_POSITIONING_MIN_TIME_REQUESTS} />
					<label>{t('label.s')}</label>
				</div>
			</LabeledField>

			<LabeledField label="message.alternative-coordinate-accuracy">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.ALTERNATIVE_SOURCES_POSITIONING_ACCURACY_COORDINATES} />
					<label>{t('label.m')}</label>
				</div>
			</LabeledField>

			<LabeledField label="label.token">
				<CFGInput uid={UIDs.ALTERNATIVE_SOURCES_POSITIONING_TOKEN} />
			</LabeledField>
		</div>
	);
};

export default AlternativeSourcesPositioning;
