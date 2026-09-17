import { useEffect, useState } from 'react';

export const useEnabled = (prop: boolean) => {
	const [enabled, setEnabled] = useState<boolean>(prop);

	useEffect(() => {
		setEnabled(prop);
	}, [prop]);

	const toggle = () => setEnabled((prev) => !prev);

	return { enabled, toggle };
};
