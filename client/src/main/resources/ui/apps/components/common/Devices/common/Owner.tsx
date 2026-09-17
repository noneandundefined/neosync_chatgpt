import React from 'react';
import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { Device_DeviceConf_Sync } from '@/rest/deviceAPI';
import { useModalContext } from '@/context/useModalContext';
import DeviceTransferModal from '../Modals/DeviceTransferModal';

interface OwnerProps {
	d: Device_DeviceConf_Sync;
	onUpdate: () => void;
}

const Owner: React.FC<OwnerProps> = ({ d, onUpdate }) => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	return (
		<React.Fragment>
			<div
				className="flex items-center justify-center gap-2"
				onClick={() =>
					open(
						<Modal title={t('label.transfer-title-device-sure')}>
							<DeviceTransferModal
								onSuccess={() => close()}
								onTransfer={(user) => {
									if (!user) {
										d.owner_uuid = null;
										d.user_email = '';
									} else {
										d.owner_uuid = user.user_uuid;
										d.user_email = user.email;
									}

									onUpdate();
								}}
								device={d}
							/>
						</Modal>
					)
				}
			>
				<p className="text-[#49525f] underline truncate max-w-[150px] hover:text-black">{d.user_email || t('message.select-no-owner')}</p>
			</div>
		</React.Fragment>
	);
};

export default Owner;
