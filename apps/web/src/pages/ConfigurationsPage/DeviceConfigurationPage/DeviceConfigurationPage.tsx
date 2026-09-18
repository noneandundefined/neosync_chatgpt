import { toast } from 'react-toastify';
import { useTranslation } from 'react-i18next';
import { useDecryptedImei } from '@/hooks/useDecryptedImei';
import React, { useCallback, useEffect, useRef, useState } from 'react';

import LoaderLine from '@/components/common/Loader/LoaderLine';
import { CONFIGURATION_DATA } from '@/constants/ConfigurationData.constant';

import { STEP_KEYS } from '@/constants/StepKeys.constant';
import { getStepsByModel } from '@/utils/StepsByModelUtils';
import { useTabParam } from '@/hooks/Configuration/useTabParam';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import DeviceConfigurationFooter from './DeviceConfigurationFooter';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';
import { useConfigurationActions } from '@/hooks/Configuration/useConfigurationActions';
import ConfigurationLayout from '@/components/common/Configurations/ConfigurationLayout';
import ConfigurationStepper from '@/components/common/Configurations/ConfigurationStepper';
import { basicConfigurationParsed, basicConfigurationResetDraft } from '@/rest/configurationAPI';
import { ConfigurationProvider, useConfigurationContext } from '@/context/useConfigurationContext';

const ConfigurationPendingApplyToast = () => {
	const { t } = useTranslation();
	const { device, imei } = useConfigurationContext();
	const warnedImeiRef = useRef<string | null>(null);

	useEffect(() => {
		if (warnedImeiRef.current === imei) return;
		warnedImeiRef.current = imei;

		if (device.configuration_sync_status !== 'pending') return;

		toast.warning(t('message.configuration-pending-apply-warning'), {
			toastId: `configuration-pending-apply-${imei}`,
			autoClose: 8000,
		});
	}, [imei, device.configuration_sync_status, t]);

	return null;
};

const DeviceConfigurationPage = () => {
	const { t } = useTranslation();
	const imei = useDecryptedImei();

	const { currentStep, setCurrentStep } = useTabParam();
	const [refreshVersion, setRefreshVersion] = useState(0);

	const stepIndex = STEP_KEYS.indexOf(currentStep);

	const {
		handleSentConfiguration,
		handleExportConfiguration,
		handleImportConfiguration,
		handleRebootDevice,
		handleRebootTelemetryDevice,
		handleEraseEepromDevice,
		handleEraseFlashDevice,
		loadingSet,
		loadingReboot,
		loadingImport,
	} = useConfigurationActions(imei ?? '');

	if (!imei) {
		return <ErrorLoaderLine title={t('message.error-get-configuration')} description={t('message.error-get-imei')} />;
	}

	const fetchConfigurationParsed = useCallback(() => basicConfigurationParsed(imei), [imei]);
	const { data: configuration, loading } = useHandleServer(['basicConfigurationParsedResp', imei], fetchConfigurationParsed);
	const reloadCurrentSection = useCallback(() => {
		setRefreshVersion((prev) => prev + 1);
	}, []);

	if (loading || !configuration) {
		return <LoaderLine title={CONFIGURATION_DATA.GET.TITLE} description={t('message.connect-to-server')} />;
	}

	return (
		<React.Fragment>
			{loadingSet && <LoaderLine title={CONFIGURATION_DATA.SET.TITLE} description={CONFIGURATION_DATA.SET.DESCRIPTION} />}
			{loadingReboot && <LoaderLine title={CONFIGURATION_DATA.REBOOT.TITLE} description={CONFIGURATION_DATA.REBOOT.DESCRIPTION} />}
			{loadingImport && <LoaderLine title={CONFIGURATION_DATA.IMPORT.TITLE} description={t('message.connect-to-server')} />}

			<ConfigurationProvider>
				<ConfigurationPendingApplyToast />

				<ConfigurationLayout
					currentStep={stepIndex}
					steps={getStepsByModel(configuration)}
					onStepClick={(index) => setCurrentStep(STEP_KEYS[index])}
					footer={
						<DeviceConfigurationFooter
							imei={imei}
							onImport={async (file: File) => {
								await handleImportConfiguration(file);
								reloadCurrentSection();
							}}
							onExport={handleExportConfiguration}
							onApply={handleSentConfiguration}
							onReboot={handleRebootDevice}
							onRebootTelemetry={handleRebootTelemetryDevice}
							onRefresh={async () => {
								await basicConfigurationResetDraft(imei);
								reloadCurrentSection();
							}}
							onReloadSection={reloadCurrentSection}
							onEraseEeprom={handleEraseEepromDevice}
							onEraseFlash={handleEraseFlashDevice}
						/>
					}
				>
					<ConfigurationStepper key={`${currentStep}-${refreshVersion}`} currentStep={stepIndex} />
				</ConfigurationLayout>
			</ConfigurationProvider>
		</React.Fragment>
	);
};

export default DeviceConfigurationPage;
