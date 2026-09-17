import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { formatRelativeTime } from '@/utils/TimeUtils';
import { ROUTES, buildRoute } from '@/constants/constants';
import { FetchFn } from '@/components/common/Table/GenericTable/GenericTable';
import CompactTable from '@/components/common/Table/CompactTable/CompactTable';
import { basicCompaniesGetWithParams, CompanyResponse } from '@/rest/companyAPI';
import { getCompanyStatusClass, getCompanyStatusLabel } from '@/utils/bulkConfigurationDetailsUtils';

const BulkConfigurationList = () => {
	const { t } = useTranslation();

	const navigate = useNavigate();

	const fetchCompanies: FetchFn<CompanyResponse> = useCallback(async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicCompaniesGetWithParams(page, limit, search, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	}, []);

	return (
		<div className="space-y-3">
			{/* <label className="inline-flex cursor-pointer border rounded px-4 py-2 text-sm hover:bg-[#E8EDF5]">
				{t('message.bulk-import-csv')}
				<input
					type="file"
					accept=".csv,text/csv"
					className="hidden"
					onChange={async (event) => {
						const file = event.target.files?.[0];
						event.target.value = '';
						if (!file) return;
						try {
							if (file.size > 10 * 1024 * 1024) throw new Error('invalid-csv');
							const state = importBulkConfigurationCsv(await file.text());
							navigate(ROUTES.CONFIGURATIONS_BULK_NEW, { state });
							toast.info(t('message.bulk-import-csv-ready'));
						} catch (error) {
							toast.error(t(error instanceof Error && error.message === 'missing-source' ? 'message.bulk-import-csv-missing-source' : 'message.bulk-import-csv-invalid'));
						}
					}}
				/>
			</label> */}
			<CompactTable<CompanyResponse>
				tableKey="bulk-configurations"
				fetchFn={fetchCompanies}
				sysToken={0}
				getRowId={(company) => company.id}
				renderTitle={(company) => company.name}
				renderSubtitle={(company) => (
					<span className="flex flex-wrap items-center gap-x-3 gap-y-0.5">
						<span>{t('message.source-value', { value: company.configuration_source })}</span>
						<span>{t('message.ttl-days', { count: company.ttl })}</span>
					</span>
				)}
				renderLabel={(company) => <span className={`px-2 rounded-[6px] border text-[13px] whitespace-nowrap ${getCompanyStatusClass(company.status)}`}>{getCompanyStatusLabel(company.status, t)}</span>}
				renderTime={(company) => formatRelativeTime(company.updated_at, t)}
				onTitleClick={(company) => navigate(buildRoute(ROUTES.CONFIGURATIONS_BULK_DETAILS, { id: company.id }))}
				emptyComponent={<div className="px-4 py-8 text-center text-sm text-[#656d76]">{t('message.bulk-configurations-not-found')}</div>}
			/>
		</div>
	);
};

export default BulkConfigurationList;
