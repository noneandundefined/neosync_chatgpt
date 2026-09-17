import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import { SOS_POLARITY } from '@/constants/Sos.constant';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import LabeledField from '@/components/ui/Form/LabeledField';

const SosButtonSettings = () => {
	const { t } = useTranslation();

	const [sosEnabled, setSosEnabled] = useState<boolean>(false);

	return (
		<div className="space-y-5">
			<div className="flex items-center justify-between">
				<p>{t('label.sos-btn')}</p>
				<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.SOS_BUTTON} vOff={1} vOn={0} onVisible={(enabled) => setSosEnabled(enabled)} />
			</div>

			<LabeledField label="message.sos-button-polarity" disabled={!sosEnabled}>
				<CFGSelect uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.SOS_POLARITY} array={SOS_POLARITY} />
			</LabeledField>
		</div>
	);
};

export default SosButtonSettings;
