import { useMemo } from 'react';
import { AnalyticQueryParams } from '@/rest/analyticAPI';
import { buildAnalyticRequest } from '@/utils/BuildAnalyticRequestUtils';
import { useAnalyticsFilters } from '@/context/useAnalyticsFiltersContext';

export const useAnalyticRequest = (type: string, query?: string): AnalyticQueryParams | null => {
	const { dateFilter, accountFilter, groupFilter, modelFilter } = useAnalyticsFilters();

	return useMemo(
		() =>
			buildAnalyticRequest(type, {
				query,
				from: dateFilter?.from,
				to: dateFilter?.to,
				user_uuid: accountFilter,
				group_id: groupFilter,
				model: modelFilter,
			}),
		[type, query, dateFilter?.from, dateFilter?.to, accountFilter, groupFilter, modelFilter]
	);
};
