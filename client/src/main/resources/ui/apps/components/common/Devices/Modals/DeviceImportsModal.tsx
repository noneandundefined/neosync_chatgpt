import { toast } from 'react-toastify';
import useIsMobile from '@/hooks/useIsMobile';
import Tooltip from '@/components/ui/Tooltip';
import Check from '@/components/@icons/check';
import Close from '@/components/@icons/close';
import { useTranslation } from 'react-i18next';
import Delete from '@/components/@icons/delete';
import { basicUsersToOwnerGet } from '@/rest/userAPI';
import { useCallback, useMemo, useState } from 'react';
import { basicGroupFieldsName } from '@/rest/groupAPI';
import { basicDevicesImportRaw } from '@/rest/deviceAPI';
import GUIButton from '@/components/ui/Button/GUIButton';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useDeviceFileUpload } from '@/hooks/useDeviceFileUpload';

interface DeviceImportsModalProps {
	onSuccess: () => void;
}

const IMPORT_BATCH_SIZE = 25;
const INITIAL_VISIBLE_ROWS = 25;
const VISIBLE_ROWS_STEP = 25;

type ImportUiDraft = {
	selectedImeis: string[];
	selectedGroup: string;
	selectedOwner: string;
	visibleRows: number;
};

const importUiDraft: ImportUiDraft = {
	selectedImeis: [],
	selectedGroup: '',
	selectedOwner: '',
	visibleRows: INITIAL_VISIBLE_ROWS,
};

const DeviceImportsModal: React.FC<DeviceImportsModalProps> = ({ onSuccess }) => {
	const { t } = useTranslation();
	const isMobile = useIsMobile();

	const { data: respUsersToOwnerGet } = useHandleServer(['respUsersToOwnerGet'], () => basicUsersToOwnerGet({ all: true }));
	const { data: respGroupFieldsName } = useHandleServer(['respGroupFieldsName'], basicGroupFieldsName);

	const { devices, setDevices, handleFileUpload } = useDeviceFileUpload();

	const [selectedImeis, setSelectedImeisState] = useState<string[]>(importUiDraft.selectedImeis);
	const [selectedGroup, setSelectedGroupState] = useState(importUiDraft.selectedGroup);
	const [selectedOwner, setSelectedOwnerState] = useState(importUiDraft.selectedOwner);
	const [visibleRows, setVisibleRowsState] = useState(importUiDraft.visibleRows);

	const setSelectedImeis = useCallback((action: React.SetStateAction<string[]>) => {
		setSelectedImeisState((prev) => {
			const next = typeof action === 'function' ? action(prev) : action;
			importUiDraft.selectedImeis = next;
			return next;
		});
	}, []);

	const setSelectedGroup = useCallback((value: string) => {
		importUiDraft.selectedGroup = value;
		setSelectedGroupState(value);
	}, []);

	const setSelectedOwner = useCallback((value: string) => {
		importUiDraft.selectedOwner = value;
		setSelectedOwnerState(value);
	}, []);

	const setVisibleRows = useCallback((action: React.SetStateAction<number>) => {
		setVisibleRowsState((prev) => {
			const next = typeof action === 'function' ? action(prev) : action;
			importUiDraft.visibleRows = next;
			return next;
		});
	}, []);
	const [importProgress, setImportProgress] = useState<{ done: number; total: number } | null>(null);

	const [fileBusy, setFileBusy] = useState(false);
	const [isImporting, setIsImporting] = useState(false);

	const busy = fileBusy || isImporting;

	const selectedDevices = useMemo(() => devices.filter((d) => selectedImeis.includes(d.imei)), [devices, selectedImeis]);
	const visibleDevices = useMemo(() => devices.slice(0, visibleRows), [devices, visibleRows]);
	const hasMoreRows = devices.length > visibleRows;

	const selectedHaveErrors = useMemo(() => selectedDevices.some((d) => Boolean(d.error)), [selectedDevices]);

	const handleDownload = (path: string, pathname: string) => {
		const link = document.createElement('a');
		link.href = path;
		link.download = pathname;
		link.click();
	};

	const handleImportDevices = async () => {
		if (!devices || devices.length === 0) return;

		const selectedForImport = devices.filter((d) => selectedImeis.includes(d.imei));

		const importDevices = selectedForImport.map((device) => ({
			imei: device.imei,
			model: device.model ?? null,
			owner_email: selectedOwner || null,
			group_name: selectedGroup || null,
		}));

		if (importDevices.length === 0) return;

		if (selectedForImport.some((d) => d.error)) {
			toast.warning(t('message.import-error'));
			return;
		}

		setIsImporting(true);
		setImportProgress({ done: 0, total: importDevices.length });

		try {
			const errorsMap = new Map<string, string>();
			let importedCount = 0;

			for (let offset = 0; offset < importDevices.length; offset += IMPORT_BATCH_SIZE) {
				const batch = importDevices.slice(offset, offset + IMPORT_BATCH_SIZE);
				const resp = await basicDevicesImportRaw(batch);

				if (resp.errors?.length) {
					resp.errors.forEach((e) => errorsMap.set(e.imei, e.error));
				}

				importedCount += batch.length;
				setImportProgress({ done: importedCount, total: importDevices.length });
			}

			if (errorsMap.size > 0) {
				setDevices((prev) =>
					prev.map((d) => {
						const foundError = errorsMap.get(d.imei);
						return foundError ? { ...d, error: foundError } : d;
					})
				);
				toast.warning(t('message.import-error'));
				return;
			}

			toast.success(t('message.successful-device-import'));

			setDevices([]);
			setSelectedImeis([]);
			setSelectedGroup('');
			setSelectedOwner('');
			setVisibleRows(INITIAL_VISIBLE_ROWS);
			setImportProgress(null);

			onSuccess();
		} finally {
			setIsImporting(false);
			setImportProgress(null);
		}
	};

	const allSelected = useMemo(() => devices.length > 0 && devices.every((d) => selectedImeis.includes(d.imei)), [devices, selectedImeis]);

	const hasSelectedForImport = selectedImeis.length > 0;

	const toggleSelectAll = () => {
		if (allSelected) {
			setSelectedImeis([]);
			return;
		}

		setSelectedImeis(devices.map((d) => d.imei));
	};

	const toggleRow = (imei: string) => {
		setSelectedImeis((prev) => (prev.includes(imei) ? prev.filter((item) => item !== imei) : [...prev, imei]));
	};

	const handleListScroll = useCallback(
		(event: React.UIEvent<HTMLDivElement>) => {
			if (!hasMoreRows) return;

			const el = event.currentTarget;
			if (el.scrollHeight - el.scrollTop - el.clientHeight > 48) return;

			setVisibleRows((prev) => Math.min(prev + VISIBLE_ROWS_STEP, devices.length));
		},
		[hasMoreRows, devices.length]
	);

	return (
		<div className="space-y-4 text-[#49525f]">
			{(!isMobile || devices.length === 0) && (
				<div className="flex flex-col space-y-1 text-[14px]">
					<div className="flex flex-row gap-2">
						<p>{t('message.file-templates')}:</p>
						<div className="flex flex-row gap-2">
							<div className="cursor-pointer text-[#296acd] hover:underline" onClick={() => handleDownload('/examples/adm/adm-devices.csv', 'adm-devices.csv')}>
								.csv
							</div>
							<div className="cursor-pointer text-[#296acd] hover:underline" onClick={() => handleDownload('/examples/adm/adm-devices.txt', 'adm-devices.txt')}>
								.txt
							</div>
						</div>
					</div>

					<p>{t('message.max-file-size-1mb')}</p>
					<p className="font-semibold text-[#9a0000]">{t('message.import-device-utf-8-desc')}</p>
				</div>
			)}

			<div className="flex items-center">
				<label id="button" className={`!pt-[0.7rem] !w-auto text-center uppercase ${busy ? 'cursor-not-allowed opacity-50 pointer-events-none' : 'cursor-pointer'}`}>
					{t('label.select-file')}
					<input
						type="file"
						className="hidden"
						accept=".csv,.xlsx,.txt"
						disabled={busy}
						onChange={async (e) => {
							const file = e.target.files?.[0];
							e.currentTarget.value = '';
							if (!file) return;

							setFileBusy(true);
							try {
								const uploaded = await handleFileUpload(file);
								if (uploaded) {
									setSelectedImeis(uploaded.filter((d) => !d.error).map((d) => d.imei));
									setVisibleRows(INITIAL_VISIBLE_ROWS);
								}
							} finally {
								setFileBusy(false);
							}
						}}
					/>
				</label>
			</div>

			{devices.length > 0 && (
				<div className="space-y-4">
					<div className="flex items-center justify-between gap-3">
						<div className={`flex flex-1 flex-wrap items-center gap-5 ${busy ? 'pointer-events-none opacity-50' : ''}`}>
							<div className="w-full sm:w-[220px]">
								<GUISelect className="relative !min-w-0 !w-full" value={selectedGroup} onChange={(e) => setSelectedGroup(e.target.value)}>
									<option value="">{t('message.assign-group')}</option>
									{respGroupFieldsName?.map((group, index) => (
										<option value={group.name} key={index}>
											{group.name}
										</option>
									))}
								</GUISelect>
							</div>

							<div className="w-full sm:w-[220px]">
								<GUISelect className="relative !min-w-0 !w-full" value={selectedOwner} onChange={(e) => setSelectedOwner(e.target.value)}>
									<option value="">{t('message.assign-owner')}</option>
									{respUsersToOwnerGet?.items?.map((user, index) => (
										<option value={user.email} key={index}>
											{user.email}
										</option>
									))}
								</GUISelect>
							</div>
						</div>

						<Tooltip title={t('message.clear-import-list')} position="bottom">
							<button
								type="button"
								onClick={() => {
									setDevices([]);
									setSelectedImeis([]);
									setSelectedGroup('');
									setSelectedOwner('');
									setVisibleRows(INITIAL_VISIBLE_ROWS);
								}}
								disabled={busy}
								className="h-[37px] w-[37px] hidden sm:flex items-center justify-center rounded border border-[#d1242f] hover:bg-[#ffecef] disabled:pointer-events-none"
							>
								<Delete fill="#d1242f" size={17} />
							</button>
						</Tooltip>
					</div>

					<div className="rounded border border-[#e6e8ec]">
						<div className="grid grid-cols-[32px_125px_1fr_1fr] items-center border-b border-[#eceff3] bg-[#f8f9fb] px-3 py-2 text-[12px] font-semibold">
							<GUICheckbox checked={allSelected} onChange={toggleSelectAll} />
							<span>{t('label.imei')}</span>
							<span>{t('label.model')}</span>
							<span>{t('label.status')}</span>
						</div>

						<div className="max-h-[200px] sm:max-h-[250px] overflow-y-auto" onScroll={handleListScroll}>
							{visibleDevices.map((d) => (
								<div key={d.imei} className={`grid grid-cols-[32px_125px_1fr_1fr] items-center border-b border-[#f1f2f5] px-3 py-2 text-[12px] last:border-b-0 ${d.error ? 'bg-[#ffe8ea]' : 'bg-white'}`}>
									<GUICheckbox checked={selectedImeis.includes(d.imei)} onChange={() => toggleRow(d.imei)} />
									<p className="truncate">{d.imei}</p>
									<span className="truncate">{d.model ? d.model.toUpperCase() : t('label.unknown')}</span>

									<span className="inline-flex items-center gap-1">
										{d.error ? (
											<>
												<Close fill="#9a0000" size={16} />
												<p className="text-[#9a0000]">{t('label.error')}</p>
											</>
										) : (
											<>
												<Check fill="#3f8f3f" size={16} />
												<p className="text-[#3f8f3f]">{t('label.status-success')}</p>
											</>
										)}
									</span>
								</div>
							))}
						</div>

						{devices.some((d) => d.error) && (
							<div className="max-h-[100px] sm:max-h-auto space-y-1 border-t border-[#eceff3] bg-[#fff8f8] px-3 py-2 overflow-y-auto">
								{devices
									.filter((d) => d.error)
									.map((d) => (
										<p key={`${d.imei}-${d.error}`} className="text-[12px] text-[#9a0000]">
											{d.imei}: {d.error}
										</p>
									))}
							</div>
						)}
					</div>
				</div>
			)}

			{devices.length > 0 && (
				<div className="space-y-2">
					{importProgress && (
						<p className="text-[12px] text-[#49525f]">
							{t('label.import')}: {importProgress.done}/{importProgress.total}
						</p>
					)}
					<GUIButton onClick={handleImportDevices} disabled={!hasSelectedForImport || selectedHaveErrors || busy} title={selectedHaveErrors ? t('message.import-error') : undefined}>
						{t('label.add')}
					</GUIButton>
				</div>
			)}
		</div>
	);
};

export default DeviceImportsModal;
