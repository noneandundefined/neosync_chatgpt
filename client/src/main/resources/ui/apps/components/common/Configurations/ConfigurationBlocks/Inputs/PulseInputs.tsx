import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { PULSE_INPUTS_DIN0, PULSE_INPUTS_DIN1 } from '@/constants/PulseInputs.constant';

const PulseInputs = () => {
	const { t } = useTranslation();

	return (
		<div>
			<div className="space-y-2 text-[14px]">
				<p>
					<span className="font-bold">{t('label.mode-frequency-meter')}</span> - {t('message.desc-frequency-meter')}
				</p>
				<p>
					<span className="font-bold">{t('label.mode-flow-meter')}</span> - {t('message.desc-flow-meter')}
				</p>
				<p>
					<span className="font-bold">{t('label.mode-instant-flow-meter')}</span> - {t('message.desc-instant-flow-meter')}
				</p>
				<p>
					<span className="font-bold">{t('label.mode-differential-flow-meter')}</span> - {t('message.desc-differential-flow-meter')}
				</p>
				<p>
					<span className="font-bold">{t('label.mode-discrete')}</span> - {t('message.desc-discrete')}
				</p>
				<p>
					<span className="font-bold">{t('label.mode-phase')}</span> - {t('message.desc-phase')}
				</p>
			</div>

			<div className="mt-5 space-y-3">
				{[0, 1].map((index) => (
					<LabeledField key={index} label="message.mode-din-label" argv={index} className="flex flex-col sm:flex-row sm:items-center sm:gap-3">
						<p className="text-sm"></p>
						<CFGSelect uid={UIDs.IMPULSE_INPUTS} index={index} className="sm:min-w-[15rem]" array={index === 0 ? PULSE_INPUTS_DIN0 : PULSE_INPUTS_DIN1} />
					</LabeledField>
				))}
			</div>
		</div>
	);
};

export default PulseInputs;
