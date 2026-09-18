import axios from 'axios';
import { useTranslation } from 'react-i18next';
import { useEffect, useRef, useState } from 'react';
import { ConfigurationManager } from '@/utils/ConfigurationManagerUtils';
import { basicConfigurationGetDraft, basicConfigurationSectionParsedSSEGet } from '@/rest/configurationAPI';

const useConfigurationManager = (imei: string, section: string) => {
	const { t } = useTranslation();

	const [manager, setManager] = useState<ConfigurationManager | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const [sseMessage, setSseMessage] = useState<string | null>(null);
	const [sseConnectionError, setSseConnectionError] = useState<string | null>(null);

	const abortRef = useRef<AbortController | null>(null);
	const eventSourceRef = useRef<EventSource | null>(null);
	const hasFinishedRef = useRef(false);

	useEffect(() => {
		if (!imei) {
			setLoading(false);
			return;
		}

		if (abortRef.current) {
			abortRef.current.abort();
		}
		if (eventSourceRef.current) {
			eventSourceRef.current.close();
		}

		const controller = new AbortController();
		abortRef.current = controller;
		let active = true;
		hasFinishedRef.current = false;

		setLoading(true);
		setError(null);
		setSseMessage(null);
		setSseConnectionError(null);
		setManager(null);

		const eventSource = basicConfigurationSectionParsedSSEGet(imei, section, {
			onOpen: (message) => {
				if (active) setSseMessage(message);
			},
			onProgress: (message) => {
				if (active) setSseMessage(message);
			},
			onConfiguration: async (baseConfig) => {
				try {
					const draft = await basicConfigurationGetDraft(imei, section, controller.signal);
					if (active) {
						setManager(new ConfigurationManager(baseConfig, draft));
						hasFinishedRef.current = true;
					}
				} catch (draftError: any) {
					if (axios.isCancel(draftError) || !active) return;
					setError(draftError?.response?.data?.message || t('message.configuration-error'));
					hasFinishedRef.current = true;
				} finally {
					if (active) setLoading(false);
				}
			},
			onError: (message) => {
				if (!active) return;
				setSseConnectionError(message);
				setError(message);
				setLoading(false);
				hasFinishedRef.current = true;
			},
			onDone: (message) => {
				if (!active) return;
				setSseMessage(message);
				setLoading(false);
				hasFinishedRef.current = true;
			},
			onClose: () => {
				if (active && !hasFinishedRef.current) {
					setSseConnectionError(t('message.sse-connection-unexpected-close'));
					setError(t('message.sse-connection-unexpected-close'));
				}
			},
		});

		eventSourceRef.current = eventSource;

		return () => {
			active = false;
			controller.abort();
			if (eventSourceRef.current) {
				eventSourceRef.current.close();
				eventSourceRef.current = null;
			}
		};
	}, [imei, section, t]);

	return { manager, loading, error, sseMessage, sseConnectionError };
};

export default useConfigurationManager;
