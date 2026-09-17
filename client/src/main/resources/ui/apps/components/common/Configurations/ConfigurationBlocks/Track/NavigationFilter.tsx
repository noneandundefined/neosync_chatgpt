import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import LabeledField from '@/components/ui/Form/LabeledField';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';

const NavigationFilter = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-3">
			<LabeledField label="message.number-of-satellites">
				<CFGInput uid={UIDs.NAV_FILTER_SATS_COUNT} />
			</LabeledField>

			<LabeledField label="message.minimum-hdop">
				<CFGInput uid={UIDs.NAV_FILTER_HDOP_TENTHS} placeholder="2.0" />
			</LabeledField>

			<LabeledField label="message.maximum-hdop">
				<CFGInput uid={UIDs.NAV_FILTER_HDOP_MAX} />
			</LabeledField>

			<LabeledField label="message.lower-height-threshold">
				<Tooltip position="bottom" title={t('message.lower-threshold-of-the-validity-zone-for-height')}>
					<div className="flex items-center gap-2">
						<CFGArrayItemInput uid={UIDs.NAV_FILTER_VALID_HEIGHT_DECA_HECTO_METERS} index={0} placeholder="0" />
						<label>{t('label.m')}</label>
					</div>
				</Tooltip>
			</LabeledField>

			<LabeledField label="message.upper-height-threshold">
				<Tooltip position="bottom" title={t('message.upper-threshold-of-the-validity-zone-for-height')}>
					<div className="flex items-center gap-2">
						<CFGArrayItemInput uid={UIDs.NAV_FILTER_VALID_HEIGHT_DECA_HECTO_METERS} index={1} placeholder="0" />
						<label>{t('label.m')}</label>
					</div>
				</Tooltip>
			</LabeledField>
		</div>
	);
};

export default NavigationFilter;
