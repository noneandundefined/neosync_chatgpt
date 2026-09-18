import { formatDate } from '@/utils/TimeUtils';
import { Firmwares_DeviceModelSource } from '@/rest/firmwaresAPI';

export const firmwaresColumns = (t: (key: string) => string) => {
	return [
		{ title: t('label.model'), key: 'device_model' },
		{ title: t('label.firmware'), key: 'firmware' },
		{ title: t('label.release_day'), key: 'release_date' },
		{
			title: t('label.type'),
			key: 'type',
			render: (fota: Firmwares_DeviceModelSource) => {
				return <span>{fota.type === 'release' ? <p>{t('label.released')}</p> : <p>{t('label.unknown')}</p>}</span>;
			},
		},
		{
			title: t('label.created_at'),
			key: 'created_at',
			render: (fota: Firmwares_DeviceModelSource) => {
				return fota.created_at ? <span>{formatDate(fota.created_at)}</span> : <span className="text-[#ccc]">{t('label.unknown')}</span>;
			},
		},
	];
};
