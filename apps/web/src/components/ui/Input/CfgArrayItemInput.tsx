import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { useEffect, useRef, useState } from 'react';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface CFGArrayItemInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	/** Уникальный идентификатор конфигурационного параметра */
	uid: string;
	/** Индекс элемента массива конфигурации, который нужно редактировать */
	index: number;
	/** Доп. функция для работы с input */
	onFn?: any;
	/** Преобразование текста при вставке из буфера */
	transformPaste?: (text: string) => string;

	min?: number;
	max?: number;
}

const clampUint8 = (n: number) => Math.max(0, Math.min(255, n));

const metersFromStoredHeight = (raw: number, index: number) => raw * (index === 0 ? 10 : 100);

const storedFromMetersHeight = (meters: number, index: number) => clampUint8(Math.floor(meters / (index === 0 ? 10 : 100)));

const CFGArrayItemInput: React.FC<CFGArrayItemInputProps> = ({ uid, index = 0, placeholder = '', style, className, onFn, value: externalValue, disabled = false, readOnly = false, min, max, transformPaste }) => {
	const { t } = useTranslation();

	const { imei, section, manager } = useConfigurationWrapperContext();
	const { value, schema, isChanged, error, handleChange } = useConfigurationField(imei, section, manager, uid, index);

	const [heightDraft, setHeightDraft] = useState('');
	const isHeightField = uid === UIDs.NAV_FILTER_VALID_HEIGHT_DECA_HECTO_METERS;
	const isHeightFocusedRef = useRef(false);

	const rawValue = externalValue !== undefined ? externalValue : Array.isArray(value) ? (value[index] ?? '') : value;

	useEffect(() => {
		if (!isHeightField || isHeightFocusedRef.current) return;

		const num = Number(rawValue);
		setHeightDraft(Number.isFinite(num) && rawValue !== '' ? String(metersFromStoredHeight(num, index)) : '');
	}, [rawValue, index, isHeightField]);

	let currentValue: number | string;

	if (uid === UIDs.SIMTIME_MIN_CONF) {
		currentValue = Math.floor(rawValue / 60);
	} else if (isHeightField) {
		currentValue = heightDraft;
	} else {
		currentValue = rawValue;
	}

	const commitHeightDraft = (input: string) => {
		const trimmed = input.trim();
		const meters = trimmed === '' ? 0 : Number(trimmed);

		if (trimmed !== '' && !Number.isFinite(meters)) return;

		const stored = storedFromMetersHeight(meters, index);
		const field = manager?.getField(uid);
		const baseArr = Array.isArray(field?.value) ? [...field!.value] : [];

		baseArr[index] = stored;
		handleChange(baseArr);
		setHeightDraft(trimmed === '' ? '' : String(metersFromStoredHeight(stored, index)));
	};

	useEffect(() => {
		if (onFn !== undefined) {
			if (Array.isArray(value)) {
				const arr = [...value];
				arr[index] = onFn;
				handleChange(arr);
			} else {
				handleChange(onFn);
			}
		}
	}, [onFn]);

	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const inputValue = e.target.value;

		if (isHeightField) {
			setHeightDraft(inputValue);
			return;
		}

		let newValue: any = inputValue;

		if (schema?.is_digits) {
			let numericValue = Number(newValue);

			if (!isNaN(numericValue)) {
				if (min !== undefined && numericValue < min) numericValue = min;
				if (max !== undefined && numericValue > max) numericValue = max;

				if (uid === UIDs.SIMTIME_MIN_CONF) {
					numericValue = numericValue * 60;
				}

				newValue = numericValue;
			}
		}

		const field = manager?.getField(uid);
		const baseArr = Array.isArray(field?.value) ? [...field!.value] : [];

		baseArr[index] = newValue;
		handleChange(baseArr);
	};

	const onPaste = (e: React.ClipboardEvent<HTMLInputElement>) => {
		if (!transformPaste || disabled || readOnly) return;

		e.preventDefault();
		const sanitized = transformPaste(e.clipboardData.getData('text'));

		if (isHeightField) {
			setHeightDraft(sanitized);
			return;
		}

		const field = manager?.getField(uid);
		const baseArr = Array.isArray(field?.value) ? [...field!.value] : [];

		baseArr[index] = sanitized;
		handleChange(baseArr);
	};

	return (
		<div className="relative flex flex-1 flex-col gap-2">
			<input
				type="text"
				inputMode={schema?.is_digits ? 'numeric' : undefined}
				pattern={schema?.is_digits ? '[0-9]*' : undefined}
				name={`${uid}_${Math.floor(Math.random() * 1000000)}_${index}`}
				id={`cfg_${imei}_${Math.floor(Math.random() * 1000000)}_${index}`}
				placeholder={placeholder}
				className={`input border p-2 rounded ${className} ${error ? 'border-red-500' : isChanged ? 'border-[#afb5c0] bg-[#f6f6f6]' : 'border-gray-300'}`}
				style={{
					...style,
					opacity: disabled ? 0.5 : 1,
					pointerEvents: disabled ? 'none' : 'auto',
				}}
				value={currentValue}
				onChange={onChange}
				onPaste={onPaste}
				onFocus={() => {
					if (isHeightField) isHeightFocusedRef.current = true;
				}}
				onBlur={() => {
					if (!isHeightField) return;
					isHeightFocusedRef.current = false;
					commitHeightDraft(heightDraft);
				}}
				disabled={disabled}
				readOnly={readOnly}
			/>

			{error && (
				<span className="absolute z-[100] top-10 rounded bg-red-500 text-white p-1 text-sm max-w-[40vw]">
					{t(error.message ?? '')} {error.argv}
				</span>
			)}
		</div>
	);
};

export default CFGArrayItemInput;
