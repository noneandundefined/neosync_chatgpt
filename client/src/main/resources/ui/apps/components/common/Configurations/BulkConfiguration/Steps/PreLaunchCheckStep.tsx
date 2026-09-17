import IndexStep from './IndexStep';
import { useTranslation } from 'react-i18next';

import ClockOutline from '@/components/@icons/clock-outline';
import CheckCircleOutline from '@/components/@icons/check-circle-outline';
import CloseCircleOutline from '@/components/@icons/close-circle-outline';
import { DeviceStatBC } from '@/pages/ConfigurationsPage/BulkConfigurationPage/BulkConfigurationPage';

interface PreLaunchCheckStepProps {
	stats: DeviceStatBC;
}

const PreLaunchCheckStep: React.FC<PreLaunchCheckStepProps> = ({ stats }) => {
	const { t } = useTranslation();

	return (
		<IndexStep step={4} title="message.provisioning-pre-launch-check">
			<div className="w-full rounded-[6px] overflow-hidden">
				<div className="flex items-start gap-4 px-5 py-5 border-b border-[#E5E7EB]">
					<CheckCircleOutline fill="#1e772e" size={20} className="mt-1" />

					<div className="flex-1">
						<p className="font-medium text-[15px] text-[#222]">{t('message.provisioning-ready-for-setup')}</p>

						<p className="mt-1 text-sm text-[#666]">{t('message.provisioning-will-be-added-to-task')}</p>
					</div>

					<span className="font-semibold text-[#222]">{stats.ready}</span>
				</div>

				<div className="flex items-start gap-4 px-5 py-5 border-b border-[#E5E7EB]">
					<ClockOutline fill="#49525f" size={20} className="mt-1" />

					<div className="flex-1">
						<p className="font-medium text-[15px] text-[#222]">{t('message.provisioning-offline-can-queue')}</p>

						<p className="mt-1 text-sm text-[#666]">{t('message.provisioning-will-setup-on-next-connection')}</p>
					</div>

					<span className="font-semibold text-[#222]">{stats.offline}</span>
				</div>

				<div className="flex items-start gap-4 px-5 py-5">
					<CloseCircleOutline fill="#b50202" size={20} className="mt-1" />

					<div className="flex-1">
						<p className="font-medium text-[15px] text-[#222]">{t('message.provisioning-incompatible-model-firmware')}</p>

						<p className="mt-1 text-sm text-[#666]">{t('message.provisioning-will-be-excluded')}</p>
					</div>

					<span className="font-semibold text-[#222]">{stats.incompatible}</span>
				</div>
			</div>
		</IndexStep>
	);
};

export default PreLaunchCheckStep;
