import { useState } from 'react';
import Table from '@/components/@icons/table';
import { Loading } from '../../Loader/Loading';
import { useTranslation } from 'react-i18next';
import ConsoleLine from '@/components/@icons/console-line';
import { DatabaseColumnResponse, basicMetaGetColumns } from '@/rest/metaAPI';
import { ANALYTIC_FUNCTIONS, ANALYTIC_TABLES, DocTab } from './analyticsDocumentation';

const tabButtonClass = (active: boolean) => ['px-3 py-1 rounded-[6px] cursor-pointer border', active ? 'bg-[#395d95] border-[#395d95]' : 'bg-white hover:bg-[#eee] border'].join(' ');

const tabButtonStyle = (active: boolean): React.CSSProperties => (active ? {} : { boxShadow: '0 0 7px rgba(0, 0, 0, 0.08)' });

const DocumentationModal = () => {
	const { t } = useTranslation();

	const [activeTab, setActiveTab] = useState<DocTab>('tables');
	const [loadingTable, setLoadingTable] = useState<string | null>(null);
	const [openTableIndex, setOpenTableIndex] = useState<number | null>(null);
	const [openFunctionIndex, setOpenFunctionIndex] = useState<number | null>(null);
	const [schemas, setSchemas] = useState<Record<string, DatabaseColumnResponse[]>>({});

	const handleTabChange = (tab: DocTab) => {
		setActiveTab(tab);
		setOpenTableIndex(null);
		setOpenFunctionIndex(null);
	};

	const handleTableToggle = async (table: string, index: number) => {
		const isCurrentlyOpen = openTableIndex === index;

		if (isCurrentlyOpen) {
			setOpenTableIndex(null);
			return;
		}

		setOpenTableIndex(index);

		if (schemas[table]) {
			return;
		}

		try {
			setLoadingTable(table);

			const response = await basicMetaGetColumns(table);

			setSchemas((prev) => ({ ...prev, [table]: response }));
		} finally {
			setLoadingTable(null);
		}
	};

	const handleFunctionToggle = (index: number) => {
		setOpenFunctionIndex((prev) => (prev === index ? null : index));
	};

	return (
		<div className="flex flex-col space-y-3">
			<div className="flex items-center gap-3">
				<button type="button" className={tabButtonClass(activeTab === 'tables')} style={tabButtonStyle(activeTab === 'tables')} onClick={() => handleTabChange('tables')}>
					<p className={activeTab === 'tables' ? 'text-white' : 'text-[#49525f]'}>{t('label.analytics-doc-tab-tables')}</p>
				</button>
				<button type="button" className={tabButtonClass(activeTab === 'functions')} style={tabButtonStyle(activeTab === 'functions')} onClick={() => handleTabChange('functions')}>
					<p className={activeTab === 'functions' ? 'text-white' : 'text-[#49525f]'}>{t('label.analytics-doc-tab-functions')}</p>
				</button>
			</div>

			{activeTab === 'tables' && (
				<div className="max-h-[500px] overflow-y-auto">
					{ANALYTIC_TABLES.map((table, index) => {
						const isOpen = openTableIndex === index;
						const columns = schemas[table.table] ?? [];

						return (
							<div className="hover:bg-[#f6f6f6]" key={table.table}>
								<div className="flex items-center justify-between border-t border-[#e4e4e4] p-3 cursor-pointer" onClick={() => handleTableToggle(table.table, index)}>
									<div className="flex items-center gap-3">
										<div className="bg-[#E8EDF5] p-2 rounded-[6px]">
											<Table fill="#395d95" />
										</div>

										<div>
											<p>{t(table.titleKey)}</p>
											<p className="text-sm">{table.table}</p>
										</div>
									</div>

									<div className="border border-[#395d95] px-3 py-2 rounded">
										<p className="text-[#395d95] text-sm">{isOpen ? t('label.hide') : t('label.more')}</p>
									</div>
								</div>

								{isOpen && (
									<div className="p-1 max-h-[400px] overflow-y-auto">
										{loadingTable === table.table ? (
											<div className="mb-5">
												<Loading />
											</div>
										) : (
											<ul className="space-y-2">
												{columns.map((column) => (
													<li key={column.column_name} className="border rounded p-2 bg-white">
														<div className="flex items-center justify-between">
															<p className="font-medium">{column.column_name}</p>

															<span className="text-xs px-2 py-1 rounded bg-[#E8EDF5] text-[#395d95]">{column.data_type}</span>
														</div>

														<div className="mt-1 text-sm text-[#6b7280]">
															{t('label.analytics-doc-nullable')}: {column.is_nullable}
														</div>

														{column.column_default && (
															<div className="text-sm text-[#6b7280]">
																{t('label.analytics-doc-default')}: {column.column_default}
															</div>
														)}
													</li>
												))}
											</ul>
										)}
									</div>
								)}
							</div>
						);
					})}
				</div>
			)}

			{activeTab === 'functions' && (
				<div>
					{ANALYTIC_FUNCTIONS.map((fn, index) => {
						const isOpen = openFunctionIndex === index;

						return (
							<div className="hover:bg-[#f6f6f6]" key={fn.name}>
								<div className="flex items-center justify-between border-t border-[#e4e4e4] p-3 cursor-pointer" onClick={() => handleFunctionToggle(index)}>
									<div className="flex items-center gap-3 min-w-0">
										<div className="bg-[#E8EDF5] p-2 rounded-[6px] shrink-0">
											<ConsoleLine fill="#395d95" />
										</div>

										<div className="min-w-0">
											<p>{fn.name}</p>
											<p className="text-sm text-[#6b7280] truncate">{fn.signature}</p>
										</div>
									</div>

									<div className="border border-[#395d95] px-3 py-2 rounded shrink-0 ml-3">
										<p className="text-[#395d95] text-sm">{isOpen ? t('label.hide') : t('label.more')}</p>
									</div>
								</div>

								{isOpen && (
									<div className="p-1 max-h-[400px] overflow-y-auto">
										<div className="border rounded p-3 bg-white mb-2">
											<p className="text-sm text-[#49525f]">{t(fn.descriptionKey)}</p>
										</div>

										{fn.parameters.length > 0 && (
											<div className="mb-2">
												<p className="text-sm font-medium text-[#49525f] px-2 mb-1">{t('label.analytics-doc-parameters')}</p>
												<ul className="space-y-2">
													{fn.parameters.map((param) => (
														<li key={param.name} className="border rounded p-2 bg-white">
															<div className="flex items-center justify-between gap-2">
																<p className="font-medium">{param.name}</p>
																<span className="text-xs px-2 py-1 rounded bg-[#E8EDF5] text-[#395d95] shrink-0">{param.type}</span>
															</div>
															{param.defaultValue && (
																<div className="mt-1 text-sm text-[#6b7280]">
																	{t('label.analytics-doc-default')}: {param.defaultValue}
																</div>
															)}
														</li>
													))}
												</ul>
											</div>
										)}

										<div>
											<p className="text-sm font-medium text-[#49525f] px-2 mb-1">{t('label.analytics-doc-returns')}</p>
											<ul className="space-y-2">
												{fn.returns.map((ret) => (
													<li key={ret.name} className="border rounded p-2 bg-white">
														<div className="flex items-center justify-between gap-2">
															<p className="font-medium">{ret.name}</p>
															<span className="text-xs px-2 py-1 rounded bg-[#E8EDF5] text-[#395d95] shrink-0">{ret.type}</span>
														</div>
													</li>
												))}
											</ul>
										</div>
									</div>
								)}
							</div>
						);
					})}
				</div>
			)}
		</div>
	);
};

export default DocumentationModal;
