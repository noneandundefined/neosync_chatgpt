import axiosClient from './axios';

const apiPath = '/analytics';

export interface WidgetResponse {
	id: number;
	type: string;
	section: string;
	title: string;
	title_key?: string;
	subtitle?: string;
	icon?: string;
	query?: string;
}

export interface AnalyticQueryResponse {
	columns?: string[];
	rows?: unknown[][];
}

export interface AnalyticQueryParams {
	type: string;
	query: string;
	from?: string;
	to?: string;
	user_uuid?: string;
	group_id?: number;
	model?: string;
	validate_only?: boolean;
}

/**
 * Запрос виджетов для аналитики
 */
export const basicAnalyticWidgets = async (): Promise<WidgetResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/widgets`);
	return response.data.message;
};

export const basicAnalyticQuery = async (params: AnalyticQueryParams): Promise<AnalyticQueryResponse> => {
	const response = await axiosClient.post(`${apiPath}/query`, params);
	return response.data.message;
};
