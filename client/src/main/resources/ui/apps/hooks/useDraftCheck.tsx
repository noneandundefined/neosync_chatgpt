import { createSSE } from '@/rest/sse';
import { useEffect, useState } from 'react';

export const EVENT_CONFIGURATION_DRAFT = 'neosync:configuration:draft:ch';

export const useDraftCheck = (imei: string) => {
	const [isChanged, setIsChanged] = useState(false);

	useEffect(() => {
		if (!imei) return;

		const url = `/devices/${imei}/configuration/draft/check`;
		const sse = createSSE(url);

		sse.onmessage = (e) => {
			if (e.data === 'true') setIsChanged(true);
			else if (e.data === 'false') setIsChanged(false);
			else setIsChanged(true);
		};

		sse.onerror = () => {
			sse.close();
		};

		return () => {
			sse.close();
		};
	}, [imei]);

	useEffect(() => {
		const handler = (e: Event) => {
			const detail = (e as CustomEvent<boolean>).detail;
			if (typeof detail === 'boolean') {
				setIsChanged(detail);
			} else {
				setIsChanged(true);
			}
		};

		window.addEventListener(EVENT_CONFIGURATION_DRAFT, handler as EventListener);

		return () => {
			window.removeEventListener(EVENT_CONFIGURATION_DRAFT, handler as EventListener);
		};
	}, []);

	return isChanged;
};
