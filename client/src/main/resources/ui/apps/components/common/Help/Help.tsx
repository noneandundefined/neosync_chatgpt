import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { HELP_SECTIONs } from '@/constants/Help.constant';
import { useModalContext } from '@/context/useModalContext';
import ChevronRight from '@/components/@icons/chevron-right';
import ModalHelpHowAddDevice from '@/components/Modal/ModalHelpHowAddDevice';
import BookOpenBlankVariantOutline from '@/components/@icons/book-open-blank-variant-outline';

interface HelpProps {
	changeSection: (id: string) => void;
}

const Help: React.FC<HelpProps> = ({ changeSection }) => {
	const { t } = useTranslation();
	const { open } = useModalContext();

	return (
		<div className="space-y-8">
			<div className="space-y-1">
				<p className="font-medium text-[18px] text-[#49525f]">{t('message.help-welcome-title')}</p>
				<p className="text-[#49525f]">{t('message.help-welcome-description')}</p>
			</div>

			<div className="flex flex-col lg:flex-row items-stretch gap-5">
				<div className="flex-1 divide-y divide-[#d1d5db]">
					{HELP_SECTIONs.map((section) => {
						const Icon = section.icon;

						return (
							<div key={section.id} onClick={() => changeSection(section.id)} className="flex items-center p-4 min-h-[100px] w-full hover:bg-white cursor-pointer">
								<div className="flex flex-1 items-center justify-between">
									<div className="flex flex-col md:flex-row md:items-center gap-5">
										<div>
											<Icon fill="#49525f" size={45} />
										</div>
										<div className="space-y-2 max-w-[40rem]">
											<p className="font-medium text-[16px]">{t(`message.${section.titleKey}`)}</p>
											<p className="text-[14px] text-[#49525f]">{t(`message.${section.descriptionKey}`)}</p>
										</div>
									</div>

									<div>
										<ChevronRight fill="#49525f" />
									</div>
								</div>
							</div>
						);
					})}
				</div>

				<div className="bg-[#f5faf6] lg:max-w-[25rem] p-5 border border-[#def2e4] rounded-[10px] space-y-8">
					<p className="font-medium text-[16px] text-[#1e772e]">{t('message.help-useful-info')}</p>

					<div className="space-y-8">
						<div
							className="flex items-start gap-5 cursor-pointer"
							onClick={() => {
								open(
									<Modal title={t('message.device-add-instruction-title')}>
										<ModalHelpHowAddDevice />
									</Modal>
								);
							}}
						>
							<div className="bg-[#e0f3e5] p-2 rounded-[10px]">
								<BookOpenBlankVariantOutline fill="#1e772e" size={31} />
							</div>

							<div>
								<p className="text-black font-medium text-[16px]">{t('message.how-create-device')}</p>
								<p className="text-[14px] font-normal text-[#49525f]">{t('message.instructions-create-device-help-desc')}</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default Help;
