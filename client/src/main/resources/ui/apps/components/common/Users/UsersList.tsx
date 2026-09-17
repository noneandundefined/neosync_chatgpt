import { userColumns } from '.';
import { useMemo, useState } from 'react';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import Delete from '@/components/@icons/delete';
import Pencil from '@/components/@icons/pencil';
import UserEditModal from './Modals/UserEditModal';
import { useModalContext } from '@/context/useModalContext';
import UserExpandedRow from './common/UserExpandedRow';
import React, { Dispatch, SetStateAction } from 'react';
import EmailConfirmModal from './Modals/EmailConfirmModal';
import Information from '@/components/@icons/information.tsx';
import EmailArrowRight from '@/components/@icons/email-arrow-right';
import ModalConfirmRemoval from '@/components/Modal/ModalConfirmRemoval';
import GenericTable, { FetchFn } from '../Table/GenericTable/GenericTable';
import { basicUserDelete, basicUserMassiveDelete, basicUserSendEmail, basicUsersGetWithParams, UserResponse } from '@/rest/userAPI';

interface UsersListProps {
	usersRole: string;
	/** Токен для обновления данных в таблице */
	sysToken: number;
	/** Токен для обновления данных в таблице */
	setSysToken: Dispatch<SetStateAction<number>>;
}

const UsersList: React.FC<UsersListProps> = ({ usersRole, sysToken, setSysToken }) => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	const [expandedRow, setExpandedRow] = useState<number | null>(null);

	const columns = useMemo(() => userColumns(t), [t]);

	const handleRowAction = (user: UserResponse, idx: number, isHovered: boolean) => (
		<div
			className="flex items-center justify-center gap-2 sm:gap-1"
			style={{
				opacity: isHovered ? 1 : 0,
			}}
		>
			<Information
				fill="#1976d2"
				size={24}
				className="cursor-pointer transition"
				onClick={(e?: React.MouseEvent<SVGSVGElement>) => {
					e?.stopPropagation();
					setExpandedRow(expandedRow === idx ? null : idx);
				}}
			/>

			<EmailArrowRight
				fill="#ff7c00"
				size={24}
				className="cursor-pointer transition"
				onClick={() => {
					open(
						<Modal title={t('message.confirm-send-email-title')}>
							<EmailConfirmModal
								user={user}
								onFn={async (token: string, language: string) => {
									await basicUserSendEmail(user.user_contact.user_uuid, token, language);
									close();
								}}
							/>
						</Modal>
					);
				}}
			/>

			<Pencil
				fill="#5ba0e8"
				size={24}
				className="cursor-pointer transition"
				onClick={() => {
					open(
						<Modal title={t('message.edit-user')}>
							<UserEditModal
								uuid={user.user_contact.user_uuid}
								onSuccess={() => {
									setSysToken((v: number) => v + 1);
									close();
								}}
							/>
						</Modal>
					);
				}}
			/>

			<Delete
				fill="#f40000"
				size={24}
				className="cursor-pointer transition"
				onClick={() => {
					open(
						<Modal title={t('message.delete-title-user-sure')}>
							<ModalConfirmRemoval
								object={user.user_contact.user_uuid}
								confirm={user.email}
								onDelete={async () => {
									await basicUserDelete(user.user_contact.user_uuid);
									setSysToken((v: number) => v + 1);
									close();
								}}
							/>
						</Modal>
					);
				}}
			/>
		</div>
	);

	const renderExpandedRow = (user: UserResponse, idx: number) => (expandedRow === idx ? <UserExpandedRow user={user} t={t} /> : null);

	const handleDelete = async (payloads: (string | number)[]) => {
		await basicUserMassiveDelete(payloads as string[]);
		setSysToken((v: number) => v + 1);
	};

	const fetchUsers: FetchFn<UserResponse> = async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicUsersGetWithParams(page, limit, search, usersRole, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	};

	return (
		<React.Fragment>
			<GenericTable<UserResponse>
				tableKey="fetchUsers"
				columns={columns}
				fetchFn={fetchUsers}
				sysToken={sysToken}
				getRowId={(user) => user.user_contact.user_uuid}
				rowActions={(user, idx, isHovered) => handleRowAction(user, idx, isHovered)}
				renderExpandedRow={renderExpandedRow}
				onDeleteFn={handleDelete}
			/>
		</React.Fragment>
	);
};

export default UsersList;
