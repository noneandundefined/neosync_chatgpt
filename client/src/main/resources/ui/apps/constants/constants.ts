export const PLACEHOLDER_ROUTE_ID = ':id';
export const PLACEHOLDER_ROUTE_DEVICE_IMEI = ':imei';
export const PLACEHOLDER_ROUTE_DEVICE_MODEL = ':model';

export const ROUTES = {
	HOME: '/',
	NOT_FOUND: '*',
	FORBIDDEN: '/403',
	UNAUTHORISED: '/unauthorised',
	SIGNUP: '/sign_up',
	SIGNIN: '/sign_in',
	HOW_GET_ACCESS: '/how_get_access',
	RESET_PASSWORD_REQ: '/password_reset',
	RESET_PASSWORD_NEW: '/password/new',

	DEVICE_SEND_COMMAND: '/devices/command',
	DEVICE_DETAILS_COMMAND: `/devices/command/${PLACEHOLDER_ROUTE_ID}`,
	DEVICE_LOGS: `/devices/${PLACEHOLDER_ROUTE_DEVICE_IMEI}/logs`,

	CONFIGURATIONS: `/configurations`,
	CONFIGURATIONS_BULK_NEW: `/configurations/bulk/new`,
	CONFIGURATIONS_BULK_DETAILS: `/configurations/bulk/${PLACEHOLDER_ROUTE_ID}`,
	CONFIGURATIONS_TEMPLATE_NEW: `/configurations/devices/template`,
	CONFIGURATIONS_TEMPLATE_EDIT: `/configurations/devices/templates/${PLACEHOLDER_ROUTE_ID}`,
	CONFIGURATIONS_DEVICE: `/configurations/devices/${PLACEHOLDER_ROUTE_DEVICE_IMEI}`,

	GROUPS: '/groups',

	FIRMWARES: '/firmwares',

	NEOSYNC_RELEASES: '/neosync/releases',
	NEOSYNC_HELP: '/neosync/help',
	NEOSYNC_ANALYTICS: '/neosync/analytics',

	DASHBOARD_USERS: '/dashboard/users',
};

type Params = Record<string, string | number>;

export const buildRoute = (route: string, params: Params): string => {
	let result = route;

	Object.entries(params).forEach(([key, value]) => {
		result = result.replace(`:${key}`, encodeURIComponent(String(value)));
	});

	return result;
};
