export const isArray = (text: string): boolean => {
	try {
		const data: unknown = JSON.parse(text);
		return Array.isArray(data);
	} catch {
		return false;
	}
};

export const arrayParseData = (data: string | string[] | number | number[] | null | undefined): string[] => {
	if (data === null || data === undefined) return [];

	if (Array.isArray(data)) return data.map(String);

	if (typeof data === 'number') {
		return [String(data)];
	}

	if (typeof data === 'string') {
		if (data.trim().startsWith('[')) {
			try {
				const parsed = JSON.parse(data);
				if (Array.isArray(parsed)) return parsed.map(String);
			} catch {}
		}

		return [data];
	}

	return [];
};

export const normalizeArrayNumber = (data: any) => {
	if (Array.isArray(data)) {
		return data.map((val) => (isNaN(Number(val)) ? val : Number(val)));
	}

	return data;
};
