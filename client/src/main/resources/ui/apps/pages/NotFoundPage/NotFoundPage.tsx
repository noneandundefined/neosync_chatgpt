import { Link } from 'react-router-dom';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';

const NotFoundPage = () => {
	const { t } = useTranslation();

	return (
		<main className="flex flex-col w-screen h-screen justify-center items-center px-4 text-center">
			<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="w-[40%] max-w-[200px] sm:max-w-[15vw] object-contain my-5 opacity-90" draggable={false} />

			<div className="w-full max-w-[300px] sm:max-w-[500px] h-[1px] bg-[#eee] my-3" />

			<p className="mb-6 text-sm sm:text-base text-[#444] max-w-[90%] sm:max-w-[500px]">{t('message.page-not-found-desc')}</p>

			<Link to={ROUTES.HOME} className="px-4 sm:px-6 py-2 mt-2 rounded bg-[#E8EDF5] font-normal text-black hover:bg-white transition text-sm sm:text-base">
				{t('message.go-to-the-main-page')}
			</Link>
		</main>
	);
};

export default NotFoundPage;
