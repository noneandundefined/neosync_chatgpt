export interface FuelAddressesResult {
	macs: { value: number[]; position: number }[];
}

type SensorType = 'RS485' | 'BLE';

export interface FuelSensorTypes {
	rs485: number[];
	ble: number[];
	all: {
		index: number;
		type: SensorType;
	}[];
}

export const FuelAddressesParse = (arr: number[]): FuelAddressesResult => {
	const macs: { value: number[]; position: number }[] = [];

	for (let i = 0; i + 6 <= arr.length; i += 6) {
		const slice6 = arr.slice(i, i + 6);
		const nonZeroCount = slice6.filter((b) => b !== 0).length;
		if (nonZeroCount >= 4) {
			macs.push({ value: slice6, position: i });
		}
	}

	return { macs };
};

export const FuelSensorTypesParse = (mask: number, fuelSensorsCount = 8): FuelSensorTypes => {
	const rs485: number[] = [];
	const ble: number[] = [];
	const all: {
		index: number;
		type: SensorType;
	}[] = [];

	const binary = mask.toString(2).padStart(fuelSensorsCount, '0');

	for (let i = fuelSensorsCount - 1; i >= 0; i--) {
		const bit = binary[i];
		const type: SensorType = bit === '1' ? 'BLE' : 'RS485';
		const index = fuelSensorsCount - 1 - i;

		all.push({ index, type });

		if (type === 'RS485') rs485.push(index);
		else ble.push(index);
	}

	return { rs485, ble, all };
};
