import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import ChevronLeft from '@/components/@icons/chevron-left';
import { HELP_USERs_CONTROLs, HELP_USERs_NOTE_KEY, HELP_USERs_STEPs } from '@/constants/Help.constant';

interface HelpUserProps {
	title: string;
	description: string;
	changeSection: () => void;
}

const HelpUser: React.FC<HelpUserProps> = ({ title, description, changeSection }) => {
	const { t } = useTranslation();

	return (
		<div className="space-y-5">
			<div className="flex items-center gap-3">
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

			<div className="flex flex-col lg:flex-row items-stretch gap-5">
				<div>
					<img src="/local/templates/neosync/help/gui-users.png" alt={t('message.help-users-image-alt')} draggable={false} />
				</div>

				<div className="bg-[#f5faf6] lg:max-w-[25rem] p-5 border border-[#def2e4] rounded-[10px] space-y-8">
					<p className="font-medium text-[16px] text-[#1e772e]">{t('message.help-how-it-works')}</p>

					<div className="space-y-8">
						<div className="space-y-8">
							{HELP_USERs_STEPs.map((step) => (
								<div key={step.id} className="flex items-start gap-5">
									<div className="flex items-center justify-center bg-[#1e772e] min-h-[20px] min-w-[20px] rounded-full">
										<p className="text-white text-sm font-medium">{step.id}</p>
									</div>

									<div>
										<p className="font-medium text-[16px]">{t(`message.${step.titleKey}`)}</p>
										<p className="text-[14px] text-[#49525f]">{t(`message.${step.descriptionKey}`)}</p>
									</div>
								</div>
							))}
						</div>
					</div>
				</div>
			</div>

			<div className="space-y-1">
				<p className="font-medium text-[16px] text-[#49525f]">{t('message.help-note-title')}</p>
				<p className="text-[#49525f]">{t(`message.${HELP_USERs_NOTE_KEY}`)}</p>
			</div>

			<div className="space-y-3">
				<p className="font-medium text-[16px] text-[#49525f]">{t('message.help-controls-title')}</p>

				<div className="bg-white border rounded-xl overflow-hidden">
					{HELP_USERs_CONTROLs.map((item, index) => {
						const Icon = item.icon;

						return (
							<div key={item.titleKey} className={`grid lg:grid-cols-[50px_280px_1fr] gap-3 px-2 py-3 lg:items-center ${index !== HELP_USERs_CONTROLs.length - 1 ? 'border-b' : ''}`}>
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

export default HelpUser;
