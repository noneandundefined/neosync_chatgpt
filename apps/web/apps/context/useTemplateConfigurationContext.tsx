import { createContext, useContext } from 'react';
import { ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';

interface TemplateConfigurationContextValue {
	isTemplate: true;
	templateId?: number;
	template?: ConfigurationTemplateResponse | null;
}

const TemplateConfigurationContext = createContext<TemplateConfigurationContextValue | null>(null);

export const useTemplateConfigurationContext = () => useContext(TemplateConfigurationContext);

export const TemplateConfigurationProvider = ({ children, templateId, template }: { children: React.ReactNode; templateId?: number; template?: ConfigurationTemplateResponse | null }) => {
	return <TemplateConfigurationContext.Provider value={{ isTemplate: true, templateId, template }}>{children}</TemplateConfigurationContext.Provider>;
};
