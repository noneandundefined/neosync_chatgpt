export const RS485_SELECT_MODBUS_BYTES_COUNT = [
	{ id: 2, value: '2' },
	{ id: 4, value: '4' },
	{ id: 6, value: '6' },
	{ id: 8, value: '8' },
];

export const RS485_SELECT_MODBUS_BD_RATE = [
	{ id: 4800, value: '4800' },
	{ id: 9600, value: '9600' },
	{ id: 19200, value: '19200' },
	{ id: 38400, value: '38400' },
	{ id: 57600, value: '57600' },
	{ id: 115200, value: '115200' },
];

export const RS485_SELECT_MODBUS_BYTES_ORDER = [
	{ id: 2, value: '2Byte-LittleEndian' },
	{ id: 1, value: 'LittleEndian' },
	{ id: 0, value: 'BigEndian' },
];

export const RS485_SELECT_MODBUS_REGISTERS_TYPE = [
	{ id: 1, value: 'Colis' },
	{ id: 2, value: 'Discrete' },
	{ id: 3, value: 'Holding' },
	{ id: 4, value: 'Input' },
];
