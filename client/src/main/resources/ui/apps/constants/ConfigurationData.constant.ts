import i18next from 'i18next';

export const CONFIGURATION_DATA = {
	GET: {
		TITLE: 'configuration-get-title',
	},
	IMPORT: {
		TITLE: 'configuration-import-title',
	},
	SET: {
		TITLE: 'configuration-set-title',
		DESCRIPTION: [
			i18next.t('message.connect-to-server'),
			i18next.t('message.configuration-set-description-2'),
			i18next.t('message.configuration-set-description-3'),
			i18next.t('message.configuration-set-description-4'),
			i18next.t('message.configuration-set-description-5'),
		],
	},
	REBOOT: {
		TITLE: 'device-reboot-title',
		DESCRIPTION: [i18next.t('message.device-reboot-description')],
	},
};
