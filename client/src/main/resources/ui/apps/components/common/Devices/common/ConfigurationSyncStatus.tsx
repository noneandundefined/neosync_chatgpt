import { useTranslation } from 'react-i18next';
import { Device_DeviceConf_Sync } from '@/rest/deviceAPI';

interface ConfigurationSyncStatusProps {
	device: Pick<Device_DeviceConf_Sync, 'configuration_sync_status' | 'configuration_sync_error'>;
	className?: string;
}

export const getConfigurationSyncStatusLabel = (t: (key: string) => string, status?: string | null, error?: string | null): string => {
	const normalized = status ?? 'idle';

	switch (normalized) {
		case 'pending':
			return t('message.configuration-sync-status-pending');
		case 'confirmed':
			return t('message.configuration-sync-status-confirmed');
		case 'failed':
			if (error === 'configuration hash mismatch') {
				return t('message.configuration-sync-error-hash-mismatch');
			}
			return error ?? t('message.configuration-sync-status-failed');
		default:
			return t('message.configuration-sync-status-idle');
	}
};

const ConfigurationSyncStatus: React.FC<ConfigurationSyncStatusProps> = ({ device, className = '' }) => {
	const { t } = useTranslation();
	const status = device.configuration_sync_status ?? 'idle';
	const label = getConfigurationSyncStatusLabel(t, status, device.configuration_sync_error);

	return (
		<span className={className} style={{ color: '#000000' }}>
			{label}
		</span>
	);
};

export default ConfigurationSyncStatus;
