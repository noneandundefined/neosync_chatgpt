import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import es from '@/locale/languages/es-es.json';
import en from '@/locale/languages/en-us.json';
import ru from '@/locale/languages/ru-ru.json';

import { LANGs } from '@/constants/Language.constant';

const resources = {
	es: {
		translation: es,
	},
	en: {
		translation: en,
	},
	ru: {
		translation: ru,
	},
};

const savedLang = localStorage.getItem('lang');
const browserLang = navigator.language.slice(0, 2);
const defaultLang = 'en';

i18n.use(initReactI18next).init({
	resources,
	lng: savedLang || (LANGs.includes(browserLang) ? browserLang : defaultLang),
	fallbackLng: 'en',
	interpolation: {
		escapeValue: false,
	},
});

export default i18n;
