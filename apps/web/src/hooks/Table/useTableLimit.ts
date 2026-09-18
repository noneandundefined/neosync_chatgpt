import { CACHEKEYs_TABLE_LIMIT_EL } from '@/constants/CacheKeys.constants';
import { useEffect, useState } from 'react';

const getInitialLimit = (key: string, defaultLimit: number): number => {
	if (typeof window === 'undefined') return defaultLimit;

	const raw = localStorage.getItem(key);
	const value = Number(raw);

	return Number.isFinite(value) && value > 0 ? value : defaultLimit;
};

export const useTableLimit = (key: string, defaultLimit = 25) => {
	const localKey = CACHEKEYs_TABLE_LIMIT_EL(key);

	const [limit, setLimit] = useState<number>(() => getInitialLimit(localKey, defaultLimit));

	useEffect(() => {
		localStorage.setItem(localKey, String(limit));
	}, [limit, key]);

	return { limit, setLimit };
};
