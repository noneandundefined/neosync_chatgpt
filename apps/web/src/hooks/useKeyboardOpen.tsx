import { useEffect, useState } from 'react';

export const useKeyboardOpen = () => {
	const [keyboardOpen, setKeyboardOpen] = useState(false);

	useEffect(() => {
		const viewport = window.visualViewport;
		if (!viewport) return;

		const onResize = () => {
			const heightDiff = window.innerHeight - viewport.height;
			setKeyboardOpen(heightDiff > 100);
		};

		viewport.addEventListener('resize', onResize);

		return () => {
			viewport.removeEventListener('resize', onResize);
		};
	}, []);

	return keyboardOpen;
};
