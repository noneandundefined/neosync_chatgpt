import Header from '@/components/Header/Header';
import Steps from '@/components/Steps';
import { getStepsByModel } from '@/utils/StepsByModelUtils';

interface ConfigurationLayoutProps {
	currentStep: number;
	steps: ReturnType<typeof getStepsByModel>;
	onStepClick: (index: number) => void;
	children: React.ReactNode;
	footer: React.ReactNode;
}

const ConfigurationLayout: React.FC<ConfigurationLayoutProps> = ({ currentStep, steps, onStepClick, children, footer }) => {
	return (
		<div className="min-h-screen flex-col md:flex-row">
			<Header show={true} arrow={true} />

			<div className="flex flex-1 flex-col items-center px-2 sm:px-6 md:px-12">
				<div className="w-full max-w-full lg:max-w-[75vw] xl:max-w-[75vw] mt-4 transition-[padding] duration-300" style={{ paddingBottom: 'var(--bottom-panel-padding)' }}>
					<Steps steps={steps} currentStep={currentStep} onClick={onStepClick} />

					<div className="max-w-full">{children}</div>
				</div>

				{footer}
			</div>
		</div>
	);
};

export default ConfigurationLayout;
