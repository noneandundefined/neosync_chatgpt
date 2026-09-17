import Plus from '@/components/@icons/plus';
import DeviceCreateModal from '@/components/common/Devices/Modals/DeviceCreateModal';
import Modal from '@/components/Modal/Modal';
import { useModalContext } from '@/context/useModalContext';
import { Dispatch, SetStateAction } from 'react';
import { useTranslation } from 'react-i18next';

interface DeviceAddButtonProps {
	setDeviceSysToken: Dispatch<SetStateAction<number>>;
}

const DeviceAddButton: React.FC<DeviceAddButtonProps> = ({ setDeviceSysToken }) => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	return (
		<button
			id="buttonhlp"
			onClick={() =>
				open(
					<Modal title={t('message.create-device')}>
						<DeviceCreateModal
							onSuccess={() => {
								setDeviceSysToken((v: number) => v + 1);
								close();
							}}
						/>
					</Modal>
				)
			}
		>
			<div className="flex items-center gap-2 w-full justify-center sm:justify-end">
				<Plus fill="#62666d" size={20} />
				<p className="text-[13px] whitespace-nowrap">{t('message.add-device')}</p>
			</div>
		</button>
	);
};

export default DeviceAddButton;
