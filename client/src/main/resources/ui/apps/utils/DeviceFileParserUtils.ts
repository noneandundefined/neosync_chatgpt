import i18next from 'i18next';
import { toast } from 'react-toastify';
import { DeviceFileParser } from '@/interface/device/deviceFileParser.interface';

// <imei>,<model?>
export const parseDevicesFromFile = (content: string): DeviceFileParser[] => {
	const lines = content
		.split('\n')
		.map((line) => line.trim())
		.filter(Boolean)
		.filter((line) => !/^sep\s*=\s*[,;]$/i.test(line))
		.filter((line) => !/^imei(\s*,\s*model)?$/i.test(line));

	return lines.map(parseLine);
};

const parseLine = (line: string): DeviceFileParser => {
	const parts = line.split(',').map((part) => part.trim());

	const imei = normalizeImei(parts[0] || '');
	if (!imei) toast.error(i18next.t('message.error-device-imports-imei'));

	let model: null | string = null;

	if (parts.length >= 2) {
		model = parts[1];
	}

	return { imei, model };
};

const normalizeImei = (raw: string): string => {
	let value = raw.trim();

	if (!value) return '';

	const excelFormula = value.match(/^=\s*"(.+)"$/);
	if (excelFormula) {
		value = excelFormula[1].trim();
	}

	value = value.replace(/^["']+|["']+$/g, '').trim();

	if (/^[+-]?\d+(\.\d+)?e[+-]?\d+$/i.test(value)) {
		const asNumber = Number(value);
		if (Number.isFinite(asNumber)) {
			value = asNumber.toFixed(0);
		}
	}

	return value;
};
