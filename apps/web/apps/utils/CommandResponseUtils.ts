import type { TFunction } from 'i18next';
import i18n from '@/utils/i18n';

export function formatCommandResponse(response: string | null | undefined, t: TFunction): string {
	const value = response ?? '';
	if (!value) {
		return '';
	}

	if (value.startsWith('message.') && i18n.exists(value)) {
		return t(value);
	}

	return value;
}
