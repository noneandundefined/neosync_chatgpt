import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import { DEVICE_MODELS, getStaticSupport } from '@/constants/DeviceModels.constant';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const Outputs = () => {
	const { t } = useTranslation();
	const { model, isTemplate } = useConfigurationWrapperContext();

	const outputCount = getStaticSupport((model ?? DEVICE_MODELS.ADM333V2) as keyof typeof DEVICE_MODELS, isTemplate).OutputCount;

	return (
		<div className="space-y-3">
			{Array.from({ length: outputCount }, (_, index) => (
				<div className="flex items-center gap-3" key={index}>
					<p className="font-medium">
						{t('label.output')} {index}
					</p>
					<CFGSwitch uid={UIDs.OUTPUT} bit={index} />
				</div>
			))}
		</div>
	);
};

export default Outputs;
