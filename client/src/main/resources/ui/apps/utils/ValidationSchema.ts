import { TFunction } from 'i18next';
import { RegisterOptions, FieldValues, Path } from 'react-hook-form';

const EMAIL_REGEX = /^\S+@\S+$/;

export const ValidationImeiSchema = <T extends FieldValues, TFieldName extends Path<T>>(t: TFunction): RegisterOptions<T, TFieldName> => ({
	required: t('message.validation-required-field'),
	minLength: { value: 15, message: t('message.validation-imei-length') },
	maxLength: { value: 15, message: t('message.validation-imei-length') },
});

export const ValidationPhoneSchema = <T extends FieldValues, TFieldName extends Path<T>>(t: TFunction): RegisterOptions<T, TFieldName> => ({
	setValueAs: (v?: string) => {
		if (!v) return '';
		return v.startsWith('+') ? v : `+${v}`;
	},
	validate: (v?: string) => !v || v.length >= 8 || t('message.validation-phone-length'),
});

export const ValidationEmailSchema = <T extends FieldValues, TFieldName extends Path<T>>(t: TFunction): RegisterOptions<T, TFieldName> => ({
	required: t('message.validation-required-field'),
	pattern: {
		value: EMAIL_REGEX,
		message: t('message.validation-email-invalid'),
	},
});

export const ValidationPasswordSchema = <T extends FieldValues, TFieldName extends Path<T>>(t: TFunction): RegisterOptions<T, TFieldName> => ({
	minLength: {
		value: 6,
		message: t('message.validation-password-min-length'),
	},
	maxLength: {
		value: 16,
		message: t('message.validation-password-max-length'),
	},
});

export const ValidationPasswordRequiredSchema = <T extends FieldValues, TFieldName extends Path<T>>(t: TFunction): RegisterOptions<T, TFieldName> => ({
	required: t('message.validation-required-field'),
	minLength: {
		value: 6,
		message: t('message.validation-password-min-length'),
	},
	maxLength: {
		value: 16,
		message: t('message.validation-password-max-length'),
	},
});
