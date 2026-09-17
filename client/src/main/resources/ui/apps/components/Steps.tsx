import '@/utils/i18n';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

interface StepsProps {
	steps: any[];
	currentStep: number;
	onClick: any;
}

const Steps: React.FC<StepsProps> = ({ steps, currentStep, onClick }) => {
	const { t } = useTranslation();

	const [hoveredStep, setHoveredStep] = useState<number | null>(null);

	return (
		<nav className="my-3">
			<ul className="flex items-center gap-6 overflow-x-auto">
				{steps.map((step, index) => {
					const isActive = currentStep == step.step;
					const isHovered = hoveredStep === step.step;

					return (
						<li key={index} className="cursor-pointer" onClick={() => onClick(step.step)} onMouseEnter={() => setHoveredStep(step.step)} onMouseLeave={() => setHoveredStep(null)}>
							<p className={`whitespace-nowrap font-medium text-[15px] ${isActive || isHovered ? 'text-[#49525f]' : 'text-[#9198a3]'}`}>{t(step.title)}</p>

							<div className={`w-full h-[2px] rounded-full mt-2 ${isActive || isHovered ? 'bg-[#49525f]' : 'bg-[#E5E8EB]'}`} />
						</li>
					);
				})}
			</ul>
		</nav>
	);
};

export default Steps;
