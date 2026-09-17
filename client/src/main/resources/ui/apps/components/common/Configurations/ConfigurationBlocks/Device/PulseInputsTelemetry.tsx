import { GUInput } from '@/components/ui/Input/GUInput';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';

const formatCounter = (value?: number) => (value === undefined || value === null ? '' : String(value));

const PulseInputsTelemetry = () => {
	const { imei } = useConfigurationWrapperContext();
	const { data: telemetry } = useConfigurationTelemetry(imei);

	const fields = [
		{ index: 0, value: telemetry?.in_d0 },
		{ index: 1, value: telemetry?.in_d1 },
	];

	return (
		<div className="space-y-3">
			{fields.map(({ index, value }) => (
				<LabeledField key={index} label={`D${index}`}>
					<GUInput type="text" value={formatCounter(value)} disabled />
				</LabeledField>
			))}
		</div>
	);
};

export default PulseInputsTelemetry;
