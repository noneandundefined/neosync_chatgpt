import { useEffect, useState } from 'react';
import GUIButton from '../ui/Button/GUIButton';
import { useTranslation } from 'react-i18next';
import { initYandexMetrika } from '@/utils/yandexMetrika';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';

const ACCEPTED = 'accepted';
const REJECTED = 'rejected';

const CookieConsentBanner = () => {
	const { t } = useTranslation();
	const [visible, setVisible] = useState(() => !localStorage.getItem(CACHEKEYs.COOKIE_CONSENT));

	useEffect(() => {
		if (localStorage.getItem(CACHEKEYs.COOKIE_CONSENT) === ACCEPTED) {
			initYandexMetrika();
		}
	}, []);

	if (!visible) {
		return null;
	}

	const acceptAll = () => {
		localStorage.setItem(CACHEKEYs.COOKIE_CONSENT, ACCEPTED);
		initYandexMetrika();
		setVisible(false);
	};

	const rejectAll = () => {
		localStorage.setItem(CACHEKEYs.COOKIE_CONSENT, REJECTED);
		setVisible(false);
	};

	return (
		<div className="fixed bottom-0 left-0 right-0 z-[200] bg-white border-t border-[#eee] shadow-[0_-4px_24px_rgba(0,0,0,0.08)]">
			<div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4 px-4 sm:px-8 py-4">
				<div className="max-w-[54rem]">
					<p className="text-[15px] text-[#1c1e22] leading-snug">
						<strong className="font-semibold">{t('message.cookie-consent-lead')}</strong> {t('message.cookie-consent-title')}
					</p>
					<p className="mt-1 text-[13px] text-[#62666d] leading-snug">{t('message.cookie-consent-description')}</p>
				</div>

				<div className="flex flex-col sm:flex-row gap-3 shrink-0 w-full sm:w-auto">
					<button type="button" onClick={rejectAll} className="h-[2.8rem] px-6 text-[0.9rem] rounded-[3px] border border-[#e4e4e4] bg-white text-[#1c1e22] hover:bg-[#f7f7f7] whitespace-nowrap">
						{t('label.cookie-reject-all')}
					</button>
					<GUIButton type="button" onClick={acceptAll} className="h-[2.8rem] px-6 text-[0.9rem] rounded-[3px] bg-[#1d3c5d] text-white hover:bg-[#163049] whitespace-nowrap">
						{t('label.cookie-accept-all')}
					</GUIButton>
				</div>
			</div>
		</div>
	);
};

export default CookieConsentBanner;
