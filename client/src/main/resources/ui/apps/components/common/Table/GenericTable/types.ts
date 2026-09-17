export interface Column<T> {
	/** Заголовок колонки (то, что будет в <th>) */
	title: string;
	/** Ключ поля в объекте данных */
	key: keyof T | string;
	/** Функция для кастомного рендера содержимого ячейки (например для иконок) */
	render?: (item: T, index: number) => React.ReactNode;
	/** Состояние видимости колонки */
	visible?: boolean;
	/** Ширина колонки */
	width?: string | number;
	/** Базовая сортировка колонки */
	sortable?: boolean;
}
