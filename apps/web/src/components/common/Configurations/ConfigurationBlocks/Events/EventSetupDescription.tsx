import { useTranslation } from 'react-i18next';

const EventSetupDescription = () => {
	const { t } = useTranslation();

	return <p className="text-sm text-justify">{t('message.alarm-sms-description')}</p>;
};

export default EventSetupDescription;
