import { useEffect } from 'react';
import Modal from '../Modal/Modal';
import GUISelect from '../ui/Select/GUISelect';
import { useTranslation } from 'react-i18next';
import ModalProfile from '../Modal/ModalProfile';
import { basicAuthSignOut } from '@/rest/authAPI';
import { LANGUAGES } from '@/constants/Language.constant';
import { useModalContext } from '@/context/useModalContext';
// import { useExportDevices } from '@/hooks/useExportDevices';
import { useRole } from '@/context/RoleContext/useRoleContext';

import Close from '../@icons/close';
import AccountBox from '../@icons/account-box';
import LogoutVariant from '../@icons/logout-variant';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import AppNavLinks from '@/components/Navigation/AppNavLinks';

interface MenuMobileProps {
	isOpen: boolean;
	onClose: () => void;
}

/* Header for mobile | in PC view this is component is hide */
const MenuMobile: React.FC<MenuMobileProps> = ({ isOpen, onClose }) => {
	const { open, close } = useModalContext();

	// const { exportTxt } = useExportDevices();
	const { t, i18n } = useTranslation();
	const { login } = useRole();

	const changeLanguage = (e: React.ChangeEvent<HTMLSelectElement>) => {
		const lang = e.target.value;
		i18n.changeLanguage(lang);
		localStorage.setItem('lang', lang);
	};

	useEffect(() => {
		const handleEscape = (event: KeyboardEvent) => {
			if (event.key === 'Escape') {
				onClose();
			}
		};

		if (isOpen) {
			document.addEventListener('keydown', handleEscape);
		} else {
			document.removeEventListener('keydown', handleEscape);
		}

		return () => {
			document.removeEventListener('keydown', handleEscape);
		};
	}, [isOpen, onClose]);

	return (
		<div className={`fixed inset-0 z-[1000] flex transition-opacity duration-300 ${isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}>
			<aside
				className={`min-w-[20rem] bg-white h-full shadow-lg transform transition-transform duration-300 ease-in-out
        					${isOpen ? 'translate-x-0' : '-translate-x-full'}`}
			>
				<nav className="absolute bg-[#fff] h-full p-6 min-w-[20rem]" style={{ boxShadow: '0 0 25px 0 rgba(0, 0, 0, 0.115)' }}>
					<div className="mb-4 cursor-pointer" onClick={onClose}>
						<Close fill="#49525f" size={20} />
					</div>

					<AppNavLinks variant="mobile" />

					<div className="absolute bottom-8 block md:hidden">
						<div
							className="flex items-center gap-3 py-2 min-w-[16rem] rounded-[8px] cursor-pointer my-2"
							onClick={() => {
								open(
									<Modal title={t('label.profile')}>
										<ModalProfile onClose={() => close()} />
									</Modal>
								);
								onClose();
							}}
						>
							<AccountBox fill="#000" size={22} />
							<p className="text-[13px]">{t('label.profile')}</p>
						</div>

						{/* <div className="flex items-center gap-3 py-2 rounded-[8px] cursor-pointer" onClick={exportTxt}>
							<Download fill="#000" size={22} />
							<p className="text-[13px]">{t('label.export-objects')}</p>
						</div> */}

						<div
							className="flex items-center gap-3 py-3 min-w-[16rem] rounded-[8px] cursor-pointer"
							onClick={() => {
								basicAuthSignOut(true);
							}}
						>
							<LogoutVariant fill="#000" size={22} />
							<p className="text-[14px]">{t('label.sign-out')}</p>
						</div>

						<div className="mt-2">
							<p className="text-[12px] text-[#777]">{t('label.language-change')}</p>
							<GUISelect onChange={changeLanguage} value={localStorage.getItem('lang') || 'ru'}>
								{LANGUAGES.map((lang, index) => (
									<option value={lang.id} key={index}>
										<p className="text-[13px]">{lang.value}</p>
									</option>
								))}
							</GUISelect>
						</div>

						<div className="h-[1px] w-full rounded-full bg-[#ccc] my-3" />

						<div className="flex items-center gap-2 !cursor-default">
							<div className="w-[1.65rem] h-[1.65rem] cursor-default rounded-[50%] flex items-center justify-center" style={{ background: '#E8EDF5' }}>
								<p className="text-[#4577c5] text-[11px]">{(login ?? '?').charAt(0).toUpperCase()}</p>
							</div>

							<p className="text-[#444] text-[13px] font-normal">{login}</p>
						</div>

						<p className="text-[11px] text-gray text-center pt-2">REQID:{localStorage.getItem(CACHEKEYs.NEOSYNC_X_REQ_ID)}</p>
					</div>
				</nav>
			</aside>

			<div className="flex-1 bg-black bg-opacity-40 transition-opacity duration-300" onClick={onClose} />
		</div>
	);
};

export default MenuMobile;
