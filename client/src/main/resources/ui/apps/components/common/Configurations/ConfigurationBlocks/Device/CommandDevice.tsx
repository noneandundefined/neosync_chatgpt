import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTcpHealth } from '@/context/useTcpHealth';
import GUIButton from '@/components/ui/Button/GUIButton';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { GUITextarea } from '@/components/ui/Textarea/GUITextarea';
import { useConfigurationContext } from '@/context/useConfigurationContext';
import CommandRecentHistory from '@/components/common/DeviceCommands/CommandRecentHistory';
import { DeviceSendCommandRequest } from '@/interface/device/deviceSendCommandRequest.interface';
import { basicDeviceCommandRecentHistory, basicDeviceSendCommand } from '@/rest/deviceCommandAPI';

const CommandDevice = () => {
	const { t } = useTranslation();

	const tcpHealth = useTcpHealth();
	const { device, imei } = useConfigurationContext();

	const { data: respDeviceCommandRecentHistory, loading: respDeviceCommandRecentHistoryLoading, reload } = useHandleServer(['respDeviceCommandRecentHistory'], basicDeviceCommandRecentHistory, { refetchInterval: 3000 });

	const [command, setCommand] = useState<string>('');
	const [openIndex, setOpenIndex] = useState<number | null>(null);

	const canSendInstant = tcpHealth && device.activated && device.status;

	const handleSubmit = async (sendMode: 'instant' | 'on_connect') => {
		if (!command.trim()) return;

		const payload: DeviceSendCommandRequest = {
			send_mode: sendMode,
			command: command,
			imeis: [imei],
		};

		await basicDeviceSendCommand(payload);

		reload();
		setCommand('');
		setOpenIndex(0);
	};

	return (
		<div className="space-y-3">
			<LabeledField label="label.command" className="flex flex-col gap-1">
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
			</LabeledField>

			<div className="flex items-center gap-2">
				<GUIButton onClick={() => handleSubmit('instant')} type="button" disabled={!canSendInstant} className={`${!canSendInstant && '!opacity-40 !cursor-not-allowed'}`}>
					{t('label.send-instant')}
				</GUIButton>

				<GUIButton onClick={() => handleSubmit('on_connect')} type="button">
					{t('label.send-on-connect')}
				</GUIButton>
			</div>

			<CommandRecentHistory data={respDeviceCommandRecentHistory} loading={respDeviceCommandRecentHistoryLoading} reload={reload} openIndex={openIndex} onOpenIndexChange={setOpenIndex} />
		</div>
	);
};

export default CommandDevice;
