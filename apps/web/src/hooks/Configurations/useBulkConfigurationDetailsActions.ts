import { toast } from 'react-toastify';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { compactTableQueryKey } from '@/components/common/Table/CompactTable/types';
import { BULK_CONFIGURATION_MAX_DEVICES } from '@/constants/BulkConfiguration.constant';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';
import { basicCompanyCancelPending, basicCompanyRetryFailed, CompanyDetailsResponse } from '@/rest/companyAPI';

export type BulkConfigurationCloneFrom = Pick<CompanyCreateRequest, 'name' | 'ttl' | 'existing_task_action' | 'launch_mode' | 'incompatible_action' | 'configuration_cource'>;

export const useBulkConfigurationDetailsActions = (companyId: number, company: CompanyDetailsResponse | undefined) => {
	const { t } = useTranslation();

	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const [actionLoading, setActionLoading] = useState<'retry' | 'cancel' | null>(null);

	const invalidateCompany = useCallback(async () => {
		await queryClient.invalidateQueries({ queryKey: ['respCompanyGetById', companyId] });
		await queryClient.invalidateQueries({ queryKey: compactTableQueryKey('bulk-configurations') });
	}, [companyId, queryClient]);

	const handleRetryFailed = useCallback(async () => {
		if (!company) {
			return;
		}

		setActionLoading('retry');
		try {
			const message = await basicCompanyRetryFailed(companyId);

			toast.success(message);

			await invalidateCompany();
		} catch {
		} finally {
			setActionLoading(null);
		}
	}, [company, companyId, invalidateCompany]);

	const handleCancelPending = useCallback(async () => {
		if (!company) {
			return;
		}

		setActionLoading('cancel');
		try {
			const message = await basicCompanyCancelPending(companyId);

			toast.success(message);

			await invalidateCompany();
		} catch {
		} finally {
			setActionLoading(null);
		}
	}, [company, companyId, invalidateCompany]);

	const handleCreateBasedOnThis = useCallback(() => {
		if (!company) {
			return;
		}

		const selectedIds = company.tasks.map((task) => task.device_imei).filter((imei): imei is string => Boolean(imei));

		if (selectedIds.length === 0) {
			toast.error(t('message.provisioning-no-devices-to-create'));
			return;
		}

		if (selectedIds.length > BULK_CONFIGURATION_MAX_DEVICES) {
			toast.error(t('message.provisioning-company-max-devices-exceeded', { count: BULK_CONFIGURATION_MAX_DEVICES }));
			return;
		}

		const cloneFrom: BulkConfigurationCloneFrom = {
			name: `${company.name} (${t('label.copy')})`,
			ttl: String(company.ttl),
			existing_task_action: company.existing_task_action,
			launch_mode: company.launch_mode,
			incompatible_action: company.incompatible_action,
			configuration_cource: company.configuration_source,
		};

		navigate(ROUTES.CONFIGURATIONS_BULK_NEW, {
			state: {
				selectedIds,
				cloneFrom,
			},
		});
	}, [company, navigate, t]);

	return {
		actionLoading,
		handleRetryFailed,
		handleCancelPending,
		handleCreateBasedOnThis,
	};
};
