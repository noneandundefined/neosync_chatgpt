import { useTranslation } from 'react-i18next';

/* Fallback technical mode */
const TechnicalWork = () => {
	const { t } = useTranslation();

	return (
		<main className="flex flex-col w-screen h-screen justify-center items-center">
			<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[15vw] min-w-[15rem] object-cover my-7 opacity-90" draggable={false} />
			<div className="w-full h-[1px] bg-[#eee]" />
			<p className="my-4 text-[2.2rem] uppercase tracking-wide font-normal max-w-[40vw] min-w-[25rem] text-balance text-center">{t('message.come-back-later')}</p>
			<p className="mb-3 text-[1rem] text-[#444] font-normal max-w-[30vw] min-w-[25rem] text-balance text-center">{t('message.come-back-later-desc')}</p>
		</main>
	);
};

export default TechnicalWork;
