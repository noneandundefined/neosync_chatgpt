import React from 'react';
import { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import Menu from '../@icons/menu.tsx';
import ChevronLeft from '../@icons/chevron-left.tsx';
import DotsVertical from '../@icons/dots-vertical.tsx';

import Tooltip from '../ui/Tooltip.tsx';
import MenuHeader from './MenuHeader.tsx';
import MenuMobile from './Menu.mobile.tsx';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants.ts';
import { getAuthState } from '@/private-route.tsx';
import Dropdown from '../ui/Dropdown/Dropdown.tsx';
import { useTcpHealth } from '@/context/useTcpHealth.tsx';
import { ROLEs_AND_NAME } from '@/constants/Roles.constant';
import { useRole } from '@/context/RoleContext/useRoleContext';

interface HeaderProps {
	show?: boolean;
	arrow?: boolean;
}

const getBackPath = (pathname: string) => {
	if (pathname === ROUTES.CONFIGURATIONS_TEMPLATE_NEW || pathname.startsWith('/configurations/devices/templates/')) {
		return `${ROUTES.CONFIGURATIONS}?tab=templates`;
	}

	if (pathname.startsWith('/configurations/devices/')) {
		return ROUTES.HOME;
	}

	return null;
};

const Header: React.FC<HeaderProps> = ({ show, arrow }) => {
	const { t } = useTranslation();
	const location = useLocation();
	const navigate = useNavigate();

	const handleBack = () => {
		const backPath = getBackPath(location.pathname);
		if (backPath) {
			navigate(backPath);
			return;
		}

		navigate(-1);
	};

	const { login, role } = useRole();

	const tcpHealth = useTcpHealth();

	const [isOpenMenu, setIsOpenMenu] = useState<boolean>(false);
	const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
	const [isMenuMobileOpen, setIsMenuMobileOpen] = useState<boolean>(false);

	useEffect(() => {
		let mounted = true;

		async function checkAuth() {
			const logged = await getAuthState();
			if (mounted) setIsLoggedIn(logged);
		}

		checkAuth();

		return () => {
			mounted = false;
		};
	}, []);

	return (
		<React.Fragment>
			{isMenuMobileOpen && <MenuMobile isOpen={isMenuMobileOpen} onClose={() => setIsMenuMobileOpen(false)} />}

			<header
				className={`sticky z-50 shrink-0 h-[60px] px-[0.4rem] lg:px-[1rem] xl:px-[2rem] w-full flex items-center justify-between text-md font-semibold text-gray-900 md:border-b border-[#e5e7eb] bg-[#fafafa] ${tcpHealth ? 'top-0' : 'top-6'}`}
			>
				<div className="flex items-center gap-3 w-[33%]">
					{isLoggedIn && !location.pathname.startsWith('/configurations/devices/') && (
						<div className="cursor-pointer hidden" id="button__menu__neomatica" onClick={() => setIsMenuMobileOpen(true)}>
							<div className="hover:bg-[#ffffff] p-2 rounded-full">
								<Menu fill="#49525f" size={23} />
							</div>
						</div>
					)}

					{arrow && (
						<Tooltip title={t('label.back')} className="!font-normal">
							<div onClick={handleBack} className="cursor-pointer">
								<div className="cursor-pointer p-1 rounded hover:bg-white">
									<ChevronLeft fill="#49525f" size={22} />
								</div>
							</div>
						</Tooltip>
					)}

					<Link to={ROUTES.HOME} className="hidden xl:flex items-center gap-2 text-[#49525f] hover:text-[#49525f] font-medium ml-3 sm:ml-0">
						Neomatica
					</Link>
				</div>

				<div className="w-[33%] hidden md:flex justify-end">
					<div className="flex items-center gap-4">
						{show && (
							<React.Fragment>
								<span>
									<p className="font-medium text-[#49525f] text-[13px] text-center whitespace-nowrap">{role ? t(`label.${ROLEs_AND_NAME[role]}`) : t('label.role-not-found')}</p>
								</span>

								<div className="w-[1px] h-[17px] rounded-full bg-[#ccc]" />

								<div className="flex items-center gap-2 !cursor-default">
									<p className="text-[#444] text-[13px] font-normal">{login ?? '?'}</p>
								</div>

								<div className="w-[1px] h-[17px] rounded-full bg-[#ccc]" />

								<div className="relative">
									<div
										className="hover:bg-white p-[5px] rounded-full transition cursor-pointer"
										onClick={(e) => {
											e.stopPropagation();
											setIsOpenMenu((prev) => !prev);
										}}
									>
										<DotsVertical size={22} fill="#49525f" />
									</div>

									<Dropdown open={isOpenMenu} close={() => setIsOpenMenu(false)} stopPropagation={true} className="absolute top-0 right-0 z-[999] font-normal">
										<MenuHeader closeMenu={() => setIsOpenMenu(false)} />
									</Dropdown>
								</div>
							</React.Fragment>
						)}
					</div>
				</div>
			</header>
		</React.Fragment>
	);
};

export default Header;
