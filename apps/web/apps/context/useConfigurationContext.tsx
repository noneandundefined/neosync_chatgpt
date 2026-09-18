import { useTranslation } from 'react-i18next';
import { useDeviceModel } from '@/hooks/useDeviceModel';
import { useDecryptedImei } from '@/hooks/useDecryptedImei';
import LoaderLine from '@/components/common/Loader/LoaderLine';
import { createContext, useCallback, useContext } from 'react';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { DEVICE_MODELS } from '@/constants/DeviceModels.constant';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';
import { CONFIGURATION_DATA } from '@/constants/ConfigurationData.constant';
import { basicDeviceGetByImei, Device_DeviceConf_Sync } from '@/rest/deviceAPI';

interface ConfigurationContext {
	model: keyof typeof DEVICE_MODELS;
	device: Device_DeviceConf_Sync;
	imei: string;
}

const ConfigurationContext = createContext<ConfigurationContext | null>(null);

export const useConfigurationContext = () => {
	const ctx = useContext(ConfigurationContext);
	if (!ctx) {
		throw new Error('useDeviceContext must be used inside DeviceProvider');
	}

	return ctx;
};

export const useOptionalConfigurationContext = () => useContext(ConfigurationContext);

export const ConfigurationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { t } = useTranslation();

	const imei = useDecryptedImei();

	if (!imei) {
		return <ErrorLoaderLine title={t('message.error-get-configuration')} description={t('message.error-get-imei')} />;
	}

	const fetchDevice = useCallback((signal?: AbortSignal) => basicDeviceGetByImei(imei, signal), [imei]);
	const { data: respDeviceGet, loading: loadingRespDeviceGet } = useHandleServer(['respDeviceGet', imei], fetchDevice, { staleTime: 0, refetchInterval: 5000 });

	if (loadingRespDeviceGet) {
		return <LoaderLine title={CONFIGURATION_DATA.GET.TITLE} description={t('message.connect-to-server')} />;
	}

	if (!respDeviceGet?.device_extended_model && !respDeviceGet?.device_model) {
		return <ErrorLoaderLine title={t('message.error-get-configuration')} description={t('message.error-get-device-model')} />;
	}

	const { model, valid } = useDeviceModel(respDeviceGet.device_extended_model ?? respDeviceGet.device_model ?? undefined);

	if (!valid) {
		return <ErrorLoaderLine title={t('message.error-get-configuration')} description={t('message.error-get-device-model')} />;
	}

	return <ConfigurationContext.Provider value={{ model, device: respDeviceGet, imei }}>{children}</ConfigurationContext.Provider>;
};
