import { useTranslation } from 'react-i18next';

const ModalHelpHowAddDevice = () => {
	const { t } = useTranslation();

	const steps = [
		t('message.device-add-instruction-step-1'),
		t('message.device-add-instruction-step-2'),
		t('message.device-add-instruction-step-3'),
		t('message.device-add-instruction-step-4'),
		t('message.device-add-instruction-step-5'),
	];

	return (
		<div className="max-h-[60vh] overflow-y-auto">
			<ol className="space-y-2 list-none">
				{steps.map((step, index) => (
					<li key={index} className="flex flex-col md:flex-row items-start gap-4 p-2">
						<div className="flex-shrink-0 w-7 h-7 flex items-center justify-center rounded-full bg-[#49525f] text-[14px] text-white font-medium">{index + 1}</div>
						<p className="leading-relaxed text-black">{step}</p>
					</li>
				))}
			</ol>
		</div>
	);
};

export default ModalHelpHowAddDevice;
