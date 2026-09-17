import { toast } from 'react-toastify';
import { useCallback, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { config } from '@/.config/config.client';
import type { SmartCaptchaHandle } from '@/components/Security/SmartCaptchaWidget';

const DEV_CAPTCHA_TOKEN = 'dev-skip';
const isDev = config.type.release === 'dev';

const useCaptcha = () => {
	const { t } = useTranslation();

	const [token, setToken] = useState<string>(isDev ? DEV_CAPTCHA_TOKEN : '');
	const widgetRef = useRef<SmartCaptchaHandle>(null);

	const reset = useCallback(() => {
		if (isDev) return;
		widgetRef.current?.reset();
		setToken('');
	}, []);

	const onVerify = useCallback((value: string) => {
		setToken(value);
	}, []);

	const validate = () => {
		if (isDev) {
			return true;
		}

		if (!token) {
			toast.error(t('message.required-captcha'));
			return false;
		}

		return true;
	};

	return {
		token,
		widgetRef,
		reset,
		onVerify,
		validate,
		isVerified: Boolean(token),
	};
};

export default useCaptcha;
