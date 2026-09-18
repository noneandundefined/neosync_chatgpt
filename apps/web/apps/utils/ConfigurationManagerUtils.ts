import { FieldType } from '@/constants/FieldType.constant';
import { ConfigurationDraftResponse, ConfigurationSectionParsedResponse } from '@/rest/configurationAPI';
import { SerializeByType } from './SerializeByTypeUtils';

export interface FieldMerged {
	uid: string;
	name: string;
	type: FieldType;
	value: any;
	default?: any;
	min?: number;
	max?: number;
	min_length?: number;
	max_length?: number;
	max_items?: number;
	is_digits?: boolean;
	section: string;
}

export class ConfigurationManager {
	private listeners: Set<() => void> = new Set();

	private base: ConfigurationSectionParsedResponse;
	private draft?: ConfigurationDraftResponse;
	private mergedMap: Record<string, FieldMerged> = {};

	constructor(base: ConfigurationSectionParsedResponse, draft?: ConfigurationDraftResponse) {
		this.base = base;
		this.draft = draft;
		this.rebuild();
	}

	public update(base?: ConfigurationSectionParsedResponse, draft?: ConfigurationDraftResponse) {
		if (base) this.base = base;
		if (draft) this.draft = draft;
		this.rebuild();
	}

	private rebuild() {
		const result: Record<string, FieldMerged> = {};

		for (const field of this.base.config_parsed) {
			const schema = this.base.schema.find((s) => s.name === field.uid);
			const draftValue = this.draft?.changes?.[field.uid];

			result[field.uid] = {
				uid: field.uid,
				name: schema?.name ?? field.uid,
				type: schema?.type ?? 'Raw',
				value: draftValue !== undefined ? draftValue : (field.value ?? schema?.default ?? null),
				default: schema?.default ?? null,
				min: schema?.min,
				max: schema?.max,
				min_length: schema?.min_length,
				max_length: schema?.max_length,
				max_items: schema?.max_items,
				is_digits: schema?.is_digits,
				section: schema?.section ?? 'unknown',
			};
		}

		this.mergedMap = result;
	}

	public getAll(): Record<string, FieldMerged> {
		return this.mergedMap;
	}

	public getField(uid: string): FieldMerged | undefined {
		return this.mergedMap[uid];
	}

	public updateField(uid: string, value: any) {
		const field = this.mergedMap[uid];
		if (!field) return;

		if (!this.draft) {
			this.draft = {
				device_imei: this.base.device_imei,
				cfg_hash: this.base.cfg_hash,
				section: this.base.section,
				changes: {},
				schema: this.base.schema,
				timestamp: Date.now(),
			};
		}

		value = SerializeByType(field.type, value);

		this.draft.changes[uid] = value;
		this.mergedMap[uid].value = value;

		this.listeners.forEach((fn) => fn());
	}

	public getDraft(): Record<string, any> {
		return this.draft?.changes ?? {};
	}

	public isFieldChanged(uid: string, index?: number): boolean {
		const field = this.mergedMap[uid];
		if (!field) return false;

		const baseField = this.base.config_parsed.find((f) => f.uid === uid);
		const schema = this.base.schema.find((s) => s.name === uid);
		const original = baseField?.value ?? schema?.default ?? null;
		const current = field.value;

		if (index !== undefined) {
			const originalArr = Array.isArray(original) ? original : [];
			const currentArr = Array.isArray(current) ? current : [];

			return JSON.stringify(currentArr[index]) !== JSON.stringify(originalArr[index]);
		}

		return JSON.stringify(current) !== JSON.stringify(original);
	}

	public getBase(): ConfigurationSectionParsedResponse {
		return this.base;
	}

	public subscribe(fn: () => void) {
		this.listeners.add(fn);
		return () => this.listeners.delete(fn);
	}
}
