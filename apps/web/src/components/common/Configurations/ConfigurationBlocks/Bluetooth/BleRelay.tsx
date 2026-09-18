import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import LabeledField from '@/components/ui/Form/LabeledField';

const BleRelay = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-4">
			<LabeledField label="label.ble-relay-status">
				<Tooltip title={t('message.tooltip-relay-state-setting')}>
					<CFGInput uid={UIDs.BLE_OUTPUT} />
				</Tooltip>
			</LabeledField>

			<LabeledField label="label.ble-broadcast-period">
				<Tooltip title={t('message.tooltip-advertising-active-period')}>
					<div className="flex items-center gap-2">
						<CFGInput uid={UIDs.BLE_BROADCAST_ACTIVE_PERIOD_SEC} />
						<label>{t('label.s')}</label>
					</div>
				</Tooltip>
			</LabeledField>

			<LabeledField label="label.ble-period-silence">
				<Tooltip title={t('message.tooltip-advertising-wait-period')}>
					<div className="flex items-center gap-2">
						<CFGInput uid={UIDs.BLE_BROADCAST_SILENT_PERIOD_SEC} />
						<label>{t('label.s')}</label>
					</div>
				</Tooltip>
			</LabeledField>

			<LabeledField label="label.address">
				<Tooltip title={t('message.tooltip-ble-relay-mac-address')}>
					<CFGInput uid={UIDs.BLE_OUTPUT_ADDR} />
				</Tooltip>
			</LabeledField>

			<LabeledField label="label.ble-encryption-key">
				<Tooltip title={t('message.tooltip-encryption-key-input')}>
					<CFGInput uid={UIDs.BLE_OUTPUT_KEY} />
				</Tooltip>
			</LabeledField>
		</div>
	);
};

export default BleRelay;
