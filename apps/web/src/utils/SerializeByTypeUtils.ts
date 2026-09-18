import {
	FieldType,
	FieldType_Ascii,
	FieldType_ByteArray,
	FieldType_MacAddressArray,
	FieldType_Phones,
	FieldType_Raw,
	FieldType_SNSArray,
	FieldType_TwoUint16,
	FieldType_TwoUint8,
	FieldType_UInt16,
	FieldType_UInt32,
	FieldType_Uint32Array,
	FieldType_UInt8,
} from '@/constants/FieldType.constant';

export const SerializeByType = (type: FieldType, value: any): any => {
	try {
		if (typeof value === 'string' && /^\[.*\]$/.test(value.trim())) {
			value = JSON.parse(value);
		}

		switch (type) {
			case FieldType_UInt8:
			case FieldType_UInt16:
			case FieldType_UInt32:
				return Number(value);

			case FieldType_TwoUint8:
			case FieldType_TwoUint16: {
				let arr: number[] = Array.isArray(value) ? value.map(Number) : [Number(value)];
				arr = arr.filter((v) => !isNaN(v));
				return arr;
			}

			case FieldType_Uint32Array:
			case FieldType_ByteArray: {
				let arr: number[] = [];

				if (Array.isArray(value)) {
					arr = value.map((v) => Number(v));
				} else if (typeof value === 'string' && /^\[.*\]$/.test(value.trim())) {
					try {
						const parsed = JSON.parse(value);
						arr = Array.isArray(parsed) ? parsed.map((v: any) => Number(v)) : [Number(parsed)];
					} catch {
						arr = [];
					}
				} else {
					arr = [Number(value)];
				}

				return arr.filter((v) => !isNaN(v));
			}

			case FieldType_Phones:
				return Array.isArray(value) ? value : String(value);

			case FieldType_Ascii:
				if (Array.isArray(value)) return value.map((v) => String(v));
				return String(value);

			case FieldType_MacAddressArray:
			case FieldType_SNSArray:
			case FieldType_Raw:
			default:
				return value;
		}
	} catch {
		return value;
	}
};
