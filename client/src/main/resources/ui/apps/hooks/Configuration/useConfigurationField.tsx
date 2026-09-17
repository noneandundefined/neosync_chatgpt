import { UIDs } from '@/constants/UID.constant';
import { EVENT_CONFIGURATION_DRAFT } from '../useDraftCheck';
import { SerializeByType } from '@/utils/SerializeByTypeUtils';
import { saveTemplateDraft } from '@/utils/TemplateDraftUtils';
import { basicConfigurationInsertDraft } from '@/rest/configurationAPI';
import { useCallback, useEffect, useRef, useState, useContext } from 'react';
import { ConfigurationManager, FieldMerged } from '@/utils/ConfigurationManagerUtils';
import { ConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { FieldType, FieldType_Ascii, FieldType_Phones, FieldType_TwoUint16, FieldType_TwoUint8, FieldType_UInt16, FieldType_UInt32, FieldType_Uint32Array, FieldType_UInt8 } from '@/constants/FieldType.constant';

export type ValidationSchemaResult = {
	message?: string;
	argv?: string;
};

export const useConfigurationField = (imei: string, section: string, manager: ConfigurationManager, uid: string, index?: number) => {
	const wrapperContext = useContext(ConfigurationWrapperContext);

	const [value, setValue] = useState<any>('');
	const [schema, setSchema] = useState<FieldMerged | null>(null);
	const [isChanged, setIsChanged] = useState(false);
	const [error, setError] = useState<ValidationSchemaResult | null>(null);
	const timerRef = useRef<NodeJS.Timeout | null>(null);

	useEffect(() => {
		if (!manager) return;

		const update = () => {
			const field = manager.getField(uid);
			if (!field || field.value === undefined) return;

			setValue((prev: any) => {
				if (JSON.stringify(prev) !== JSON.stringify(field.value)) {
					return field.value;
				}

				return prev;
			});

			setSchema(field);
			setIsChanged(manager.isFieldChanged(uid, index));
		};

		const unsubscribe = manager.subscribe(update);
		update();

		return () => {
			unsubscribe();
		};
	}, [manager, uid, index]);

	const validate = (value: string, schema: FieldMerged): ValidationSchemaResult | null => {
		if (!schema) return null;

		if (uid === UIDs.SERVER_HOST) {
			const host = Array.isArray(value) ? String(value[index ?? 0] ?? '') : String(value ?? '');

			if (/^https?:\/\//i.test(host.trim())) {
				return { message: 'message.rule-remove-http-prefix' };
			}
		}

		const type = schema.type as FieldType;

		if ([FieldType_UInt8, FieldType_UInt16, FieldType_UInt32, FieldType_TwoUint8, FieldType_TwoUint16, FieldType_Uint32Array].includes(type)) {
			const numberValue = Array.isArray(value) ? value.map(Number) : Number(value);

			if (Array.isArray(numberValue)) {
				const min = schema.min;
				const max = schema.max;

				if (min != null && numberValue.some((val) => val < min)) {
					return {
						message: 'message.rule-value-numbers',
						argv: `${min}...${max}`,
					};
				}

				if (max != null && numberValue.some((val) => val > max)) {
					return {
						message: 'message.rule-value-numbers',
						argv: `${min}...${max}`,
					};
				}
			} else {
				const min = schema.min;
				const max = schema.max;

				if (min != null && numberValue < min) {
					return {
						message: 'message.rule-value-numbers',
						argv: `${min}...${max}`,
					};
				}

				if (max != null && numberValue > max) {
					return {
						message: 'message.rule-value-numbers',
						argv: `${min}...${max}`,
					};
				}
			}
		}

		if ([FieldType_Ascii, FieldType_Phones].includes(type) && schema.max_length && value.length > schema.max_length) {
			return {
				message: 'message.rule-value-max-length',
				argv: `${schema.max_length}`,
			};
		}

		return null;
	};

	const handleChange = useCallback(
		(newValue: any) => {
			if (!manager || !schema) return;
			setValue(newValue);

			const err = validate(newValue, schema);
			setError(err);

			const serialized = SerializeByType(schema.type as FieldType, newValue);
			manager.updateField(uid, serialized);

			if (timerRef.current) clearTimeout(timerRef.current);

			if (!err) {
				timerRef.current = setTimeout(async () => {
					await saveDraft(serialized);
				}, 300);
			}
		},
		[manager, schema, uid, index]
	);

	const saveDraft = useCallback(
		async (serializedValue: any) => {
			if (!manager || !schema || error) return;

			const changes = manager.getDraft();

			if (serializedValue !== undefined) {
				changes[uid] = serializedValue;
			}

			if (Object.keys(changes).length > 0) {
				const cfg_hash = manager.getBase()?.cfg_hash ?? 0;

				if (wrapperContext?.isTemplate) {
					saveTemplateDraft(cfg_hash, changes, wrapperContext?.templateId);
					return;
				}

				window.dispatchEvent(new CustomEvent(EVENT_CONFIGURATION_DRAFT, { detail: true }));
				await basicConfigurationInsertDraft(imei, section, cfg_hash, changes);
			}
		},
		[manager, schema, uid, imei, section, error, wrapperContext]
	);

	return { value, schema, isChanged, error, setError, handleChange, saveDraft };
};
