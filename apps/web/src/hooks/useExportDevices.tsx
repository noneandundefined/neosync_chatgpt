import { basicDevicesGetWithParams } from '@/rest/deviceAPI';
import { EXPORT_ADM_DEVICES_PATHNAME } from '@/constants/Files.constants';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';

export const useExportDevices = () => {
	const exportTxt = async () => {
		const columnSort: ColumnSortState = { key: 'created_at', direction: 'desc' };

		const first = await basicDevicesGetWithParams(1, 1, '', columnSort);
		const totalAll = first.total_all ?? first.total ?? 0;

		if (!totalAll) {
			const blob = new Blob([''], { type: 'text/plain;charset=utf-8' });
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');

			a.href = url;
			a.download = EXPORT_ADM_DEVICES_PATHNAME;
			a.click();

			URL.revokeObjectURL(url);
			return;
		}

		const chunkSize = 1000;
		const pages = Math.max(1, Math.ceil(totalAll / chunkSize));
		const allItems: any[] = [];

		for (let page = 1; page <= pages; page += 1) {
			const resp = await basicDevicesGetWithParams(page, chunkSize, '', columnSort);
			allItems.push(...(resp.items ?? []));
		}

		const lines = allItems.map((d: any) => {
			const imei = d.imei;
			const model = d.device_model ?? '';

			const arr = [imei, model];

			while (arr[arr.length - 1] === '') arr.pop();

			return arr.join(',');
		});

		const text = lines.join('\n');

		const blob = new Blob([text], { type: 'text/plain;charset=utf-8' });
		const url = URL.createObjectURL(blob);

		const a = document.createElement('a');
		a.href = url;
		a.download = EXPORT_ADM_DEVICES_PATHNAME;
		a.click();

		URL.revokeObjectURL(url);
	};

	return { exportTxt };
};
