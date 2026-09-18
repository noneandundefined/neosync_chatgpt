import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface LoaderLineProps {
	title: string;
	description: string[] | string | null;
}

const LoaderLine: React.FC<LoaderLineProps> = ({ title, description }) => {
	const { t } = useTranslation();
	const [dots, setDots] = useState('');
	const [index, setIndex] = useState(0);

	const isArray = Array.isArray(description);

	useEffect(() => {
		const interval = setInterval(() => {
			setDots((prev) => (prev.length >= 3 ? '' : prev + '.'));
		}, 400);

		return () => clearInterval(interval);
	}, []);

	useEffect(() => {
		if (!isArray) return;

		const interval = setInterval(() => {
			setIndex((prev) => (description && Array.isArray(description) ? (prev + 1) % description.length : 0));
		}, 5000);

		return () => clearInterval(interval);
	}, [description, isArray]);

	const currentText = isArray ? description[index] : (description ?? t('message.connect-to-server'));

	return (
		<div className="fixed inset-0 h-screen w-screen bg-[#00000070] z-[1000] flex flex-col items-center justify-center">
			<div className="bg-[#fff] w-screen h-auto pt-[3rem] pb-0 flex flex-col justify-center" style={{ boxShadow: '0 0 20px 0 rgba(0, 0, 0, 0.15' }}>
				<div className="w-full text-center sm:text-left sm:ml-[25%] px-3">
					<p className="font-medium text-[1.45rem]">{t(`message.${title}`)}</p>
					<p className="text-[14px] text-[#555] my-3">
						{currentText} <span className="dots">{dots}</span>
					</p>
				</div>
				<div id="loader-line" className="mt-[7rem]"></div>
			</div>
		</div>
	);
};

export default LoaderLine;
