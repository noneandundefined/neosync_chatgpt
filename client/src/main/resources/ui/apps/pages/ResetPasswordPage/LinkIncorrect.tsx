import Header from '@/components/Header/Header';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

const LinkIncorrect = () => {
	const { t } = useTranslation();

	return (
		<main className="flex h-screen w-screen min-h-0 flex-col bg-[#f7f9fc]" id="signin__main">
			<Header />

			<div className="flex min-h-0 flex-1 items-center justify-center">
				<div className="w-[50vw]" id="signin__resize__width">
					<div className="flex flex-col items-start space-y-8" id="signin__position">
						<div className="w-[50%] p-5 shadow-[0_0_7px_0_rgba(0,0,0,0.08)] text-center space-y-3">
							<p className="text-[16px]">{t('message.incorrect-link-password-reset')}</p>
							<div className="mt-6 text-sm" id="get__access__signin">
								<Link to={ROUTES.SIGNIN} className="text-[#1d3c5d] hover:underline">
									{t('message.back-to-signin')}
								</Link>
							</div>
						</div>
					</div>
				</div>

				<div id="signin__right__info">
					<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[38vw] object-cover -ml-[10vw]" draggable={false} />
				</div>
			</div>
		</main>
	);
};

export default LinkIncorrect;
