import type { ComponentType } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation } from 'react-router-dom';

import Car from '../@icons/car';
import Poll from '../@icons/poll';
import TrainCar from '@/components/@icons/train-car';
import HelpCircle from '@/components/@icons/help-circle';
import AccountGroup from '@/components/@icons/account-group';
import CardBulletedOutline from '@/components/@icons/card-bulleted-outline';

import { ROUTES } from '@/constants/constants';
import CogOutline from '../@icons/cog-outline';
import { hasPermission } from '@/constants/Roles.constant';
import { useRole } from '@/context/RoleContext/useRoleContext';

export type AppNavVariant = 'sidebar' | 'mobile';

type NavIconProps = { fill: string; size?: number; className?: string };

type AppNavItemConfig = {
	route: string;
	labelKey: string;
	Icon: ComponentType<NavIconProps>;
	permission?: string;
	/** Extra classes on <li> per layout */
	liExtra?: Partial<Record<AppNavVariant, string>>;
};

const APP_NAV_ITEMS: AppNavItemConfig[] = [
	{ route: ROUTES.HOME, labelKey: 'label.devices', Icon: Car },
	{ route: ROUTES.DEVICE_SEND_COMMAND, labelKey: 'label.commands', Icon: CardBulletedOutline },
	{ route: ROUTES.GROUPS, labelKey: 'label.groups', Icon: TrainCar },
	{ route: ROUTES.CONFIGURATIONS, labelKey: 'label.configurations', Icon: CogOutline },
	{
		route: ROUTES.DASHBOARD_USERS,
		labelKey: 'label.users',
		Icon: AccountGroup,
		permission: 'view:users',
		liExtra: { mobile: 'min-w-[14rem]' },
	},
	{
		route: ROUTES.NEOSYNC_ANALYTICS,
		labelKey: 'label.analytics',
		Icon: Poll,
		permission: 'view:analytics',
		liExtra: { mobile: 'min-w-[14rem]' },
	},
	{
		route: ROUTES.NEOSYNC_HELP,
		labelKey: 'label.help',
		Icon: HelpCircle,
		liExtra: { mobile: 'min-w-[14rem]' },
	},
];

const activeLi = 'border-l-2 border-[#395d95]';
const inactiveLi = 'bg-transparent hover:bg-[#fff] border-l-2 border-transparent hover:border-[#395d95]';

const variantClasses: Record<AppNavVariant, { li: string; link: string }> = {
	sidebar: {
		li: `px-3 py-[10px] rounded-r-[8px] cursor-pointer`,
		link: 'flex items-center gap-5 text-[#000] hover:text-[#000] font-normal text-[14px]',
	},
	mobile: {
		li: `px-3 py-[7px] rounded-r-[8px] cursor-pointer`,
		link: 'flex items-center gap-3 text-[#000] hover:text-[#000] font-normal text-[14px]',
	},
};

interface AppNavLinksProps {
	variant: AppNavVariant;
	collapsed?: boolean;
}

const AppNavLinks: React.FC<AppNavLinksProps> = ({ variant, collapsed }) => {
	const { role } = useRole();
	const { t } = useTranslation();
	const location = useLocation();

	const { li: liBase, link: linkClass } = variantClasses[variant];

	return (
		<ul className="flex flex-col gap-2">
			{APP_NAV_ITEMS.map(({ route, labelKey, Icon, permission, liExtra }) => {
				if (permission !== undefined) {
					if (role === null || !hasPermission(role, permission)) {
						return null;
					}
				}

				const isActive = location.pathname === route;
				const extra = liExtra?.[variant];

				return (
					<li key={route} className={`${isActive ? activeLi : inactiveLi} ${liBase} ${extra ?? ''}`.trim()}>
						<Link to={route} className={linkClass}>
							<Icon fill={isActive ? '#395d95' : '#49525f'} size={22} className="cursor-pointer shrink-0" />
							{!collapsed && <p className={isActive ? 'font-medium text-[#395d95]' : 'font-medium text-[#49525f]'}>{t(labelKey)}</p>}
						</Link>
					</li>
				);
			})}
		</ul>
	);
};

export default AppNavLinks;
