import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import ChevronLeft from '@/components/@icons/chevron-left';
import { HELP_DEVICEs_CONTROLs } from '@/constants/Help.constant';

interface HelpDeviceProps {
	title: string;
	description: string;
	changeSection: () => void;
}

const HelpDevice: React.FC<HelpDeviceProps> = ({ title, description, changeSection }) => {
	const { t } = useTranslation();

	return (
		<div className="space-y-5">
			<div className="sticky top-0 z-10 bg-[#fafafa] flex items-center gap-3">
				<Tooltip title={t('label.back')}>
					<div onClick={changeSection} className="cursor-pointer p-1 rounded hover:bg-white">
						<ChevronLeft fill="#49525f" size={22} />
					</div>
				</Tooltip>

				<p className="text-[#555] font-medium text-[16px]">{t('message.help-breadcrumb')}</p>
			</div>

			<div className="space-y-1">
				<p className="font-medium text-[18px] text-[#49525f]">{title}</p>
				<p className="text-[#49525f]">{description}</p>
			</div>

			<div className="space-y-3">
				<p className="font-medium text-[16px] text-[#49525f]">{t('message.help-device-columns-title')}</p>

				<ul className="list-disc pl-6 space-y-2 text-[15px] text-[#333]">
					<li>
						<span className="font-semibold">{t('label.imei')}</span> — {t('message.help-device-column-imei-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.model')}</span> — {t('message.help-device-column-model-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.firmware')}</span> — {t('message.help-device-column-firmware-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.name')}</span> — {t('message.help-device-column-name-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.phone')}</span> — {t('message.help-device-column-phone-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.organization')}</span> — {t('message.help-device-column-organization-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('message.help-device-column-date-sync-label')}</span> — {t('message.help-device-column-date-sync-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.owner')}</span> — {t('message.help-device-column-owner-desc')}
					</li>

					<li>
						<span className="font-semibold">{t('label.connection')}</span> — {t('message.help-device-column-connection-desc')}
					</li>
				</ul>
			</div>

			<div className="space-y-3">
				<p className="font-medium text-[16px] text-[#49525f]">{t('message.help-controls-title')}</p>

				<div className="bg-white border rounded-xl overflow-hidden">
					{HELP_DEVICEs_CONTROLs.map((item, index) => {
						const Icon = item.icon;

						return (
							<div key={item.titleKey} className={`grid lg:grid-cols-[50px_280px_1fr] gap-3 px-2 py-3 lg:items-center ${index !== HELP_DEVICEs_CONTROLs.length - 1 ? 'border-b' : ''}`}>
								<div className="flex lg:justify-center">
									<div className="bg-[#f7f7f7] rounded-lg p-2">
										<Icon fill="#49525f" size={19} />
									</div>
								</div>

								<div className="font-medium text-[14px]">{t(`message.${item.titleKey}`)}</div>

								<div className="text-[#62666d] text-[14px] whitespace-pre-line">{t(`message.${item.descriptionKey}`)}</div>
							</div>
						);
					})}
				</div>
			</div>
		</div>
	);
};

export default HelpDevice;
