import { useTranslation } from 'react-i18next';

/* Fallback for unexpected errors on the client */
const ErrorFallback = () => {
	const { t } = useTranslation();

	return (
		<main className="fixed top-0 left-0 bg-white z-[100] flex flex-col items-center justify-center h-screen w-screen px-4 text-center">
			<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="w-[40%] max-w-[200px] sm:max-w-[15vw] object-contain my-5 opacity-90" draggable={false} />
			<div className="w-full max-w-[300px] sm:max-w-[500px] h-[1px] bg-[#eee]" />
			<h1 className="my-5 text-xl sm:text-3xl">{t('message.something-went-wrong')}</h1>
			<p className="text-sm sm:text-lg text-gray-700 max-w-[90%] sm:max-w-[500px]">{t('message.something-went-wrong-desc')}</p>
		</main>
	);
};

export default ErrorFallback;
