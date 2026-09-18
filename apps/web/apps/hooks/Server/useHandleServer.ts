import { QueryKey, QueryObserverResult, RefetchOptions, useQuery, useQueryClient, UseQueryOptions } from '@tanstack/react-query';

// Общий тип, определяющий контракт возврата нашего хука.
type UseHandleServerType<T, K extends string = 'data'> = {
	[P in K]: T | null;
} & {
	loading: boolean;
	reload: (options?: RefetchOptions) => Promise<QueryObserverResult<T | null, Error>>; // Обновляем тип reload
	updateHServer: (updatedFields: Partial<T>) => void;
};

// Универсальный хук для отправки запросов к серверу и отслеживания состояния загрузки.
export const useHandleServer = <T>(queryKey: QueryKey, fn: (signal?: AbortSignal) => Promise<T>, options?: Omit<UseQueryOptions<T, Error, T, QueryKey>, 'queryKey' | 'queryFn'>) => {
	const queryClient = useQueryClient();

	const {
		data: queryData,
		isLoading,
		refetch,
	} = useQuery({
		queryKey: queryKey,
		queryFn: ({ signal }) => fn(signal),
		...options,
	});

	const update = (updatedFields: Partial<T>) => {
		queryClient.setQueryData(queryKey, (old: T | undefined) => ({
			...old,
			...updatedFields,
		}));
	};

	return {
		data: queryData,
		loading: isLoading,
		reload: refetch,
		updateHServer: update,
	} as UseHandleServerType<T>;
};
