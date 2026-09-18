export const safeDecodeBase64 = (input: string | number) => {
	if (typeof input !== 'string') return input;

	const base64regex = /^[A-Za-z0-9+/]+={0,2}$/;

	if (input.length % 4 === 0 && base64regex.test(input)) {
		try {
			return Uint8Array.from(atob(input), (c) => c.charCodeAt(0));
		} catch {
			return input;
		}
	}

	return input;
};
