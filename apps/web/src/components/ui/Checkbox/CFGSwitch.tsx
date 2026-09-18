import GUISwitch from './GUISwitch';
import { useEffect, useState } from 'react';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { BitHelperUtils } from '@/utils/BitHelperUtils';

interface CFGSwitchProps {
	/** Уникальный идентификатор конфигурационного параметра */
	uid: string;
	/** Номер бита в маске */
	bit?: number;
	/** Номера битов в маске */
	bits?: number[];
	/** Индекс элемента массива конфигурации, который нужно редактировать */
	index?: number;
	/** Флаги значений при off */
	vOff?: number | string;
	/** Флаги значений при on */
	vOn?: number | string;
	/** Заблокирован ли элемент */
	disabled?: boolean;
	/** Доп. функция для работы с кастомными переключениями */
	onChange?: (enabled: boolean, rawValue: any) => void;
	/** Доп. функция для получения значения при первом render */
	onVisible?: (enabled: boolean) => void;
}

const CFGSwitch: React.FC<CFGSwitchProps> = ({ uid, bit, bits, index, vOff, vOn, disabled, onChange, onVisible }) => {
	const { imei, section, manager } = useConfigurationWrapperContext();
	const { value, handleChange, saveDraft } = useConfigurationField(imei, section, manager, uid);

	const [checked, setChecked] = useState<boolean>(false);

	useEffect(() => {
		let enabled = false;

		if (Array.isArray(value) && index != null) {
			enabled = value[index] === vOn;
		} else if (bit != null) {
			if (value != null && vOff !== undefined && vOn !== undefined) {
				const bitValue = BitHelperUtils.checkBit(value, bit) ? 1 : 0;
				enabled = bitValue === Number(vOn);
			} else {
				enabled = value != null ? BitHelperUtils.checkBit(value, bit) : false;
			}
		} else if (bits?.length) {
			enabled = value != null ? bits.some((b) => BitHelperUtils.checkBit(value, b)) : false;
		} else {
			enabled = value === vOn;
		}

		setChecked(enabled);
		onVisible?.(enabled);
	}, [value, bit, index, vOn]);

	const handleToggle = async () => {
		let newValue;

		if (Array.isArray(value) && index != null) {
			const copy = [...value];
			copy[index] = checked ? vOff! : vOn!;
			newValue = copy;
		} else if (bit != null) {
			const numericValue = Number(value) || 0;
			newValue = numericValue ^ (1 << bit);
		} else {
			newValue = checked ? vOff! : vOn!;
		}

		setChecked(!checked);

		handleChange(newValue);
		await saveDraft(newValue);

		onChange?.(!checked, newValue);
	};

	return <GUISwitch checked={checked} onChange={handleToggle} disabled={disabled} />;
};

export default CFGSwitch;
