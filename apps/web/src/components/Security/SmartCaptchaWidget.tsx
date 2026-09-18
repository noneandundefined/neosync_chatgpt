import { config } from '@/.config/config.client';
import { YANDEX_SMARTCAPTCHA_LINK } from '@/constants/SmartCaptcha.constant';
import { forwardRef, useEffect, useId, useImperativeHandle, useRef } from 'react';
import { useTranslation } from 'react-i18next';

interface SmartCaptchaWidgetProps {
	onVerify: (token: string) => void;
	className?: string;
}

export interface SmartCaptchaHandle {
	reset: () => void;
}

type SmartCaptchaLang = 'ru' | 'en' | 'be' | 'kk' | 'tt' | 'uk' | 'uz' | 'tr';

interface SmartCaptchaApi {
	render: (
		container: HTMLElement | string,
		params: {
			sitekey: string;
			callback?: (token: string) => void;
			hl?: SmartCaptchaLang;
			test?: boolean;
		}
	) => number;
	reset: (widgetId?: number) => void;
	destroy: (widgetId?: number) => void;
	subscribe?: (widgetId: number, event: string, callback: () => void) => void;
}

interface SmartCaptchaWindow extends Window {
	smartCaptcha?: SmartCaptchaApi;
}

const isDev = config.type.release === 'dev';
const DEV_CAPTCHA_TOKEN = 'dev-skip';

const getSiteKey = () => {
	const fromEnv = import.meta.env.VITE_YANDEX_SMARTCAPTCHA as string | undefined;
	if (fromEnv) {
		return fromEnv;
	}

	if (isDev) {
		return config.yandex.SMARTCAPTCHA_SITE_KEY_DEV;
	}

	if (config.type.release == 'prod') {
		return config.yandex.SMARTCAPTCHA_SITE_KEY_PROD;
	}

	return config.yandex.SMARTCAPTCHA_SITE_KEY_PREPROD;
};

const getWidgetLang = (language: string): SmartCaptchaLang => {
	const code = language.split('-')[0];
	if (code === 'ru') {
		return 'ru';
	}

	return 'en';
};

const readSmartToken = (container: HTMLElement) => {
	return container.querySelector<HTMLInputElement>('input[name="smart-token"]')?.value ?? '';
};

const SmartCaptchaWidget = forwardRef<SmartCaptchaHandle, SmartCaptchaWidgetProps>(({ onVerify, className }, ref) => {
	const { i18n } = useTranslation();
	const reactId = useId();
	const containerRef = useRef<HTMLDivElement | null>(null);
	const widgetId = useRef<number | null>(null);
	const onVerifyRef = useRef(onVerify);
	const lastTokenRef = useRef('');
	const callbackName = `neosyncSmartCaptcha_${reactId.replace(/[^a-zA-Z0-9_]/g, '')}`;

	useImperativeHandle(ref, () => ({
		reset: () => {
			const w = window as unknown as SmartCaptchaWindow;
			if (w.smartCaptcha && widgetId.current != null) {
				w.smartCaptcha.reset(widgetId.current);
			}
			lastTokenRef.current = '';
			onVerifyRef.current('');
		},
	}));

	useEffect(() => {
		onVerifyRef.current = onVerify;
	}, [onVerify]);

	useEffect(() => {
		if (isDev) {
			onVerifyRef.current(DEV_CAPTCHA_TOKEN);
			return;
		}

		const siteKey = getSiteKey();
		const container = containerRef.current;

		if (!siteKey || !container) {
			return;
		}

		const emitToken = (token: string) => {
			if (token === lastTokenRef.current) {
				return;
			}

			lastTokenRef.current = token;
			onVerifyRef.current(token);
		};

		(window as unknown as Record<string, unknown>)[callbackName] = (token: string) => emitToken(token ?? '');

		const observer = new MutationObserver(() => emitToken(readSmartToken(container)));
		observer.observe(container, { childList: true, subtree: true, attributes: true, attributeFilter: ['value'] });

		if (!document.querySelector(`script[src="${YANDEX_SMARTCAPTCHA_LINK}"]`)) {
			const script = document.createElement('script');
			script.src = YANDEX_SMARTCAPTCHA_LINK;
			script.defer = true;
			document.head.appendChild(script);
		}

		const w = window as unknown as SmartCaptchaWindow;

		const interval = setInterval(() => {
			if (!w.smartCaptcha || widgetId.current != null) {
				return;
			}

			widgetId.current = w.smartCaptcha.render(container, {
				sitekey: siteKey,
				hl: getWidgetLang(i18n.language),
				callback: (token: string) => emitToken(token),
			});

			w.smartCaptcha.subscribe?.(widgetId.current, 'token-expired', () => emitToken(''));
			w.smartCaptcha.subscribe?.(widgetId.current, 'network-error', () => emitToken(''));
			clearInterval(interval);
		}, 100);

		return () => {
			clearInterval(interval);
			observer.disconnect();
			delete (window as unknown as Record<string, unknown>)[callbackName];

			if (w.smartCaptcha && widgetId.current != null) {
				w.smartCaptcha.destroy(widgetId.current);
				widgetId.current = null;
			}

			lastTokenRef.current = '';
			onVerifyRef.current('');
		};
	}, [callbackName, i18n.language]);

	if (isDev) {
		return null;
	}

	const siteKey = getSiteKey();
	if (!siteKey) {
		return null;
	}

	return (
		<div
			ref={containerRef}
			id={`captcha-container-${reactId.replace(/:/g, '')}`}
			className={['smart-captcha', className].filter(Boolean).join(' ')}
			data-sitekey={siteKey}
			data-hl={getWidgetLang(i18n.language)}
			data-callback={callbackName}
			style={{ height: 100 }}
		></div>
	);
});

export default SmartCaptchaWidget;
