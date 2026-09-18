import { useTranslation } from 'react-i18next';
import { GUInput } from '@/components/ui/Input/GUInput';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';

const formatMillivolts = (value?: number) => (value === undefined || value === null ? '' : String(value));

const AnalogInputsTelemetry = () => {
	const { t } = useTranslation();

	const { imei } = useConfigurationWrapperContext();
	const { data: telemetry } = useConfigurationTelemetry(imei);

	const fields = [
		{ index: 0, value: telemetry?.in_a0 },
		{ index: 1, value: telemetry?.in_a1 },
	];

	return (
		<div className="space-y-3">
			{fields.map(({ index, value }) => (
				<LabeledField key={index} label={`A${index}`}>
					<div className="flex items-center gap-2">
						<GUInput type="text" value={formatMillivolts(value)} disabled />
						<p className="text-[#666] text-[13px]">{t('label.mV')}</p>
					</div>
				</LabeledField>
			))}
		</div>
	);
};

export default AnalogInputsTelemetry;
