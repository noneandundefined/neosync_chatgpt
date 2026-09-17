import { useTranslation } from 'react-i18next';

const Spinner = () => {
	const { t } = useTranslation();

	return (
		<div className="relative flex justify-center items-center">
			<p className="font-medium">{t('message.please-wait')}</p>
			<div className={`absolute right-0 w-[1.24rem] h-[1.24rem] border-b-2 border-[#000] rounded-full animate-spin`}></div>
		</div>
	);
};

export default Spinner;
