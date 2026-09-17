import { config as configClient } from '@/.config/config.client';

export type ProductAnalyticsProperties = Record<string, string | number | boolean | null | undefined>;

export interface ProductAnalyticsEvent {
	event_id: string;
	occurred_at: string;
	session_id: string;
	event_name: string;
	category: string;
	path?: string;
	entity_type?: string;
	entity_id?: string;
	device_id?: number;
	success?: boolean;
	duration_ms?: number;
	error_code?: string;
	app_version?: string;
	properties?: ProductAnalyticsProperties;
}

const FLUSH_INTERVAL = 15000;
const MAX_BATCH_SIZE = 20;
const SESSION_KEY = 'neosync_analytics_session_id';
const APP_VERSION = '1.1.0';

let queue: ProductAnalyticsEvent[] = [];
let flushTimer: number | null = null;
const recentSuccessfulEvents: Array<{ name: string; at: number }> = [];

const analyticsBaseUrl = () => {
	if (configClient.type.release == 'dev') return configClient.links.URL_BACKEND_DEV;
	if (configClient.type.release == 'prod') return configClient.links.URL_BACKEND_PROD;
	return configClient.links.URL_BACKEND_PREPROD;
};

const createEventID = () => {
	if (typeof crypto !== 'undefined' && crypto.randomUUID) return crypto.randomUUID();
	return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
};

export const getAnalyticsSessionID = () => {
	let sessionID = sessionStorage.getItem(SESSION_KEY);
	if (!sessionID) {
		sessionID = createEventID();
		sessionStorage.setItem(SESSION_KEY, sessionID);
	}
	return sessionID;
};

const cleanProperties = (properties?: ProductAnalyticsProperties) => {
	if (!properties) return undefined;

	return Object.fromEntries(
		Object.entries(properties)
			.filter(([, value]) => ['string', 'number', 'boolean'].includes(typeof value) || value === null)
			.map(([key, value]) => [key.slice(0, 64), typeof value === 'string' ? value.slice(0, 256) : value])
	);
};

export const flushProductAnalytics = async () => {
	if (!queue.length) return;

	const events = queue.splice(0, MAX_BATCH_SIZE);
	let retryDelay = 1000;
	const publicPaths = ['/sign_in', '/how_get_access', '/password_reset', '/password/new'];
	const publicEvents = events.filter(
		(event) =>
			event.event_name === 'session_started' ||
			(publicPaths.some((path) => event.path?.startsWith(path)) && event.event_name !== 'login_success' && event.event_name !== 'password_reset_success')
	);
	const authenticatedEvents = events.filter((event) => !publicEvents.includes(event));
	const send = async (batch: ProductAnalyticsEvent[], endpoint: string) => {
		if (!batch.length) return;
		try {
			const response = await fetch(`${analyticsBaseUrl()}${endpoint}`, {
				method: 'POST',
				credentials: 'include',
				keepalive: true,
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ events: batch }),
			});

			if (!response.ok) {
				queue = [...batch, ...queue].slice(0, 200);
				if (response.status === 401) retryDelay = 30000;
			}
		} catch {
			queue = [...batch, ...queue].slice(0, 200);
		}
	};

	await Promise.all([send(publicEvents, '/analytics/public-events'), send(authenticatedEvents, '/analytics/events')]);

	if (queue.length) scheduleFlush(retryDelay);
};

const scheduleFlush = (delay = FLUSH_INTERVAL) => {
	if (flushTimer !== null) return;
	flushTimer = window.setTimeout(() => {
		flushTimer = null;
		void flushProductAnalytics();
	}, delay);
};

export const trackProductEvent = (eventName: string, category: string, options: Partial<Omit<ProductAnalyticsEvent, 'event_id' | 'occurred_at' | 'session_id' | 'event_name' | 'category'>> = {}) => {
	if (options.success === true) {
		recentSuccessfulEvents.push({ name: eventName, at: Date.now() });
		if (recentSuccessfulEvents.length > 50) recentSuccessfulEvents.shift();
	}
	queue.push({
		event_id: createEventID(),
		occurred_at: new Date().toISOString(),
		session_id: getAnalyticsSessionID(),
		event_name: eventName,
		category,
		path: window.location.pathname,
		app_version: APP_VERSION,
		...options,
		properties: cleanProperties(options.properties),
	});

	if (queue.length >= MAX_BATCH_SIZE) void flushProductAnalytics();
	else scheduleFlush();
};

export const hasSuccessfulProductEventSince = (eventName: string, timestamp: number) => recentSuccessfulEvents.some((event) => event.at >= timestamp && event.name === eventName);

const normalizedApiPath = (url = '') =>
	url
		.replace(/\b\d{15}\b/g, ':imei')
		.replace(/\/\d+(?=\/|$)/g, '/:id')
		.split('?')[0];

const resolveApiEvent = (method: string, url: string) => {
	const path = normalizedApiPath(url);
	if (path.includes('/auth/signin')) return ['login', 'authentication'];
	if (path.includes('/auth/check')) return ['session_authenticated', 'session'];
	if (path.includes('/password/reset')) return ['password_reset', 'authentication'];
	if (path.includes('/configuration/apply-template')) return ['template_apply', 'configuration'];
	if (path.includes('/configuration/apply')) return ['configuration_apply', 'configuration'];
	if (path.includes('/configuration/import')) return ['configuration_import', 'configuration'];
	if (path.includes('/configuration/draft') && method !== 'get') return ['configuration_draft_save', 'configuration'];
	if (path.includes('/configuration-templates')) return [`template_${method}`, 'configuration'];
	if (path.includes('/firmware')) return ['firmware_update', 'firmware'];
	if (path.includes('/cmd')) return ['command_send', 'command'];
	if (path.includes('/companies')) return ['bulk_configuration', 'configuration'];
	if (path.includes('/devices/import')) return ['device_import', 'device'];
	if (path === '/devices' && method === 'post') return ['device_create', 'device'];
	if (path.includes('/groups') && method !== 'get') return [`group_${method}`, 'group'];
	return method === 'get' ? null : ['api_mutation', 'api'];
};

export const trackApiOutcome = (method: string, url: string, success: boolean, durationMs: number, status?: number, errorCode?: string) => {
	if (url.includes('/analytics/events')) return;
	const resolved = resolveApiEvent(method, url);
	if (!resolved) return;

	const [baseName, category] = resolved;
	const deviceImei = url.match(/\b\d{15}\b/)?.[0];
	trackProductEvent(`${baseName}_${success ? 'success' : 'failed'}`, category, {
		success,
		duration_ms: Math.max(0, Math.round(durationMs)),
		error_code: errorCode,
		entity_type: deviceImei ? 'device' : undefined,
		entity_id: deviceImei,
		properties: { api_path: normalizedApiPath(url), method, status: status ?? 0 },
	});
};

export const trackSearchOutcome = (url: string, resultCount: number) => {
	if (!/[?&]search=[^&]+/i.test(url)) return;
	trackProductEvent(resultCount === 0 ? 'search_zero_results' : 'search_used', 'search', {
		properties: { api_path: normalizedApiPath(url), result_count: Math.max(0, resultCount) },
	});
};
