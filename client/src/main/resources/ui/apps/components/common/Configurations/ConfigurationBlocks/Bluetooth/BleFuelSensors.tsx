import Sensor from './Sensor';
import { MAX_HEIGHT } from '.';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { useEffect, useRef, useState } from 'react';
import { FuelAddressesParse } from '@/utils/FuelAddressesParse.Utils';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { TelemetryResponse } from '@/interface/configuration/configurationTelemetryResponse.interface';

interface ParsedSensor {
	address: string;
	typeLabel: string;
	index: number;
}

interface BleFuelSensorsProps {
	telemetry: TelemetryResponse | null;
}

const BleFuelSensors: React.FC<BleFuelSensorsProps> = ({ telemetry }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	// wrapper
	const wrapperRef = useRef<HTMLDivElement | null>(null);
	const [isOverflow, setIsOverflow] = useState<boolean>(false);

	// configurations
	const { value: bleSensorsFuel, handleChange: handleBleSensorsFuelChange, saveDraft: saveBleSensorsFuelDraft } = useConfigurationField(imei, section, manager, UIDs.FUEL_SENSOR_ADDRESS_LIST);

	const bleSensorTypesMask = useConfigurationVal(manager, UIDs.BLE_SENSOR_TYPE);

	const [sensors, setSensors] = useState<ParsedSensor[]>([]);

	useEffect(() => {
		if (!bleSensorsFuel || bleSensorsFuel.length === 0) {
			setSensors([]);
			return;
		}

		const { macs } = FuelAddressesParse(Array.from(bleSensorsFuel));

		const parsedSensors: ParsedSensor[] = macs.map((mac) => {
			const address = mac.value.map((b) => b.toString(16).padStart(2, '0').toUpperCase()).join(':');

			return {
				address,
				typeLabel: t('label.adm-unknown'),
				index: mac.position,
			};
		});

		setSensors(parsedSensors);
	}, [bleSensorsFuel, bleSensorTypesMask]);

	const handleRemove = async (idx: number) => {
		const arr = Array.from(bleSensorsFuel);

		for (let i = 0; i < 6; i++) {
			arr[idx + i] = 0;
		}

		handleBleSensorsFuelChange(arr);
		await saveBleSensorsFuelDraft(arr);
	};

	useEffect(() => {
		if (!wrapperRef.current) return;

		const h = wrapperRef.current.scrollHeight;
		setIsOverflow(h > MAX_HEIGHT);
	}, [sensors]);

	return (
		<div className="w-full" ref={wrapperRef}>
			{sensors && sensors.length > 0 ? (
				<div className={`${isOverflow ? 'flex flex-wrap flex-row' : 'flex flex-col'} gap-6 md:max-h-[600px] items-start content-start`}>
					{sensors.map((sensor, index) => {
						const sensorTelemetry = telemetry?.fuelinfo?.[sensor.address];

						return (
							<div key={index} className="w-full md:w-auto">
								<Sensor
									index={index}
									sensor={sensor.address}
									typeLabel="label.adm-unknown"
									onRemove={() => handleRemove(sensor.index)}
									telemetry={sensorTelemetry}
									fields={[
										{
											label: 'label.fuel-level',
											value: sensorTelemetry?.fuel_level ?? t('label.no-data'),
											unit: 'label.units',
										},
										{
											label: 'label.fuel-temperature',
											value: sensorTelemetry?.temp ?? t('label.no-data'),
											unit: 'label.degrees-celsius',
										},
										{
											label: 'label.sensor-battery-level',
											value: sensorTelemetry?.voltage != null ? sensorTelemetry.voltage.toFixed(1) : t('label.no-data'),
											unit: 'label.v',
										},
										{
											label: 'label.rssi',
											value: sensorTelemetry?.rssi ?? t('label.no-data'),
											unit: 'label.dBm',
										},
									]}
								/>
							</div>
						);
					})}
				</div>
			) : (
				<p className="font-medium text-[#ccc]">{t('message.ble-fuel-none')}</p>
			)}
		</div>
	);
};

export default BleFuelSensors;
