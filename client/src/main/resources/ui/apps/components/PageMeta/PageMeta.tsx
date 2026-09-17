import React from 'react';
import { Helmet } from 'react-helmet-async';
import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router-dom';

export const OG_LOCALE_MAP: Record<string, string> = {
	ru: 'ru_RU',
	en: 'en_US',
	es: 'es_ES',
};

export interface PageMetaProps {
	descriptionKey: string;
	ogTitleKey: string;
	ogDescriptionKey: string;
	path?: string;
	title?: React.ReactNode;
	noIndex?: boolean;
}

const PageMeta: React.FC<PageMetaProps> = ({ descriptionKey, ogTitleKey, ogDescriptionKey, path, title, noIndex }) => {
	const location = useLocation();
	const { t, i18n } = useTranslation();

	const description = t(descriptionKey);
	const ogTitle = t(ogTitleKey);
	const ogDescription = t(ogDescriptionKey);
	const canonicalPath = path ?? location.pathname;
	const pageUrl = canonicalPath === '/' ? `https://neosync.neomatica.ru/` : `https://neosync.neomatica.ru/${canonicalPath}`;
	const ogLocale = OG_LOCALE_MAP[i18n.language] ?? 'en_US';

	return (
		<Helmet>
			<html lang={i18n.language} />

			{title && <title>{title}</title>}

			<meta name="description" content={description} />

			{noIndex && <meta name="robots" content="noindex, nofollow" />}

			<meta property="og:type" content="website" />
			<meta property="og:site_name" content="NeoSync" />
			<meta property="og:title" content={ogTitle} />
			<meta property="og:description" content={ogDescription} />
			<meta property="og:image" content="https://neosync.neomatica.ru/local/templates/neomatica/images/neomatica-logo.png" />
			<meta property="og:url" content={pageUrl} />
			<meta property="og:locale" content={ogLocale} />
			<meta name="twitter:card" content="summary" />
			<meta name="twitter:title" content={ogTitle} />
			<meta name="twitter:description" content={ogDescription} />
			<meta name="twitter:image" content="https://neosync.neomatica.ru/local/templates/neomatica/images/neomatica-logo.png" />
			<link rel="canonical" href={pageUrl} />
		</Helmet>
	);
};

export default PageMeta;
