import type { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';
import { BULK_CONFIGURATION_MAX_DEVICES } from '@/constants/BulkConfiguration.constant';

// Prefix spreadsheet formulas (and literal apostrophes) reversibly.
export const encodeBulkCsv = (rows: Record<string, unknown>[]): string => {
	const cell = (value: unknown) => {
		let text = String(value ?? '');
		if (/^[=+\-@\t\r\n']/.test(text)) text = `'${text}`;
		return `"${text.replace(/"/g, '""')}"`;
	};
	const headers = Object.keys(rows[0]);
	return '\uFEFF' + [headers, ...rows.map((row) => headers.map((key) => row[key]))].map((row) => row.map(cell).join(';')).join('\r\n');
};

export const decodeBulkCsv = (input: string): Record<string, string>[] => {
	const rows: string[][] = [];
	let row: string[] = [],
		value = '',
		quoted = false;
	const text = input.replace(/^\uFEFF/, '');
	for (let i = 0; i < text.length; i++) {
		const ch = text[i];
		if (ch === '"') {
			if (quoted && text[i + 1] === '"') {
				value += '"';
				i++;
			} else quoted = !quoted;
		} else if (!quoted && (ch === ';' || ch === '\n' || ch === '\r')) {
			row.push(value);
			value = '';
			if (ch !== ';') {
				rows.push(row);
				row = [];
				if (ch === '\r' && text[i + 1] === '\n') i++;
			}
		} else value += ch;
	}
	if (quoted) throw new Error('invalid-csv');
	if (value || row.length) {
		row.push(value);
		rows.push(row);
	}
	const headers = rows.shift();
	if (!headers?.length || new Set(headers).size !== headers.length) throw new Error('invalid-csv');
	return rows
		.filter((r) => r.some(Boolean))
		.map((r) => {
			if (r.length !== headers.length) throw new Error('invalid-csv');
			return Object.fromEntries(headers.map((key, index) => [key, /^'[=+\-@\t\r\n']/.test(r[index]) ? r[index].slice(1) : r[index]]));
		});
};

export const importBulkConfigurationCsv = (text: string): { selectedIds: string[]; importedConfiguration: CompanyCreateRequest } => {
	const rows = decodeBulkCsv(text);
	const first = rows[0];
	if (!first || first.csv_version !== '1') throw new Error('invalid-csv');
	const sharedKeys = [
		'company_id',
		'name',
		'priority_model',
		'ttl',
		'existing_task_action',
		'launch_mode',
		'incompatible_action',
		'configuration_cource',
		'configuration_template_id',
		'configuration_device_cource',
		'configuration_file_name',
		'configuration_file_base64',
		'selected_imeis_json',
	];
	if (rows.some((row) => sharedKeys.some((key) => row[key] !== first[key]))) throw new Error('invalid-csv');
	const imeis: unknown = JSON.parse(first.selected_imeis_json);
	if (!Array.isArray(imeis) || imeis.some((imei) => typeof imei !== 'string' || !/^\d{15}$/.test(imei))) throw new Error('invalid-csv');
	const selectedIds = [...new Set<string>(imeis)];
	if (!selectedIds.length || selectedIds.length > BULK_CONFIGURATION_MAX_DEVICES) throw new Error('invalid-csv');
	if (
		!first.name?.trim() ||
		first.name.length > 255 ||
		!['1', '7', '30'].includes(first.ttl) ||
		!['skip', 'replace', 'overwrite'].includes(first.existing_task_action) ||
		!['all', 'selected'].includes(first.launch_mode) ||
		first.incompatible_action !== 'exclude'
	)
		throw new Error('invalid-csv');
	const data: CompanyCreateRequest = {
		name: first.name,
		ttl: first.ttl,
		priority_model: first.priority_model || undefined,
		existing_task_action: first.existing_task_action === 'overwrite' ? 'replace' : first.existing_task_action,
		launch_mode: first.launch_mode,
		incompatible_action: 'exclude',
		configuration_cource: first.configuration_cource as CompanyCreateRequest['configuration_cource'],
	};
	switch (data.configuration_cource) {
		case 'device_sources':
			if (!selectedIds.includes(first.configuration_device_cource)) throw new Error('missing-source');
			data.configuration_device_cource = first.configuration_device_cource;
			break;
		case 'template_sources':
			data.configuration_template_id = Number(first.configuration_template_id);
			if (!Number.isSafeInteger(data.configuration_template_id) || data.configuration_template_id <= 0) throw new Error('missing-source');
			break;
		case 'file_sources': {
			const base64 = first.configuration_file_base64;
			if (!base64 || base64.length > 6830 || !/^[A-Za-z0-9+/]+={0,2}$/.test(base64)) throw new Error('missing-source');
			const binary = atob(base64);
			if (!binary.length || binary.length > 5120 || btoa(binary) !== base64 || first.configuration_file_name.length > 255) throw new Error('invalid-csv');
			data.configuration_file_name = first.configuration_file_name || `company_${first.company_id}.neosync`;
			data.configuration_file_base64 = base64;
			break;
		}
		default:
			throw new Error('invalid-csv');
	}
	return { selectedIds, importedConfiguration: data };
};
