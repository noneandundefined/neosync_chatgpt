import { toast } from 'react-toastify';
import Modal from '@/components/Modal/Modal';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/constants';
import { useTableContext } from './TableContext';
import SearchInputTable from './SearchInputTable';
import React, { useEffect, useState } from 'react';
import CommonViewColumn from './common/CommonViewColumn';
import Dropdown from '@/components/ui/Dropdown/Dropdown';
import { hasPermission } from '@/constants/Roles.constant';
import { useModalContext } from '@/context/useModalContext';
import { useRole } from '@/context/RoleContext/useRoleContext';
import ModalConfirmUnknownRemoval from '@/components/Modal/ModalConfirmUnknownRemoval';
import { BULK_CONFIGURATION_MAX_DEVICES } from '@/constants/BulkConfiguration.constant';

import Delete from '@/components/@icons/delete';
import ViewColumn from '@/components/@icons/view-column';
import CogOutline from '@/components/@icons/cog-outline';
import AccountSwitch from '@/components/@icons/account-switch';

const TableToolbar = () => {
	const { role } = useRole();
	const { t } = useTranslation();
	const { open, close } = useModalContext();

	const navigate = useNavigate();

	const [openViewColumn, setOpenViewColumn] = useState<boolean>(false);

	const { search, setPage, totalAll, selectedCount, setSearch, hasSelected, selectedIds, onDeleteFn, enableTransferSelectedAction, onTransferSelected } = useTableContext<any>(); // resetSelection
	const [localSearch, setLocalSearch] = useState(search);

	const canTransferSelected = Boolean(enableTransferSelectedAction && onTransferSelected && role !== null && hasPermission(role, 'transfer:users'));

	const handleDeleteSelected = async () => {
		await onDeleteFn(selectedIds);
		close();
	};

	useEffect(() => {
		const timeout = setTimeout(() => {
			if (localSearch === search) {
				return;
			}

			setSearch(localSearch);
			setPage(1);
		}, 300);

		return () => clearTimeout(timeout);
	}, [localSearch, search, setSearch, setPage]);

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setLocalSearch(e.target.value);
	};

	return (
		<React.Fragment>
			<div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 md:gap-0 w-full">
				{/* Фильтры */}
				<div className="flex flex-col md:flex-row md:items-center gap-3 w-full">
					<div className="flex flex-wrap items-center gap-3">
						<p className="text-sm whitespace-nowrap">
							{t('message.marked-objects')}: {selectedCount}
						</p>

						<div className="w-[1px] h-[17px] rounded-full bg-[#ccc]" />

						<p className="text-sm whitespace-nowrap">
							{t('label.total')}: {totalAll}
						</p>
					</div>

					{hasSelected && (
						<div className="flex gap-3 md:ml-3">
							<Tooltip title={t('message.delete-selected')} position="bottom">
								<button
									className="p-2 hover:bg-white rounded-full"
									onClick={() => {
										open(
											<Modal title={t('message.delete-title-confirm')}>
												<ModalConfirmUnknownRemoval onDelete={handleDeleteSelected} />
											</Modal>
										);
									}}
								>
									<Delete fill="#49525f" size={22} />
								</button>
							</Tooltip>

							{canTransferSelected && (
								<div className="flex items-center">
									<div className="h-[17px] w-[1px] rounded-full bg-[#ccc] mr-3" />

									<div className='flex gap-2'>
										<Tooltip title={t('label.transfer-selected-device')} position="bottom">
											<button className="p-2 hover:bg-white rounded-full" onClick={() => onTransferSelected?.(selectedIds)}>
												<AccountSwitch fill="#49525f" size={22} />
											</button>
										</Tooltip>

										<Tooltip title={t('label.bulk-configuration')} position="bottom">
											<button
												className="p-2 hover:bg-white rounded-full"
												onClick={() => {
													if (selectedIds.length > BULK_CONFIGURATION_MAX_DEVICES) {
														toast.error(t('message.provisioning-company-max-devices-exceeded', { count: BULK_CONFIGURATION_MAX_DEVICES }));
														return;
													}

													navigate(ROUTES.CONFIGURATIONS_BULK_NEW, { state: { selectedIds } });
												}}
											>
												<CogOutline fill="#49525f" size={22} />
											</button>
										</Tooltip>
									</div>
								</div>
							)}
						</div>
					)}
				</div>

				{/* Table setting */}
				<div className="relative w-full flex gap-3 justify-end">
					<SearchInputTable value={localSearch} onChange={handleChange} placeholder={t('label.search')} />

					<Tooltip title={t('message.managing-column-creation')}>
						<button
							id="buttonhlp"
							className="!p-[7px] md:!p-[5px] !h-auto"
							onClick={(e) => {
								e.stopPropagation();
								setOpenViewColumn((prev) => !prev);
							}}
						>
							<ViewColumn size={20} fill="#49525f" />
						</button>
					</Tooltip>

					<Dropdown open={openViewColumn} close={() => setOpenViewColumn(false)} stopPropagation={true} className="absolute top-0 z-[999] font-normal">
						<CommonViewColumn />
					</Dropdown>
				</div>
			</div>
		</React.Fragment>
	);
};

export default TableToolbar;
