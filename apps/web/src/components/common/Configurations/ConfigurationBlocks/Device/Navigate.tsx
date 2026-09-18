import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationContext } from '@/context/useConfigurationContext';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';

// const formatCoord = (value?: number) => (value === undefined || value === null ? '' : value.toFixed(6));

const Navigate = () => {
	const { imei } = useConfigurationContext();
	const { data: telemetry } = useConfigurationTelemetry(imei);

	return (
		<div className="space-y-3">
			<LabeledField label="label.gps-satellites" className="flex flex-col gap-1">
				<input type="text" value={telemetry?.gps_satellites ?? ''} disabled />
			</LabeledField>

			<LabeledField label="label.glonass-satellites" className="flex flex-col gap-1">
				<input type="text" value={telemetry?.glonass_satellites ?? ''} disabled />
			</LabeledField>

			{/* <LabeledField label="label.latitude" className="flex flex-col gap-1">
				<input type="text" value={formatCoord(telemetry?.lat)} disabled />
			</LabeledField>

			<LabeledField label="label.longitude" className="flex flex-col gap-1">
				<input type="text" value={formatCoord(telemetry?.lon)} disabled />
			</LabeledField> */}
		</div>
	);
};

export default Navigate;
