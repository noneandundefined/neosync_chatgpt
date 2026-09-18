const TEMPLATE_DRAFT_PREFIX = 'neosync:template-draft:';

export interface TemplateDraftData {
	cfg_hash: number;
	changes: Record<string, any>;
}

export const getTemplateDraftKey = (model: string) => `${TEMPLATE_DRAFT_PREFIX}${model}`;

export const loadTemplateDraft = (model: string): TemplateDraftData | null => {
	const raw = sessionStorage.getItem(getTemplateDraftKey(model));
	if (!raw) return null;

	try {
		return JSON.parse(raw) as TemplateDraftData;
	} catch {
		return null;
	}
};

export const saveTemplateDraft = (model: string, draft: TemplateDraftData) => {
	sessionStorage.setItem(getTemplateDraftKey(model), JSON.stringify(draft));
};

export const clearTemplateDraft = (model: string) => {
	sessionStorage.removeItem(getTemplateDraftKey(model));
};

export const hasTemplateDraftChanges = (model: string): boolean => {
	const draft = loadTemplateDraft(model);
	return !!draft && Object.keys(draft.changes).length > 0;
};

export const mergeTemplateDraftChanges = (model: string, sectionChanges: Record<string, any>, cfgHash: number) => {
	const existing = loadTemplateDraft(model) ?? { cfg_hash: cfgHash, changes: {} };

	existing.cfg_hash = cfgHash;

	for (const [key, value] of Object.entries(sectionChanges)) {
		existing.changes[key] = value;
	}

	saveTemplateDraft(model, existing);
	return existing;
};
