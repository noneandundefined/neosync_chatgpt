import { useTranslation } from 'react-i18next';
import ArrowRight from '../@icons/arrow-right';
import ContentCopy from '../@icons/content-copy';
import EmailOutline from '../@icons/email-outline';
import { copyToBuffer } from '@/utils/CopyBufferUtils';

const SOCIAL_MEDIA_RU = [
	{
		ico: '/local/templates/social-media/tg-ico.png',
		title: 'Telegram',
		description: 'message.support-telegram-description',
		href: 'https://t.me/Neomatica_Support',
	},
	{
		ico: '/local/templates/social-media/wa-ico.png',
		title: 'WhatsApp',
		description: 'message.support-whatsapp-description',
		href: 'https://wa.me/89068770418',
	},
	{
		ico: '/local/templates/social-media/max-ico.png',
		title: 'MAX',
		description: 'message.support-max-description',
		href: 'https://max.ru/u/f9LHodD0cOL6QV4-PmHr3LD-YxyOzzsE8UaGewxT2_pw1ItMRIE3lEjzdpU',
	},
];

const SOCIAL_MEDIA_EU = [
	{
		ico: '/local/templates/social-media/tg-ico.png',
		title: 'Telegram',
		description: 'message.support-telegram-description',
		href: 'https://t.me/Neomatica_Support',
	},
	{
		ico: '/local/templates/social-media/wa-ico.png',
		title: 'WhatsApp',
		description: 'message.support-whatsapp-description',
		href: 'https://wa.me/89068770418 ',
	},
];

const ModalSupport = () => {
	const { t, i18n } = useTranslation();

	const socialMediaArray = i18n.language === 'ru' ? SOCIAL_MEDIA_RU : SOCIAL_MEDIA_EU;

	return (
		<div className="space-y-3">
			<div className="bg-[#fafafc] border rounded-[12px] p-2 flex items-center justify-between">
				<div className="flex items-center gap-5">
					<div className="bg-[#eff5fc] p-2 rounded-full">
						<EmailOutline fill="#395d95" size={25} />
					</div>

					<div className="space-y-1">
						<p className="text-[#49525f] font-medium">{t('label.email')}</p>
						<p className="text-[#395d95] font-medium">{i18n.language === 'ru' ? 'support@neomatica.ru' : 'support@neomatica.com'}</p>
					</div>
				</div>

				<div className="bg-white p-2 rounded-full cursor-pointer" onClick={() => copyToBuffer(i18n.language === 'ru' ? 'support@neomatica.ru' : 'support@neomatica.com')}>
					<ContentCopy fill="#49525f" size={21} />
				</div>
			</div>

			<div className="space-y-3">
				<p className="font-medium text-[#333]">{t('label.messengers')}</p>

				<div className="flex flex-wrap gap-3">
					{socialMediaArray.map((media, index) => (
						<div
							key={index}
							className="bg-[#fafafc] flex-1 border rounded-[12px] p-2 flex items-center justify-between min-w-[15rem] cursor-pointer"
							onClick={() => window.open(media.href, '_blank', 'noopener,noreferrer')}
						>
							<div className="flex items-center gap-3">
								<img src={media.ico} alt={media.title} draggable={false} width={40} />

								<div className="space-y-1">
									<p className="text-[#49525f] font-medium">{t(media.title)}</p>
									<p className="text-[14px] text-[#49525f]">{t(media.description)}</p>
								</div>
							</div>

							<div className="cursor-pointer">
								<ArrowRight fill="#49525f" />
							</div>
						</div>
					))}
				</div>
			</div>
		</div>
	);
};

export default ModalSupport;
