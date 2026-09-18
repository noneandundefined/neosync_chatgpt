import Modal from '../Modal/Modal';
import Download from '../@icons/download';
import GUISelect from '../ui/Select/GUISelect';
import { useTranslation } from 'react-i18next';
import AccountBox from '../@icons/account-box';
import ModalProfile from '../Modal/ModalProfile';
import { basicAuthSignOut } from '@/rest/authAPI';
import LogoutVariant from '../@icons/logout-variant';
import { LANGUAGES } from '@/constants/Language.constant';
import { useModalContext } from '@/context/useModalContext';
import { useExportDevices } from '@/hooks/useExportDevices';

interface MenuHeaderProps {
	closeMenu: () => void;
}

const MenuHeader: React.FC<MenuHeaderProps> = ({ closeMenu }) => {
	const { t, i18n } = useTranslation();

	const { open, close } = useModalContext();

	const { exportTxt } = useExportDevices();

	const changeLanguage = (e: React.ChangeEvent<HTMLSelectElement>) => {
		const lang = e.target.value;
		i18n.changeLanguage(lang);
		localStorage.setItem('lang', lang);
	};

	return (
		<div className="bg-white min-w-[15rem] flex flex-col" style={{ boxShadow: '0 9px 17px 0 rgba(0, 0, 0, 0.3)' }}>
			<div className="flex flex-col p-2 px-3">
				<p className="text-[12px] text-[#777]">{t('label.language-change')}</p>
				<GUISelect onChange={changeLanguage} className="relative !min-w-auto" value={localStorage.getItem('lang') || 'ru'}>
					{LANGUAGES.map((lang, index) => (
						<option value={lang.id} key={index}>
							<p className="text-[13px]">{lang.value}</p>
						</option>
					))}
				</GUISelect>
			</div>

			<div
				className="flex items-center gap-3 py-[9px] pl-3 cursor-pointer hover:bg-[#f3f3f3] transition"
				onClick={() => {
					open(
						<Modal title={t('label.profile')}>
							<ModalProfile onClose={() => close()} />
						</Modal>
					);
					closeMenu();
				}}
			>
				<AccountBox fill="#555" size={21} />
				<p className="text-[13px]">{t('label.profile')}</p>
			</div>

			<div className="flex items-center gap-3 py-[9px] pl-3 cursor-pointer hover:bg-[#f3f3f3] transition" onClick={exportTxt}>
				<Download fill="#555" size={21} />
				<p className="text-[13px]">{t('label.export-objects')}</p>
			</div>

			<div className="flex items-center gap-3 py-[9px] pl-3 cursor-pointer hover:bg-[#f3f3f3] transition" onClick={() => basicAuthSignOut(true)}>
				<LogoutVariant fill="#555" size={21} />
				<p className="text-[13px]">{t('label.sign-out')}</p>
			</div>
		</div>
	);
};

export default MenuHeader;
