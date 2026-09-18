import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import axiosClient from './axios';

const apiPath = '/meta';

export interface Releases {
	releases: Release[];
}

export interface Release {
	version: string;
	date: string;
	changes: Record<string, ChangeSection>;
}

export interface ChangeSection {
	added: string[];
	fixed: string[];
}

export interface DatabaseColumnResponse {
	column_name: string;
	data_type: string;
	is_nullable: string;
	column_default: string | null;
}

export const basicMetaTcpHealth = async (): Promise<boolean> => {
	const response = await axiosClient.get(`${apiPath}/health/tcp`);
	return response.data.message.status;
};

export const basicMetaGetReleases = async (): Promise<Releases> => {
	const response = await axiosClient.get(`${apiPath}/s/releases`);
	return response.data.message;
};

export const basicMetaGetVersion = async (): Promise<any> => {
	const versionCache = sessionStorage.getItem(CACHEKEYs.NEOSYNC_VERSION);
	if (versionCache) {
		return versionCache;
	}

	const response = await axiosClient.get(`${apiPath}/s/version`);

	sessionStorage.setItem(CACHEKEYs.NEOSYNC_VERSION, response.data.message);
	return response.data.message;
};

/**
 * Получение колонок всех таблиц
 */
export const basicMetaGetColumns = async (table: string): Promise<DatabaseColumnResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/schema?table=${table}`);
	return response.data.message;
};
