import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useEnabled } from '@/hooks/Configuration/useEnabled';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { RS485_SELECT_MODBUS_BD_RATE, RS485_SELECT_MODBUS_BYTES_COUNT, RS485_SELECT_MODBUS_BYTES_ORDER, RS485_SELECT_MODBUS_REGISTERS_TYPE } from '@/constants/RS485.constant';

const Modbus = () => {
	const { t } = useTranslation();
	const { manager } = useConfigurationWrapperContext();

	const modbusRegisterAddrValue = useConfigurationVal(manager, UIDs.MODBUS_DEVICES_ADDR);
	const enabledStates = modbusRegisterAddrValue.map((val: any) => useEnabled(val !== 255));

	return (
		<div className="w-full flex items-center gap-4 flex-wrap my-3">
			{Array.from({ length: 5 }, (_, index) => {
				const { enabled, toggle } = enabledStates[index];

				return (
					<div className="flex-1 border border-[#e5e8eb] p-2 space-y-2 min-w-[51vw] sm:min-w-[20rem] max-w-[30rem]" key={index}>
						<div className="mb-2 flex items-center justify-between">
							<p className="font-medium">
								{t('label.sensor-short')} {index}
							</p>
							<GUISwitch checked={enabled} onChange={toggle} />
						</div>

						<LabeledField label="message.sensor-address" className="flex flex-col gap-1" disabled={!enabled}>
							<CFGArrayItemInput uid={UIDs.MODBUS_DEVICES_ADDR} index={index} className="max-h-[1.75rem]" disabled={!enabled} onFn={!enabled ? 255 : undefined} min={0} max={255} />
						</LabeledField>

						<LabeledField label="message.register-address" className="flex flex-col gap-1" disabled={!enabled}>
							<CFGArrayItemInput uid={UIDs.MODBUS_BYTES_ORDER} index={index} className="max-h-[1.75rem]" disabled={!enabled} min={0} max={65535} />
						</LabeledField>

						<LabeledField label="message.number-of-bytes" className="flex flex-col" disabled={!enabled}>
							<CFGSelect uid={UIDs.MODBUS_BYTES_COUNT} index={index} array={RS485_SELECT_MODBUS_BYTES_COUNT} className="mini min-w-full" disabled={!enabled} />
						</LabeledField>

						<LabeledField label="message.transmission-speed" className="flex flex-col" disabled={!enabled}>
							<CFGSelect uid={UIDs.MODBUS_DEVICES_BD_RATE} index={index} array={RS485_SELECT_MODBUS_BD_RATE} className="mini min-w-full" disabled={!enabled} />
						</LabeledField>

						<LabeledField label="message.register-type" className="flex flex-col" disabled={!enabled}>
							<CFGSelect uid={UIDs.MODBUS_REGISTERS_TYPE} index={index} array={RS485_SELECT_MODBUS_REGISTERS_TYPE} className="mini min-w-full" disabled={!enabled} />
						</LabeledField>

						<LabeledField label="message.byte-order" className="flex flex-col" disabled={!enabled}>
							<CFGSelect uid={UIDs.MODBUS_REGISTERS_ADDR_2} index={index} array={RS485_SELECT_MODBUS_BYTES_ORDER} className="mini min-w-full" disabled={!enabled} />
						</LabeledField>
					</div>
				);
			})}
		</div>
	);
};

export default Modbus;
