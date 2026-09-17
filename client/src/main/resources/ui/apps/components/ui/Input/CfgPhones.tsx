import React, { useEffect } from 'react';
import { useCallback, useState } from 'react';
import Delete from '@/components/@icons/delete';
import { toast } from 'react-toastify';
import Plus from '@/components/@icons/plus';
import PlaylistRemove from '@/components/@icons/playlist-remove';
import Colors from '@/constants/Color.constant';
import { arrayParseData } from '@/utils/ArrayUtils';
import { useTranslation } from 'react-i18next';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import Tooltip from '../Tooltip';

interface PhonesInputProps {
	/** Уникальный идентификатор конфигурационного параметра */
	uid: string;
	/** Текст placeholder для input */
	placeholder?: string;
	/** Дополнительные inline-стили */
	style?: React.CSSProperties;
	/** Дополнительные CSS-классы для input */
	className?: string;
	/** Флаг блокировки поля ввода */
	disabled?: boolean;
}

const MAX_PHONES = 4;

const CfgPhonesInput: React.FC<PhonesInputProps> = ({ uid, placeholder = '', style, className, disabled = false }) => {
	const { t } = useTranslation();

	const { imei, section, manager } = useConfigurationWrapperContext();
	const { value, handleChange } = useConfigurationField(imei, section, manager, uid);

	const [localValue, setLocalValue] = useState<string[]>([]);
	const [errors, setErrors] = useState<(string | null)[]>([]);

	useEffect(() => {
		const parsed = arrayParseData(value);
		setLocalValue(parsed);
		setErrors(parsed.map(() => null));
	}, [value]);

	const updatePhones = (newPhones: string[], newErrors?: (string | null)[]) => {
		setLocalValue(newPhones);
		if (newErrors) setErrors(newErrors);
		handleChange(newPhones);
	};

	const handlePhoneChange = useCallback(
		(index: number, e: React.ChangeEvent<HTMLInputElement>) => {
			const input = e.target.value.replace(/\s+/g, '');
			const phone = input.startsWith('+') ? input : '+' + input;

			const updated = [...localValue];
			updated[index] = phone === '+' ? '' : phone;

			let currentError: string | null = null;

			const updatedErrors = [...errors];
			updatedErrors[index] = currentError;

			updatePhones(updated, updatedErrors);
		},
		[localValue, errors, updatePhones]
	);

	const handleDelete = useCallback(
		(index: number) => {
			const updated = [...localValue];
			updated.splice(index, 1);

			const updatedErrors = [...errors];
			updatedErrors.splice(index, 1);

			updatePhones(updated, updatedErrors);
		},
		[localValue, updatePhones]
	);

	const handleCleanPhones = useCallback(() => {
		updatePhones([]);
	}, [updatePhones]);

	const handleAddPhone = () => {
		if (localValue.length >= MAX_PHONES) {
			toast.warning(t('message.max-phone-numbers-reached'));
			return;
		}

		updatePhones([...localValue, ''], [...errors, null]);
	};

	return (
		<React.Fragment>
			<div className="flex flex-col sm:flex-row sm:items-center gap-3">
				<p className="font-mornal">{t('message.list-phone-numbers')}:</p>
				<div className="flex items-center gap-3">
					<Tooltip title={t('message.add-phone-number')}>
						<div id="buttonhlp" className="!p-[6px] !h-auto">
							<Plus size={20} fill="#008c25" onClick={handleAddPhone} />
						</div>
					</Tooltip>

					<Tooltip title={t('message.clear-list')}>
						<div id="buttonhlp" className="!p-[6px] !h-auto">
							<PlaylistRemove size={20} fill={Colors.color_error_base} onClick={handleCleanPhones} />
						</div>
					</Tooltip>
				</div>
			</div>

			<div className="mt-4 sm:mt-2 space-y-2">
				{localValue.length > 0 ? (
					localValue.map((phone, index) => (
						<div key={index} className="relative">
							<div className="flex items-center gap-3">
								<input
									type="tel"
									name={uid}
									id={`cfg_${imei}_${Math.floor(Math.random() * 1000000)}_${index}`}
									value={phone}
									style={style}
									className={className}
									placeholder={placeholder}
									onChange={(e) => handlePhoneChange(index, e)}
									disabled={disabled}
								/>

								<div id="buttonhlp" className="!p-[8px] !h-auto" onClick={() => handleDelete(index)}>
									<Delete size={20} fill="#b50202" />
								</div>
							</div>

							{errors[index] && <span className="absolute z-[100] top-11 rounded bg-red-500 text-white p-1 text-sm max-w-[40vw]">{t(errors[index]!)}</span>}
						</div>
					))
				) : (
					<p className="flex justify-center my-[2rem] text-sm sm:text-base">{t('message.none-phone-numbers')}</p>
				)}
			</div>
		</React.Fragment>
	);
};

export default CfgPhonesInput;
