import { AnalyticQueryParams } from '@/rest/analyticAPI';

type AnalyticsFilters = {
	from?: string;
	to?: string;
	user_uuid?: string | null;
	group_id?: number | null;
	model?: string | null;
	query?: string;
};

export const buildAnalyticRequest = (type: string, options: AnalyticsFilters): AnalyticQueryParams | null => {
	if (!options.query?.trim()) {
		return null;
	}

	const base: AnalyticQueryParams = {
		type,
		query: options.query,
	};

	if (type === 'stat_live' || type === 'sql') {
		return {
			...base,
			user_uuid: options.user_uuid || undefined,
			group_id: options.group_id && options.group_id > 0 ? options.group_id : undefined,
			model: options.model || undefined,
		};
	}

	if (!options.from || !options.to) {
		return null;
	}

	return {
		...base,
		from: options.from,
		to: options.to,
		user_uuid: options.user_uuid || undefined,
		group_id: options.group_id && options.group_id > 0 ? options.group_id : undefined,
		model: options.model || undefined,
	};
};
