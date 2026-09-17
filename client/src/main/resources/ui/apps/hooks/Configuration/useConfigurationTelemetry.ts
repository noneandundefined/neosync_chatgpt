import { useCallback } from 'react';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicConfigurationTelemetry } from '@/rest/configurationAPI';

const TELEMETRY_REFETCH_MS = 5000;

export const useConfigurationTelemetry = (imei: string, enabled = true) => {
	const fetchConfigurationTelemetry = useCallback(() => basicConfigurationTelemetry(imei), [imei]);

	return useHandleServer(['respConfigurationTelemetry', imei], fetchConfigurationTelemetry, {
		enabled: enabled && !!imei,
		refetchInterval: TELEMETRY_REFETCH_MS,
	});
};
