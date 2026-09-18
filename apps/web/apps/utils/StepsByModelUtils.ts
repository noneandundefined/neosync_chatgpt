import { STEPS } from '@/constants/ConfigurationSteps.constant';
import { SUPPORT_CONFIG } from '@/constants/DeviceModels.constant';
import { ConfigurationParsedResponse } from '@/rest/configurationAPI';

interface Step {
	title: string;
	step: number;
}

export const getStepsByModel = (configuration: ConfigurationParsedResponse | null): Step[] => {
	if (!configuration) {
		const disabled = ['label.bluetooth', 'label.one-wire'];
		return STEPS.filter((step) => !disabled.includes(step.title));
	}

	const support = SUPPORT_CONFIG(configuration);

	const supportMap: Record<string, boolean> = {
		'label.sim': support.Sim_Tab,
		'label.server': support.Server_Tab,
		'label.events': support.Event_Tab,
		'label.track': support.Track_Tab,
		'label.inputs': support.ModeHybrid ? false : support.Input_Tab,
		'label.outputs': support.Output_Tab,
		'label.rs-485': support.RS485_Tab,
		'label.bluetooth': support.Bluetooth_Tab,
		'label.one-wire': support.OneWire_Tab,
	};

	return STEPS.filter((step) => supportMap[step.title] ?? true);
};
