import { useCallback, useState } from 'react';
import PageLayout from '../PageLayout';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import { useNavigate } from 'react-router-dom';
import { basicDeviceLogs, basicDeviceLogsExport } from '@/rest/deviceAPI';
import ChevronLeft from '@/components/@icons/chevron-left';
import { EVENT_VIEW_MAP } from '@/constants/Event.constant';
import { useDecryptedImei } from '@/hooks/useDecryptedImei';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';
import { basicConfigurationGetHistories, basicConfigurationHistoryExport, ConfigurationHistory } from '@/rest/configurationAPI';

function formatDate(ts: number): string {
	return new Date(ts * 1000).toLocaleString('ru-RU', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
	});
}

function formatHistoryDate(value: string): string {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return value;

	return date.toLocaleString('ru-RU', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
	});
}

function renderPayload(payload: any): string {
	if (!payload) return '';

	if (typeof payload !== 'object') {
		return String(payload);
	}

	return Object.entries(payload)
		.filter(([_, v]) => v !== null && v !== '' && v !== '0001-01-01T00:00:00Z')
		.map(([k, v]) => {
			if (k === 'last_mod_time') {
				const ts = Number(v);
				return `${k}=${Number.isFinite(ts) && ts > 0 ? formatDate(ts) : v}`;
			}

			return `${k}=${typeof v === 'object' ? JSON.stringify(v) : v}`;
		})
		.join(' ');
}

function renderConfigValue(value: any): string {
	if (value === undefined) return '—';
	if (value === null) return 'null';
	if (typeof value === 'string') return value;
	if (typeof value === 'object') return JSON.stringify(value);
	return String(value);
}

const originLabel = (origin: string, t: TFunction) => {
	switch (origin) {
		case 'neosync':
			return t('label.config-origin-neosync');
		case 'tracker':
			return t('label.config-origin-tracker');
		case 'baseline':
			return t('label.config-origin-baseline');
		default:
			return t('label.config-origin-unknown');
	}
};

const DeviceLogsPage = () => {
	const { t } = useTranslation();
	const imei = useDecryptedImei();
	const safeImei = imei ?? '';
	const navigate = useNavigate();

	const [tab, setTab] = useState<'events' | 'history'>('events');
	const [expandedHistoryId, setExpandedHistoryId] = useState<number | null>(null);

	const fetchDeviceLogs = useCallback((signal?: AbortSignal) => basicDeviceLogs(safeImei, signal), [safeImei]);
	const { data: respDeviceLogs, loading: respDeviceLogsLoading } = useHandleServer(['respDeviceLogs', safeImei], fetchDeviceLogs, {
		enabled: !!imei,
	});

	const fetchConfigurationHistory = useCallback((signal?: AbortSignal) => basicConfigurationGetHistories(safeImei, signal), [safeImei]);
	const { data: configurationHistory, loading: configurationHistoryLoading } = useHandleServer(['respConfigurationHistory', safeImei], fetchConfigurationHistory, {
		enabled: !!imei && tab === 'history',
		staleTime: 15_000,
	});

	if (!imei) {
		return <ErrorLoaderLine title={t('label.error')} description={t('message.error-get-imei')} />;
	}

	const renderHistoryItem = (item: ConfigurationHistory) => {
		const expanded = expandedHistoryId === item.id;

		return (
			<div key={item.id} className="border-b border-[#e9e9e9] last:border-b-0">
				<div className="flex flex-wrap items-center gap-x-4 gap-y-2 px-3 py-3 hover:bg-[#f8f8f8]">
					<button className="flex-1 min-w-[280px] text-left" onClick={() => setExpandedHistoryId(expanded ? null : item.id)}>
						<div className="flex flex-wrap items-center gap-2">
							<span className="font-medium text-[#49525f]">{formatHistoryDate(item.apply_at)}</span>
							<span className="text-xs text-[#656d76]">hash: {item.cfg_hash}</span>
							<span className="text-xs px-2 py-0.5 rounded bg-[#f1f1f1]">{originLabel(item.origin, t)}</span>
							{item.cfg_sync_status && <span className="text-xs px-2 py-0.5 rounded bg-[#f1f1f1]">{item.cfg_sync_status}</span>}
						</div>

						<div className="mt-1 text-xs text-[#656d76]">
							{t('label.changed-fields')}: {item.change_count}
							{item.cfg_sync_error ? <span className="ml-3 text-red-600">{item.cfg_sync_error}</span> : null}
						</div>
					</button>

					<button
						className="px-3 py-1.5 text-xs border rounded hover:bg-white"
						onClick={(e) => {
							e.stopPropagation();
							void basicConfigurationHistoryExport(imei, item.id);
						}}
					>
						{t('label.download-version')}
					</button>
				</div>

				{expanded && (
					<div className="px-3 pb-3">
						{item.parse_error ? (
							<div className="text-sm text-red-600">{item.parse_error}</div>
						) : item.changes.length === 0 ? (
							<div className="text-sm text-[#656d76]">{t('message.configuration-history-no-diff')}</div>
						) : (
							<div className="border rounded overflow-hidden">
								{item.changes.map((change) => (
									<div key={change.uid} className="grid grid-cols-[minmax(170px,0.7fr)_1fr_32px_1fr] gap-2 px-3 py-2 text-xs border-b last:border-b-0 items-start">
										<div className="min-w-0">
											<div className="font-medium text-[#49525f] break-all">{change.name}</div>
											<div className="text-[#8b949e]">UID {change.uid}{change.section ? ` · ${change.section}` : ''}</div>
										</div>
										<div className="break-all text-[#656d76]">{renderConfigValue(change.old_value)}</div>
										<div className="text-center text-[#8b949e]">→</div>
										<div className="break-all text-[#24292f]">{renderConfigValue(change.new_value)}</div>
									</div>
								))}
							</div>
						)}
					</div>
				)}
			</div>
		);
	};

	return (
		<PageLayout>
			<div className="flex items-center gap-3 shrink-0">
				<Tooltip title={t('label.back')}>
					<div onClick={() => navigate(-1)} className="cursor-pointer p-1 rounded hover:bg-white">
						<ChevronLeft fill="#49525f" size={22} />
					</div>
				</Tooltip>

				<p className="text-[#555] font-medium text-[16px]">
					{t('label.terminal-events')} <span className="text-black ml-2">({imei})</span>
				</p>
			</div>

			<div className="flex items-center justify-between gap-3 mt-4 border-b">
				<div className="flex">
					<button
						className={`px-3 py-2 text-sm border-b-2 ${tab === 'events' ? 'border-[#49525f] text-[#24292f]' : 'border-transparent text-[#656d76]'}`}
						onClick={() => setTab('events')}
					>
						{t('label.terminal-events')}
					</button>
					<button
						className={`px-3 py-2 text-sm border-b-2 ${tab === 'history' ? 'border-[#49525f] text-[#24292f]' : 'border-transparent text-[#656d76]'}`}
						onClick={() => setTab('history')}
					>
						{t('label.configuration-history')}
					</button>
				</div>

				{tab === 'events' && (
					<button className="px-3 py-1.5 text-xs border rounded hover:bg-white" onClick={() => void basicDeviceLogsExport(imei)}>
						{t('label.download-events')}
					</button>
				)}
			</div>

			<div className="relative flex-1 min-h-0 min-w-0 my-3 border-b">
				<div className="absolute inset-0 overflow-auto py-2 text-sm">
					{tab === 'events' ? (
						respDeviceLogsLoading ? (
							<Loading />
						) : respDeviceLogs?.length === 0 ? (
							<p className="text-[#555]">{t('message.terminal-events-here')}</p>
						) : (
							respDeviceLogs?.map((log) => {
								const view = EVENT_VIEW_MAP[log.event];

								return (
									<div className="flex items-center gap-4 w-max min-w-full whitespace-nowrap" key={log.id}>
										<div className="text-[#555] shrink-0">{formatDate(log.ts)}</div>
										<div className="shrink-0">{view ? t(view.text) : log.event}</div>
										<div className="text-[#555]">{renderPayload(log.payload)}</div>
									</div>
								);
							})
						)
					) : configurationHistoryLoading ? (
						<Loading />
					) : configurationHistory?.length === 0 ? (
						<p className="text-[#555]">{t('message.configuration-history-empty')}</p>
					) : (
						<div className="border rounded bg-white">{configurationHistory?.map(renderHistoryItem)}</div>
					)}
				</div>
			</div>
		</PageLayout>
	);
};

export default DeviceLogsPage;
