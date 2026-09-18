import { useState } from 'react';
import PulseInputs from './PulseInputs';
import AnalogInputs from './AnalogInputs';
import { useTranslation } from 'react-i18next';
import IndexConfiguration from '../IndexConfiguration';
import GUISelect from '@/components/ui/Select/GUISelect';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigInputInner = () => {
	const { t } = useTranslation();
	const { support } = useConfigurationWrapperContext();

	const [inputsType, setInputsType] = useState<string>('mode-analog-inputs');

	return (
		<>
			<div className="flex flex-col sm:flex-row gap-3 sm:items-center">
				<p className="whitespace-nowrap">{t('label.mode-type-configurable')}</p>
				<GUISelect onChange={(e) => setInputsType(e.target.value)}>
					<option value="mode-analog-inputs">{t('label.mode-analog-inputs')}</option>
					{support.InputPulse && <option value="mode-pulse-inputs">{t('label.mode-pulse-inputs')}</option>}
				</GUISelect>
			</div>

			{inputsType === 'mode-analog-inputs' ? <AnalogInputs /> : <IndexConfiguration title={t('label.mode-pulse-inputs')} content={<PulseInputs />} />}
		</>
	);
};

const ConfigInput = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigInputInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigInput;
