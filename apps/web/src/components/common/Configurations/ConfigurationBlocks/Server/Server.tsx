import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { arrayParseData } from '@/utils/ArrayUtils.ts';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useEnabled } from '@/hooks/Configuration/useEnabled';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField.tsx';
import { buildTrackingServers, TRACKING_SERVERS_PROTO, TrackingServer } from '@/constants/TrackingServers.constant';

interface ServerProps {
	index: number;
	hostValue: string;
	isLastActive: boolean;
}

const sanitizeServerHostPaste = (text: string) =>
	text
		.trim()
		.replace(/^https?:\/\//i, '')
		.replace(/\/+$/, '');

const detectPresetId = (servers: TrackingServer[], host: string, port: number): string => {
	const id = `${host}:${port}`;

	const found = servers.find((s) => s.id === id);
	if (found) return found.id;

	return 'custom';
};

export const Server: React.FC<ServerProps> = ({ index, hostValue, isLastActive }) => {
	const { t } = useTranslation();
	const { imei, model, section, manager, isTemplate } = useConfigurationWrapperContext();

	const hostField = useConfigurationField(imei, section, manager, UIDs.SERVER_HOST);
	const portField = useConfigurationField(imei, section, manager, UIDs.SERVER_PORT);

	const host = arrayParseData(hostField.value)[index] ?? '';
	const port = Number(arrayParseData(portField.value)[index]);

	const trackingServers = useMemo(() => buildTrackingServers(isTemplate ? undefined : model), [model, isTemplate]);
	const presetValue = useMemo(() => detectPresetId(trackingServers, host, port), [trackingServers, host, port]);

	const { enabled, toggle } = useEnabled(hostValue !== '');

	const handleToggle = () => {
		if (isLastActive && enabled) return;
		toggle();
	};

	return (
		<>
			<div className="w-full sm:min-w-[25rem] space-y-3">
				<div className="flex flex-row items-center justify-between">
					<label>{t('label.connection-type')}: </label>
					<GUISwitch className={`${isLastActive && enabled ? 'opacity-50 !cursor-not-allowed' : 'cursor-pointer'}`} checked={enabled} onChange={handleToggle} disabled={isLastActive && enabled} />
				</div>

				<LabeledField label="message.server-presets" disabled={!enabled}>
					<CFGSelect uid={UIDs.SERVER_HOST} index={index} array={trackingServers} presetType="server" _value={presetValue} writeFromExtra="address" disabled={!enabled} />
				</LabeledField>

				<LabeledField label="label.address" disabled={!enabled}>
					<CFGArrayItemInput uid={UIDs.SERVER_HOST} index={index} onFn={!enabled ? '' : undefined} disabled={!enabled} transformPaste={sanitizeServerHostPaste} />
				</LabeledField>

				<LabeledField label="label.port" disabled={!enabled}>
					<CFGArrayItemInput uid={UIDs.SERVER_PORT} index={index} onFn={!enabled ? '' : undefined} disabled={!enabled} />
				</LabeledField>

				<LabeledField label="label.protocol" disabled={!enabled}>
					<CFGSelect uid={UIDs.TRAFFIC_PROTOCOL} index={index} array={TRACKING_SERVERS_PROTO} disabled={!enabled} />
				</LabeledField>
			</div>
		</>
	);
};
export default Server;
