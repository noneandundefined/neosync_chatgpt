import Widget from './Widget';
import { useState } from 'react';
import Editor from '@monaco-editor/react';
import { useTranslation } from 'react-i18next';
import GUIButton from '@/components/ui/Button/GUIButton';
import { buildAnalyticRequest } from '@/utils/BuildAnalyticRequestUtils';
import { AnalyticQueryResponse, basicAnalyticQuery } from '@/rest/analyticAPI';

const formatCell = (value: unknown) => {
	if (value === null || value === undefined) {
		return '—';
	}

	if (typeof value === 'object') {
		return JSON.stringify(value);
	}

	return String(value);
};

const SQLWidget = () => {
	const { t } = useTranslation();

	const [draftQuery, setDraftQuery] = useState('');
	const [result, setResult] = useState<AnalyticQueryResponse | null>(null);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleRun = async () => {
		const request = buildAnalyticRequest('sql', { query: draftQuery });

		if (!request) {
			setError(t('message.fields-not-filled'));
			return;
		}

		setLoading(true);
		setError(null);

		try {
			const response = await basicAnalyticQuery(request);
			setResult(response);
		} catch {
			setResult(null);
			setError(t('message.error-occurred'));
		} finally {
			setLoading(false);
		}
	};

	const columns = result?.columns ?? [];
	const rows = result?.rows ?? [];

	return (
		<Widget title={t('label.analytics-sql-console')} loading={loading} editable={false}>
			<div className="flex h-full flex-col gap-3 min-h-[420px]">
				<div className="flex flex-col gap-3">
					<label className="block text-sm text-[#49525f]">SQL</label>
					<Editor
						height={200}
						value={draftQuery}
						onChange={(val) => setDraftQuery(val ?? '')}
						defaultLanguage="sql"
						theme="vs-light"
						options={{
							minimap: { enabled: false },
							fontSize: 14,
							lineNumbers: 'on',
							lineNumbersMinChars: 3,
							scrollBeyondLastLine: false,
							wordWrap: 'on',
							scrollbar: {
								vertical: 'auto',
								horizontal: 'auto',
							},
						}}
					/>
					<div className="flex justify-end">
						<GUIButton type="button" onClick={handleRun} disabled={loading || !draftQuery.trim()}>
							{t('label.analytics-sql-run')}
						</GUIButton>
					</div>
				</div>

				<div className="flex flex-col flex-1 min-h-0 gap-2">
					<p className="text-sm font-medium text-[#49525f]">{t('label.output')}</p>

					{error && <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{error}</div>}

					{!error && !loading && rows.length === 0 && <div className="rounded-lg border border-[#e4e4e4] bg-[#fafafa] px-3 py-6 text-sm text-[#6b7280] text-center">{t('label.no-data')}</div>}

					{rows.length > 0 && (
						<div className="flex-1 max-h-[500px] overflow-auto rounded-lg border border-[#dbdbdb]">
							<table className="w-full text-sm border-collapse">
								<thead className="sticky top-0 z-[1] bg-white">
									<tr>
										{columns.map((col) => (
											<th key={col} className="px-3 py-2 text-left font-medium text-[#49525f] border-b border-[#dbdbdb] whitespace-nowrap">
												{col}
											</th>
										))}
									</tr>
								</thead>
								<tbody className="">
									{rows.map((row, rowIndex) => (
										<tr key={rowIndex} className={rowIndex % 2 === 0 ? 'bg-white' : 'bg-[#f7f7f7]'}>
											{columns.map((_, colIndex) => (
												<td key={`${rowIndex}-${colIndex}`} className="px-3 py-2 text-[#49525f] border-b border-[#dbdbdb] align-top font-mono text-xs max-w-[280px] break-words">
													{formatCell(row[colIndex])}
												</td>
											))}
										</tr>
									))}
								</tbody>
							</table>
							<div className="sticky bottom-0 bg-[#f7f7f7] border-t border-[#dbdbdb] px-3 py-1.5 text-xs text-[#49525f]">
								{t('label.total')}: {rows.length}
							</div>
						</div>
					)}
				</div>
			</div>
		</Widget>
	);
};

export default SQLWidget;
