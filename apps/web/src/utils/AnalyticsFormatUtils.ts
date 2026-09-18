import { TFunction } from 'i18next';
import { WidgetResponse } from '@/rest/analyticAPI';

const routeLabels: Array<[RegExp, string]> = [
	[/^\/$/, 'analytics-route-devices'],
	[/^\/sign_in/, 'analytics-route-sign-in'],
	[/^\/dashboard\/users/, 'analytics-route-users'],
	[/^\/devices\/command/, 'analytics-route-commands'],
	[/^\/devices\/[^/]+\/logs/, 'analytics-route-device-logs'],
	[/^\/configurations\/devices\//, 'analytics-route-device-configuration'],
	[/^\/configurations\/bulk/, 'analytics-route-bulk-configuration'],
	[/^\/configurations/, 'analytics-route-configurations'],
	[/^\/groups/, 'analytics-route-groups'],
	[/^\/firmwares/, 'analytics-route-firmwares'],
	[/^\/neosync\/analytics/, 'analytics-route-analytics'],
	[/^\/neosync\/help/, 'analytics-route-help'],
	[/^\/neosync\/releases/, 'analytics-route-releases'],
];

const apiLabels: Array<[string, string]> = [
	['/auth/signin', 'analytics-api-sign-in'],
	['/auth/check', 'analytics-api-session-check'],
	['/devices/:imei/configuration/apply-template/:id', 'analytics-api-template-apply'],
	['/devices/:imei/configuration/apply', 'analytics-api-configuration-apply'],
	['/devices/:imei/configuration/import', 'analytics-api-configuration-import'],
	['/devices/:imei/configuration/draft/:id', 'analytics-api-draft-save'],
	['/devices/:imei/firmware', 'analytics-api-firmware-update'],
	['/devices/:imei/cmd/send-and-wait', 'analytics-api-command-send'],
	['/devices/cmd', 'analytics-api-command-send'],
	['/devices/import', 'analytics-api-device-import'],
	['/devices', 'analytics-api-devices'],
	['/companies', 'analytics-api-bulk-configuration'],
	['/groups', 'analytics-api-groups'],
];

const translate = (t: TFunction, key: string, fallback: string) => t(`label.${key}`, { defaultValue: fallback });
const readableFallback = (value: string) => value.replace(/[_-]+/g, ' ').replace(/^./, (character) => character.toUpperCase());

export const getAnalyticsErrorReasonLabel = (t: TFunction, reason: string) => {
	const lowered = reason.toLocaleLowerCase();
	if (lowered.includes('tracker offline')) return translate(t, 'analytics-reason-tracker_offline', 'Device is offline');
	if (lowered.includes('timeout')) return translate(t, 'analytics-reason-delivery_timeout', 'Delivery took too long');
	if (lowered.includes('firmware')) return translate(t, 'analytics-reason-unsupported_firmware', 'Firmware is not supported');
	return translate(t, 'analytics-reason-config_error', 'Configuration error');
};

export const getAnalyticsWidgetTitle = (t: TFunction, widget: WidgetResponse) => (widget.title_key ? translate(t, widget.title_key, widget.title) : widget.title);

export const getAnalyticsColumnLabel = (t: TFunction, column: string) => translate(t, `analytics-column-${column}`, readableFallback(column));

export const getAnalyticsRouteLabel = (t: TFunction, path: string) => {
	const route = routeLabels.find(([pattern]) => pattern.test(path));
	return route ? translate(t, route[1], path) : path;
};

export const getAnalyticsApiLabel = (t: TFunction, path: string) => {
	const endpoint = apiLabels.find(([value]) => path.includes(value));
	return endpoint ? translate(t, endpoint[1], path) : translate(t, 'analytics-api-other', 'Other operation');
};

export const formatAnalyticsValue = (t: TFunction, language: string, column: string, value: unknown) => {
	if (value === null || value === undefined || value === '') return '—';
	if (typeof value === 'number') return value.toLocaleString(language);

	const text = String(value);
	if (['page', 'path', 'landing_page'].includes(column)) return getAnalyticsRouteLabel(t, text);
	if (column === 'api_path') return getAnalyticsApiLabel(t, text);
	if (column === 'field') return translate(t, text, readableFallback(text));
	if (column === 'control') return translate(t, `analytics-value-${text.toLocaleLowerCase()}`, text);
	if (column === 'error' || column === 'error_code') {
		const translated = t(`label.analytics-value-${text.toLocaleLowerCase()}`, { defaultValue: '' });
		return translated || translate(t, 'analytics-error-other', 'The action could not be completed');
	}
	if (column === 'cohort' || column === 'date' || column.endsWith('_time')) {
		const date = new Date(text);
		if (!Number.isNaN(date.getTime())) return new Intl.DateTimeFormat(language, { dateStyle: 'medium', timeStyle: column.endsWith('_time') ? 'short' : undefined }).format(date);
	}
	if (['event', 'event_name', 'failed_event', 'category', 'status', 'role', 'step', 'input_type', 'method', 'source', 'medium'].includes(column)) {
		return translate(t, `analytics-value-${text.toLocaleLowerCase()}`, readableFallback(text));
	}

	return text;
};
