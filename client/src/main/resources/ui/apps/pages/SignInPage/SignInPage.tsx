import '@/utils/i18n';
import { Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import useCaptcha from '@/hooks/useCaptcha';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { basicAuthSignIn } from '@/rest/authAPI';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import InputPassword from '@/components/ui/Form/InputPassword';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';
import { SignInRequestWithToken } from '@/interface/auth/signInRequest.interface';
import { ValidationEmailSchema, ValidationPasswordRequiredSchema } from '@/utils/ValidationSchema';

const SignInPage = () => {
	const captcha = useCaptcha();
	const { t } = useTranslation();
	const remember = localStorage.getItem(CACHEKEYs.AUTH_REMEMBER_ME) === 'true';

	const {
		register,
		handleSubmit,
		watch,
		setValue,
		formState: { errors },
	} = useForm<SignInRequestWithToken>({
		mode: 'onChange',
		defaultValues: {
			login: '',
			password: '',
			remember_me: remember,
			turnstile_token: '',
		},
	});

	const onSubmit = async (data: SignInRequestWithToken) => {
		await basicAuthSignIn(data);

		localStorage.setItem(CACHEKEYs.AUTH_REMEMBER_ME, String(data.remember_me));

		window.location.replace(ROUTES.HOME);
	};

	return (
		<main className="w-screen h-screen">
			<div className="flex justify-between items-center h-screen p-[1rem] sm:p-[4rem] lg:px-[7rem] xl:px-[12rem]">
				<div className="w-full sm:w-[55vw] lg:w-[35vw] xl:w-[28vw]">
					<h1 className="text-center sm:text-left text-[19px] sm:text-2xl font-semibold text-gray-900 mb-8">{t('message.signin-title')}</h1>

					<form className="flex flex-col items-stretch w-full space-y-8">
						<div className="flex flex-col space-y-6">
							<div>
								<label className="block text-[15px] font-medium text-gray-800 mb-1">{t('label.login')}</label>
								<GUInput
									type="email"
									{...register('login', ValidationEmailSchema<SignInRequestWithToken, 'login'>(t))}
									placeholder={t('label.login-placeholder')}
									className="!min-h-[2.8rem]"
									error={errors.login?.message}
								/>
							</div>

							<div>
								<InputPassword
									label={t('label.password')}
									placeholder={t('label.password-placeholder')}
									className="!min-h-[2.8rem]"
									{...register('password', ValidationPasswordRequiredSchema<SignInRequestWithToken, 'password'>(t))}
									error={errors.password?.message}
								/>
							</div>
						</div>

						<div className="flex flex-col gap-6 justify-between text-sm">
							<label className="flex items-center text-[15px] cursor-pointer">
								<GUICheckbox checked={watch('remember_me')} onChange={(v) => setValue('remember_me', v)} />
								{t('label.remember-me')}
							</label>
							<div className="flex justify-between">
								<Link to={ROUTES.RESET_PASSWORD_REQ} className="text-[14px] text-[#1d3c5d] hover:underline">
									{t('label.forgot-password')}
								</Link>
							</div>
						</div>

						{/* Yandex SmartCaptcha */}
						<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} />

						<GUIButton
							type="submit"
							onClick={handleSubmit(async (data) => {
								if (!captcha.validate()) return;

								await onSubmit({
									...data,
									turnstile_token: captcha.token,
								}).finally(captcha.reset);
							})}
							disabled={!watch('login') || !watch('password')}
						>
							{t('label.signin')}
						</GUIButton>
					</form>

					<div className="text-center sm:text-left mt-6 text-sm">
						<Link to={ROUTES.HOW_GET_ACCESS} className="text-[#1d3c5d] hover:underline">
							{t('label.get-access')}
						</Link>
					</div>
				</div>

				<div className="hidden lg:block">
					<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[35vw] object-cover" draggable={false} />
				</div>
			</div>
		</main>
	);
};

export default SignInPage;
