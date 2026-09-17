type QueryParams = Record<string, string | number>;

export const interpolateQuery = (template: string, params: QueryParams) => {
	return template.replace(/{{\s*(\w+)\s*}}/g, (_, key) => {
		const value = params[key];

		if (value === null || value === undefined || value === '') {
			return 'NULL';
		}

		if (typeof value === 'string') {
			return `'${value}'`;
		}

		return String(value);
	});
};
