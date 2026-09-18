import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import GUISelect from '@/components/ui/Select/GUISelect';
import { MULTI_APN_TYPE } from '@/constants/Sim.constant';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

export function GetModeMultiApn(value: number): number {
	const bit3 = BitHelperUtils.checkBit(value, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_0) ? 1 : 0;
	const bit4 = BitHelperUtils.checkBit(value, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_1) ? 1 : 0;

	return (bit4 << 1) | bit3;
}

export function SetModeMultiApn(value: number, mode: number): number {
	const bit3 = mode & 1;
	const bit4 = (mode >> 1) & 1;

	let updated = value;

	updated = BitHelperUtils.clearBit(updated, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_0);
	updated = BitHelperUtils.clearBit(updated, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_1);

	if (bit3) updated = BitHelperUtils.setBit(updated, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_0);
	if (bit4) updated = BitHelperUtils.setBit(updated, BITs.DEVICE_FUNCTION_3.MULTI_APN_MODE_BIT_1);

	return updated;
}

const MultiApn = () => {
	const { t } = useTranslation();
	const { manager, imei, section } = useConfigurationWrapperContext();

	const { value: deviceFunction3, handleChange: handleDeviceFunction3Change, saveDraft: saveDeviceFunction3Draft } = useConfigurationField(imei, section, manager, UIDs.DEVICE_FUNCTION_3);

	const currentMode = GetModeMultiApn(deviceFunction3);
	const isEnabled = currentMode !== 0;

	const handleCheckboxToggle = (e: React.ChangeEvent<HTMLInputElement>) => {
		const checked = e.target.checked;
		const newMode = checked ? 1 : 0;

		const updated = SetModeMultiApn(deviceFunction3, newMode);

		handleDeviceFunction3Change(updated);
		saveDeviceFunction3Draft(updated);
	};

	const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
		const newMode = Number(e.target.value);
		const updated = SetModeMultiApn(deviceFunction3, newMode);

		handleDeviceFunction3Change(updated);
		saveDeviceFunction3Draft(updated);
	};

	return (
		<div className="w-full sm:min-w-[22rem] space-y-4">
			<LabeledField label="message.advanced-apn" className="flex justify-between sm:justify-between items-center gap-5">
				<GUISwitch checked={isEnabled} onChange={handleCheckboxToggle} />
			</LabeledField>

			<LabeledField label="message.multiapn-operating-mode" className="flex flex-col space-y-1">
				<GUISelect value={currentMode} onChange={handleChange}>
					{MULTI_APN_TYPE.map((type, index) => (
						<option key={index} value={type.id}>
							{t(type.value)}
						</option>
					))}
				</GUISelect>
			</LabeledField>

			{Array.from({ length: 2 }, (_, index) => (
				<LabeledField label="message.working-time-apn" argv={`APN${index}`} className="flex-col" disabled={currentMode !== 1}>
					<div className="flex items-center gap-2">
						<CFGArrayItemInput index={index} uid={UIDs.SIMTIME_MIN_CONF} disabled={currentMode !== 1} />
						<label>{t('label.second-short')}</label>
					</div>
				</LabeledField>
			))}
		</div>
	);
};

export default MultiApn;
