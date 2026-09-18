import React from 'react';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { clearMacBlock } from '@/utils/СlearMacUtils';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { LLS_COUNT } from '@/constants/DeviceModels.constant';
import { FuelSensorTypesParse } from '@/utils/FuelAddressesParse.Utils';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { TelemetryResponse } from '@/interface/configuration/configurationTelemetryResponse.interface';

interface LLSProps {
	telemetry: TelemetryResponse | null;
}

const LLS: React.FC<LLSProps> = ({ telemetry }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	const getValueOrNoData = (v: any) => (v === 0 ? 0 : (v ?? t('label.no-data')));

	const { value: fuelSensorTypeVal, handleChange: handleFuelSensorTypeVal } = useConfigurationField(imei, section, manager, UIDs.FUEL_SENSOR_TYPE);

	const { value: fuelAddressList, schema: fuelAddressListSchema, handleChange: handleChangeFuelAddressList, saveDraft: saveFuelAddressList } = useConfigurationField(imei, section, manager, UIDs.FUEL_SENSOR_ADDRESS_LIST);

	const normalizeFuelArray = (): number[] => {
		if (Array.isArray(fuelAddressList)) {
			return fuelAddressList.map((val) => Number(val) || 0);
		}

		if (typeof fuelAddressList === 'string') {
			try {
				const parsed = JSON.parse(fuelAddressList);
				if (Array.isArray(parsed)) {
					return parsed.map((val) => Number(val) || 0);
				}
			} catch {
				return [];
			}
		}

		return Array.isArray(fuelAddressList) ? fuelAddressList : [];
	};

	const baseArray = normalizeFuelArray();
	const fuelSensorType = FuelSensorTypesParse(fuelSensorTypeVal.toString(2));

	const getLLSValue = (rs485Index: number) => {
		const pos = rs485Index * 6;
		return baseArray[pos] ?? 0;
	};

	const onChangeLLS = (rs485Index: number, val: number) => {
		const arr = [...baseArray];
		const pos = rs485Index * 6;
		arr[pos] = val;
		handleChangeFuelAddressList(arr);
	};

	const onToggleRS485 = async (index: number, enabled: boolean) => {
		const newMask = BitHelperUtils.changeBit(fuelSensorTypeVal, index, enabled ? 0 : 1);
		handleFuelSensorTypeVal(newMask);

		const arr = [...baseArray];
		const pos = index * 6;

		clearMacBlock(arr, pos);

		arr[pos] = 0;
		handleChangeFuelAddressList(arr);
		await saveFuelAddressList(arr);
	};

	return (
		<div className="my-3">
			<p className="w-full font-semibold text-[#9a0000]">{t('message.rs485-ble-priority-warning')}</p>

			<p className="font-medium my-5">{t('label.setting-short')}</p>
			<div className="w-full flex items-center gap-4 flex-wrap">
				{Array.from({ length: LLS_COUNT }, (_, index) => {
					const isBle = BitHelperUtils.checkBit(fuelSensorTypeVal, index);
					const isEnabled = !isBle;

					const isRS485 = fuelSensorType.rs485.includes(index);

					const sensorKey = (index + 1).toString();
					const sensorTelemetry = telemetry?.fuelinfo?.[sensorKey];
					const hasIncompleteData = !sensorTelemetry || sensorTelemetry.temp == null || sensorTelemetry.fuel_level == null;
					const temp = hasIncompleteData ? t('label.no-data') : getValueOrNoData(sensorTelemetry.temp);
					const fuel_level = hasIncompleteData ? t('label.no-data') : getValueOrNoData(sensorTelemetry.fuel_level);

					return (
						<div key={index} className="flex flex-1 flex-col gap-1 min-w-[51vw] sm:min-w-[20rem] max-w-[40rem] border border-[#e5e8eb] p-1">
							<div className="w-full sm:w-auto p-2 space-y-3">
								<p className="font-medium text-[14px]">
									{t('label.sensor-short')} {index}
								</p>
							</div>

							<div className={`w-full items-center gap-4 ${!isEnabled ? 'opacity-40' : ''}`}>
								<div className="p-2">
									<div className="flex items-center justify-between gap-3">
										<p className="text-[13px]">
											{t('message.fuel-level')} {index}
										</p>
										<div className="flex items-center gap-2">
											<input type="text" value={fuel_level} className="mini min-w-[90%]" disabled />
											<label>{t('label.units')}</label>
										</div>
									</div>
									<div className="flex items-center justify-between gap-3">
										<p className="text-[13px]">
											{t('label.temperature')} {index}
										</p>
										<div className="flex items-center gap-2">
											<input type="text" value={temp} className="mini min-w-[90%]" disabled />
											<label>{t('label.degrees-celsius')}</label>
										</div>
									</div>
								</div>
							</div>

							<div className="w-full sm:w-auto border border-[#e5e8eb] p-2 space-y-3">
								<div className="flex items-center justify-between">
									<p className="text-[14px]">{t('message.lls-usage')}</p>
									<GUISwitch checked={isEnabled} onChange={(e) => onToggleRS485(index, e.target.checked)} />
								</div>

								<LabeledField label="label.address" argv={index} className="flex flex-col sm:flex-row sm:items-center gap-5" disabled={!isEnabled}>
									<input
										type="text"
										inputMode={fuelAddressListSchema?.is_digits ? 'numeric' : undefined}
										pattern={fuelAddressListSchema?.is_digits ? '[0-9]*' : undefined}
										name={`cfg_${imei}_${Math.floor(Math.random() * 1000000)}`}
										id={`cfg_${imei}_${Math.floor(Math.random() * 1000000)}`}
										className={`input border p-2 rounded mini border-gray-300`}
										value={isRS485 ? getLLSValue(index) : 0}
										style={{
											opacity: !isEnabled ? 0.5 : 1,
											pointerEvents: !isEnabled ? 'none' : 'auto',
										}}
										onChange={(e) => {
											if (!isRS485) return;

											let val = Number(e.target.value);
											if (isNaN(val)) val = 0;

											if (val < 0 || val > 255) return;
											onChangeLLS(index, val);
										}}
										disabled={!isEnabled}
									/>
								</LabeledField>
							</div>
						</div>
					);
				})}
			</div>
		</div>
	);
};

export default LLS;
