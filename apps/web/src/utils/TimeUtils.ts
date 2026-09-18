import { TFunction } from 'i18next';

export const formatRelativeTime = (iso: string, t: TFunction): string => {
	const diffMs = Date.now() - new Date(iso).getTime();
	const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));

	if (days <= 0) return t('label.today');
	if (days === 1) return t('label.yesterday');
	if (days < 30) return t('message.days-ago', { count: days });

	const months = Math.floor(days / 30);
	if (months < 12) return t('message.months-ago', { count: months });

	return t('message.years-ago', { count: Math.floor(months / 12) });
};

export const getDurationTime = (dateString: string) => {
	if (!dateString) return '-';

	const now = new Date();
	const target = new Date(dateString);

	let diff = Math.abs(now.getTime() - target.getTime());

	const days = Math.floor(diff / (1000 * 60 * 60 * 24));
	diff -= days * (1000 * 60 * 60 * 24);

	const hours = Math.floor(diff / (1000 * 60 * 60));
	diff -= hours * (1000 * 60 * 60);

	const minutes = Math.floor(diff / (1000 * 60));

	return `${days} д. ${hours} ч. ${minutes} мин.`;
};

export const formatRuDate = (date: Date): string => {
	return date
		.toLocaleString('ru-RU', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit',
		})
		.replace(',', '');
};

export const formatDate = (date: string): string => {
	return new Intl.DateTimeFormat('ru-RU', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
		hour12: false,
	}).format(new Date(date));
};

export const formatLocalDate = (iso: string): string => {
	if (!iso) return '-';

	const hasTimezone = /(?:Z|[+-]\d{2}:?\d{2})$/i.test(iso);
	const normalizedIso = hasTimezone ? iso : `${iso}Z`;
	const date = new Date(normalizedIso);
	if (Number.isNaN(date.getTime())) return '-';

	return date.toLocaleString('ru-RU', {
		day: '2-digit',
		month: '2-digit',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
		hour12: false,
	});
};

export const formatDiff = (isoDate: string): string => {
	const now = new Date();
	const target = new Date(isoDate);

	let diffMs = now.getTime() - target.getTime();
	if (diffMs < 0) diffMs = 0;

	const diffMinutes = Math.floor(diffMs / (1000 * 60));
	const days = Math.floor(diffMinutes / (60 * 24));
	const hours = Math.floor((diffMinutes % (60 * 24)) / 60);
	const minutes = diffMinutes % 60;

	return `${days} д. ${hours} ч. ${minutes} мин.`;
};
