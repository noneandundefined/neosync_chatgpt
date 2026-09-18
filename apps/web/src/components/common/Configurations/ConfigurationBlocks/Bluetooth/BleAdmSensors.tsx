import Sensor from './Sensor';
import { MAX_HEIGHT } from '.';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { useEffect, useRef, useState } from 'react';
import { BLE_SENSOR_TYPES } from '@/constants/BleSensorTypes.constant';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { TelemetryResponse } from '@/interface/configuration/configurationTelemetryResponse.interface';

interface BleSensorsAdmProps {
	telemetry: TelemetryResponse | null;
}

const BleSensorsAdm: React.FC<BleSensorsAdmProps> = ({ telemetry }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	const wrapperRef = useRef<HTMLDivElement | null>(null);
	const [isOverflow, setIsOverflow] = useState<boolean>(false);

	const bleSensorTypes = useConfigurationVal(manager, UIDs.BLE_SENSOR_TYPE);
	const { value: bleSensorsAdm, handleChange: handleBleSensorsAdmChange, saveDraft: saveBleSensorsAdmDraft } = useConfigurationField(imei, section, manager, UIDs.BLE_ADM_SENSOR_ADDRESS_LIST);

	const [sensors, setSensors] = useState<string[] | null>(bleSensorsAdm);

	useEffect(() => {
		if (bleSensorsAdm) {
			setSensors(bleSensorsAdm);
		}
	}, [bleSensorsAdm]);

	const handleRemove = async (index: number) => {
		setSensors((prev) => {
			if (!prev) return null;

			const updatedSensors = prev.filter((_, i) => i !== index);
			handleBleSensorsAdmChange(updatedSensors);
			saveBleSensorsAdmDraft(updatedSensors);
			return updatedSensors;
		});
	};

	useEffect(() => {
		if (!wrapperRef.current) return;

		const height = wrapperRef.current.scrollHeight;
		setIsOverflow(height > MAX_HEIGHT);
	}, [sensors]);

	return (
		<div className="w-full min-h-[50px]" ref={wrapperRef}>
			{sensors && sensors.length > 0 ? (
				<div className={`${isOverflow ? 'flex flex-wrap flex-row' : 'flex flex-col'} gap-6 md:max-h-[600px] max-h-auto items-start content-start`}>
					{sensors.map((sensor, index) => {
						const typeCode = bleSensorTypes?.[index] ?? 255;
						const typeLabel = BLE_SENSOR_TYPES[typeCode];

						const sensorTelemetry = telemetry?.blesensorinfo?.[sensor];

						return (
							<div key={index} className="w-full md:w-auto">
								<Sensor
									index={index}
									sensor={sensor}
									typeLabel={typeLabel}
									onRemove={() => handleRemove(index)}
									telemetry={sensorTelemetry}
									fields={[
										{
											label: 'label.voltage',
											value: sensorTelemetry?.voltage ?? t('label.no-data'),
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
				<p className="font-medium text-[#ccc]">{t('message.ble-adm-none')}</p>
			)}
		</div>
	);
};

export default BleSensorsAdm;
