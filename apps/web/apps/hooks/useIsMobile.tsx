import { useEffect, useState } from 'react';

export const useIsMobile = (breakpoint = 800) => {
	const [isMobile, setIsMobile] = useState(false);

	useEffect(() => {
		const media = window.matchMedia(`(max-width: ${breakpoint}px)`);
		const update = () => setIsMobile(media.matches);

		update();
		media.addEventListener('change', update);
		return () => media.removeEventListener('change', update);
	}, [breakpoint]);

	return isMobile;
};

export default useIsMobile;
