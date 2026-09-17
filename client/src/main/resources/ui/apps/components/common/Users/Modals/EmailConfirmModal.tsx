import useCaptcha from '@/hooks/useCaptcha';
import { UserResponse } from '@/rest/userAPI';
import { useTranslation } from 'react-i18next';
import React, { useEffect, useState } from 'react';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUIButton from '@/components/ui/Button/GUIButton';
import { LANGUAGES } from '@/constants/Language.constant';
import LabeledField from '@/components/ui/Form/LabeledField';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';

interface EmailConfirmModalProps {
	user: UserResponse | null;
	onFn: (token: string, language: string) => Promise<void>;
}

const resolveUserLanguage = (user: UserResponse | null) => {
	const language = user?.user_contact.language?.trim();
	if (language && LANGUAGES.some((item) => item.id === language)) {
		return language;
	}

	return 'ru';
};

const EmailConfirmModal: React.FC<EmailConfirmModalProps> = ({ user, onFn }) => {
	const { t } = useTranslation();
	const captcha = useCaptcha();
	const [language, setLanguage] = useState(() => resolveUserLanguage(user));

	useEffect(() => {
		setLanguage(resolveUserLanguage(user));
	}, [user?.user_contact.user_uuid, user?.user_contact.language]);

	const handleSend = async () => {
		if (!captcha.validate()) return;
		await onFn(captcha.token, language).finally(captcha.reset);
	};

	return (
		<React.Fragment>
			<div className="mt-4 sm:px-0 space-y-3">
				<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">
					{t('message.confirm-send-email-info', {
						login: user?.email,
					})}
				</p>

				<LabeledField label="message.send-email-recipient-language">
					<GUISelect value={language} onChange={(e) => setLanguage(e.target.value)}>
						{LANGUAGES.map((lang) => (
							<option value={lang.id} key={lang.id}>
								{lang.value}
							</option>
						))}
					</GUISelect>
				</LabeledField>

				<div className="border-l-2 border-l-[#395d95] p-2 px-4 space-y-2">
					<p className="text-sm sm:text-base text-[#49525f] max-w-full sm:max-w-[40rem]">
						<span className="font-medium text-[14px]">{t('label.login')}</span>: <span className="underline text-[14px] text-[#2847f3]">{user?.email}</span>
					</p>
					<p className="text-sm sm:text-base text-[#49525f] max-w-full sm:max-w-[40rem]">
						<span className="font-medium text-[14px]">{t('label.password')}</span>: <span className="text-[14px]">********</span>
					</p>
				</div>
			</div>

			{/* Yandex SmartCaptcha */}
			<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} className="mt-6" />

			<GUIButton className="mt-5" onClick={handleSend}>
				{t('label.send')}
			</GUIButton>
		</React.Fragment>
	);
};

export default EmailConfirmModal;
