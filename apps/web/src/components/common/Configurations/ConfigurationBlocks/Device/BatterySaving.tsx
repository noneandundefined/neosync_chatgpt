import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';

const BatterySaving = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-2">
			<div className="flex items-center justify-between">
				<p className="text-[14px]">{t('message.charging-limit')}</p>
				<Tooltip title={t('message.tooltip-battery-charge-limit-description')}>
					<CFGSwitch uid={UIDs.BATTERY_SAFE_MODE} bit={2} />
				</Tooltip>
			</div>

			<p className="text-sm text-[#666] max-w-[90%]"></p>
		</div>
	);
};

export default BatterySaving;
