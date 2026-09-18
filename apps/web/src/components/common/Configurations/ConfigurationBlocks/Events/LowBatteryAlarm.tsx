import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import LabeledField from '@/components/ui/Form/LabeledField';

const LowBatteryAlarm = () => {
	const { t } = useTranslation();

	const [smsSendEnabled, setSmsSendEnabled] = useState<boolean>(false);

	return (
		<div className="space-y-4">
			<div className="flex items-center justify-between">
				<p>{t('message.send-sms-messages')}</p>
				<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_TRUE_1} bit={BITs.CUSTOM_MASK_DEFAULT_TRUE_1.SEND_SMS_MESSAGES} vOff={0} vOn={1} onVisible={(enabled) => setSmsSendEnabled(enabled)} />
			</div>

			<LabeledField label="message.trigger-threshold" disabled={!smsSendEnabled}>
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.LOW_BATTERY_LEVEL} disabled={!smsSendEnabled} />
					<label>{t('label.mV')}</label>
				</div>
			</LabeledField>
		</div>
	);
};

export default LowBatteryAlarm;
