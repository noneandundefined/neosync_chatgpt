export const prettyFormatPhone = (phone: string): string => {
	const digits = phone.replace(/\D/g, '');

	if (digits.length <= 1) return `+${digits}`;
	if (digits.length <= 4) return `+${digits[0]} (${digits.slice(1)}`;
	if (digits.length <= 7) return `+${digits[0]} (${digits.slice(1, 4)}) ${digits.slice(4)}`;
	if (digits.length <= 9) return `+${digits[0]} (${digits.slice(1, 4)}) ${digits.slice(4, 7)}-${digits.slice(7)}`;
	if (digits.length <= 11) return `+${digits[0]} (${digits.slice(1, 4)}) ${digits.slice(4, 7)}-${digits.slice(7, 9)}-${digits.slice(9)}`;

	return `+${digits}`;
};

export const cleanPhone = (formatted: string): string => {
	return '+' + formatted.replace(/\D/g, '');
};
