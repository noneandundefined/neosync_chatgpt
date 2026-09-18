import { useState } from 'react';
import PageLayout from '../../PageLayout';
import { DeviceShort } from '@/rest/deviceCommandAPI';
import CommandSend from '@/components/common/DeviceCommands/CommandSend';
import CommandDevices from '@/components/common/DeviceCommands/CommandDevices';

const DeviceCommandPage = () => {
	const [selectedDevice, setSelectedDevice] = useState<DeviceShort[]>([]);

	return (
		<PageLayout>
			<div className="flex flex-col md:flex-row flex-1 min-h-0 md:h-full gap-5 max-md:overflow-y-auto">
				{/* Devices for send command */}
				<CommandDevices selectedDevice={selectedDevice} setSelectedDevice={setSelectedDevice} />

				{/* Send command */}
				<CommandSend selectedDevice={selectedDevice} />
			</div>
		</PageLayout>
	);
};

export default DeviceCommandPage;
