import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';

const IButton = () => {
	const { t } = useTranslation();

	return (
		<div className="min-w-[10rem]">
			<div className="flex items-center gap-10">
				<p className="font-normal">{t('message.ibutton-reader-polling')}</p>
				<CFGSwitch uid={UIDs.IBUTTON_ENABLED} vOff={0} vOn={1} />
			</div>
		</div>
	);
};

export default IButton;
