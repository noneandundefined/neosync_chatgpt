import { useTranslation } from 'react-i18next';
import { STEP_KEYS } from '@/constants/StepKeys.constant';
import LoaderLine from '@/components/common/Loader/LoaderLine';
import { useTabParam } from '@/hooks/Configuration/useTabParam';
import { STEPS } from '@/constants/ConfigurationSteps.constant';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import ConfigurationTemplateFooter from './ConfigurationTemplateFooter';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';
import { CONFIGURATION_DATA } from '@/constants/ConfigurationData.constant';
import useTemplateExitGuard from '@/hooks/Configuration/useTemplateExitGuard';
import ConfigurationLayout from '@/components/common/Configurations/ConfigurationLayout';
import { TemplateConfigurationProvider } from '@/context/useTemplateConfigurationContext';
import ConfigurationStepper from '@/components/common/Configurations/ConfigurationStepper';
import useTemplateConfigurationManager from '@/hooks/Configuration/useTemplateConfigurationManager';
import { basicConfigurationTemplateGetById, ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';

interface ConfigurationTemplateEditorPageProps {
	templateId?: number;
}

const ConfigurationTemplateEditorPage = ({ templateId }: ConfigurationTemplateEditorPageProps) => {
	const { t } = useTranslation();

	useTemplateExitGuard(templateId);

	const { currentStep, setCurrentStep } = useTabParam();
	const stepIndex = STEP_KEYS.indexOf(currentStep);

	const { data: template, loading: loadingTemplate } = useHandleServer<ConfigurationTemplateResponse | null>(['configurationTemplate', templateId], () => basicConfigurationTemplateGetById(templateId!), {
		enabled: !!templateId,
	});

	const { manager: deviceSectionManager, loading: loadingSection } = useTemplateConfigurationManager('device', true, templateId);

	if (templateId && loadingTemplate) {
		return <LoaderLine title={CONFIGURATION_DATA.GET.TITLE} description={t('message.connect-to-server')} />;
	}

	if (templateId && !template) {
		return <ErrorLoaderLine title={t('message.error-get-configuration')} description={t('message.configuration-template-not-found')} />;
	}

	if (loadingSection || !deviceSectionManager) {
		return <LoaderLine title={CONFIGURATION_DATA.GET.TITLE} description={t('message.connect-to-server')} />;
	}

	return (
		<TemplateConfigurationProvider templateId={templateId} template={template ?? null}>
			<ConfigurationLayout currentStep={stepIndex} steps={STEPS} onStepClick={(index) => setCurrentStep(STEP_KEYS[index])} footer={<ConfigurationTemplateFooter />}>
				<ConfigurationStepper currentStep={stepIndex} />
			</ConfigurationLayout>
		</TemplateConfigurationProvider>
	);
};

export default ConfigurationTemplateEditorPage;
