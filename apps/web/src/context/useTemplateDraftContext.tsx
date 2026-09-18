import { createContext, useContext } from 'react';

interface TemplateDraftContextValue {
	model: string;
	hasUnsavedChanges: boolean;
	setHasUnsavedChanges: (value: boolean) => void;
}

const TemplateDraftContext = createContext<TemplateDraftContextValue | null>(null);

export const useTemplateDraftContext = () => useContext(TemplateDraftContext);

export const TemplateDraftProvider = TemplateDraftContext.Provider;
