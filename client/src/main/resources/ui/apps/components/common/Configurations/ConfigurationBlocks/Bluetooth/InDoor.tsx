import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';

const InDoor = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-4">
			<div className="flex items-center justify-between">
				<p>{t('message.receive-tag-identifier')}</p>
				<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.RECEIVE_TAG_IDENTIFIER} vOff={0} vOn={1} />
			</div>

			<div className="flex items-center justify-between">
				<p className="max-w-[80%]">{t('message.replace-coordinates-with-ble-tags')}</p>
				<Tooltip title={t('message.tooltip-tag-coordinates-protocol-setting')}>
					<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.REPLACE_COORDINATES_WITH_BLE_TAGS} vOff={0} vOn={1} />
				</Tooltip>
			</div>
		</div>
	);
};

export default InDoor;
