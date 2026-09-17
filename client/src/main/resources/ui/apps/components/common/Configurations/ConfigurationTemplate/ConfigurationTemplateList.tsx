import { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { formatRelativeTime } from '@/utils/TimeUtils';
import { ROUTES, buildRoute } from '@/constants/constants';
import CompactTable from '@/components/common/Table/CompactTable/CompactTable';
import { FetchFn } from '@/components/common/Table/GenericTable/GenericTable';
import OfficeBuildingCogOutline from '@/components/@icons/office-building-cog-outline';
import { basicConfigurationTemplateDelete, basicConfigurationTemplatesGetWithParams, ConfigurationTemplateResponse } from '@/rest/configurationTemplateAPI';

const ConfigurationTemplateList = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	const [sysToken, setSysToken] = useState(0);

	const typeOfSavingLabel = (type: ConfigurationTemplateResponse['type_of_saving']) => (type === 'full' ? t('message.configuration-template-type-full') : t('message.configuration-template-type-modified'));

	const handleDeleteMany = async (ids: (string | number)[]) => {
		await Promise.all(ids.map((id) => basicConfigurationTemplateDelete(Number(id))));
		setSysToken((v) => v + 1);
	};

	const fetchTemplates: FetchFn<ConfigurationTemplateResponse> = useCallback(async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicConfigurationTemplatesGetWithParams(page, limit, search, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	}, []);

	return (
		<div className="space-y-4">
			<div className="flex justify-end">
				<button id="buttonhlp" onClick={() => navigate(ROUTES.CONFIGURATIONS_TEMPLATE_NEW)}>
					<div className="flex items-center gap-2 w-full justify-center sm:justify-end">
						<OfficeBuildingCogOutline fill="#62666d" size={20} />
						<p className="text-[13px] whitespace-nowrap">{t('label.create-configuration-template')}</p>
					</div>
				</button>
			</div>

			<CompactTable<ConfigurationTemplateResponse>
				tableKey="configuration-templates"
				fetchFn={fetchTemplates}
				sysToken={sysToken}
				getRowId={(template) => template.id}
				renderTitle={(template) => template.name}
				renderSubtitle={(template) => (template.model ? `${t('label.model')}: ${template.model}` : undefined)}
				renderLabel={(template) => <span className="px-3 py-[2px] rounded-full border text-[#49525f] text-[13px] whitespace-nowrap">{typeOfSavingLabel(template.type_of_saving)}</span>}
				renderTime={(template) => formatRelativeTime(template.updated_at, t)}
				onTitleClick={(template) => navigate(buildRoute(ROUTES.CONFIGURATIONS_TEMPLATE_EDIT, { id: template.id }))}
				onDeleteFn={handleDeleteMany}
				emptyComponent={<div className="px-4 py-8 text-center text-sm text-[#656d76]">{t('message.configuration-templates-not-found')}</div>}
			/>
		</div>
	);
};

export default ConfigurationTemplateList;
