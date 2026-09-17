import React, { useEffect } from 'react';
import Modal from '@/components/Modal/Modal';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import Pencil from '@/components/@icons/pencil';
import { ROLES } from '@/constants/Roles.constant';
import { Loading } from '@/components/common/Loader/Loading';
import { useModalContext } from '@/context/useModalContext';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { useRole } from '@/context/RoleContext/useRoleContext';
import ModalAnalyticEdit from '@/components/Modal/ModalAnalyticEdit';

interface WidgetProps {
	title?: string;
	editableTitle?: string;
	className?: string;
	loading?: boolean;
	children: React.ReactNode;
	widgetId?: number;
	editableQuery?: string;
	onApplyQuery?: (nextQuery: string) => void;
	onApplyTitle?: (nextTitle: string) => void;
	editable?: boolean;
}

type CachedWidgetValue = {
	title: string;
	query: string;
};

type CachedWidgetMap = Record<string, CachedWidgetValue>;

const normalizeQuery = (value: string) => {
	return (value ?? '')
		.replace(/\\r\\n/g, '\n')
		.replace(/\\n/g, '\n')
		.replace(/\\t/g, '\t')
		.trim();
};

const Widget: React.FC<WidgetProps> = ({ title, editableTitle, className = '', loading = false, children, widgetId, editableQuery, onApplyQuery, onApplyTitle, editable = true }) => {
	const { role } = useRole();
	const { t } = useTranslation();
	const { open } = useModalContext();

	const displayTitle = title ?? '';
	const modalTitle = editableTitle ?? displayTitle;
	const queryValue = normalizeQuery(editableQuery ?? '');

	const canUseLocalAnalyticsQuery = role === ROLES.SUPERADMIN || role === ROLES.SUPPORT;

	const canEdit = editable && Boolean(widgetId) && Boolean(onApplyQuery) && canUseLocalAnalyticsQuery;

	useEffect(() => {
		if (!widgetId || !onApplyQuery || !canUseLocalAnalyticsQuery) return;

		const raw = localStorage.getItem(CACHEKEYs.ANALYTIC_WIDGETS);
		if (!raw) return;

		try {
			const parsed = JSON.parse(raw) as CachedWidgetMap;

			const saved = parsed[String(widgetId)];
			if (!saved) return;

			if (saved.query) onApplyQuery(normalizeQuery(saved.query));
			if (saved.title && onApplyTitle) onApplyTitle(saved.title);
		} catch {
			// Ignore an invalid local cache and use the server widget settings.
		}
	}, [widgetId, onApplyQuery, onApplyTitle, canUseLocalAnalyticsQuery]);

	return (
		<section className={`bg-[#f7f7f7] rounded-xl p-2 sm:p-3 h-full border min-w-0 ${className}`}>
			<React.Fragment>
				<div className="mb-3 flex items-start justify-between gap-3">
					<h3 className="text-base font-medium text-gray-800 min-w-0 break-words">{title && displayTitle}</h3>
					{canEdit && (
						<Tooltip title={t('message.editing-widget')}>
							<button
								type="button"
								onClick={() => {
									open(
										<Modal title={t('message.editing-widget')} width="780px">
											<ModalAnalyticEdit
												title={modalTitle}
												query={queryValue}
												onSave={(nextTitle, nextQuery) => {
													const cleanQuery = normalizeQuery(nextQuery);

													if (onApplyTitle) onApplyTitle(nextTitle);
													if (onApplyQuery) onApplyQuery(cleanQuery);

													if (!canUseLocalAnalyticsQuery) return;

													const raw = localStorage.getItem(CACHEKEYs.ANALYTIC_WIDGETS);
													let parsed: CachedWidgetMap = {};
													if (raw) {
														try {
															parsed = JSON.parse(raw) as CachedWidgetMap;
														} catch {
															parsed = {};
														}
													}

													parsed[String(widgetId)] = {
														title: nextTitle,
														query: cleanQuery,
													};

													localStorage.setItem(CACHEKEYs.ANALYTIC_WIDGETS, JSON.stringify(parsed));
												}}
											/>
										</Modal>
									);
								}}
							>
								<Pencil fill="#6b7280" size={15} />
							</button>
						</Tooltip>
					)}
				</div>

				{loading ? (
					<div className="flex justify-center items-center mb-3">
						<Loading />
					</div>
				) : (
					<div className="bg-white p-2 sm:p-3 rounded-xl min-w-0" style={{ boxShadow: '0 0 7px rgba(0, 0, 0, 0.08)' }}>
						{children}
					</div>
				)}
			</React.Fragment>
		</section>
	);
};

export default Widget;
