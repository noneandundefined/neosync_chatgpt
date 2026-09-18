import { basicConfigurationApply, basicConfigurationExport, basicConfigurationImport, basicRebootConfigurationTelemetry } from '@/rest/configurationAPI';
import { basicDeviceCommandEraseEeprom, basicDeviceCommandEraseFlash, basicDeviceCommandReboot } from '@/rest/deviceCommandAPI';
import { useState } from 'react';

export const useConfigurationActions = (imei: string) => {
	const [loadingSet, setLoadingSet] = useState(false);
	const [loadingImport, setLoadingImport] = useState(false);
	const [loadingReboot, setLoadingReboot] = useState(false);
	const [loadingEraseFlash, setLoadingEraseFlash] = useState(false);
	const [loadingEraseEeprom, setLoadingEraseEeprom] = useState(false);
	const [loadingRebootTelemetry, setLoadingRebootTelemetry] = useState(false);

	const handleSentConfiguration = async () => {
		try {
			setLoadingSet(true);
			await basicConfigurationApply(imei);
		} finally {
			setLoadingSet(false);
		}
	};

	const handleExportConfiguration = async () => {
		await basicConfigurationExport(imei);
	};

	const handleImportConfiguration = async (file: File) => {
		try {
			setLoadingImport(true);
			await basicConfigurationImport(imei, file);
		} finally {
			setLoadingImport(false);
		}
	};

	const handleRebootTelemetryDevice = async () => {
		try {
			setLoadingRebootTelemetry(true);
			await basicRebootConfigurationTelemetry(imei);
		} finally {
			setLoadingRebootTelemetry(false);
		}
	};

	const handleRebootDevice = async () => {
		try {
			setLoadingReboot(true);
			await basicDeviceCommandReboot(imei);
		} finally {
			setLoadingReboot(false);
		}
	};

	const handleEraseEepromDevice = async () => {
		try {
			setLoadingEraseEeprom(true);
			await basicDeviceCommandEraseEeprom(imei);
		} finally {
			setLoadingEraseEeprom(false);
		}
	};

	const handleEraseFlashDevice = async () => {
		try {
			setLoadingEraseFlash(true);
			await basicDeviceCommandEraseFlash(imei);
		} finally {
			setLoadingEraseFlash(false);
		}
	};

	return {
		handleSentConfiguration,
		handleExportConfiguration,
		handleImportConfiguration,
		handleRebootDevice,
		handleRebootTelemetryDevice,
		handleEraseEepromDevice,
		handleEraseFlashDevice,
		loadingSet,
		loadingImport,
		loadingReboot,
		loadingRebootTelemetry,
		loadingEraseEeprom,
		loadingEraseFlash,
	};
};
