import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import LabeledField from '@/components/ui/Form/LabeledField';

const ENABLED_VALUE = 0x04000001;
const DISABLED_VALUE = 0xfffffdff;

const Period = () => {
	const { t } = useTranslation();

	return (
		<div className="min-w-auto sm:min-w-[20rem] space-y-4">
			<LabeledField label="message.period-in-motion">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.POINTS_PERIOD_MOVE} />
					<label>{t('label.s')}</label>
				</div>
			</LabeledField>

			<LabeledField label="message.period-at-standstill">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.POINTS_PERIOD_HOLD} />
					<label>{t('label.s')}</label>
				</div>
			</LabeledField>

			<div className="flex items-center justify-between">
				<label className="max-w-[50%]">{t('message.setting-points-only-end-period')}</label>
				<CFGSwitch uid={UIDs.STAT_MASK} vOff={DISABLED_VALUE} vOn={ENABLED_VALUE} />
			</div>
		</div>
	);
};

export default Period;
