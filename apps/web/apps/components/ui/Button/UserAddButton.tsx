import Plus from '@/components/@icons/plus';
import { useTranslation } from 'react-i18next';

interface UserAddButtonProps {
	onClick: () => void;
}

const UserAddButton: React.FC<UserAddButtonProps> = ({ onClick }) => {
	const { t } = useTranslation();

	return (
		<button id="buttonhlp" onClick={onClick}>
			<div className="flex items-center gap-2 w-full justify-center sm:justify-end">
				<Plus fill="#62666d" size={20} />
				<p className="text-[13px]">{t('message.add-user')}</p>
			</div>
		</button>
	);
};

export default UserAddButton;
