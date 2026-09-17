import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

const HowGetAccessPage = () => {
	const { t } = useTranslation();

	return (
		<main className="w-screen h-screen">
			<div className="flex justify-between items-center h-screen p-[2rem] sm:p-[4rem] lg:px-[7rem] xl:px-[12rem]">
				<div className="w-full sm:w-[55vw] lg:w-[35vw] xl:w-[28vw]">
					<div>
						<h1 className="text-center sm:text-left text-[19px] sm:text-2xl font-semibold text-gray-900 mb-8">{t('message.title-how-get-access')}?</h1>

						<div className="flex flex-col items-stretch space-y-8">
							<div className="w-full flex flex-col space-y-6 items-center sm:items-start">
								<p>{t('message.need-to-get-access')}:</p>

								<div className="mt-5 space-y-3">
									<ul className="list-decimal list-outside ml-0 sm:ml-4 space-y-2 text-[#334155] inline-block sm:text-left">
										<li>
											<span className="font-semibold">{t('label.name-company')}.</span>
										</li>
										<li>
											<span className="font-semibold">E-mail</span>, {t('message.email-body-req-access')}.
										</li>
										<li>
											<span className="font-semibold">{t('label.contact-data')}</span> {t('message.responsible-person-access')}.
										</li>
									</ul>
								</div>

								<p className="text-[14px]">
									{t('message.request-template')}:{' '}
									<a href="https://disk.yandex.ru/i/mNO-3WBBog2zkA" target="_blank" rel="noopener noreferrer" className="text-[#1d3c5d] hover:underline">
										https://disk.yandex.ru/i/mNO-3WBBog2zkA
									</a>
								</p>

								<div>
									<p className="font-semibold text-[14px] text-[#9a0000]">{t('message.howgetaccess-warning')}</p>
								</div>
							</div>
						</div>

						<div className="text-center sm:text-left mt-6 text-sm">
							<Link to={ROUTES.SIGNIN} className="text-[#1d3c5d] hover:underline">
								{t('message.back-to-signin')}
							</Link>
						</div>
					</div>
				</div>

				<div className="hidden lg:block">
					<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[35vw] object-cover -ml-[10vw]" draggable={false} />
				</div>
			</div>
		</main>
	);
};

export default HowGetAccessPage;
