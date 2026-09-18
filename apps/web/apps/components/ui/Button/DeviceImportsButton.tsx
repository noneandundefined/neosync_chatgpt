import FileImportOutline from '@/components/@icons/file-import-outline';
import DeviceImportsModal from '@/components/common/Devices/Modals/DeviceImportsModal';
import Modal from '@/components/Modal/Modal';
import { useModalContext } from '@/context/useModalContext';
import { Dispatch, SetStateAction } from 'react';
import { useTranslation } from 'react-i18next';

interface DeviceImportsButtonProps {
	setDeviceSysToken: Dispatch<SetStateAction<number>>;
}

const DeviceImportsButton: React.FC<DeviceImportsButtonProps> = ({ setDeviceSysToken }) => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	return (
		<button
			id="buttonhlp"
			onClick={() =>
				open(
					<Modal title={t('message.import-devices')}>
						<DeviceImportsModal
							onSuccess={() => {
								setDeviceSysToken((v: number) => v + 1);
								close();
							}}
						/>
					</Modal>
				)
			}
			className="flex-1"
		>
			<div className="flex items-center gap-2 w-full justify-center sm:justify-end">
				<FileImportOutline fill="#62666d" size={20} />
				<p className="text-[13px] whitespace-nowrap">{t('message.import-devices')}</p>
			</div>
		</button>
	);
};

export default DeviceImportsButton;
