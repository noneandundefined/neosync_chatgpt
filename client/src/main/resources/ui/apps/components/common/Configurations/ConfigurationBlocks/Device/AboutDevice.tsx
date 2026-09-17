import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { copyToBuffer } from '@/utils/CopyBufferUtils';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationContext } from '@/context/useConfigurationContext';

const AboutDevice = () => {
	const { t } = useTranslation();
	const { model, device, imei } = useConfigurationContext();

	const [time, setTime] = useState<string>('');

	useEffect(() => {
		if (!device.status && device.updated_at) {
			const updatedAt = new Date(device.updated_at);

			const formatted = updatedAt.toLocaleString('ru-RU', {
				timeZone: 'UTC',
				hour12: false,
				year: 'numeric',
				month: '2-digit',
				day: '2-digit',
				hour: '2-digit',
				minute: '2-digit',
				second: '2-digit',
			});

			setTime(formatted + ' UTC');
			return;
		}

		const updateTime = () => {
			const now = new Date();

			const formatted = now.toLocaleString('ru-RU', {
				timeZone: 'UTC',
				hour12: false,
				year: 'numeric',
				month: '2-digit',
				day: '2-digit',
				hour: '2-digit',
				minute: '2-digit',
				second: '2-digit',
			});

			setTime(formatted + ' UTC');
		};

		updateTime();

		const interval = setInterval(updateTime, 3000);

		return () => clearInterval(interval);
	}, [device.status, device.updated_at]);

	/** Firmwares */
	const firmwareVersion = device.firmware_version ?? 0;
	const firmwareVersionUpd = device.firmware_version_upd ?? 0;
	const firmwareUpdating = device.firmware_update_status === 'pending';

	let colorFirmware = device.activated ? (firmwareVersion < firmwareVersionUpd ? '#b50202' : 'black') : '#ccc';
	if (firmwareUpdating) colorFirmware = '#dd6300';

	const firmwareValue = firmwareUpdating ? t('message.firmware-update-in-progress') : `0x${firmwareVersion.toString(16).toUpperCase().padStart(2, '0')}`;

	return (
		<div className="space-y-3">
			<LabeledField label="label.time-utc" className="flex flex-col gap-1">
				<input type="text" value={time} disabled />
			</LabeledField>

			<LabeledField label="message.remote-configuration-server" className="flex flex-col gap-1">
				<input type="text" value={device.status ? t('message.packet-has-been-sent-successfully') : t('label.disabled')} className={device.status ? 'text-[#067b00]' : 'text-[#b50202]'} disabled />
			</LabeledField>

			<LabeledField label="label.device-type" className="flex flex-col gap-1">
				<input type="text" value={device.device_extended_model?.toUpperCase() ?? model} disabled />
			</LabeledField>

			<LabeledField label="message.tracker-firmware-version" className="flex flex-col gap-1">
				<input type="text" style={{ color: colorFirmware }} value={firmwareValue} disabled />
			</LabeledField>

			<LabeledField label="IMEI" className="flex flex-col gap-1">
				<div className="flex w-full items-center gap-2">
					<input type="text" value={device.imei} disabled />
					<button id="buttonhlp" className="!text-sm" onClick={() => copyToBuffer(imei)}>
						{t('label.copy')}
					</button>
				</div>
			</LabeledField>
		</div>
	);
};

export default AboutDevice;
