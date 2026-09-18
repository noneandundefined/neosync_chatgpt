import { Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import { basicMetaGetVersion } from '@/rest/metaAPI';
import SwapHorizontal from '../@icons/swap-horizontal';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import AppNavLinks from '@/components/Navigation/AppNavLinks';
import { useHandleServer } from '@/hooks/Server/useHandleServer';

/* Left sidebar */
const SideBar = () => {
	const { t } = useTranslation();

	const [collapsed, setCollapsed] = useState<boolean>(() => {
		const saved = sessionStorage.getItem(CACHEKEYs.SIDEBAR_COLLAPSED);
		return saved ? JSON.parse(saved) : false;
	});

	useEffect(() => {
		sessionStorage.setItem(CACHEKEYs.SIDEBAR_COLLAPSED, JSON.stringify(collapsed));
	}, [collapsed]);

	const { data: basicMetaGetVersionResp } = useHandleServer(['basicMetaGetVersionResp'], basicMetaGetVersion);

	return (
		<aside id="side__bar__neomatica" className={`min-h-0 shrink-0 overflow-hidden border-r pr-5 py-4 transition-all duration-300 ${collapsed ? 'pr-2' : 'min-w-[17rem]'}`}>
			<nav className="flex h-full flex-col justify-between">
				<AppNavLinks variant="sidebar" collapsed={collapsed} />

				<div className="space-y-5">
					{/* VERSION */}
					{!collapsed && (
						<Link to={ROUTES.NEOSYNC_RELEASES} className="text-sm">
							<p className="font-medium text-[#49525f]">NeoSync {basicMetaGetVersionResp}</p>
						</Link>
					)}

					<div className="px-3 py-[10px] rounded-r-[8px] cursor-pointer bg-transparent hover:bg-[#fff] border-l-2 border-transparent hover:border-[#395d95]" onClick={() => setCollapsed(!collapsed)}>
						<div className="flex items-center gap-5 text-[#000] hover:text-[#000] font-normal text-[14px]">
							<SwapHorizontal fill="#49525f" className="shrink-0" />
							{!collapsed && <p className="font-medium text-[#49525f]">{t('label.collapse')}</p>}
						</div>
					</div>
				</div>
			</nav>
		</aside>
	);
};

export default SideBar;
