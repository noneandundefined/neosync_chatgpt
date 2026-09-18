export interface TemplateDraftData {
	cfg_hash: number;
	changes: Record<string, any>;
	timestamp: number;
}

export const EVENT_TEMPLATE_DRAFT = 'neosync:template:draft:ch';

const TEMPLATE_DRAFT_KEY = 'template-draft:current';

export const getTemplateDraftStorageKey = (templateId?: number | null) => (templateId ? `template-draft:${templateId}` : TEMPLATE_DRAFT_KEY);

export const getTemplateDraft = (templateId?: number | null): TemplateDraftData | null => {
	const raw = sessionStorage.getItem(getTemplateDraftStorageKey(templateId));
	if (!raw) return null;

	try {
		return JSON.parse(raw) as TemplateDraftData;
	} catch {
		return null;
	}
};

export const saveTemplateDraft = (cfgHash: number, changes: Record<string, any>, templateId?: number | null) => {
	const existing = getTemplateDraft(templateId);

	const merged: TemplateDraftData = {
		cfg_hash: cfgHash || existing?.cfg_hash || 0,
		changes: { ...(existing?.changes ?? {}), ...changes },
		timestamp: Math.floor(Date.now() / 1000),
	};

	sessionStorage.setItem(getTemplateDraftStorageKey(templateId), JSON.stringify(merged));
	window.dispatchEvent(new CustomEvent(EVENT_TEMPLATE_DRAFT, { detail: true }));
};

export const clearTemplateDraft = (templateId?: number | null) => {
	sessionStorage.removeItem(getTemplateDraftStorageKey(templateId));
	window.dispatchEvent(new CustomEvent(EVENT_TEMPLATE_DRAFT, { detail: false }));
};

export const hasTemplateDraft = (templateId?: number | null): boolean => {
	const draft = getTemplateDraft(templateId);
	return !!draft && Object.keys(draft.changes).length > 0;
};

export const getTemplateDraftChanges = (templateId?: number | null): Record<string, any> => {
	return getTemplateDraft(templateId)?.changes ?? {};
};

export const getTemplateDraftCfgHash = (templateId?: number | null): number => {
	return getTemplateDraft(templateId)?.cfg_hash ?? 0;
};
