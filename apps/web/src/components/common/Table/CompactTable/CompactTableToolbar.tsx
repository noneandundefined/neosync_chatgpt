import { useState } from 'react';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import Dropdown from '@/components/ui/Dropdown/Dropdown';
import { useModalContext } from '@/context/useModalContext';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useCompactTableContext } from './CompactTableContext';
import DotsHorizontal from '@/components/@icons/dots-horizontal';
import ModalConfirmUnknownRemoval from '@/components/Modal/ModalConfirmUnknownRemoval';

const CompactTableToolbar = () => {
	const { t } = useTranslation();
	const { open, close } = useModalContext();
	const [openMenu, setOpenMenu] = useState(false);

	const { allChecked, toggleAll, selectedCount, hasSelected, selectedIds, onDeleteFn, bulkActions, selectionMode } = useCompactTableContext<any>();

	if (selectionMode === 'single') {
		return null;
	}

	const handleDeleteSelected = async () => {
		if (!onDeleteFn) return;
		await onDeleteFn(selectedIds);

		close();

		setOpenMenu(false);
	};

	return (
		<div className="flex flex-wrap items-center gap-3 px-4 py-3 border-b border-[#e9e9e9] bg-[#f6f6f6] min-h-[3.7rem]">
			<GUICheckbox checked={allChecked} onChange={() => toggleAll()} />

			{hasSelected && (
				<div className="flex items-center gap-6">
					<p className="text-sm text-[#49525f] whitespace-nowrap">
						{t('label.selected')}: {selectedCount}
					</p>

					{onDeleteFn && (
						<div className="relative mt-1">
							<button
								type="button"
								className="px-2 rounded-md border border-[#d1d9e0]"
								onClick={(e) => {
									e.stopPropagation();
									setOpenMenu((prev) => !prev);
								}}
							>
								<DotsHorizontal fill="#49525f" size={20} />
							</button>

							<Dropdown open={openMenu} close={() => setOpenMenu(false)} className="absolute top-8 left-0 bg-white py-1 border border-[#e9e9e9] min-w-[11rem] z-[200]" stopPropagation>
								<p
									className="text-sm hover:bg-[#f7f7f7] px-3 py-[9px] md:py-[5px] cursor-pointer"
									onClick={() => {
										setOpenMenu(false);
										open(
											<Modal title={t('message.delete-title-confirm')}>
												<ModalConfirmUnknownRemoval onDelete={handleDeleteSelected} />
											</Modal>
										);
									}}
								>
									{t('label.delete')}
								</p>
							</Dropdown>
						</div>
					)}

					{bulkActions}
				</div>
			)}
		</div>
	);
};

export default CompactTableToolbar;
