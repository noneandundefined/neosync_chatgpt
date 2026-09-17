import useIsMobile from '@/hooks/useIsMobile';
import { useTranslation } from 'react-i18next';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useCompactTableContext } from './CompactTableContext';

function CompactTableBody<T>() {
	const { t } = useTranslation();
	const isMobile = useIsMobile();

	const { pageData, getRowId, checkedItems, toggleOne, hoveredRow, setHoveredRow, onRowClick, onTitleClick, renderTitle, renderSubtitle, renderLabel, renderMeta, renderTime, isUnread, selectionMode } =
		useCompactTableContext<T>();

	if (pageData.length === 0) {
		return (
			<div className="px-4 py-10 text-center">
				<p className="text-[13px] text-[#656d76]">{t('message.data-available')}</p>
			</div>
		);
	}

	return (
		<div>
			{pageData.map((item, index) => {
				const id = getRowId(item);
				const isSelected = !!checkedItems[id];
				const isHovered = hoveredRow === index;
				const unread = isUnread?.(item) ?? false;

				return (
					<div
						key={id}
						className={`flex items-start gap-3 px-4 py-4 border-b border-[#e9e9e9] last:border-b-0 transition-colors ${onRowClick || selectionMode === 'single' ? 'cursor-pointer' : ''} ${isSelected && selectionMode === 'single' ? 'bg-[#f6f6f6]' : isHovered ? 'bg-[#f6f6f6]' : 'bg-white'}`}
						onClick={() => {
							if (selectionMode === 'single') {
								toggleOne(id);
							}

							onRowClick?.(item, index);
						}}
						onMouseEnter={!isMobile ? () => setHoveredRow(index) : undefined}
						onMouseLeave={!isMobile ? () => setHoveredRow(null) : undefined}
					>
						<div className="flex items-center gap-2 pt-1 shrink-0">
							{unread ? <span className="h-2 w-2 rounded-full bg-[#0969da] shrink-0" /> : <span className="h-2 w-2 shrink-0" />}

							<div onClick={(e) => e.stopPropagation()}>
								<GUICheckbox checked={isSelected} onChange={() => toggleOne(id)} />
							</div>
						</div>

						<div className="flex-1 min-w-0">
							<div
								className={`font-medium text-[#49525f] ${onTitleClick ? 'cursor-pointer hover:underline' : ''}`}
								onClick={
									onTitleClick
										? (e) => {
												e.stopPropagation();
												onTitleClick(item, index);
											}
										: undefined
								}
							>
								{renderTitle(item, index)}
							</div>
							{renderSubtitle && <div className="mt-0.5 text-sm text-[#656d76] line-clamp-2">{renderSubtitle(item, index)}</div>}
						</div>

						{renderLabel && <div className="hidden sm:flex items-center shrink-0 text-sm text-[#656d76] px-2">{renderLabel(item, index)}</div>}

						<div className="flex items-center gap-3 shrink-0">
							{renderMeta?.(item, index)}
							{renderTime && <div className="text-sm text-[#656d76] whitespace-nowrap min-w-[5rem] text-right">{renderTime(item, index)}</div>}
						</div>
					</div>
				);
			})}
		</div>
	);
}

export default CompactTableBody;
