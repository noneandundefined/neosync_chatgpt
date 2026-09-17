import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import GUIButton from '@/components/ui/Button/GUIButton';
import { GUITextarea } from '@/components/ui/Textarea/GUITextarea';

import Modal from '@/components/Modal/Modal';
import { useTcpHealth } from '@/context/useTcpHealth';
import CommandRecentHistory from './CommandRecentHistory';
import { useModalContext } from '@/context/useModalContext';
import CommandHistoryModal from './Modals/CommandHistoryModal';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { DeviceSendCommandRequest } from '@/interface/device/deviceSendCommandRequest.interface';
import { basicDeviceCommandRecentHistory, basicDeviceSendCommand, DeviceShort } from '@/rest/deviceCommandAPI';

interface CommandSendProps {
	selectedDevice: DeviceShort[];
}

const CommandSend: React.FC<CommandSendProps> = ({ selectedDevice }) => {
	const { t } = useTranslation();

	const tcpHealth = useTcpHealth();

	const { data: respDeviceCommandRecentHistory, loading: respDeviceCommandRecentHistoryLoading, reload } = useHandleServer(['respDeviceCommandRecentHistory'], basicDeviceCommandRecentHistory, { refetchInterval: 3000 });

	const { open, close } = useModalContext();

	const [openIndex, setOpenIndex] = useState<number | null>(null);
	const [command, setCommand] = useState<string>('');

	const handleSubmit = async (sendMode: 'instant' | 'on_connect') => {
		if (!command || selectedDevice.length <= 0) return;

		const payload: DeviceSendCommandRequest = {
			send_mode: sendMode,
			command: command,
			imeis: selectedDevice.map((d) => d.imei),
		};

		await basicDeviceSendCommand(payload);

		reload();
		setCommand('');
		setOpenIndex(0);
	};

	const hasActivatedDevice = selectedDevice.some((d) => d.activated);

	const hasOnlineDevice = selectedDevice.some((d) => d.status);

	const canSendInstant = tcpHealth && hasActivatedDevice && hasOnlineDevice;

	return (
		<React.Fragment>
			<div className="flex-1 w-full h-auto md:min-h-0 md:h-full bg-white border border-[#e4e4e4] py-4 px-5 rounded space-y-3 flex flex-col">
				<div className="flex shrink-0">
					<p className="text-center font-medium text-[17px] text-[#49525f]">{t('label.command')}</p>
				</div>

				<form className="flex flex-col flex-1 md:min-h-0">
					<div className="flex-1 space-y-3 py-3">
						<div>
							<GUITextarea
								placeholder={t('message.enter-the-command')}
								className="text-[14px]"
								value={command}
								onChange={(e) => setCommand(e.target.value)}
								onKeyDown={(e) => {
									if (e.key === 'Enter' && !e.shiftKey) {
										e.preventDefault();
										handleSubmit('instant');
									}
								}}
							/>
						</div>

						<div className="flex items-center gap-2">
							<GUIButton onClick={() => handleSubmit('instant')} type="submit" disabled={!canSendInstant} className={`${!canSendInstant && '!opacity-40 !cursor-not-allowed'}`}>
								{t('label.send-instant')}
							</GUIButton>

							<GUIButton onClick={() => handleSubmit('on_connect')} type="submit">
								{t('label.send-on-connect')}
							</GUIButton>
						</div>

						<CommandRecentHistory data={respDeviceCommandRecentHistory} loading={respDeviceCommandRecentHistoryLoading} reload={reload} openIndex={openIndex} onOpenIndexChange={setOpenIndex} />
					</div>

					<div
						className="flex justify-center cursor-pointer hover:underline shrink-0"
						onClick={() =>
							open(
								<Modal title={t('message.command-history')}>
									<CommandHistoryModal closeModal={() => close()} />
								</Modal>
							)
						}
					>
						<p className="text-[14px] text-[#1d3c5d]">{t('message.command-history')}</p>
					</div>
				</form>
			</div>
		</React.Fragment>
	);
};

export default CommandSend;
