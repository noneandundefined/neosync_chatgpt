import { useTranslation } from 'react-i18next';
import { STEP_KEYS } from '@/constants/StepKeys.constant';
import LoaderLine from '@/components/common/Loader/LoaderLine';
import { useTabParam } from '@/hooks/Configuration/useTabParam';
import { STEPS } from '@/constants/ConfigurationSteps.constant';
import ConfigurationTemplateFooter from './ConfigurationTemplateFooter';
import { CONFIGURATION_DATA } from '@/constants/ConfigurationData.constant';
import useTemplateExitGuard from '@/hooks/Configuration/useTemplateExitGuard';
import ConfigurationLayout from '@/components/common/Configurations/ConfigurationLayout';
import { TemplateConfigurationProvider } from '@/context/useTemplateConfigurationContext';
import ConfigurationStepper from '@/components/common/Configurations/ConfigurationStepper';
import useTemplateConfigurationManager from '@/hooks/Configuration/useTemplateConfigurationManager';

const ConfigurationTemplateCreatePage = () => {
	const { t } = useTranslation();

	useTemplateExitGuard();

	const { currentStep, setCurrentStep } = useTabParam();

	const stepIndex = STEP_KEYS.indexOf(currentStep);

	const { manager: deviceSectionManager, loading: loadingSection } = useTemplateConfigurationManager('device');

	if (loadingSection || !deviceSectionManager) {
		return <LoaderLine title={CONFIGURATION_DATA.GET.TITLE} description={t('message.connect-to-server')} />;
	}

	return (
		<TemplateConfigurationProvider>
			<ConfigurationLayout currentStep={stepIndex} steps={STEPS} onStepClick={(index) => setCurrentStep(STEP_KEYS[index])} footer={<ConfigurationTemplateFooter />}>
				<ConfigurationStepper currentStep={stepIndex} />
			</ConfigurationLayout>
		</TemplateConfigurationProvider>
	);
};

export default ConfigurationTemplateCreatePage;
