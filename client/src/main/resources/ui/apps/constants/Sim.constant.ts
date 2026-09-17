export const OPERATORS = [
	{
		id: '',
		value: 'label.no',
		extra: {
			apn: '',
			login: '',
			password: '',
		},
	},
	{
		id: 'mts',
		value: 'MTS',
		extra: {
			apn: 'internet.mts.ru',
			login: 'mts',
			password: 'mts',
		},
	},
	{
		id: 'beeline',
		value: 'Beeline',
		extra: {
			apn: 'm2m.beeline.ru',
			login: 'beeline',
			password: 'beeline',
		},
	},
	{
		id: 'gdata',
		value: 'Megafon',
		extra: {
			apn: 'internet',
			login: 'gdata',
			password: 'gdata',
		},
	},
	{
		id: 'tele2',
		value: 'Tele2',
		extra: {
			apn: 'internet.tele2.ru',
			login: '',
			password: '',
		},
	},
	{
		id: 'emt',
		value: 'M2M Express',
		extra: {
			apn: 'internet.emt.ee',
			login: 'emt',
			password: 'emt',
		},
	},
];

export const detectOperatorId = (apn: string, login: string, password: string): string => {
	const exact = OPERATORS.find((op) => op.extra.apn === apn && op.extra.login === login && op.extra.password === password);
	if (exact) return String(exact.id);

	if (apn) {
		const byApn = OPERATORS.find((op) => op.extra.apn === apn);
		if (byApn) return String(byApn.id);
	}

	if (apn || login || password) {
		return 'custom';
	}

	return '';
};

export const SIM_PRIORITY = [
	{ id: 0, value: 'label.no-priority' },
	{ id: 1, value: 'SIM0' },
	{ id: 2, value: 'SIM1' },
];

export const NET_SIM = [
	{ id: 0, value: 'label.auto' },
	{ id: 1, value: '2G' },
	{ id: 2, value: 'LTE' },
];

export const MULTI_APN_TYPE = [
	{
		id: 0,
		value: 'message.multiapn-not-specified',
	},
	{
		id: 1,
		value: 'message.multiapn-periodic',
	},
	{
		id: 2,
		value: 'message.multiapn-in-parallel',
	},
];
