import PageLayout from '../PageLayout';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';

import BulkConfigurationList from '@/components/common/Configurations/BulkConfiguration/BulkConfigurationList';
import ConfigurationTemplateList from '@/components/common/Configurations/ConfigurationTemplate/ConfigurationTemplateList';

type Tab = 'templates' | 'companies' | 'tasks';

const ConfigurationsPage = () => {
	const { t } = useTranslation();

	const [searchParams, setSearchParams] = useSearchParams();

	const tab = (searchParams.get('tab') as Tab) ?? 'templates';

	const changeTab = (newTab: Tab) => {
		setSearchParams({ tab: newTab });
	};

	return (
		<PageLayout>
			<div className="space-y-5">
				<div className="flex border-b">
					<div className="flex gap-[3rem]">
						<p className={`text-center py-3 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'templates' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => changeTab('templates')}>
							{t('label.configuration-templates')}
						</p>

						<p className={`text-center py-3 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'companies' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => changeTab('companies')}>
							{t('label.bulk-configurations')}
						</p>
					</div>
				</div>

				{tab === 'templates' && <ConfigurationTemplateList />}

				{tab === 'companies' && <BulkConfigurationList />}
			</div>
		</PageLayout>
	);
};

export default ConfigurationsPage;
