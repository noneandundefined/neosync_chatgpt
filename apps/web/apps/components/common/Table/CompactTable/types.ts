import { FetchFn } from '../GenericTable/GenericTable';
import type { ReactNode } from 'react';

export interface CompactTableProps<T> {
	/** Идентификатор таблицы */
	tableKey: string;
	/** Запрос к серверу с пагинацией */
	fetchFn: FetchFn<T>;
	/** Токен для обновления данных */
	sysToken: number;
	/** Уникальный идентификатор строки */
	getRowId: (item: T) => string | number;
	/** Основная строка (жирный заголовок) */
	renderTitle: (item: T, index: number) => ReactNode;
	/** Дополнительное описание под заголовком */
	renderSubtitle?: (item: T, index: number) => ReactNode;
	/** Метка/категория справа от текста */
	renderLabel?: (item: T, index: number) => ReactNode;
	/** Аватар или иконка справа */
	renderMeta?: (item: T, index: number) => ReactNode;
	/** Время/дата справа */
	renderTime?: (item: T, index: number) => ReactNode;
	/** Непрочитанный элемент — синяя точка слева */
	isUnread?: (item: T) => boolean;
	/** Клик по строке */
	onRowClick?: (item: T, index: number) => void;
	/** Клик по заголовку строки */
	onTitleClick?: (item: T, index: number) => void;
	/** Массовое удаление выбранных */
	onDeleteFn?: (ids: (string | number)[]) => Promise<void>;
	/** Дополнительные bulk-действия при выборе строк */
	bulkActions?: ReactNode;
	/** Пустое состояние */
	emptyComponent?: ReactNode;
	/** Нижняя подсказка слева от пагинации */
	footer?: ReactNode;
	/** Поиск над таблицей */
	enableSearch?: boolean;
	/** Режим выбора строк */
	selectionMode?: 'multiple' | 'single';
	/** Контролируемый выбор строк */
	controlledSelectedIds?: (string | number)[];
	/** Изменение выбранных строк */
	onSelectionChange?: (ids: (string | number)[]) => void;
	/** Синхронизировать поиск с URL */
	persistSearchInUrl?: boolean;
}

export const compactTableQueryKey = (tableKey: string) => [`compact-${tableKey}`] as const;
