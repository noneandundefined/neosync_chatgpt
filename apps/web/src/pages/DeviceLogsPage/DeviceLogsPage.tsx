import { useCallback } from 'react';
import PageLayout from '../PageLayout';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { basicDeviceLogs } from '@/rest/deviceAPI';
import ChevronLeft from '@/components/@icons/chevron-left';
import { EVENT_VIEW_MAP } from '@/constants/Event.constant';
import { useDecryptedImei } from '@/hooks/useDecryptedImei';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import ErrorLoaderLine from '@/components/common/Loader/ErrorLoaderLine';

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

			return `${k}=${v}`;
		})
		.join(' ');
}

const DeviceLogsPage = () => {
	const { t } = useTranslation();
	const imei = useDecryptedImei();

	const navigate = useNavigate();

	if (!imei) {
		return <ErrorLoaderLine title={t('label.error')} description={t('message.error-get-imei')} />;
	}

	const fetchDeviceLogs = useCallback(() => basicDeviceLogs(imei), [imei]);
	const { data: respDeviceLogs, loading: respDeviceLogsLoading } = useHandleServer(['respDeviceLogs', imei], fetchDeviceLogs);

	return (
		<PageLayout>
			<div className="flex items-center gap-3 shrink-0">
				<Tooltip title={t('label.back')}>
					<div onClick={() => navigate(-1)} className="cursor-pointer p-1 rounded hover:bg-white">
						{/* <p className="text-sm hover:underline">{t('label.output')}</p> */}
						<ChevronLeft fill="#49525f" size={22} />
					</div>
				</Tooltip>

				<p className="text-[#555] font-medium text-[16px]">
					{t('label.terminal-events')} <span className="text-black ml-2">({imei})</span>
				</p>
			</div>

			<div className="relative flex-1 min-h-0 min-w-0 my-5 border-t border-b">
				<div className="absolute inset-0 overflow-auto py-3 text-sm">
					{respDeviceLogsLoading ? (
						<Loading />
					) : respDeviceLogs?.length === 0 ? (
						<p className="text-[#555]">{t('message.terminal-events-here')}</p>
					) : (
						respDeviceLogs?.map((log, index) => {
							const view = EVENT_VIEW_MAP[log.event];

							return (
								<div className="flex items-center gap-4 w-max min-w-full whitespace-nowrap" key={index}>
									<div className="text-[#555] shrink-0">{formatDate(log.ts)}</div>
									<div className="shrink-0">{t(view.text)}</div>
									<div className="text-[#555]">{renderPayload(log.payload)}</div>
								</div>
							);
						})
					)}
				</div>
			</div>
		</PageLayout>
	);
};

export default DeviceLogsPage;
