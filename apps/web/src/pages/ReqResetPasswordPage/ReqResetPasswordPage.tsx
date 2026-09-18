import { Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { basicRequestPasswordReset } from '@/rest/authAPI';
import { ValidationEmailSchema } from '@/utils/ValidationSchema';

const ReqResetPasswordPage = () => {
	const { t } = useTranslation();

	const {
		register,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<{ email: string }>({
		mode: 'onChange',
		defaultValues: {
			email: '',
		},
	});

	const onSubmit = async (data: { email: string }) => {
		await basicRequestPasswordReset(data.email);
	};

	return (
		<main className="w-screen h-screen">
			<div className="flex justify-between items-center h-screen p-[1rem] sm:p-[4rem] lg:px-[7rem] xl:px-[12rem]">
				<div className="w-full sm:w-[55vw] lg:w-[35vw] xl:w-[28vw]">
					<h1 className="text-center sm:text-left text-[19px] sm:text-2xl font-semibold text-gray-900 mb-8">{t('label.reset-password')}</h1>

					<form className="flex flex-col items-stretch space-y-8">
						<div className="w-full flex flex-col space-y-6">
							<div className="w-full">
								<label className="block text-[15px] font-medium text-gray-800 mb-1">{t('label.login')}</label>
								<GUInput
									type="email"
									placeholder={t('label.login-placeholder')}
									{...register('email', ValidationEmailSchema<{ email: string }, 'email'>(t))}
									className="!min-h-[2.8rem]"
									error={errors.email?.message}
								/>
							</div>
						</div>

						<GUIButton type="submit" onClick={handleSubmit(onSubmit)} disabled={!watch('email')}>
							{t('label.reset-password')}
						</GUIButton>
					</form>

					<div className="text-center sm:text-left mt-6 text-sm" id="get__access__signin">
						<Link to={ROUTES.SIGNIN} className="text-[#1d3c5d] hover:underline">
							{t('message.back-to-signin')}
						</Link>
					</div>
				</div>

				<div className="hidden lg:block">
					<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[35vw] object-cover -ml-[10vw]" draggable={false} />
				</div>
			</div>
		</main>
	);
};

export default ReqResetPasswordPage;
