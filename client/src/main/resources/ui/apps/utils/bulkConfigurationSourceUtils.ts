import { TFunction } from 'i18next';
import { format } from 'date-fns';

import { DeviceDetailedStatusResponse } from '@/rest/deviceAPI';
import { ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';

export type BulkConfigurationSourceType = CompanyCreateRequest['configuration_cource'];

export type BulkDeviceSourceEligibility = { eligible: true } | { eligible: false; reasonKey: string };

export const getDominantFirmwareVersion = (devices: DeviceDetailedStatusResponse[]): number | null => {
	const counts = new Map<number, number>();

	for (const device of devices) {
		if (device.firmware_version == null || device.firmware_version === 0) {
			continue;
		}

		counts.set(device.firmware_version, (counts.get(device.firmware_version) ?? 0) + 1);
	}

	let dominant: number | null = null;
	let maxCount = 0;

	for (const [firmware, count] of counts) {
		if (count > maxCount || (count === maxCount && (dominant === null || firmware < dominant))) {
			maxCount = count;
			dominant = firmware;
		}
	}

	return dominant;
};

export const getBulkDeviceSourceEligibility = (device: DeviceDetailedStatusResponse, priorityModel: string, dominantFirmware: number | null): BulkDeviceSourceEligibility => {
	if (priorityModel && device.device_model !== priorityModel) {
		return { eligible: false, reasonKey: 'message.provisioning-device-source-reason-model' };
	}

	if (device.configuration_sync_status !== 'confirmed') {
		return { eligible: false, reasonKey: 'message.provisioning-device-source-reason-not-confirmed' };
	}

	if (dominantFirmware != null && dominantFirmware > 0) {
		if (device.firmware_version == null || device.firmware_version === 0) {
			return { eligible: false, reasonKey: 'message.provisioning-device-source-reason-firmware-unknown' };
		}

		if (device.firmware_version !== dominantFirmware) {
			return { eligible: false, reasonKey: 'message.provisioning-device-source-reason-firmware' };
		}
	}

	return { eligible: true };
};

export const isBulkTemplateCompatible = (template: ConfigurationTemplateResponse, priorityModel: string): boolean => {
	if (!priorityModel || !template.model) {
		return true;
	}

	return template.model === priorityModel;
};

export const getBulkDeviceSourceEligibilityLabel = (t: TFunction, eligibility: BulkDeviceSourceEligibility): string => {
	if (eligibility.eligible) {
		return '';
	}

	return t(eligibility.reasonKey);
};

export const getBulkConfigurationSourceTypeLabel = (source: BulkConfigurationSourceType, t: TFunction): string => {
	switch (source) {
		case 'template_sources':
			return t('message.provisioning-source-type-template');
		case 'file_sources':
			return t('message.provisioning-source-type-file');
		case 'device_sources':
			return t('message.provisioning-source-type-device');
		default:
			return source;
	}
};

export const buildBulkProvisioningDefaultName = (
	t: TFunction,
	params: {
		model: string;
		source: BulkConfigurationSourceType;
		sourceName?: string | null;
		selectedCount: number;
		createdAt: Date;
	}
): string => {
	const model = params.model.trim() || t('message.provisioning-default-name-unknown-model');
	const sourceType = getBulkConfigurationSourceTypeLabel(params.source, t);
	const sourceName = params.sourceName?.trim();
	const nameParams = { model, sourceType, sourceName, selectedCount: params.selectedCount, date: format(params.createdAt, 'dd.MM.yyyy HH:mm') };

	if (!sourceName) {
		return t('message.provisioning-default-name-partial', nameParams);
	}

	return t('message.provisioning-default-name', nameParams);
};
