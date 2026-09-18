export const config = {
	/* Base URLs */
	links: {
		URL_BACKEND_DEV: 'http://localhost:8080/api/v1',
		URL_FRONTEND_DEV: 'http://localhost:5173',
		URL_BACKEND_PROD: 'https://neosync.neomatica.ru/api/v1',
		URL_FRONTEND_PROD: 'https://neosync.neomatica.ru',
		URL_BACKEND_PREPROD: 'http://77.222.42.127/api/v1',
		URL_FRONTEND_PREPROD: 'http://77.222.42.127',
	},
	/* Yandex SmartCaptcha client keys */
	yandex: {
		SMARTCAPTCHA_SITE_KEY_DEV: '',
		SMARTCAPTCHA_SITE_KEY_PROD: 'ysc1_uMcAbp9sZiuOgf8iflcrbJf2PDEMgKTKFGGk5XEce95acf3c',
		SMARTCAPTCHA_SITE_KEY_PREPROD: 'ysc1_uMcAbp9sZiuOgf8iflcrbJf2PDEMgKTKFGGk5XEce95acf3c',
	},
	type: {
		release: 'dev', // prod | preprod
	},
};
