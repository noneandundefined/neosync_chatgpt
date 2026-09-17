import React, { createContext, useContext, useState } from 'react';

type DateFilter = { from: string; to: string } | null;

interface AnalyticsFiltersState {
	dateFilter: DateFilter;
	setDateFilter: (v: DateFilter) => void;

	accountFilter: string | null;
	setAccountFilter: (v: string | null) => void;

	groupFilter: number | null;
	setGroupFilter: (v: number | null) => void;

	modelFilter: string | null;
	setModelFilter: (v: string | null) => void;

	reset: () => void;
}

const AnalyticsFiltersContext = createContext<AnalyticsFiltersState | null>(null);

export const useAnalyticsFilters = () => {
	const ctx = useContext(AnalyticsFiltersContext);
	if (!ctx) throw new Error('useAnalyticsFilters outside provider');

	return ctx;
};

export const AnalyticsFiltersProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const [dateFilter, setDateFilter] = useState<DateFilter>(null);
	const [accountFilter, setAccountFilter] = useState<string | null>(null);
	const [groupFilter, setGroupFilter] = useState<number | null>(null);
	const [modelFilter, setModelFilter] = useState<string | null>(null);

	const reset = () => {
		setDateFilter(null);
		setAccountFilter(null);
		setGroupFilter(null);
		setModelFilter(null);
	};

	return (
		<AnalyticsFiltersContext.Provider
			value={{
				dateFilter,
				setDateFilter,
				accountFilter,
				setAccountFilter,
				groupFilter,
				setGroupFilter,
				modelFilter,
				setModelFilter,
				reset,
			}}
		>
			{children}
		</AnalyticsFiltersContext.Provider>
	);
};
