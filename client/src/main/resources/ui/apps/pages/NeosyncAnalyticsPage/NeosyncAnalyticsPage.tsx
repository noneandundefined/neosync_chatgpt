import PageLayout from '../PageLayout';
import Modal from '@/components/Modal/Modal';
import { getPeriods } from '@/utils/DateUtils';
import { useTranslation } from 'react-i18next';
import { basicUsersGetFull } from '@/rest/userAPI';
import { useEffect, useMemo, useState } from 'react';
import { basicGroupsGetFull } from '@/rest/groupAPI';
import { basicFirmwareDeviceModelsGet } from '@/rest/firmwaresAPI';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUIButton from '@/components/ui/Button/GUIButton';
import { basicAnalyticWidgets } from '@/rest/analyticAPI';
import { useModalContext } from '@/context/useModalContext';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import WidgetRenderer from '@/components/common/Analytics/WidgetRenderer';
import { Loading, SpinnerLoading } from '@/components/common/Loader/Loading';
import DocumentationModal from '@/components/common/Analytics/Modals/DocumentationModal';
import { getAnalyticsWidgetTitle } from '@/utils/AnalyticsFormatUtils';
import { AnalyticsFiltersProvider, useAnalyticsFilters } from '@/context/useAnalyticsFiltersContext';
import { GUInput } from '@/components/ui/Input/GUInput';

const sections = [
	{ id: 'overview', titleKey: 'analytics-section-overview', descriptionKey: 'analytics-section-overview-description' },
	{ id: 'funnels', titleKey: 'analytics-section-funnels', descriptionKey: 'analytics-section-funnels-description' },
	{ id: 'behavior', titleKey: 'analytics-section-behavior', descriptionKey: 'analytics-section-behavior-description' },
	{ id: 'configurations', titleKey: 'analytics-section-configurations', descriptionKey: 'analytics-section-configurations-description' },
	{ id: 'devices', titleKey: 'analytics-section-devices', descriptionKey: 'analytics-section-devices-description' },
	{ id: 'errors', titleKey: 'analytics-section-errors', descriptionKey: 'analytics-section-errors-description' },
	{ id: 'retention', titleKey: 'analytics-section-retention', descriptionKey: 'analytics-section-retention-description' },
	{ id: 'acquisition', titleKey: 'analytics-section-acquisition', descriptionKey: 'analytics-section-acquisition-description' },
	{ id: 'performance', titleKey: 'analytics-section-performance', descriptionKey: 'analytics-section-performance-description' },
	{ id: 'tools', titleKey: 'analytics-section-tools', descriptionKey: 'analytics-section-tools-description' },
] as const;

const NeosyncAnalyticsPage = () => (
	<PageLayout>
		<AnalyticsFiltersProvider>
			<AnalyticsWorkspace />
		</AnalyticsFiltersProvider>
	</PageLayout>
);

const AnalyticsWorkspace = () => {
	const { t } = useTranslation();
	const { open } = useModalContext();
	const { setDateFilter, setAccountFilter, accountFilter, setGroupFilter, groupFilter, setModelFilter, modelFilter, reset } = useAnalyticsFilters();
	const [activeSection, setActiveSection] = useState<string>('overview');
	const [search, setSearch] = useState('');
	const { data: widgets, loading: widgetsLoading } = useHandleServer(['respAnalyticWidgets'], basicAnalyticWidgets);
	const { data: users, loading: usersLoading } = useHandleServer(['respUsersGetFull'], basicUsersGetFull);
	const { data: groups, loading: groupsLoading } = useHandleServer(['respGroupsGetFull'], basicGroupsGetFull);
	const { data: models, loading: modelsLoading } = useHandleServer(['respAnalyticsDeviceModels'], basicFirmwareDeviceModelsGet);
	const periods = useMemo(() => getPeriods(), []);
	const [periodIndex, setPeriodIndex] = useState(0);
	const currentSection = sections.find((section) => section.id === activeSection) || sections[0];
	const visibleWidgets = useMemo(() => {
		const normalizedSearch = search.trim().toLocaleLowerCase();
		return (widgets || []).filter((widget) => widget.section === activeSection && (!normalizedSearch || getAnalyticsWidgetTitle(t, widget).toLocaleLowerCase().includes(normalizedSearch)));
	}, [widgets, activeSection, search, t]);

	useEffect(() => {
		const today = periods[0];
		setDateFilter({ from: today.from.toISOString(), to: today.to.toISOString() });
	}, [periods, setDateFilter]);

	const resetFilters = () => {
		reset();
		setPeriodIndex(0);
		const today = periods[0];
		setDateFilter({ from: today.from.toISOString(), to: today.to.toISOString() });
		localStorage.removeItem(CACHEKEYs.ANALYTIC_WIDGETS);
	};

	return (
		<div className="min-w-0 w-full overflow-x-hidden overflow-y-auto pb-5">
			<div className="mb-4">
				<h1 className="text-xl font-medium text-[#30343b]">{t('label.analytics-title')}</h1>
			</div>

			<section className="bg-[#f7f7f7] border rounded-xl p-3 sm:p-4 mb-4">
				<div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-x-4 gap-y-3">
					<div className="min-w-0">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-report-section')}</p>
						<GUISelect className="relative !min-w-0 !w-full" value={activeSection} onChange={(event) => setActiveSection(String(event.target.value))}>
							{sections.map((section) => (
								<option key={section.id} value={section.id}>
									{t(`label.${section.titleKey}`)}
								</option>
							))}
						</GUISelect>
					</div>
					<div className="min-w-0">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-period')}</p>
						<GUISelect
							className="relative !min-w-0 !w-full"
							value={periodIndex}
							onChange={(event) => {
								const index = Number(event.target.value);
								setPeriodIndex(index);
								const period = periods[index];
								setDateFilter({ from: period.from.toISOString(), to: period.to.toISOString() });
							}}
						>
							{periods.map((period, index) => (
								<option key={index} value={index}>
									{t(`label.${period.label}`)}
								</option>
							))}
						</GUISelect>
					</div>
					<div className="min-w-0">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-account-filter')}</p>
						<GUISelect className="relative !min-w-0 !w-full" onChange={(event) => setAccountFilter(String(event.target.value))} value={accountFilter ?? ''}>
							{!usersLoading && <option value="">{t('message.all-accounts')}</option>}
							{usersLoading ? (
								<option value="">
									<Loading />
								</option>
							) : (
								users?.map((user) => (
									<option value={user.user_uuid} key={user.user_uuid}>
										{user.email}
									</option>
								))
							)}
						</GUISelect>
					</div>
					<div className="min-w-0">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-group-filter')}</p>
						<GUISelect className="relative !min-w-0 !w-full" onChange={(event) => setGroupFilter(event.target.value ? Number(event.target.value) : null)} value={groupFilter ?? ''}>
							{!groupsLoading && <option value="">{t('message.all-groups')}</option>}
							{groupsLoading ? (
								<option value="">
									<Loading />
								</option>
							) : (
								groups?.map((group) => (
									<option value={group.id} key={group.id}>
										{group.name}
									</option>
								))
							)}
						</GUISelect>
					</div>
					<div className="min-w-0">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-model-filter')}</p>
						<GUISelect className="relative !min-w-0 !w-full" onChange={(event) => setModelFilter(String(event.target.value) || null)} value={modelFilter ?? ''}>
							{!modelsLoading && <option value="">{t('message.analytics-all-models')}</option>}
							{modelsLoading ? (
								<option value="">
									<Loading />
								</option>
							) : (
								models?.map((model) => (
									<option value={model} key={model}>
										{model}
									</option>
								))
							)}
						</GUISelect>
					</div>
					<div className="min-w-0 sm:col-span-2 xl:col-span-2">
						<p className="text-xs text-[#7b838d] mb-1">{t('label.analytics-report-search')}</p>
						<GUInput
							value={search}
							onChange={(event) => setSearch(event.target.value)}
							placeholder={t('message.analytics-report-search-placeholder')}
						/>
					</div>
					<div className="flex flex-col gap-2 xl:justify-end xl:col-span-1">
						<GUIButton className="w-full sm:w-auto" onClick={resetFilters}>
							{t('label.analytics-reset-filters')}
						</GUIButton>
						<button
							type="button"
							className="text-sm text-[#1976d2] hover:text-[#00479b] py-2"
							onClick={() =>
								open(
									<Modal title={t('label.documentation')}>
										<DocumentationModal />
									</Modal>
								)
							}
						>
							{t('label.documentation')}
						</button>
					</div>
				</div>
			</section>

			<div className="mb-4">
				<h2 className="text-lg font-medium text-[#30343b]">{t(`label.${currentSection.titleKey}`)}</h2>
				<p className="text-xs sm:text-sm text-[#7b838d] mt-1">{t(`message.${currentSection.descriptionKey}`)}</p>
			</div>

			{widgetsLoading && (
				<div className="py-16">
					<SpinnerLoading />
				</div>
			)}
			{activeSection === 'tools' ? (
				<div className="grid grid-cols-12">
					<WidgetRenderer widget={{ id: 0, section: 'tools', type: 'sql', title: '', query: '' }} />
				</div>
			) : (
				<div className="grid grid-cols-12 gap-3 sm:gap-4 items-stretch">
					{visibleWidgets.map((widget) => (
						<WidgetRenderer key={widget.id} widget={widget} allWidgets={widgets || undefined} />
					))}
				</div>
			)}
			{!widgetsLoading && activeSection !== 'tools' && !visibleWidgets.length && (
				<div className="bg-[#f7f7f7] border rounded-xl p-8 sm:p-12 text-center text-sm text-[#7b838d]">{t('message.analytics-reports-not-found')}</div>
			)}
		</div>
	);
};

export default NeosyncAnalyticsPage;
