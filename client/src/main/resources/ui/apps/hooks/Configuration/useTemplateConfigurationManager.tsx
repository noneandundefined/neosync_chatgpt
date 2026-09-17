import { useTranslation } from 'react-i18next';
import { useEffect, useRef, useState } from 'react';
import { ConfigurationManager } from '@/utils/ConfigurationManagerUtils';
import { getTemplateDraft } from '@/utils/TemplateDraftUtils';
import { basicConfigurationTemplateSectionParsedSSEGet, basicConfigurationTemplateSectionParsedSSEGetById, ConfigurationTemplateSectionParsedResponse } from '@/rest/configurationTemplateAPI';
import { ConfigurationDraftResponse, ConfigurationSectionParsedResponse } from '@/rest/configurationAPI';

const TEMPLATE_IMEI_PLACEHOLDER = 'template';

const toSectionResponse = (data: ConfigurationTemplateSectionParsedResponse): ConfigurationSectionParsedResponse => ({
	device_imei: TEMPLATE_IMEI_PLACEHOLDER,
	cfg_hash: data.cfg_hash,
	section: data.section,
	config_parsed: data.config_parsed,
	schema: data.schema,
	timestamp: data.timestamp,
});

const toDraftResponse = (section: string, draft: ReturnType<typeof getTemplateDraft>): ConfigurationDraftResponse | undefined => {
	if (!draft || Object.keys(draft.changes).length === 0) return undefined;

	return {
		device_imei: TEMPLATE_IMEI_PLACEHOLDER,
		cfg_hash: draft.cfg_hash,
		section,
		changes: draft.changes,
		schema: [],
		timestamp: draft.timestamp,
	};
};

const useTemplateConfigurationManager = (section: string, enabled = true, templateId?: number) => {
	const { t } = useTranslation();

	const [manager, setManager] = useState<ConfigurationManager | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const [sseMessage, setSseMessage] = useState<string | null>(null);
	const [sseConnectionError, setSseConnectionError] = useState<string | null>(null);

	const eventSourceRef = useRef<EventSource | null>(null);
	const hasFinishedRef = useRef(false);

	useEffect(() => {
		if (!enabled) {
			setLoading(false);
			return;
		}

		if (eventSourceRef.current) {
			eventSourceRef.current.close();
		}

		let active = true;
		hasFinishedRef.current = false;

		setLoading(true);
		setError(null);
		setSseMessage(null);
		setSseConnectionError(null);
		setManager(null);

		const callbacks = {
			onOpen: (message?: string) => setSseMessage(message ?? null),
			onProgress: (message?: string) => setSseMessage(message ?? null),
			onConfiguration: (baseConfig: ConfigurationTemplateSectionParsedResponse) => {
				if (!active) return;

				const draft = getTemplateDraft(templateId);
				const sectionDraft = toDraftResponse(section, draft);
				setManager(new ConfigurationManager(toSectionResponse(baseConfig), sectionDraft));
				hasFinishedRef.current = true;
				setLoading(false);
			},
			onError: (message?: string) => {
				setSseConnectionError(message ?? t('message.error-get-configuration'));
				setError(message ?? t('message.error-get-configuration'));
				setLoading(false);
				hasFinishedRef.current = true;
			},
			onDone: (message?: string) => {
				setSseMessage(message ?? null);
				setLoading(false);
				hasFinishedRef.current = true;
			},
			onClose: () => {
				if (active && loading && !hasFinishedRef.current) {
					setSseConnectionError(t('message.sse-connection-unexpected-close'));
					setError(t('message.sse-connection-unexpected-close'));
				}
			},
		};

		const eventSource = templateId ? basicConfigurationTemplateSectionParsedSSEGetById(templateId, section, callbacks) : basicConfigurationTemplateSectionParsedSSEGet(section, callbacks);

		eventSourceRef.current = eventSource;

		return () => {
			active = false;
			if (eventSourceRef.current) {
				eventSourceRef.current.close();
				eventSourceRef.current = null;
			}
		};
	}, [section, enabled, templateId, t]);

	return { manager, loading, error, sseMessage, sseConnectionError };
};

export default useTemplateConfigurationManager;
