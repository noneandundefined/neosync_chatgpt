import { toast } from 'react-toastify';
import { useForm } from 'react-hook-form';
import useCaptcha from '@/hooks/useCaptcha';
import LinkIncorrect from './LinkIncorrect';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { basicPasswordReset } from '@/rest/authAPI';
import { Link, useLocation } from 'react-router-dom';
import GUIButton from '@/components/ui/Button/GUIButton';
import InputPassword from '@/components/ui/Form/InputPassword';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';
import { ValidationPasswordRequiredSchema } from '@/utils/ValidationSchema';

const ResetPasswordPage = () => {
	const { t } = useTranslation();
	const captcha = useCaptcha();

	const location = useLocation();
	const params = new URLSearchParams(location.search);

	const uuid = params.get('uuid');
	const exp = params.get('exp');
	const sig = params.get('sig');

	const {
		register,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<{
		password: string;
		password_repeat: string;
		turnstile_token: string;
	}>({
		mode: 'onChange',
		defaultValues: {
			password: '',
			password_repeat: '',
			turnstile_token: '',
		},
	});

	if (!uuid || !exp || !sig) return <LinkIncorrect />;

	const onSubmit = async (data: { password: string; password_repeat: string; turnstile_token: string }) => {
		if (data.password !== data.password_repeat) {
			toast.error(t('message.passwords-not-match'));
			return;
		}

		await basicPasswordReset(uuid, exp, sig, data.password, data.turnstile_token);
	};

	return (
		<main className="w-screen h-screen">
			<div className="flex justify-between items-center h-screen p-[1rem] sm:p-[4rem] lg:px-[7rem] xl:px-[12rem]">
				<div className="w-full sm:w-[55vw] lg:w-[35vw] xl:w-[28vw]">
					<h1 className="text-center sm:text-left text-[19px] sm:text-2xl font-semibold text-gray-900 mb-8">{t('label.reset-password')}</h1>

					<form className="flex flex-col items-stretch w-full space-y-8">
						<div>
							<label className="block text-[15px] font-medium text-gray-800 mb-1">{t('label.new-password')}</label>
							<InputPassword
								placeholder={t('label.password-placeholder')}
								className="!min-h-[2.8rem]"
								{...register(
									'password',
									ValidationPasswordRequiredSchema<
										{
											password: string;
											password_repeat: string;
											turnstile_token: string;
										},
										'password'
									>(t)
								)}
								error={errors.password?.message}
								onCopy={(e) => e.preventDefault()}
								onPaste={(e) => e.preventDefault()}
								onCut={(e) => e.preventDefault()}
							/>
						</div>

						<div>
							<label className="block text-[15px] font-medium text-gray-800 mb-1">{t('message.repeat-password')}</label>
							<InputPassword
								placeholder={t('message.repeat-password')}
								className="!min-h-[2.8rem]"
								{...register(
									'password_repeat',
									ValidationPasswordRequiredSchema<
										{
											password: string;
											password_repeat: string;
											turnstile_token: string;
										},
										'password_repeat'
									>(t)
								)}
								error={errors.password_repeat?.message}
								onCopy={(e) => e.preventDefault()}
								onPaste={(e) => e.preventDefault()}
								onCut={(e) => e.preventDefault()}
							/>
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
							disabled={!watch('password') || !watch('password_repeat')}
						>
							{t('label.reset-password')}
						</GUIButton>
					</form>

					<div className="text-left mt-6 text-sm" id="get__access__signin">
						<Link to={ROUTES.SIGNIN} className="text-[#1d3c5d] hover:underline">
							{t('message.back-to-signin')}
						</Link>
					</div>
				</div>

				<div className="hidden lg:block">
					<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[38vw] object-cover -ml-[10vw]" draggable={false} />
				</div>
			</div>
		</main>
	);
};

export default ResetPasswordPage;
