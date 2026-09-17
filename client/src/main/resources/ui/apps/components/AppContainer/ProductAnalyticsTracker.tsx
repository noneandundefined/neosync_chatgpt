import { useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import { flushProductAnalytics, hasSuccessfulProductEventSince, trackProductEvent } from '@/utils/ProductAnalytics';

const ProductAnalyticsTracker = () => {
	const location = useLocation();
	const pageStartedAt = useRef(Date.now());
	const previousPath = useRef(location.pathname);

	useEffect(() => {
		const query = new URLSearchParams(window.location.search);
		trackProductEvent('session_started', 'session', {
			properties: {
				language: navigator.language,
				viewport: `${window.innerWidth}x${window.innerHeight}`,
				referrer_host: document.referrer ? new URL(document.referrer).host : '',
				utm_source: query.get('utm_source') || '',
				utm_medium: query.get('utm_medium') || '',
				utm_campaign: query.get('utm_campaign') || '',
			},
		});

		const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming | undefined;
		if (navigation) {
			trackProductEvent('page_performance', 'performance', {
				duration_ms: Math.round(navigation.loadEventEnd || navigation.duration),
				properties: {
					dom_ms: Math.round(navigation.domContentLoadedEventEnd),
					ttfb_ms: Math.round(navigation.responseStart),
				},
			});
		}
	}, []);

	useEffect(() => {
		const now = Date.now();
		if (previousPath.current !== location.pathname) {
			if (previousPath.current.startsWith('/configurations') && !hasSuccessfulProductEventSince('configuration_apply_success', pageStartedAt.current)) {
				trackProductEvent('configuration_abandoned', 'configuration', { success: false, error_code: 'navigation_away', path: previousPath.current, duration_ms: now - pageStartedAt.current });
			}
			trackProductEvent('page_duration', 'navigation', {
				path: previousPath.current,
				duration_ms: now - pageStartedAt.current,
			});
		}

		previousPath.current = location.pathname;
		pageStartedAt.current = now;
		trackProductEvent('page_view', 'navigation', {
			properties: { query_present: Boolean(location.search) },
		});
	}, [location.pathname, location.search]);

	useEffect(() => {
		const startedForms = new WeakSet<HTMLFormElement>();
		const handleClick = (event: MouseEvent) => {
			const element = (event.target as HTMLElement | null)?.closest<HTMLElement>('button, a, [data-analytics-event]');
			if (!element) return;

			trackProductEvent(element.dataset.analyticsEvent || 'control_clicked', 'interaction', {
				properties: {
					control:
						element.dataset.analyticsControl ||
						element.getAttribute('aria-label') ||
						element.getAttribute('title') ||
						element.getAttribute('name') ||
						element.textContent?.trim().replace(/\s+/g, ' ').slice(0, 80) ||
						element.tagName.toLowerCase(),
				},
			});
		};
		const handleFocus = (event: FocusEvent) => {
			const form = (event.target as HTMLElement | null)?.closest<HTMLFormElement>('form');
			if (!form || startedForms.has(form)) return;
			startedForms.add(form);
			trackProductEvent('form_started', 'interaction', { properties: { form: form.id || form.getAttribute('name') || window.location.pathname } });
		};
		const handleSubmit = (event: SubmitEvent) => {
			const form = event.target as HTMLFormElement;
			trackProductEvent('form_submitted', 'interaction', { properties: { form: form.id || form.getAttribute('name') || window.location.pathname } });
		};

		const handleInvalid = (event: Event) => {
			const input = event.target as HTMLInputElement;
			trackProductEvent('form_validation_failed', 'error', {
				success: false,
				error_code: input.validity.valueMissing ? 'required' : input.validity.typeMismatch ? 'type_mismatch' : 'invalid',
				properties: { field: input.name || input.id || 'unknown', input_type: input.type || 'unknown' },
			});
		};

		const handleError = (event: ErrorEvent) => {
			trackProductEvent('client_error', 'error', { success: false, error_code: event.error?.name || 'javascript_error' });
		};

		const handleRejection = () => trackProductEvent('client_error', 'error', { success: false, error_code: 'unhandled_rejection' });
		const handleVisibility = () => {
			if (document.visibilityState === 'hidden') {
				trackProductEvent('page_duration', 'navigation', { duration_ms: Date.now() - pageStartedAt.current });
				void flushProductAnalytics();
			}
		};

		document.addEventListener('click', handleClick, true);
		document.addEventListener('focusin', handleFocus, true);
		document.addEventListener('submit', handleSubmit, true);
		document.addEventListener('invalid', handleInvalid, true);
		document.addEventListener('visibilitychange', handleVisibility);
		window.addEventListener('error', handleError);
		window.addEventListener('unhandledrejection', handleRejection);

		return () => {
			document.removeEventListener('click', handleClick, true);
			document.removeEventListener('focusin', handleFocus, true);
			document.removeEventListener('submit', handleSubmit, true);
			document.removeEventListener('invalid', handleInvalid, true);
			document.removeEventListener('visibilitychange', handleVisibility);
			window.removeEventListener('error', handleError);
			window.removeEventListener('unhandledrejection', handleRejection);
		};
	}, []);

	return null;
};

export default ProductAnalyticsTracker;
