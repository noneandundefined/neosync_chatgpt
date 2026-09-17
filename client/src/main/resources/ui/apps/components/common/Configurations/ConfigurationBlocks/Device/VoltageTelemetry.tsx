import { useTranslation } from 'react-i18next';
import { GUInput } from '@/components/ui/Input/GUInput';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationContext } from '@/context/useConfigurationContext';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';

const formatMillivolts = (value?: number) => (value === undefined || value === null ? '' : String(value));

const VoltageTelemetry = () => {
	const { t } = useTranslation();

	const { imei } = useConfigurationContext();
	const { data: telemetry } = useConfigurationTelemetry(imei);

	return (
		<div className="space-y-3">
			<LabeledField label="label.nutrition">
				<div className="flex items-center gap-2">
					<GUInput type="text" value={formatMillivolts(telemetry?.v_power)} disabled />
					<p className="text-[#666] text-[13px]">{t('label.mV')}</p>
				</div>
			</LabeledField>

			<LabeledField label="label.battery">
				<div className="flex items-center gap-2">
					<GUInput type="text" value={formatMillivolts(telemetry?.v_battery)} disabled />
					<p className="text-[#666] text-[13px]">{t('label.mV')}</p>
				</div>
			</LabeledField>
		</div>
	);
};

export default VoltageTelemetry;
