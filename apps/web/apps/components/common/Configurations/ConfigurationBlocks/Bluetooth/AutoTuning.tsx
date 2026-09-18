import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { basicDeviceCommandFindBleSensors } from '@/rest/deviceCommandAPI';
import { useConfigurationContext } from '@/context/useConfigurationContext';
import { BLEAUTOCATCH_FIND_ALL_COMMAND, BLEAUTOCATCH_FIND_NEARBY_COMMAND, BLEAUTOCATCH_STOP_FIND_COMMAND } from '@/constants/Command.constant';

const AutoTuning = () => {
	const { imei } = useConfigurationContext();
	const { t } = useTranslation();

	const [searchStart, setSearchStart] = useState(false);
	const [searchSelected, setSearchSelected] = useState<string>(BLEAUTOCATCH_FIND_ALL_COMMAND);

	const handleBleSearch = async () => {
		try {
			setSearchStart(!searchStart);

			if (!searchStart) {
				await basicDeviceCommandFindBleSensors(imei, searchSelected);
			} else {
				await basicDeviceCommandFindBleSensors(imei, BLEAUTOCATCH_STOP_FIND_COMMAND);
			}
		} catch {
			setSearchStart(false);
		}
	};

	useEffect(() => {
		let timer: NodeJS.Timeout;

		if (searchStart) {
			timer = setTimeout(async () => {
				setSearchStart(false);
				try {
					await basicDeviceCommandFindBleSensors(imei, BLEAUTOCATCH_STOP_FIND_COMMAND);
				} catch (err: any) {
					console.error(err);
				}
			}, 150000);
		}

		return () => clearTimeout(timer);
	}, [searchStart, imei]);

	return (
		<>
			<div className="w-full sm:max-w-[20rem] space-y-3">
				<p>{t('message.select-sensor-search-radius')}</p>
				<div className="space-y-3">
					<div>
						<label className="flex items-center gap-3">
							<input type="radio" className="w-auto" value={BLEAUTOCATCH_FIND_ALL_COMMAND} checked={searchSelected === BLEAUTOCATCH_FIND_ALL_COMMAND} onChange={(e) => setSearchSelected(e.target.value)} />
							{t('message.search-all')}
						</label>
					</div>

					<div>
						<label className="flex items-center gap-3">
							<input type="radio" className="w-auto" value={BLEAUTOCATCH_FIND_NEARBY_COMMAND} checked={searchSelected === BLEAUTOCATCH_FIND_NEARBY_COMMAND} onChange={(e) => setSearchSelected(e.target.value)} />
							{t('message.search-nearby')}
						</label>
					</div>

					<button className="bg-[#eee] text-[12px] text-black py-[4px] px-5 rounded-[6px] uppercase text-center cursor-pointer hover:bg-[#e3e3e3] transition border border-[#ccc]" onClick={handleBleSearch}>
						{searchStart ? t('message.ble-btn-search-stop') : t('message.ble-btn-search-start')}
					</button>

					<div className="flex items-center gap-4">
						<p className="whitespace-nowrap text-[13px]">{t('message.ble-search-progress')}:</p>
						<div className="w-full h-[4px] bg-gray-300 overflow-hidden">
							<div className={`h-[4px] bg-blue-500 ${searchStart ? 'transition-all duration-[150000ms]' : ''}`} style={{ width: searchStart ? '100%' : '0%' }} />
						</div>
					</div>
				</div>
			</div>

			<p className="text-sm w-full font-semibold text-[#9a0000] mt-3">{t('message.autosetting-bluetooth')}</p>
			<p className="w-full font-semibold text-[#9a0000] mt-3">{t('message.rs485-ble-priority-warning')}</p>
		</>
	);
};

export default AutoTuning;
