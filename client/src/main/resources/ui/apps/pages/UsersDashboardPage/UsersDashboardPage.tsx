import { useState } from 'react';
import PageLayout from '@/pages/PageLayout';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import GUISelect from '@/components/ui/Select/GUISelect';
import { useModalContext } from '@/context/useModalContext';
import UsersList from '@/components/common/Users/UsersList';
import { useRole } from '@/context/RoleContext/useRoleContext';
import UserAddButton from '@/components/ui/Button/UserAddButton';
import UserCreateModal from '@/components/common/Users/Modals/UserCreateModal';
import { getAvailableRoles, ROLEs_AND_NAME } from '@/constants/Roles.constant';

const UsersDashboardPage = () => {
	const { t } = useTranslation();
	const { role } = useRole();

	const { open, close } = useModalContext();

	const [userSysToken, setUserSysToken] = useState<number>(0);
	const [usersRole, setUsersRole] = useState<string>('');

	const availableRoles = getAvailableRoles(role);

	return (
		<PageLayout>
			<div className="flex flex-col sm:flex-row justify-end mb-4 gap-3">
				<div>
					<GUISelect
						onChange={(e) => {
							setUsersRole(e.target.value);
							setUserSysToken((v: number) => v + 1);
						}}
					>
						<option value="">{t('label.all')}</option>
						{availableRoles.map((role, index) => (
							<option value={role} key={index}>
								{t(`label.${ROLEs_AND_NAME[role]}`)}
							</option>
						))}
					</GUISelect>
				</div>

				<UserAddButton
					onClick={() =>
						open(
							<Modal title={t('message.add-user')}>
								<UserCreateModal
									onSuccess={() => {
										setUserSysToken((v: number) => v + 1);
										close();
									}}
								/>
							</Modal>
						)
					}
				/>
			</div>

			<UsersList usersRole={usersRole} sysToken={userSysToken} setSysToken={setUserSysToken} />
		</PageLayout>
	);
};

export default UsersDashboardPage;
