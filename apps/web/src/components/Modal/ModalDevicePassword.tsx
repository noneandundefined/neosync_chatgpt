import Tooltip from '../ui/Tooltip';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import GUIButton from '../ui/Button/GUIButton';
import LabeledField from '../ui/Form/LabeledField';
import InputPassword from '../ui/Form/InputPassword';
import { Device_DeviceConf_Sync } from '@/rest/deviceAPI';
import { basicConfigurationChangePass } from '@/rest/configurationAPI';
import { ConfigurationChangePassRequest } from '@/interface/configuration/configurationChangePassRequest.interface';

interface ModalDevicePasswordProps {
	device: Device_DeviceConf_Sync;
}

const ModalDevicePassword: React.FC<ModalDevicePasswordProps> = ({ device }) => {
	const { t } = useTranslation();

	const {
		register,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<ConfigurationChangePassRequest>({
		mode: 'onChange',
		defaultValues: {
			old_password: '',
			new_password: '',
		},
	});

	const onSubmit = async (data: ConfigurationChangePassRequest) => {
		await basicConfigurationChangePass(device.imei, data);
	};

	return (
		<form className="space-y-3">
			<p className="text-sm sm:text-base max-w-full sm:max-w-[40rem]">{t('message.change-password-desc')}</p>

			<LabeledField label="label.old-password">
				<InputPassword
					{...register('old_password', {
						required: t('message.validation-required-field'),
					})}
					error={errors.old_password?.message}
				/>
			</LabeledField>

			<LabeledField label="label.new-password">
				<div className="relative">
					<Tooltip title={t('message.requirements-pass')} position="bottom">
						<InputPassword
							{...register('new_password', {
								required: t('message.validation-required-field'),
							})}
							error={errors.new_password?.message}
						/>
					</Tooltip>
				</div>
			</LabeledField>

			<GUIButton type="submit" disabled={!watch('old_password') || !watch('new_password')} onClick={handleSubmit(onSubmit)}>
				{t('label.change-password')}
			</GUIButton>
		</form>
	);
};

export default ModalDevicePassword;
