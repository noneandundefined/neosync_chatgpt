export const emptyToNull = (value?: string | null): string | null => {
	const trimmed = value?.trim() ?? '';

	if (trimmed === '' || trimmed === '+') {
		return null;
	}

	return trimmed;
};
