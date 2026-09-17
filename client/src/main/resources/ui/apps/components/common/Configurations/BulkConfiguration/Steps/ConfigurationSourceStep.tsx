import IndexStep from './IndexStep';
import { toast } from 'react-toastify';
import FileHelper from '@/utils/FilesUtils';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { UseFormRegister } from 'react-hook-form';
import GUIRadio from '@/components/ui/Radio/GUIRadio';
import { formatRelativeTime } from '@/utils/TimeUtils';
import LabeledField from '@/components/ui/Form/LabeledField';
import { DeviceDetailedStatusResponse } from '@/rest/deviceAPI';
import { FetchFn } from '@/components/common/Table/GenericTable/GenericTable';
import CompactTable from '@/components/common/Table/CompactTable/CompactTable';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';
import { getConfigurationSyncStatusLabel } from '@/components/common/Devices/common/ConfigurationSyncStatus';
import { getBulkDeviceSourceEligibility, getBulkDeviceSourceEligibilityLabel, isBulkTemplateCompatible } from '@/utils/bulkConfigurationSourceUtils';
import { basicConfigurationTemplatesGetWithParams, ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';

interface ConfigurationSourceStepProps {
	register: UseFormRegister<CompanyCreateRequest>;
	configurationSource: CompanyCreateRequest['configuration_cource'];
	devices?: DeviceDetailedStatusResponse[] | null;
	priorityModel: string;
	dominantFirmware: number | null;
	selectedImei: string | null;
	onSelectDevice: (imei: string) => void;
	selectedTemplateId: number | null;
	onSelectTemplate: (templateId: number) => void;
	configurationFile: File | null;
	onSelectConfigurationFile: (file: File) => void | Promise<void>;
}

const ConfigurationSourceStep: React.FC<ConfigurationSourceStepProps> = ({
	register,
	configurationSource,
	devices = [],
	priorityModel,
	dominantFirmware,
	selectedImei,
	onSelectDevice,
	selectedTemplateId,
	onSelectTemplate,
	configurationFile,
	onSelectConfigurationFile,
}) => {
	const { t } = useTranslation();

	const [fileBusy, setFileBusy] = useState(false);

	const fetchDevices: FetchFn<DeviceDetailedStatusResponse> = useCallback(
		async ({ page, limit, search }) => {
			const all = devices ?? [];
			const query = search.trim().toLowerCase();
			const filtered = query ? all.filter((device) => device.imei.toLowerCase().includes(query) || (device.device_model ?? '').toLowerCase().includes(query)) : all;
			const start = (page - 1) * limit;

			return {
				items: filtered.slice(start, start + limit),
				total: filtered.length,
				totalPages: Math.max(1, Math.ceil(filtered.length / limit)),
			};
		},
		[devices]
	);

	const fetchTemplates: FetchFn<ConfigurationTemplateResponse> = useCallback(async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicConfigurationTemplatesGetWithParams(page, limit, search, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	}, []);

	const typeOfSavingLabel = (type: ConfigurationTemplateResponse['type_of_saving']) => (type === 'full' ? t('message.configuration-template-type-full') : t('message.configuration-template-type-modified'));

	const handleDeviceSelectionChange = (ids: (string | number)[]) => {
		if (ids.length === 0) {
			if (selectedImei) {
				onSelectDevice(selectedImei);
			}

			return;
		}

		const imei = String(ids[0]);
		const device = devices?.find((item) => item.imei === imei);
		if (!device) {
			return;
		}

		const eligibility = getBulkDeviceSourceEligibility(device, priorityModel, dominantFirmware);
		if (!eligibility.eligible) {
			toast.error(getBulkDeviceSourceEligibilityLabel(t, eligibility));
			return;
		}

		if (selectedImei !== imei) {
			onSelectDevice(imei);
		}
	};

	const handleTemplateSelectionChange = (ids: (string | number)[]) => {
		if (ids.length === 0) {
			if (selectedTemplateId) {
				onSelectTemplate(selectedTemplateId);
			}

			return;
		}

		const templateId = Number(ids[0]);

		if (selectedTemplateId !== templateId) {
			onSelectTemplate(templateId);
		}
	};

	const handleFileUpload = async (file: File) => {
		const isBinary = await FileHelper.isBin(file);
		if (!isBinary) {
			toast.error(t('message.file-not-binary'));
			return;
		}

		if (!FileHelper.fileMax5Kb(file)) {
			toast.error(t('message.file-too-large-5kb'));
			return;
		}

		await onSelectConfigurationFile(file);
	};

	return (
		<IndexStep step={3} title="message.provisioning-configuration-source">
			<LabeledField label="message.provisioning-configuration-source-select">
				<div className="flex flex-col gap-3">
					<GUIRadio label={t('message.provisioning-configuration-source-template-recommended')} value="template_sources" register={register('configuration_cource')} />
					<GUIRadio label={t('message.provisioning-configuration-source-file')} value="file_sources" register={register('configuration_cource')} />
					<GUIRadio label={t('message.provisioning-configuration-source-device-advanced')} value="device_sources" register={register('configuration_cource')} />
				</div>
			</LabeledField>

			{configurationSource === 'device_sources' && (
				<>
					<p className="font-semibold text-[#9a0000]">{t('message.provisioning-configuration-source-device-warning')}</p>

					<CompactTable<DeviceDetailedStatusResponse>
						tableKey="bulk-configuration-source-devices"
						fetchFn={fetchDevices}
						sysToken={devices?.length ?? 0}
						persistSearchInUrl={false}
						selectionMode="single"
						controlledSelectedIds={selectedImei ? [selectedImei] : []}
						onSelectionChange={handleDeviceSelectionChange}
						getRowId={(device) => device.imei}
						renderTitle={(device) => device.imei}
						renderSubtitle={(device) => {
							const parts = [device.device_model ?? '—'];
							if (device.firmware_version != null && device.firmware_version > 0) {
								parts.push(`${t('message.tracker-firmware-version')}: ${device.firmware_version}`);
							}

							const eligibility = getBulkDeviceSourceEligibility(device, priorityModel, dominantFirmware);
							if (!eligibility.eligible) {
								parts.push(getBulkDeviceSourceEligibilityLabel(t, eligibility));
							}

							return parts.join(' · ');
						}}
						renderLabel={(device) => <span className="px-3 py-[2px] rounded-full border text-[#49525f] text-[13px] whitespace-nowrap">{getConfigurationSyncStatusLabel(t, device.configuration_sync_status, null)}</span>}
						emptyComponent={<div className="px-4 py-8 text-center text-sm text-[#656d76]">{t('message.data-available')}</div>}
					/>
				</>
			)}

			{configurationSource === 'template_sources' && (
				<CompactTable<ConfigurationTemplateResponse>
					tableKey="bulk-configuration-source-templates"
					fetchFn={fetchTemplates}
					sysToken={0}
					persistSearchInUrl={false}
					selectionMode="single"
					controlledSelectedIds={selectedTemplateId ? [selectedTemplateId] : []}
					onSelectionChange={handleTemplateSelectionChange}
					getRowId={(template) => template.id}
					renderTitle={(template) => template.name}
					renderSubtitle={(template) => {
						const parts: string[] = [];
						if (template.model) {
							parts.push(`${t('label.model')}: ${template.model}`);
						}

						if (priorityModel && template.model && !isBulkTemplateCompatible(template, priorityModel)) {
							parts.push(t('message.provisioning-configuration-template-model-mismatch-short'));
						}

						return parts.length ? parts.join(' · ') : undefined;
					}}
					renderLabel={(template) => {
						const compatible = isBulkTemplateCompatible(template, priorityModel);
						return (
							<span className={`px-3 py-[2px] rounded-full border text-[13px] whitespace-nowrap ${compatible ? 'text-[#49525f]' : 'text-[#b42318] border-[#f5c2c7]'}`}>
								{compatible ? typeOfSavingLabel(template.type_of_saving) : t('message.provisioning-incompatible-model-firmware')}
							</span>
						);
					}}
					renderTime={(template) => formatRelativeTime(template.updated_at, t)}
					emptyComponent={<div className="px-4 py-8 text-center text-sm text-[#656d76]">{t('message.configuration-templates-not-found')}</div>}
				/>
			)}

			{configurationSource === 'file_sources' && (
				<div className="flex flex-col items-center gap-3 pt-2">
					<p className="font-semibold text-[#9a0000]">{t('message.provisioning-configuration-file-neosync-hint')}</p>

					<label id="button" className={`!pt-[0.7rem] !w-auto text-center uppercase ${fileBusy ? 'cursor-not-allowed opacity-50 pointer-events-none' : 'cursor-pointer'}`}>
						{t('label.select-file')}
						<input
							type="file"
							className="hidden"
							accept=".neosync"
							disabled={fileBusy}
							onChange={async (e) => {
								const file = e.target.files?.[0];
								e.currentTarget.value = '';
								if (!file) return;

								setFileBusy(true);
								try {
									await handleFileUpload(file);
								} finally {
									setFileBusy(false);
								}
							}}
						/>
					</label>

					{configurationFile && <p className="text-sm text-[#656d76]">{configurationFile.name}</p>}
				</div>
			)}
		</IndexStep>
	);
};

export default ConfigurationSourceStep;
