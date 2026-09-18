const METRIKA_ID = 111693443;
const METRIKA_SCRIPT = `https://mc.yandex.ru/metrika/tag.js?id=${METRIKA_ID}`;

interface YandexMetrika {
	(...args: unknown[]): void;
	a?: unknown[];
	l?: number;
}

interface YandexMetrikaWindow extends Window {
	ym?: YandexMetrika;
}

export const initYandexMetrika = () => {
	const w = window as YandexMetrikaWindow;

	if (document.querySelector(`script[src="${METRIKA_SCRIPT}"]`)) {
		return;
	}

	w.ym =
		w.ym ||
		function (...args: unknown[]) {
			(w.ym!.a = w.ym!.a || []).push(args);
		};
	w.ym.l = Date.now();

	const script = document.createElement('script');
	script.async = true;
	script.src = METRIKA_SCRIPT;
	document.head.appendChild(script);

	w.ym(METRIKA_ID, 'init', {
		ssr: true,
		webvisor: true,
		clickmap: true,
		ecommerce: 'dataLayer',
		referrer: document.referrer,
		url: location.href,
		accurateTrackBounce: true,
		trackLinks: true,
	});
};
