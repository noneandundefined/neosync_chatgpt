import React from 'react';
import InDoor from './InDoor';
import BleRelay from './BleRelay';
import AutoTuning from './AutoTuning';
import { useRef, useState } from 'react';
import BleSensorsAdm from './BleAdmSensors';
import BleFuelSensors from './BleFuelSensors';
import { useTranslation } from 'react-i18next';
import IndexConfiguration from '../IndexConfiguration';
import GUISelect from '@/components/ui/Select/GUISelect';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

export const MAX_HEIGHT = 600;

const ConfigBluetoothInner = () => {
	const { t } = useTranslation();
	const { imei, support, isTemplate } = useConfigurationWrapperContext();

	const { data: respConfigurationTelemetry } = useConfigurationTelemetry(imei, !isTemplate);

	const [bleType, setBleType] = useState<string>('mode-ble-sensor');

	const wrapperRef = useRef<HTMLDivElement | null>(null);

	return (
		<div>
			{support.InDoor && <IndexConfiguration title={t('message.indoor-navigation')} content={<InDoor />} />}

			<div className="flex flex-col sm:flex-row gap-3 sm:items-center">
				<p className="whitespace-nowrap">{t('label.mode-type-configurable')}</p>
				<GUISelect onChange={(e) => setBleType(e.target.value)}>
					<option value="mode-ble-sensor">{t('label.ble-sensors')}</option>
					{support.BleRele && <option value="mode-ble-relay">{t('label.ble-relay')}</option>}
				</GUISelect>
			</div>

			{bleType == 'mode-ble-sensor' ? (
				<React.Fragment>
					{!isTemplate && <IndexConfiguration title={t('label.auto-tuning')} content={<AutoTuning />} />}

					<div ref={wrapperRef} className={`flex flex-col gap-4`}>
						<div>
							<p className="font-medium my-3">{t('label.ble-adm-sensor')}</p>
							<BleSensorsAdm telemetry={respConfigurationTelemetry} />
						</div>
						<div>
							<p className="font-medium my-3">{t('label.ble-fuel-sensor')}</p>
							<BleFuelSensors telemetry={respConfigurationTelemetry} />
						</div>
					</div>
				</React.Fragment>
			) : (
				<IndexConfiguration title={t('label.ble-relay')} content={<BleRelay />} maxWidth="30rem" />
			)}
		</div>
	);
};

const ConfigBluetooth = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigBluetoothInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigBluetooth;
