import { config } from '@/.config/config.client';
import { ConfigurationSectionParsedResponse } from './configurationAPI';
import i18next from 'i18next';

export interface SseCallbacks {
	onOpen?: (message: string) => void;
	onProgress?: (message: string) => void;
	onConfiguration?: (data: ConfigurationSectionParsedResponse) => void;
	onError?: (message: string) => void;
	onDone?: (message: string) => void;
	onClose?: () => void;
}

export interface DraftSSECallbacks {
	onInitial?: (hasDraft: boolean) => void;
	onUpdate?: () => void;
	onError?: (err: any) => void;
	onClose?: () => void;
}

export function createSSE(url: string): EventSource {
	const lang = i18next.language || 'en';

	const base = config.type.release == 'dev' ? config.links.URL_BACKEND_DEV : config.type.release == 'prod' ? config.links.URL_BACKEND_PROD : config.links.URL_BACKEND_PREPROD;

	const fullUrl = `${base}${url}${url.includes('?') ? '&' : '?'}lang=${lang}`;

	const es = new EventSource(fullUrl, {
		withCredentials: true,
	} as any);

	return es;
}
