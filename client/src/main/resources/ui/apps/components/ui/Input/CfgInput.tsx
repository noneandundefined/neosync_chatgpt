import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface CFGInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	/** Уникальный идентификатор конфигурационного параметра */
	uid: string;
	valueParam?: string;
	min?: number;
	max?: number;
}

const CFGInput: React.FC<CFGInputProps> = ({ uid, valueParam, className, placeholder, style, disabled, readOnly, min, max }) => {
	const { t } = useTranslation();

	const { imei, section, manager } = useConfigurationWrapperContext();
	const { value, schema, isChanged, error, setError, handleChange } = useConfigurationField(imei, section, manager, uid);

	const cleanedValue = value !== null && value !== undefined ? String(value).replace(/\x00/g, '') : '';

	const isTenthsUint8 = uid === UIDs.NAV_FILTER_HDOP_TENTHS;
	const clampUint8 = (n: number) => Math.max(0, Math.min(255, n));

	const formatTenths = (raw: unknown) => {
		const num = Number(raw);
		if (!Number.isFinite(num)) return '';
		return (num / 10).toFixed(1);
	};

	const parseTenthsToUint8 = (s: string) => {
		const normalized = s.replace(',', '.').trim();
		if (normalized === '') return null;
		const num = Number(normalized);
		if (!Number.isFinite(num)) return null;
		return clampUint8(Math.round(num * 10));
	};

	return (
		<div className="flex flex-1 flex-col gap-2 relative">
			<input
				type="text"
				name={uid}
				id={`cfg_${imei}_${Math.floor(Math.random() * 1000000)}`}
				inputMode={isTenthsUint8 ? 'decimal' : schema?.is_digits ? 'numeric' : undefined}
				pattern={isTenthsUint8 ? undefined : schema?.is_digits ? '[0-9]*' : undefined}
				placeholder={placeholder}
				className={`input border p-2 rounded ${className} ${error ? 'border-red-500' : isChanged ? 'border-[#afb5c0] bg-[#f6f6f6]' : 'border-gray-300'}`}
				style={{
					...style,
					opacity: disabled ? 0.5 : 1,
					pointerEvents: disabled ? 'none' : 'auto',
				}}
				value={valueParam ? valueParam : uid === UIDs.AUTH_GLOBAL_PASS && cleanedValue === '0' ? '' : isTenthsUint8 ? (cleanedValue === '' ? '' : formatTenths(cleanedValue)) : cleanedValue}
				onChange={(e) => {
					let val = e.target.value;

					if (uid === UIDs.AUTH_GLOBAL_PASS) {
						const regex = /^[A-Za-z0-9!"#$%&'()*+\-./:;<=>?@[\\\]^_{|}~]{0,8}$/;
						if (!regex.test(val)) {
							setError({ message: 'message.requirements-pass' });
							return;
						}

						handleChange(val);
						return;
					}

					if (isTenthsUint8) {
						if (val.trim() === '') {
							handleChange('');
							return;
						}

						const parsed = parseTenthsToUint8(val);
						if (parsed === null) return;
						handleChange(parsed);
						return;
					}

					if (schema?.is_digits) {
						const num = Number(val);
						if (!isNaN(num)) {
							if (min !== undefined && num < min) val = String(min);
							if (max !== undefined && num > max) val = String(max);
						}
					}

					handleChange(val);
				}}
				disabled={disabled}
				readOnly={readOnly}
			/>

			{error && (
				<span className="absolute z-[100] top-10 rounded bg-red-500 text-white p-1 text-sm min-w-[12vw] max-w-[40vw]">
					{t(error.message ?? '')} {error.argv}
				</span>
			)}
		</div>
	);
};

export default CFGInput;
