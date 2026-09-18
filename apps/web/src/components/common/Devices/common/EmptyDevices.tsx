import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { Dispatch, SetStateAction } from 'react';
import { UserAccessResponse } from '@/rest/userAPI';
import GUIButton from '@/components/ui/Button/GUIButton';
import { hasPermission } from '@/constants/Roles.constant';
import { useModalContext } from '@/context/useModalContext';
import DeviceCreateModal from '../Modals/DeviceCreateModal';
import { useRole } from '@/context/RoleContext/useRoleContext';
import ModalHelpHowAddDevice from '@/components/Modal/ModalHelpHowAddDevice';

interface EmptyDevicesProps {
	/** Профиль пользователя для доступов */
	userAccesses: UserAccessResponse | null;
	/** Токен для обновления данных в таблице */
	setSysToken: Dispatch<SetStateAction<number>>;
}

const EmptyDevices: React.FC<EmptyDevicesProps> = ({ userAccesses, setSysToken }) => {
	const { t } = useTranslation();

	const { role } = useRole();
	const { open, close } = useModalContext();

	return (
		<div className="flex-1 bg-white rounded-md border border-[#e4e4e4]">
			<div className="flex flex-col space-y-7 items-center justify-center h-full">
				<img src="/local/templates/neomatica/images/neomatica-with-text-logo.png" alt="neomatica" className="max-w-[15vw] min-w-[15rem] object-cover my-7 opacity-90" draggable={false} />
				<div className="w-full h-[1px] bg-[#eee]" />
				<p className="min-w-[20rem] my-4 text-[1.5rem] md:text-[2rem] uppercase tracking-wide font-normal max-w-[50vw] text-center">{t('message.empty-device-added')}</p>
				<p className="mb-3 text-[1rem] text-[#444] font-normal max-w-[30vw] text-balance text-center">{t('message.to-get-started')}:</p>

				{(userAccesses?.access_treker_create ?? true) && role !== null && hasPermission(role, 'create:device') && (
					<GUIButton
						className="!w-auto"
						onClick={() =>
							open(
								<Modal title={t('message.create-device')}>
									<DeviceCreateModal
										onSuccess={() => {
											setSysToken((v: number) => v + 1);
											close();
										}}
									/>
								</Modal>
							)
						}
					>
						{t('message.add-device')}
					</GUIButton>
				)}

				<div
					className="text-[14px] text-[#1d3c5d] cursor-pointer underline"
					onClick={() => {
						open(
							<Modal title={t('message.device-add-instruction-title')}>
								<ModalHelpHowAddDevice />
							</Modal>
						);
					}}
				>
					<p>{t('message.how-add-device-to-neosync')}</p>
				</div>
			</div>
		</div>
	);
};

export default EmptyDevices;
