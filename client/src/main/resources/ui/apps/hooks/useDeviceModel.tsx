import { DEVICE_MODELS } from '@/constants/DeviceModels.constant';

const DEFAULT_MODEL = 'ADM007BLE' as keyof typeof DEVICE_MODELS;

type UseDeviceModelResult = {
	model: keyof typeof DEVICE_MODELS;
	valid: boolean;
};

export const useDeviceModel = (model?: string): UseDeviceModelResult => {
	if (!model) {
		return { model: DEFAULT_MODEL, valid: false };
	}

	const normalized = model.trim().toUpperCase();

	const deviceModels = Object.keys(DEVICE_MODELS) as Array<keyof typeof DEVICE_MODELS>;

	const exactMatch = deviceModels.find((deviceModel) => normalized === deviceModel);
	if (exactMatch) {
		return {
			model: exactMatch,
			valid: true,
		};
	}

	const matched = [...deviceModels].sort((a, b) => b.length - a.length).find((deviceModel) => normalized.includes(deviceModel));

	if (matched) {
		return {
			model: matched,
			valid: true,
		};
	}

	return {
		model: DEFAULT_MODEL,
		valid: false,
	};
};
