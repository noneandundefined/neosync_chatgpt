import React from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';
import { basicFirmwareUpdate } from '@/rest/firmwaresAPI';
import { Device_DeviceConf_Sync } from '@/rest/deviceAPI.ts';
import Tooltip from '@/components/ui/Tooltip.tsx';

interface ModalUpdateFirmwareProps {
	device: Device_DeviceConf_Sync;
	onSuccess?: (imei: string) => void;
}

/* Modal for update firmware device */
const ModalUpdateFirmware: React.FC<ModalUpdateFirmwareProps> = ({ device, onSuccess }) => {
	const { t } = useTranslation();

	const hex = `0x${device.firmware_version_upd?.toString(16).toUpperCase().padStart(2, '0')}`;

	const handleUpdateFirmware = async () => {
		await basicFirmwareUpdate(device.imei);
		onSuccess?.(device.imei);
	};

	return (
		<React.Fragment>
			<div className="flex flex-col gap-2 mt-3">
				<p className="text-sm">
					{t('message.tracker-firmware-version')} <span className="text-[14px]">{device.imei}</span> {t('message.outdated-firmware')}.
				</p>
				<p className="text-sm mb-2">
					{t('message.new-version')}: <span className="text-[14px]">{hex}</span>
				</p>
				<p className="mb-6 text-gray-600 text-sm max-w-[45rem]">{t('message.firmware-update-description')}</p>

				<Tooltip title={device.status ? '' : t('message.device-not-connected-to-server')} position="bottom">
					<GUIButton onClick={handleUpdateFirmware} disabled={!device.status}>
						{t('label.update')}
					</GUIButton>
				</Tooltip>
			</div>
		</React.Fragment>
	);
};

export default ModalUpdateFirmware;
