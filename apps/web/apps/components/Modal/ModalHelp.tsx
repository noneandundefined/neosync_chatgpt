import { SUPPORT_EMAIL_EN, SUPPORT_EMAIL_RU, SUPPORT_PHONE_NUMBER, SUPPORT_SKYPE } from '@/constants/Support.constant';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicFirmwareDeviceModelsGet } from '@/rest/firmwaresAPI';
import { useState } from 'react';
import GUISelect from '../ui/Select/GUISelect';

const ModalHelp = () => {
	const { t, i18n } = useTranslation();

	const [model, setModel] = useState<string>('');

	const { data: basicFotaDeviceModelsGetResp } = useHandleServer(['basicFotaDeviceModelsGetResp'], basicFirmwareDeviceModelsGet);

	const handleOpenManual = () => {
		if (!model) return;

		const lang = i18n.language === 'ru' ? 'manual_ru' : 'manual_en';

		const url = `https://www.neomatica.com/upload/${lang}/${model}.pdf`;
		if (url) window.open(url);
	};

	return (
		<div className="space-y-4">
			<div className="space-y-3">
				<p className="font-medium">{t('label.manuals')}</p>
				<p className="text-[14px]">{t('message.manuals-info')}</p>

				<div className="flex flex-col sm:flex-row sm:items-center gap-2">
					<p className="text-[14px]">{t('message.manuals-select-device')}</p>

					<GUISelect value={model} onChange={(e) => setModel(e.target.value)}>
						{basicFotaDeviceModelsGetResp?.map((modelSource, index) => (
							<option key={index} value={modelSource}>
								{modelSource}
							</option>
						))}
					</GUISelect>
				</div>

				<GUIButton onClick={handleOpenManual}>{t('message.manuals-open-button')}</GUIButton>

				<div className="space-y-2">
					<p className="font-medium">{t('message.support-contacts')}</p>
					<p className="text-[14px]">
						{t('label.phone')}: {SUPPORT_PHONE_NUMBER}
					</p>
					<p className="text-[14px]">
						{t('label.email')}: {i18n.language === 'ru' ? SUPPORT_EMAIL_RU : SUPPORT_EMAIL_EN}
					</p>
					<p className="text-[14px]">
						{t('label.skype')}: {SUPPORT_SKYPE}
					</p>
				</div>
			</div>
		</div>
	);
};

export default ModalHelp;
