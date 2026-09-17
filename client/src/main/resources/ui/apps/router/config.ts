import { ROUTES } from '@/constants/constants';
import { lazy, LazyExoticComponent } from 'react';

/* Auth */
const AuthSignin = lazy(() => import('@/pages/SignInPage/index'));
const AuthReqResetPassword = lazy(() => import('@/pages/ReqResetPasswordPage/index'));
const AuthResetPassword = lazy(() => import('@/pages/ResetPasswordPage/index'));
const AuthHowGetAccess = lazy(() => import('@/pages/HowGetAccessPage/index'));

/* Device */
const DeviceCommands = lazy(() => import('@/pages/DeviceCommandDashboardPage/DeviceCommandPage/index'));
const DeviceCommandDetails = lazy(() => import('@/pages/DeviceCommandDashboardPage/DeviceCommandDetailsPage/index'));
const DeviceLogs = lazy(() => import('@/pages/DeviceLogsPage/index'));

/* Configurations */
const Configurations = lazy(() => import('@/pages/ConfigurationsPage/index'));
const BulkConfiguration = lazy(() => import('@/pages/ConfigurationsPage/BulkConfigurationPage/index'));
const BulkConfigurationDetails = lazy(() => import('@/pages/ConfigurationsPage/BulkConfigurationDetailsPage/index'));
const ConfigurationTemplate = lazy(() => import('@/pages/ConfigurationsPage/ConfigurationTemplatePage/index'));
const ConfigurationTemplateEdit = lazy(() => import('@/pages/ConfigurationsPage/ConfigurationTemplatePage/ConfigurationTemplateEditPage'));
const DeviceConfiguration = lazy(() => import('@/pages/ConfigurationsPage/DeviceConfigurationPage/index'));

/* Groups */
const Group = lazy(() => import('@/pages/GroupPage/index'));

/* Firmwares */
const Firmwares = lazy(() => import('@/pages/FirmwarePage/index'));

/* Dashboard */
const DashboardUsers = lazy(() => import('@/pages/UsersDashboardPage/index'));

/* Neosync */
const NeosyncReleases = lazy(() => import('@/pages/NeosyncReleasesPage/index'));
const NeosyncHelp = lazy(() => import('@/pages/NeosyncHelpPage/index'));
const NeosyncAnalytics = lazy(() => import('@/pages/NeosyncAnalyticsPage/index'));

const NotFound = lazy(() => import('@/pages/NotFoundPage/index'));
const Home = lazy(() => import('@/pages/ApplicationPage/index'));

export interface CustomRouteConfig {
	path: string;
	title?: string;
	loginRequired?: boolean;
	redirectIfLogged?: boolean;
	component: LazyExoticComponent<() => JSX.Element>;
}

const config: CustomRouteConfig[] = [
	/* Auth */
	{
		path: ROUTES.SIGNIN,
		loginRequired: false,
		redirectIfLogged: true,
		component: AuthSignin,
		title: 'label.page-signin',
	},
	{
		path: ROUTES.RESET_PASSWORD_REQ,
		loginRequired: false,
		redirectIfLogged: true,
		component: AuthReqResetPassword,
		title: 'label.page-reset-password',
	},
	{
		path: ROUTES.RESET_PASSWORD_NEW,
		loginRequired: false,
		redirectIfLogged: true,
		component: AuthResetPassword,
		title: 'label.page-reset-password',
	},
	{
		path: ROUTES.HOW_GET_ACCESS,
		loginRequired: false,
		redirectIfLogged: true,
		component: AuthHowGetAccess,
		title: 'label.page-how-get-access',
	},
	/* Home */
	{
		path: ROUTES.HOME,
		loginRequired: true,
		component: Home,
	},
	/* Dashboard */
	{
		path: ROUTES.DASHBOARD_USERS,
		loginRequired: true,
		component: DashboardUsers,
		title: 'label.page-user-management',
	},
	/* Devices */
	{
		path: ROUTES.DEVICE_LOGS,
		loginRequired: true,
		component: DeviceLogs,
		title: 'label.page-device-logs',
	},
	{
		path: ROUTES.DEVICE_SEND_COMMAND,
		loginRequired: true,
		component: DeviceCommands,
		title: 'label.page-send-commands',
	},
	{
		path: ROUTES.DEVICE_DETAILS_COMMAND,
		loginRequired: true,
		component: DeviceCommandDetails,
		title: 'label.page-details-command',
	},
	{
		path: ROUTES.CONFIGURATIONS_DEVICE,
		loginRequired: true,
		component: DeviceConfiguration,
		title: 'label.page-terminal-config',
	},
	{
		path: ROUTES.CONFIGURATIONS,
		loginRequired: true,
		component: Configurations,
		title: 'label.configurations',
	},
	{
		path: ROUTES.CONFIGURATIONS_BULK_NEW,
		loginRequired: true,
		component: BulkConfiguration,
		title: 'message.og-title-provisioning-wizard',
	},
	{
		path: ROUTES.CONFIGURATIONS_BULK_DETAILS,
		loginRequired: true,
		component: BulkConfigurationDetails,
		title: 'message.og-title-bulk-configuration-details',
	},
	{
		path: ROUTES.CONFIGURATIONS_TEMPLATE_NEW,
		loginRequired: true,
		component: ConfigurationTemplate,
		title: 'message.og-title-configuration-template-create',
	},
	{
		path: ROUTES.CONFIGURATIONS_TEMPLATE_EDIT,
		loginRequired: true,
		component: ConfigurationTemplateEdit,
		title: 'message.og-title-provisioning-wizard',
	},
	/* Groups */
	{
		path: ROUTES.GROUPS,
		loginRequired: true,
		component: Group,
		title: 'label.page-group-management',
	},
	/* Firmwares */
	{
		path: ROUTES.FIRMWARES,
		loginRequired: true,
		component: Firmwares,
		title: 'label.page-firmwares-devices',
	},
	/* Neosync */
	{
		path: ROUTES.NEOSYNC_RELEASES,
		loginRequired: true,
		component: NeosyncReleases,
		title: 'label.page-neosync-releases',
	},
	{
		path: ROUTES.NEOSYNC_HELP,
		loginRequired: true,
		component: NeosyncHelp,
		title: 'label.page-neosync-help',
	},
	{
		path: ROUTES.NEOSYNC_ANALYTICS,
		loginRequired: true,
		component: NeosyncAnalytics,
		title: 'message.analytics-og-title',
	},
	/* NotFound */
	{
		path: ROUTES.NOT_FOUND,
		loginRequired: false,
		redirectIfLogged: false,
		component: NotFound,
		title: 'label.page-not-found',
	},
];

export default config;
