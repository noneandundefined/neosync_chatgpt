import { useEffect } from 'react';
import { config } from '@/.config/config.client';

interface AntiDebugProps {
	redirectUrl?: string;
	checkInterval?: number;
	enableInDev?: boolean;
}

export const useAntiDebug = ({ redirectUrl = '/', checkInterval = 1000, enableInDev = false }: AntiDebugProps = {}) => {
	useEffect(() => {
		const isProd = config.type.release === 'prod';

		if (!isProd && !enableInDev) {
			return;
		}

		const devtoolsDetector = new Image();
		Object.defineProperty(devtoolsDetector, 'id', {
			get() {
				window.location.href = redirectUrl;
			},
		});

		const sizeCheck = setInterval(() => {
			if (window.outerWidth - window.innerWidth > 160 || window.outerHeight - window.innerHeight > 160) {
				window.location.href = redirectUrl;
			}
		}, checkInterval);

		const keyHandler = (e: KeyboardEvent) => {
			if (e.key === 'F12' || (e.ctrlKey && e.shiftKey && ['I', 'J', 'C'].includes(e.key)) || (e.ctrlKey && e.key === 'U')) {
				e.preventDefault();
				return false;
			}
		};

		document.addEventListener('keydown', keyHandler);

		const debuggerTrap = setInterval(() => {
			debugger;
		}, 500);

		return () => {
			clearInterval(sizeCheck);
			clearInterval(debuggerTrap);
			document.removeEventListener('keydown', keyHandler);
		};
	}, [redirectUrl, checkInterval, enableInDev]);
};
