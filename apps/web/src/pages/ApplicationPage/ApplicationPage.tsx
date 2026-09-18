import { useState } from 'react';
import PageLayout from '../PageLayout';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { basicUserGetMeAccesses } from '@/rest/userAPI';
import HelpCircle from '@/components/@icons/help-circle';
import { hasPermission } from '@/constants/Roles.constant';
import { useModalContext } from '@/context/useModalContext';
import { useRole } from '@/context/RoleContext/useRoleContext';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import DevicesList from '@/components/common/Devices/DevicesList';
import DeviceAddButton from '@/components/ui/Button/DeviceAddButton';
import ModalHelpHowAddDevice from '@/components/Modal/ModalHelpHowAddDevice';
import DeviceImportsButton from '@/components/ui/Button/DeviceImportsButton';

const ApplicationPage = () => {
	const { role } = useRole();
	const { t } = useTranslation();
	const { open } = useModalContext();

	const [deviceSysToken, setDeviceSysToken] = useState<number>(0);

	const { data: respUserGetMeAccesses } = useHandleServer(['respUserGetMeAccesses'], basicUserGetMeAccesses);

	return (
		<PageLayout>
			<div className="flex flex-col sm:flex-row justify-end mb-4 gap-2">
				{(respUserGetMeAccesses?.access_treker_create ?? true) && role !== null && hasPermission(role, 'create:device') && <DeviceAddButton setDeviceSysToken={setDeviceSysToken} />}
				<div className="flex gap-2">
					{(respUserGetMeAccesses?.access_treker_create ?? true) && role !== null && hasPermission(role, 'import:device') && <DeviceImportsButton setDeviceSysToken={setDeviceSysToken} />}

					{(respUserGetMeAccesses?.access_treker_create ?? true) && role !== null && hasPermission(role, 'create:device') && (
						<button
							id="buttonhlp"
							onClick={() => {
								open(
									<Modal title={t('message.device-add-instruction-title')}>
										<ModalHelpHowAddDevice />
									</Modal>
								);
							}}
						>
							<div className="flex items-center w-full">
								<HelpCircle fill="#62666d" size={18} />
							</div>
						</button>
					)}
				</div>
			</div>

			<DevicesList userAccesses={respUserGetMeAccesses} sysToken={deviceSysToken} setSysToken={setDeviceSysToken} />
		</PageLayout>
	);
};

export default ApplicationPage;
