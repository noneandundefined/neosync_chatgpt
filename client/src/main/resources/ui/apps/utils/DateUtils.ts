import { startOfDay, endOfDay, subDays, startOfMonth, endOfMonth } from 'date-fns';

export const formatAnalytic = (date: Date) => date.toLocaleDateString('ru-RU');

export const getPeriods = () => {
	const now = new Date();

	return [
		{
			label: 'last-14',
			from: startOfDay(subDays(now, 13)),
			to: endOfDay(now),
		},
		{
			label: 'today',
			from: startOfDay(now),
			to: endOfDay(now),
		},
		{
			label: 'yesterday',
			from: startOfDay(subDays(now, 1)),
			to: endOfDay(subDays(now, 1)),
		},
		{
			label: 'last-7',
			from: startOfDay(subDays(now, 6)),
			to: endOfDay(now),
		},
		{
			label: 'last-30',
			from: startOfDay(subDays(now, 29)),
			to: endOfDay(now),
		},
		{
			label: 'this-month',
			from: startOfMonth(now),
			to: endOfDay(now),
		},
		{
			label: 'last-month',
			from: startOfMonth(subDays(startOfMonth(now), 1)),
			to: endOfMonth(subDays(startOfMonth(now), 1)),
		},
	];
};
