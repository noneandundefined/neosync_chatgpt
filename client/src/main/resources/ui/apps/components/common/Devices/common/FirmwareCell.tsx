import { Device_DeviceConf_Sync } from '@/rest/deviceAPI';

interface FirmwareCellProps {
	d: Device_DeviceConf_Sync;
	t: (key: string) => string;
	onClick: any;
	updatedDevices: Set<string>;
}

const FirmwareCell: React.FC<FirmwareCellProps> = ({ d, t, onClick, updatedDevices }) => {
	const current = d.firmware_version;
	const updated = d.firmware_version_upd ?? 0;
	const hex = `0x${current.toString(16).toUpperCase().padStart(2, '0')}`;

	const isUpdating = d.firmware_update_status === 'pending';
	const isUpdated = updatedDevices.has(d.imei);

	let color = d.activated ? (current < updated ? '#b50202' : 'black') : '#ccc';
	if (isUpdated) color = 'green';
	if (isUpdating) color = '#dd6300';

	return (
		<span
			style={{
				color: d.activated ? color : '#ccc',
				cursor: d.activated && current < updated && !isUpdating ? 'pointer' : 'default',
				textDecoration: d.activated && current < updated && !isUpdating ? 'underline' : 'none',
			}}
			className={isUpdated || isUpdating ? 'animate-pulse' : ''}
			onClick={() => d.activated && current < updated && !isUpdating && onClick(d)}
		>
			{d.activated ? (isUpdating ? t('message.firmware-update-in-progress') : hex) : t('label.unknown')}
		</span>
	);
};

export default FirmwareCell;
