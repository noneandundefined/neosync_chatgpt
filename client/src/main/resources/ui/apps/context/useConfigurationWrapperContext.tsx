import { useTranslation } from 'react-i18next';
import { createContext, useContext } from 'react';
import LoaderLine from '@/components/common/Loader/LoaderLine';
import { useTabParam } from '@/hooks/Configuration/useTabParam';
import { SUPPORT_CONFIG } from '@/constants/DeviceModels.constant';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';
import { CONFIGURATION_DATA } from '@/constants/ConfigurationData.constant';
import { useOptionalConfigurationContext } from './useConfigurationContext';
import useConfigurationManager from '@/hooks/Configuration/useConfigurationManager';
import { useTemplateConfigurationContext } from './useTemplateConfigurationContext';
import useTemplateConfigurationManager from '@/hooks/Configuration/useTemplateConfigurationManager';

const ConfigurationWrapperContext = createContext<any>(null);
export const useConfigurationWrapperContext = () => useContext(ConfigurationWrapperContext);
export { ConfigurationWrapperContext };

export const ConfigurationWrapperProvider = ({ children }: { children: React.ReactNode }) => {
	const { t } = useTranslation();
	const templateCtx = useTemplateConfigurationContext();

	const deviceCtx = useOptionalConfigurationContext();
	const imei = deviceCtx?.imei ?? '';

	const { currentStep: section } = useTabParam();

	const isTemplate = !!templateCtx?.isTemplate;
	const templateId = templateCtx?.templateId;

	const deviceManager = useConfigurationManager(imei, section);
	const templateManager = useTemplateConfigurationManager(section, isTemplate, templateId);

	const { manager, loading, error, sseMessage, sseConnectionError } = isTemplate ? templateManager : deviceManager;

	if (loading) {
		const title = sseConnectionError ? t('message.sse-connection-error-title') : CONFIGURATION_DATA.GET.TITLE;

		return <LoaderLine title={title} description={sseMessage} />;
	}

	if (error || sseConnectionError) {
		const displayError = error || sseConnectionError;
		return <ErrorLoaderLine imei={isTemplate ? undefined : imei} title={t('message.error-get-configuration')} description={t(displayError || 'message.error-get-configuration')} />;
	}

	if (!manager) {
		return <LoaderLine title={t(CONFIGURATION_DATA.GET.TITLE)} description={t('message.initial-manager')} />;
	}

	const base = manager.getBase();
	const support = SUPPORT_CONFIG(base);
	const model = deviceCtx?.model;

	const value = isTemplate
		? {
				imei: 'template',
				section,
				support,
				manager,
				isTemplate: true,
				templateId,
			}
		: { imei, section, support, manager, model };

	return <ConfigurationWrapperContext.Provider value={value}>{children}</ConfigurationWrapperContext.Provider>;
};
