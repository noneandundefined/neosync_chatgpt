import PageLayout from '../PageLayout';
import { useState } from 'react';

import Modal from '@/components/Modal/Modal';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import GroupsList from '@/components/common/Groups/GroupsList.tsx';
import GroupCreateModal from '@/components/common/Groups/Modals/GroupCreateModal';

import { useTranslation } from 'react-i18next';

import Plus from '@/components/@icons/plus.tsx';
import { basicUserGetMeAccesses } from '@/rest/userAPI';
import { useModalContext } from '@/context/useModalContext';

const GroupPage = () => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	const [groupSysToken, setGroupSysToken] = useState<number>(0);

	const { data: respUserGetMeAccesses } = useHandleServer(['respUserGetMeAccesses'], basicUserGetMeAccesses);

	return (
		<PageLayout>
			{(respUserGetMeAccesses?.access_group_manage ?? true) && (
				<div className="flex flex-col sm:flex-row justify-end mb-4">
					<button
						id="buttonhlp"
						onClick={() =>
							open(
								<Modal title={t('message.create-group')}>
									<GroupCreateModal
										onSuccess={() => {
											setGroupSysToken((v: number) => v + 1);
											close();
										}}
									/>
								</Modal>
							)
						}
					>
						<div className="flex items-center gap-2 w-full justify-center sm:justify-end">
							<Plus fill="#555" size={20} />
							<p className="text-[13px]">{t('message.add-group')}</p>
						</div>
					</button>
				</div>
			)}

			<GroupsList userAccesses={respUserGetMeAccesses} sysToken={groupSysToken} setSysToken={setGroupSysToken} />
		</PageLayout>
	);
};

export default GroupPage;
