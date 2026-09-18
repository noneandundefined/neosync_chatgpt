import { useTranslation } from 'react-i18next';

export const groupColumns = () => {
	const { t } = useTranslation();

	return [
		{ title: t('label.title'), key: 'name', sortable: true },
		{ title: t('label.description'), key: 'description', sortable: true },
		{ title: t('label.objects'), key: 'objects', sortable: true },
	];
};
