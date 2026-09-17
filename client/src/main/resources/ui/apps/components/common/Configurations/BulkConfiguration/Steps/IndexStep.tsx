import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import ChevronUp from '@/components/@icons/chevron-up';
import ChevronDown from '@/components/@icons/chevron-down';
import BulkConfigurationStepper from './BulkConfigurationStepper';

interface IndexStepProps {
	step: number;
	title: string;
	children: React.ReactNode;
	defaultExpanded?: boolean;
}

const IndexStep: React.FC<IndexStepProps> = ({ step, title, children, defaultExpanded }) => {
	const { t } = useTranslation();

	const [expanded, setExpanded] = useState(defaultExpanded ?? step === 1);

	return (
		<div className="flex gap-6 items-stretch">
			<BulkConfigurationStepper number={step} />

			<section className="flex-1 min-w-0">
				<button type="button" className="w-full flex items-center justify-between gap-3 text-left" onClick={() => setExpanded((prev) => !prev)}>
					<p className="font-medium text-[#49525f]">{t(title)}</p>

					{expanded ? <ChevronUp fill="#49525f" size={20} /> : <ChevronDown fill="#49525f" size={20} />}
				</button>

				{expanded && <div className="space-y-5 mt-5">{children}</div>}
			</section>
		</div>
	);
};

export default IndexStep;
