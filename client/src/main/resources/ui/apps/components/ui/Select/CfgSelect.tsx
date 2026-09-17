import GUISelect from './GUISelect';
import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { arrayParseData } from '@/utils/ArrayUtils';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import { PRESET_MAPPING } from '@/constants/PresetMapping.constant';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

function getPresetStoredValue(item: { id?: unknown; extra?: Record<string, unknown> }, writeFromExtra?: string): string {
	if (writeFromExtra != null && item?.extra != null && writeFromExtra in item.extra) {
		return String(item.extra[writeFromExtra]);
	}

	return String(item.id);
}

interface CFGSelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
	/** Уникальный идентификатор конфигурационного параметра */
	uid: string;
	/** Номер бита в маске */
	bit?: number;
	/** Индекс элемента массива конфигурации, который нужно редактировать */
	index?: number;
	/** Массив значений для отображения в списке */
	array: any[];
	/** Тип пресета для автоматического маппинга дополнительных параметров */
	presetType?: keyof typeof PRESET_MAPPING;
	/** Колбэк, вызываемый при изменении значения */
	onChange?: any;
	/** Значение для select */
	_value?: string;

	writeFromExtra?: string;
}

const CFGSelect: React.FC<CFGSelectProps> = ({ uid, bit, index = 0, array, presetType, onChange, _value, writeFromExtra }) => {
	const { t } = useTranslation();

	const { imei, section, manager } = useConfigurationWrapperContext();
	const { value, isChanged, handleChange } = useConfigurationField(imei, section, manager, uid);

	const valueArray: string[] = useMemo(() => arrayParseData(value).map(String), [value]);

	/** Get local value */
	const localCfgValue = useMemo(() => {
		if (bit != null) {
			return BitHelperUtils.checkBit(Number(value) || 0, bit) ? '1' : '0';
		}

		return _value ?? valueArray[index] ?? '';
	}, [value, bit, index, _value]);

	/** Handle select value */
	const handleSelectChange = useCallback(
		(e: React.ChangeEvent<HTMLSelectElement>) => {
			const selectedId = e.target.value;
			const selected = array.find((item) => String(item.id) === selectedId);

			if (!selected) return;

			/** Bit changed */
			if (bit != null) {
				const numericValue = Number(value) || 0;
				const selectedBitValue = Number(selectedId) as 0 | 1;

				const newValue = BitHelperUtils.changeBit(numericValue, bit, selectedBitValue);

				/** Save changes */
				handleChange(newValue);
				onChange?.(selectedId);
				return;
			}

			const updatedArray = [...valueArray];

			if (writeFromExtra && selected.extra?.[writeFromExtra] !== undefined) {
				updatedArray[index] = String(selected.extra[writeFromExtra]);
			} else {
				updatedArray[index] = selectedId;
			}

			handleChange(updatedArray);
			onChange?.(selectedId);

			if (presetType && selected.extra) {
				const mapping = PRESET_MAPPING[presetType];
				if (!mapping) return;

				Object.entries(selected.extra).forEach(([extraKey, extraValue]) => {
					const targetUid = mapping[extraKey];
					if (!targetUid) return;

					const targetField = manager.getField(targetUid);
					if (!targetField) return;

					const arr = Array.isArray(targetField.value) ? [...targetField.value] : arrayParseData(targetField.value);

					arr[index] = String(extraValue);
					manager.updateField(targetUid, arr);
				});
			}
		},
		[value, bit, index, uid, array, presetType, writeFromExtra, handleChange, onChange, valueArray]
	);

	const valueMatchesPreset = useMemo(() => {
		if (bit != null) return true;
		return array.some((item) => String(item.id) === String(localCfgValue) || getPresetStoredValue(item, writeFromExtra) === localCfgValue);
	}, [array, localCfgValue, bit, writeFromExtra]);

	const showCustomPresetOption = bit == null && writeFromExtra != null && localCfgValue !== '' && !valueMatchesPreset;

	const option = useMemo(() => {
		const base = array.map((item, optIdx) => (
			<option value={String(item.id)} key={String(item.id) || `opt-${optIdx}`}>
				{t(item.value)}
			</option>
		));
		if (!showCustomPresetOption) return base;
		return [
			...base,
			<option value={localCfgValue} key={`preset-custom-${index}-${localCfgValue}`}>
				{t('label.custom')}
			</option>,
		];
	}, [array, t, showCustomPresetOption, localCfgValue, index]);

	return (
		<div className={isChanged ? '[&>div>div]:!border-[#afb5c0]' : undefined}>
			<GUISelect value={localCfgValue} onChange={handleSelectChange}>
				{option}
			</GUISelect>
		</div>
	);
};

export default CFGSelect;
