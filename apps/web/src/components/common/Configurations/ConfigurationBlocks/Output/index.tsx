import Outputs from './Outputs';
import { useTranslation } from 'react-i18next';
import IndexConfiguration from '../IndexConfiguration';
import { ConfigurationWrapperProvider } from '@/context/useConfigurationWrapperContext';

const ConfigOutputInner = () => {
	const { t } = useTranslation();

	return <IndexConfiguration title={t('label.outputs')} content={<Outputs />} maxWidth="30rem" />;
};

const ConfigOutput = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigOutputInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigOutput;
