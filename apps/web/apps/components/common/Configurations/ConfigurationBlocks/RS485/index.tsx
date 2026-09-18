import LLS from './LLS';
import React from 'react';
import Modbus from './Modbus';
import { useState } from 'react';
import ADM20Reader from './ADM20Reader';
import LLSTariration from './LLSTariration';
import { useTranslation } from 'react-i18next';
import IndexConfiguration from '../IndexConfiguration';
import GUISelect from '@/components/ui/Select/GUISelect';
import { useConfigurationTelemetry } from '@/hooks/Configuration/useConfigurationTelemetry';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigRS485Inner = () => {
	const { t } = useTranslation();
	const { imei, support, isTemplate } = useConfigurationWrapperContext();

	const [rs485Type, setRs485Type] = useState<string>('lls');

	const { data: respConfigurationTelemetry } = useConfigurationTelemetry(imei, !isTemplate);

	return (
		<div>
			<div className="flex flex-col sm:flex-row gap-3 sm:items-center">
				<p className="whitespace-nowrap">{t('message.type-configurable-sensor')}</p>
				<GUISelect onChange={(e) => setRs485Type(e.target.value)}>
					{support.LLS && <option value="lls">LLS</option>}
					{support.ADM20 && <option value="adm20-reader">{t('label.adm20-reader')}</option>}
					{support.Modbus && <option value="modbus">Modbus</option>}
				</GUISelect>
			</div>

			{rs485Type === 'lls' ? (
				<React.Fragment>
					<LLS telemetry={respConfigurationTelemetry} />
					<IndexConfiguration title={t('label.tariration')} content={<LLSTariration />} />
				</React.Fragment>
			) : rs485Type === 'adm20-reader' ? (
				<IndexConfiguration title="ADM20" content={<ADM20Reader telemetry={respConfigurationTelemetry} />} maxWidth="60rem" />
			) : (
				<Modbus />
			)}
		</div>
	);
};

const ConfigRS485 = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigRS485Inner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigRS485;
