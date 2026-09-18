import { UserAccessResponse } from '@/rest/userAPI';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import Delete from '@/components/@icons/delete.tsx';
import Pencil from '@/components/@icons/pencil.tsx';
import GroupEditModal from './Modals/GroupEditModal';
import React, { Dispatch, SetStateAction } from 'react';
import { useModalContext } from '@/context/useModalContext';
import AccountSwitch from '@/components/@icons/account-switch';
import CopyToAccountsModal from './Modals/CopyToAccountsModal';
// import { useRole } from '@/context/RoleContext/useRoleContext.tsx';
import { groupColumns } from '@/components/common/Groups/index.tsx';
import ModalConfirmRemoval from '@/components/Modal/ModalConfirmRemoval';
import GenericTable, { FetchFn } from '@/components/common/Table/GenericTable/GenericTable.tsx';
import { basicGroupDelete, basicGroupMassiveDelete, basicGroupsGetWithParams, GroupResponse } from '@/rest/groupAPI.ts';

interface GroupsListProps {
	/** Профиль пользователя для доступов */
	userAccesses: UserAccessResponse | null;
	/** Токен для обновления данных в таблице */
	sysToken: number;
	/** Токен для обновления данных в таблице */
	setSysToken: Dispatch<SetStateAction<number>>;
}

const GroupsList: React.FC<GroupsListProps> = ({ userAccesses, sysToken, setSysToken }) => {
	// const { role } = useRole();
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	const handleRowAction = (group: GroupResponse, isHovered: boolean) => (
		<div
			className="flex items-center justify-center gap-2 sm:gap-2"
			style={{
				opacity: isHovered ? 1 : 0,
			}}
		>
			{(userAccesses?.access_group_manage ?? true) && (
				<React.Fragment>
					{userAccesses?.user_uuid === group.user_uuid && (
						<AccountSwitch
							fill="#666666"
							size={24}
							className="cursor-pointer transition"
							onClick={() => {
								open(
									<Modal title={t('message.copy-to-account')}>
										<CopyToAccountsModal groupId={group.id} onSuccess={() => setSysToken((v: number) => v + 1)} />
									</Modal>
								);
							}}
						/>
					)}

					<Pencil
						fill="#5ba0e8"
						size={24}
						className="cursor-pointer transition"
						onClick={() => {
							open(
								<Modal title={t('message.edit-group')} width="700px">
									<GroupEditModal
										groupId={group.id}
										userUuid={userAccesses?.user_uuid}
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
								<Modal title={t('message.delete-title-group-sure')}>
									<ModalConfirmRemoval
										object={group.id}
										confirm={group.name}
										onDelete={async () => {
											await basicGroupDelete(group.id);
											setSysToken((v: number) => v + 1);
											close();
										}}
									/>
								</Modal>
							);
						}}
					/>
				</React.Fragment>
			)}
		</div>
	);

	const columns = groupColumns();

	const handleDelete = async (payloads: (string | number)[]) => {
		const numberIds = payloads.map((id) => Number(id)).filter((id) => !isNaN(id));
		await basicGroupMassiveDelete(numberIds);
		setSysToken((v: number) => v + 1);
	};

	const fetchGroups: FetchFn<GroupResponse> = async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicGroupsGetWithParams(page, limit, search, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	};

	return (
		<React.Fragment>
			<GenericTable<GroupResponse>
				tableKey="fetchGroups"
				columns={columns}
				fetchFn={fetchGroups}
				sysToken={sysToken}
				getRowId={(group) => group.id}
				rowActions={(group, _, isHovered) => handleRowAction(group, isHovered)}
				onDeleteFn={handleDelete}
			/>
		</React.Fragment>
	);
};

export default GroupsList;
