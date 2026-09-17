import Alert from '../@icons/alert';
import { useTranslation } from 'react-i18next';

const TerminalTemporarilyLimited = () => {
	const { t } = useTranslation();

	return (
		<div className={`fixed top-0 left-0 w-full z-[103] flex items-center justify-center bg-white`}>
			<div className="py-[3px] px-3 rounded-full cursor-default">
				<div className="flex items-center gap-2">
					<Alert fill="#cd9a00" size={17} />
					<p className="text-[12px] text-[#000]">{t('message.terminal-temporarily-limited-updates')}</p>
				</div>
			</div>
		</div>
	);
};

export default TerminalTemporarilyLimited;
