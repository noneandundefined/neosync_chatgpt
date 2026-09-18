import { ROUTES } from '@/constants/constants';
import { useTabParam } from '@/hooks/Configuration/useTabParam';
import { basicConfigurationExport, basicConfigurationResetDraft } from '@/rest/configurationAPI';
import { useTranslation } from 'react-i18next';

interface ErrorLoaderLineProps {
	imei?: string;
	title: string;
	description: string;
}

const ErrorLoaderLine: React.FC<ErrorLoaderLineProps> = ({ imei = '', title, description }) => {
	const { t } = useTranslation();
	const { currentStep, setCurrentStep } = useTabParam();

	return (
		<div className="fixed inset-0 h-screen w-screen bg-[#00000070] z-[1000] flex flex-col items-center justify-center">
			<div className="bg-[#fff] w-screen h-auto pt-[3rem] pb-0 flex flex-col justify-center" style={{ boxShadow: '0 0 20px 0 rgba(0, 0, 0, 0.15' }}>
				<div className="w-full text-center sm:text-left sm:ml-[25%] px-3">
					<p className="font-medium text-[1.45rem] text-[#921414]">{title}</p>
					<p className="text-[14px] text-[#555] my-3">{description}</p>

					<div className="space-y-2 mt-6">
						{imei && (
							<>
								<p
									className="text-[#1976d2] hover:text-[#00479b] font-medium text-md md:text-[12px] cursor-pointer"
									onClick={async () => {
										await basicConfigurationExport(imei);
									}}
								>
									{t('message.save-settings-file')}
								</p>
								<p
									className="text-[#1976d2] hover:text-[#00479b] font-medium text-md md:text-[12px] cursor-pointer"
									onClick={async () => {
										await basicConfigurationResetDraft(imei);
										if (currentStep == 'device') {
											setCurrentStep('sim');
										} else {
											setCurrentStep('device');
										}
									}}
								>
									{t('message.repeat-get-configuration')}
								</p>
							</>
						)}
						<p className="text-[#1976d2] hover:text-[#00479b] font-medium text-md md:text-[12px] cursor-pointer" onClick={() => (window.location.href = ROUTES.HOME)}>
							{t('message.back-to-main-page')}
						</p>
					</div>
				</div>
				<div id="loader-line-error" className="mt-[7rem]"></div>
			</div>
		</div>
	);
};

export default ErrorLoaderLine;
