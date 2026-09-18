import * as XLSX from 'xlsx';
import i18next from 'i18next';
import { toast } from 'react-toastify';
import FileHelper from '@/utils/FilesUtils';
import { useHandleServer } from './Server/useHandleServer';
import { basicDevicesImportCheck } from '@/rest/deviceAPI';
import { basicFirmwareDeviceModelsGet } from '@/rest/firmwaresAPI';
import { parseDevicesFromFile } from '@/utils/DeviceFileParserUtils';
import { Dispatch, SetStateAction, useCallback, useState } from 'react';
import { DeviceFileParser } from '@/interface/device/deviceFileParser.interface';

const MAX_DEVICES = 100;

let cachedDevices: DeviceFileParser[] = [];

export const useDeviceFileUpload = () => {
	const [devices, setDevicesState] = useState<DeviceFileParser[]>(cachedDevices);

	const setDevices = useCallback<Dispatch<SetStateAction<DeviceFileParser[]>>>((action) => {
		setDevicesState((prev) => {
			const next = typeof action === 'function' ? action(prev) : action;
			cachedDevices = next;
			return next;
		});
	}, []);

	const { data: basicFotaDeviceModelsGetResp } = useHandleServer(['basicFotaDeviceModelsGetResp'], basicFirmwareDeviceModelsGet);

	const handleFileUpload = async (file?: File): Promise<DeviceFileParser[] | undefined> => {
		if (!file) return;

		const allowedExtensions = ['csv', 'xlsx', 'txt'];
		const ext = file.name.split('.').pop()?.toLowerCase();

		if (!ext || !allowedExtensions.includes(ext)) {
			toast.error(i18next.t('message.error-ext-device-imports'));
			return;
		}

		if (!FileHelper.fileMax1Mb(file)) {
			toast.error(i18next.t('message.file-too-large-1mb'));
			return;
		}

		let parsed: DeviceFileParser[] = [];

		try {
			if (ext === 'xlsx') {
				parsed = await parseXlsx(file);
			} else {
				const text = await file.text();

				if (text.includes('�')) {
					toast.error(i18next.t('message.error-utf8-device-imports'));
					return;
				}

				parsed = parseDevicesFromFile(text);
			}
		} catch {
			toast.error(i18next.t('message.error-read-device-imports'));
			return;
		}

		if (parsed.length > MAX_DEVICES) {
			toast.error(i18next.t('message.error-imports-100-devices'));
		}

		const parsedLimited = parsed.slice(0, MAX_DEVICES);

		const uniqueDevices: DeviceFileParser[] = [];
		const seenImei = new Set<string>();

		let hasError = false;

		for (const device of parsedLimited) {
			device.error = undefined;

			if (!device.imei) {
				device.error = i18next.t('message.error-device-imports-imei');
			} else if (device.imei.length !== 15) {
				device.error = i18next.t('message.error-device-imports-imei-validate');
			}

			if (seenImei.has(device.imei)) {
				continue;
			}

			if (device.model && basicFotaDeviceModelsGetResp && !basicFotaDeviceModelsGetResp.includes(device.model.toUpperCase())) {
				device.error = i18next.t('message.error-device-imports-model-validate').toLowerCase();
			}

			if (device.error) hasError = true;

			seenImei.add(device.imei);
			uniqueDevices.push(device);
		}

		const imeisToCheck = uniqueDevices.filter((device) => !device.error && device.imei.length === 15).map((device) => device.imei);

		if (imeisToCheck.length > 0) {
			try {
				const serverErrors = await basicDevicesImportCheck(imeisToCheck);
				const errorsByImei = new Map(serverErrors.map((item) => [item.imei, item.error]));

				for (const device of uniqueDevices) {
					const serverError = errorsByImei.get(device.imei);
					if (!serverError) continue;

					device.error = serverError;
					hasError = true;
				}
			} catch {
				/* Local parse result stays; import will still return server errors. */
			}
		}

		if (hasError) toast.error(i18next.t('message.import-error'));

		setDevices(uniqueDevices);

		return uniqueDevices;
	};

	const removeDevice = (imei: string) => {
		setDevices((prev) => prev.filter((d) => d.imei !== imei));
	};

	return { devices, setDevices, handleFileUpload, removeDevice };
};

const parseXlsx = async (file: File): Promise<DeviceFileParser[]> => {
	const data = await file.arrayBuffer();

	const workbook = XLSX.read(data, { type: 'array' });

	const sheetName = workbook.SheetNames[0];
	const sheet = workbook.Sheets[sheetName];

	const rows: any[] = XLSX.utils.sheet_to_json(sheet, {
		defval: '',
	});

	return rows.map((row) => {
		const imei = String(row.imei || row.IMEI || '').trim();
		const model = String(row.model || row.MODEL || '').trim();

		return {
			imei,
			model: model || null,
			error: undefined,
		};
	});
};
