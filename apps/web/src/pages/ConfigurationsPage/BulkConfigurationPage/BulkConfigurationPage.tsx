import { toast } from 'react-toastify';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import Header from '@/components/Header/Header';
import { useQueryClient } from '@tanstack/react-query';
import { basicCompanyCreate } from '@/rest/companyAPI';
import { basicDeviceStatusDetailed } from '@/rest/deviceAPI';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { compactTableQueryKey } from '@/components/common/Table/CompactTable/types';
import { basicConfigurationTemplateGetById } from '@/rest/configurationTemplateAPI';
import { BULK_CONFIGURATION_MAX_DEVICES } from '@/constants/BulkConfiguration.constant';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';
import { BulkConfigurationCloneFrom } from '@/hooks/Configurations/useBulkConfigurationDetailsActions';
import { getBulkDeviceSourceEligibility, getDominantFirmwareVersion, isBulkTemplateCompatible, buildBulkProvisioningDefaultName } from '@/utils/bulkConfigurationSourceUtils';

import GUIButton from '@/components/ui/Button/GUIButton';

import PreLaunchCheckStep from '@/components/common/Configurations/BulkConfiguration/Steps/PreLaunchCheckStep';
import GeneralSettingsStep from '@/components/common/Configurations/BulkConfiguration/Steps/GeneralSettingsStep';
import ScopeOfApplicationStep from '@/components/common/Configurations/BulkConfiguration/Steps/ScopeOfApplicationStep';
import ConfigurationSourceStep from '@/components/common/Configurations/BulkConfiguration/Steps/ConfigurationSourceStep';

export type DeviceStatBC = {
	compatible: number;
	incompatible: number;
	ready: number;
	offline: number;
	priorityModel: string;
	models: string[];
	modelCounts: Record<string, number>;
	dominantFirmware: number | null;
};

const fileToBase64 = (file: File) =>
	new Promise<string>((resolve, reject) => {
		const reader = new FileReader();

		reader.onload = () => {
			const result = reader.result;

			if (typeof result !== 'string') {
				reject(new Error('invalid-file'));
				return;
			}

			const base64 = result.includes(',') ? result.split(',')[1] : result;
			resolve(base64);
		};

		reader.onerror = () => reject(reader.error ?? new Error('invalid-file'));
		reader.readAsDataURL(file);
	});

const BulkConfiguration = () => {
	const { t } = useTranslation();
	const { state } = useLocation();

	const navigate = useNavigate();

	const queryClient = useQueryClient();

	const selectedIdsRef = useRef<string[] | null>(null);
	const cloneFromRef = useRef<BulkConfigurationCloneFrom | null>(null);

	if (state?.selectedIds) {
		selectedIdsRef.current = state.selectedIds as string[];
	}

	if (state?.cloneFrom) {
		cloneFromRef.current = state.cloneFrom as BulkConfigurationCloneFrom;
	}

	const selectedIds = selectedIdsRef.current ?? [];
	const cloneFrom = cloneFromRef.current;
	const importedConfiguration = useRef<CompanyCreateRequest | undefined>(state?.importedConfiguration).current;
	const previousSourceRef = useRef(importedConfiguration?.configuration_cource ?? cloneFrom?.configuration_cource ?? 'template_sources');
	const exceedsDeviceLimit = selectedIds.length > BULK_CONFIGURATION_MAX_DEVICES;
	const lastAutoNameRef = useRef<string | null>(null);
	const preserveCloneNameRef = useRef(Boolean(cloneFrom?.name || importedConfiguration?.name));
	const [nameCreatedAt] = useState(() => new Date());

	const [configurationFile, setConfigurationFile] = useState<File | null>(() => {
		if (!importedConfiguration?.configuration_file_base64 || !importedConfiguration.configuration_file_name) return null;
		const bytes = Uint8Array.from(atob(importedConfiguration.configuration_file_base64), (char) => char.charCodeAt(0));
		return new File([bytes], importedConfiguration.configuration_file_name);
	});
	const [selectedModel, setSelectedModel] = useState<string | null>(importedConfiguration?.priority_model ?? null);
	const fileReadVersionRef = useRef(0);

	const fetchDeviceStatusDetailed = useCallback(() => {
		if (!selectedIds.length) {
			return Promise.resolve([]);
		}

		return basicDeviceStatusDetailed(selectedIds);
	}, [selectedIds]);
	const { data: respDeviceStatusDetailed } = useHandleServer(['respDeviceStatusDetailed', selectedIds], fetchDeviceStatusDetailed);

	const stats: DeviceStatBC = useMemo(() => {
		if (!respDeviceStatusDetailed?.length) {
			return {
				compatible: 0,
				incompatible: 0,
				ready: 0,
				offline: 0,
				priorityModel: '',
				models: [] as string[],
				modelCounts: {},
				dominantFirmware: null,
			};
		}

		const modelCount = new Map<string, number>();

		for (const device of respDeviceStatusDetailed) {
			const model = device.device_model;
			if (!model) continue;

			modelCount.set(model, (modelCount.get(model) ?? 0) + 1);
		}

		let priorityModel = '';
		let compatible = 0;

		for (const [model, count] of [...modelCount].sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))) {
			if (count > compatible) {
				compatible = count;
				priorityModel = model;
			}
		}

		if (selectedModel && modelCount.has(selectedModel)) {
			priorityModel = selectedModel;
			compatible = modelCount.get(selectedModel)!;
		}

		let ready = 0;
		let offline = 0;
		let incompatible = 0;

		for (const device of respDeviceStatusDetailed) {
			if (device.device_model !== priorityModel) {
				incompatible++;
				continue;
			}

			if (device.status) {
				ready++;
			} else {
				offline++;
			}
		}

		return {
			compatible,
			incompatible,
			ready,
			offline,
			priorityModel,
			models: [...modelCount.keys()],
			modelCounts: Object.fromEntries(modelCount),
			dominantFirmware: getDominantFirmwareVersion(respDeviceStatusDetailed.filter((device) => device.device_model === priorityModel)),
		};
	}, [respDeviceStatusDetailed, selectedModel]);

	const {
		register,
		control,
		setValue,
		getValues,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<CompanyCreateRequest>({
		defaultValues: {
			name: cloneFrom?.name ?? '',
			ttl: cloneFrom?.ttl ?? '1',
			existing_task_action: cloneFrom?.existing_task_action ?? 'skip',
			launch_mode: cloneFrom?.launch_mode ?? 'all',
			incompatible_action: cloneFrom?.incompatible_action ?? 'exclude',
			configuration_cource: cloneFrom?.configuration_cource ?? 'template_sources',
			configuration_device_cource: null,
			configuration_template_id: null,
			configuration_file_name: null,
			configuration_file_base64: null,
			...importedConfiguration,
		},
	});

	const configurationSource = watch('configuration_cource');
	const selectedImei = watch('configuration_device_cource');
	const selectedTemplateId = watch('configuration_template_id');
	const configurationFileName = watch('configuration_file_name');

	const fetchSelectedTemplate = useCallback(() => {
		if (!selectedTemplateId) {
			return Promise.resolve(null);
		}

		return basicConfigurationTemplateGetById(selectedTemplateId);
	}, [selectedTemplateId]);

	const { data: selectedTemplate } = useHandleServer(['bulkConfigurationTemplate', selectedTemplateId], fetchSelectedTemplate, {
		enabled: configurationSource === 'template_sources' && !!selectedTemplateId,
	});

	const provisioningSourceName = useMemo(() => {
		switch (configurationSource) {
			case 'device_sources':
				return selectedImei ?? null;
			case 'file_sources':
				return configurationFileName;
			case 'template_sources':
				return selectedTemplate?.name ?? null;
			default:
				return null;
		}
	}, [configurationSource, selectedImei, configurationFileName, selectedTemplate?.name]);

	const autoProvisioningName = useMemo(
		() =>
			buildBulkProvisioningDefaultName(t, {
				model: stats.priorityModel,
				source: configurationSource,
				sourceName: provisioningSourceName,
				selectedCount: selectedIds.length,
				createdAt: nameCreatedAt,
			}),
		[t, stats.priorityModel, configurationSource, provisioningSourceName, selectedIds.length, nameCreatedAt]
	);

	const handleGenerateName = () => {
		preserveCloneNameRef.current = false;
		lastAutoNameRef.current = autoProvisioningName;
		setValue('name', autoProvisioningName, { shouldDirty: true, shouldValidate: true });
	};

	useEffect(() => {
		if (previousSourceRef.current === configurationSource) return;
		previousSourceRef.current = configurationSource;
		fileReadVersionRef.current++;
		setValue('configuration_device_cource', null);
		setValue('configuration_template_id', null);
		setValue('configuration_file_name', null);
		setValue('configuration_file_base64', null);
		setConfigurationFile(null);
	}, [configurationSource, setValue]);

	const handleSelectModel = (model: string) => {
		if (model === stats.priorityModel || !stats.models.includes(model)) return;
		setSelectedModel(model);
		fileReadVersionRef.current++;
		if (configurationFileName) {
			setValue('configuration_file_name', null, { shouldDirty: true });
			setValue('configuration_file_base64', null, { shouldDirty: true });
			setConfigurationFile(null);
			toast.info(t('message.provisioning-model-change-file-reset'));
		}
	};

	useEffect(() => {
		if (!stats.priorityModel) return;
		if (configurationSource === 'device_sources' && selectedImei) {
			const device = respDeviceStatusDetailed?.find((item) => item.imei === selectedImei);
			if (!device || !getBulkDeviceSourceEligibility(device, stats.priorityModel, stats.dominantFirmware).eligible) {
				setValue('configuration_device_cource', null, { shouldDirty: true });
				toast.info(t('message.provisioning-model-change-device-reset'));
			}
		}
		if (configurationSource === 'template_sources' && selectedTemplateId && selectedTemplate?.id === selectedTemplateId && !isBulkTemplateCompatible(selectedTemplate, stats.priorityModel)) {
			setValue('configuration_template_id', null, { shouldDirty: true });
			toast.info(t('message.provisioning-model-change-template-reset'));
		}
	}, [configurationSource, selectedImei, selectedTemplateId, selectedTemplate, stats.priorityModel, stats.dominantFirmware, respDeviceStatusDetailed, setValue, t]);

	const handleSelectDevice = (imei: string) => {
		if (configurationSource !== 'device_sources') {
			return;
		}

		setValue('configuration_device_cource', selectedImei === imei ? null : imei, {
			shouldDirty: true,
		});
	};

	const handleSelectTemplate = (templateId: number) => {
		if (configurationSource !== 'template_sources') {
			return;
		}

		setValue('configuration_template_id', selectedTemplateId === templateId ? null : templateId, {
			shouldDirty: true,
		});
	};

	const handleSelectConfigurationFile = async (file: File) => {
		const readVersion = ++fileReadVersionRef.current;
		setConfigurationFile(file);
		setValue('configuration_file_name', file.name, { shouldDirty: true });
		setValue('configuration_file_base64', null, { shouldDirty: true });
		const base64 = await fileToBase64(file);
		if (readVersion !== fileReadVersionRef.current) return;
		setValue('configuration_file_base64', base64, { shouldDirty: true });
	};

	useEffect(() => {
		if (preserveCloneNameRef.current) {
			return;
		}

		const currentName = getValues('name');
		if (currentName && lastAutoNameRef.current !== null && currentName !== lastAutoNameRef.current) {
			return;
		}

		lastAutoNameRef.current = autoProvisioningName;
		setValue('name', autoProvisioningName);
	}, [cloneFrom?.name, autoProvisioningName, getValues, setValue]);

	const onSubmit = async (data: CompanyCreateRequest) => {
		if (exceedsDeviceLimit) {
			toast.error(t('message.provisioning-company-max-devices-exceeded', { count: BULK_CONFIGURATION_MAX_DEVICES }));
			return;
		}

		if (data.configuration_cource === 'device_sources' && !data.configuration_device_cource) {
			toast.error(t('message.provisioning-configuration-device-source-required'));
			return;
		}

		if (data.configuration_cource === 'device_sources' && data.configuration_device_cource) {
			const sourceDevice = respDeviceStatusDetailed?.find((device) => device.imei === data.configuration_device_cource);
			if (!sourceDevice) {
				toast.error(t('message.provisioning-configuration-device-source-not-found'));
				return;
			}

			const eligibility = getBulkDeviceSourceEligibility(sourceDevice, stats.priorityModel, stats.dominantFirmware);
			if (!eligibility.eligible) {
				toast.error(t(eligibility.reasonKey));
				return;
			}
		}

		if (data.configuration_cource === 'template_sources' && !data.configuration_template_id) {
			toast.error(t('message.provisioning-configuration-template-source-required'));
			return;
		}

		if (data.configuration_cource === 'template_sources' && data.configuration_template_id) {
			const template = await basicConfigurationTemplateGetById(data.configuration_template_id);
			if (!isBulkTemplateCompatible(template, stats.priorityModel)) {
				toast.error(t('message.provisioning-configuration-template-model-mismatch-short'));
				return;
			}
		}

		if (data.configuration_cource === 'file_sources' && (!data.configuration_file_base64 || !data.configuration_file_name)) {
			toast.error(t('message.provisioning-configuration-file-source-required'));
			return;
		}

		await basicCompanyCreate({
			...data,
			priority_model: stats.priorityModel || undefined,
			imeis: selectedIds,
		});

		await queryClient.invalidateQueries({ queryKey: compactTableQueryKey('bulk-configurations') });
		navigate(`${ROUTES.CONFIGURATIONS}?tab=companies`);
	};

	if (!selectedIds.length) {
		return <Navigate to={`${ROUTES.CONFIGURATIONS}?tab=companies`} replace />;
	}

	return (
		<div className={`flex flex-col h-[calc(100vh-24px)]`}>
			<Header show={true} />

			<main className="flex-1 mt-4 pb-3 px-3">
				<div className="max-w-[50rem] mx-auto space-y-10">
					<div className="space-y-1">
						<p className="font-medium text-[18px] text-[#49525f]">{t('message.provisioning-wizard-title')}</p>
						<p className="text-[#49525f]">{t('message.provisioning-wizard-description')}</p>
					</div>

					{exceedsDeviceLimit && (
						<div className="rounded-[6px] border border-[#f5c2c7] bg-[#fff5f5] px-4 py-3 text-sm text-[#b42318]">
							{t('message.provisioning-company-max-devices-warning', { count: BULK_CONFIGURATION_MAX_DEVICES, selected: selectedIds.length })}
						</div>
					)}

					<div className="space-y-5">
						<GeneralSettingsStep register={register} control={control} errors={errors} onGenerateName={handleGenerateName} />

						<ScopeOfApplicationStep register={register} stats={stats} state={{ selectedIds }} exceedsDeviceLimit={exceedsDeviceLimit} onSelectModel={handleSelectModel} />

						<ConfigurationSourceStep
							register={register}
							configurationSource={configurationSource}
							devices={respDeviceStatusDetailed}
							priorityModel={stats.priorityModel}
							dominantFirmware={stats.dominantFirmware}
							selectedImei={selectedImei ?? null}
							onSelectDevice={handleSelectDevice}
							selectedTemplateId={selectedTemplateId ?? null}
							onSelectTemplate={handleSelectTemplate}
							configurationFile={configurationFile}
							onSelectConfigurationFile={handleSelectConfigurationFile}
						/>

						<PreLaunchCheckStep stats={stats} />
					</div>

					<GUIButton onClick={handleSubmit(onSubmit)} disabled={exceedsDeviceLimit}>
						{t('label.create')}
					</GUIButton>
				</div>
			</main>
		</div>
	);
};

export default BulkConfiguration;
