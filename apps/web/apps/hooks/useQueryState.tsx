import { useEffect, useState } from 'react';
import { useLocation, useSearchParams } from 'react-router-dom';

function useQueryState<T = string>(key: string, defaultValue: T, enabled = true): [T, (value: T) => void] {
	const [searchParams, setSearchParams] = useSearchParams();
	const location = useLocation();

	const initialValue = enabled ? ((searchParams.get(key) as T) ?? defaultValue) : defaultValue;
	const [state, setState] = useState<T>(initialValue);

	useEffect(() => {
		if (!enabled) return;

		setSearchParams(
			(prev) => {
				const next = new URLSearchParams(prev);

				if (state === defaultValue || state === '' || state === null) {
					next.delete(key);
				} else {
					next.set(key, String(state));
				}

				return next;
			},
			{ replace: true, state: location.state }
		);
	}, [enabled, state, key, defaultValue, setSearchParams, location.state]);

	return [state, setState];
}

export default useQueryState;
