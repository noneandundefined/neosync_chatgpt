import { useTranslation } from 'react-i18next';

export const Loading = () => {
	const { t } = useTranslation();

	return (
		<div className="flex gap-5 justify-center items-center cursor-default">
			<p className="text-[13px] text-[#49525f]">{t('message.please-wait')}</p>
			<div className={`w-[1rem] h-[1rem] border-b-2 border-[#49525f] rounded-full animate-spin`}></div>
		</div>
	);
};

export const SpinnerLoading = () => {
	return (
		<div className="flex gap-5 justify-center items-center cursor-default my-5">
			<div className={`w-[1rem] h-[1rem] border-b-2 border-[#49525f] rounded-full animate-spin`}></div>
		</div>
	);
};
