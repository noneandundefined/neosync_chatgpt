import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { basicDeviceCommandOWireSensors } from '@/rest/deviceCommandAPI';
import { useConfigurationContext } from '@/context/useConfigurationContext';

const AutoSettingTemp = () => {
	const { t } = useTranslation();
	const { imei } = useConfigurationContext();

	const [loading, setLoading] = useState<boolean>(false);

	const handleFind = async () => {
		try {
			setLoading(true);
			await basicDeviceCommandOWireSensors(imei);
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="max-w-[25rem] my-4">
			<div className="min-w-[10rem]">
				<button
					className="bg-[#eee] text-[12px] text-black font-medium py-[5px] px-5 rounded-[6px] uppercase text-center cursor-pointer hover:bg-[#e3e3e3] transition border border-[#ccc]"
					disabled={loading}
					onClick={handleFind}
				>
					{t('message.start-temperature-sensor-search')}
				</button>
			</div>
		</div>
	);
};

export default AutoSettingTemp;
