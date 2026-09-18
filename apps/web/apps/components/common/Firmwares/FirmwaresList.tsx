import { Dispatch, SetStateAction, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { firmwaresColumns } from '.';
import GenericTable, { FetchFn } from '../Table/GenericTable/GenericTable';
import { basicFirmwareSources, Firmwares_DeviceModelSource } from '@/rest/firmwaresAPI';
import Information from '@/components/@icons/information';
import History from '@/components/@icons/history';

export interface FirmwaresListProps {
	/** Токен для обновления данных в таблице */
	sysToken: number;
	/** Токен для обновления данных в таблице */
	setSysToken: Dispatch<SetStateAction<number>>;
}

const FirmwaresList: React.FC<FirmwaresListProps> = ({ sysToken }) => {
	const { t } = useTranslation();

	const handleRowAction = (_: Firmwares_DeviceModelSource, _1: number, isHovered: boolean) => (
		<div
			className="flex items-center justify-center gap-2 sm:gap-1"
			style={{
				opacity: isHovered ? 1 : 0,
			}}
		>
			<span title={t('message.firmware-history-title')}>
				<History fill="#5ba0e8" size={24} className="cursor-pointer transition" />
			</span>

			<span title={t('message.firmware-info-title')}>
				<Information fill="#1976d2" size={24} className="cursor-pointer transition" />
			</span>
		</div>
	);

	const columns = useMemo(() => firmwaresColumns(t), [t]);

	const fetchFirmwares: FetchFn<Firmwares_DeviceModelSource> = async ({ page, limit, search }, signal) => {
		const res = await basicFirmwareSources(page, limit, search, signal);
		return { items: res.items, total: 0, totalPages: res.total_pages };
	};

	return (
		<GenericTable<Firmwares_DeviceModelSource>
			tableKey="fetchFirmwares"
			columns={columns}
			fetchFn={fetchFirmwares}
			sysToken={sysToken}
			getRowId={(fota) => fota.device_model}
			rowActions={(fota, idx, isHovered) => handleRowAction(fota, idx, isHovered)}
			onDeleteFn={async (_) => {
				console.log('deleted');
			}}
		/>
	);
};

export default FirmwaresList;
