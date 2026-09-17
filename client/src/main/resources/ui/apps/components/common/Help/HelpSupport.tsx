import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import FaceAgent from '@/components/@icons/face-agent';
import OpenInNew from '@/components/@icons/open-in-new';
import ModalSupport from '@/components/Modal/ModalSupport';
import { useModalContext } from '@/context/useModalContext';

const HelpSupport = () => {
	const { t } = useTranslation();

	const { open } = useModalContext();

	return (
		<div className="bg-[#eff5fc] rounded-[6px] p-4 border border-[#e3eefa]">
			<div className="flex items-center justify-between">
				<div className="flex flex-col md:flex-row md:items-center gap-5">
					<FaceAgent fill="#395d95" size={35} />

					<div>
						<p className="font-medium text-[14px]">{t('message.help-support-title')}</p>
						<p className="text-[14px]">{t('message.help-support-description')}</p>
					</div>
				</div>

				<div
					className="bg-[#395d95] px-4 py-2 rounded-[6px] flex items-center gap-3 cursor-pointer"
					onClick={() =>
						open(
							<Modal title={t('message.help-support-contact')}>
								<ModalSupport />
							</Modal>
						)
					}
				>
					<p className="hidden md:block text-[14px] text-white">{t('message.help-support-contact')}</p>
					<OpenInNew fill="#fff" size={15} />
				</div>
			</div>
		</div>
	);
};

export default HelpSupport;
