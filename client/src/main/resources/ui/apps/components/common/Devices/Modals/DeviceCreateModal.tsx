import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import useCaptcha from '@/hooks/useCaptcha';
import { useTranslation } from 'react-i18next';
import { basicDeviceCreate } from '@/rest/deviceAPI';
import { emptyToNull } from '@/utils/emptyToNullUtils';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUIButton from '@/components/ui/Button/GUIButton';
import LabeledField from '@/components/ui/Form/LabeledField';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { DEVICE_MODELS } from '@/constants/DeviceModels.constant';
import { basicFirmwareDeviceModelsGet } from '@/rest/firmwaresAPI';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';
import { DeviceCreateRequest } from '@/interface/device/deviceCreateRequest.interface';
import { ValidationImeiSchema, ValidationPhoneSchema } from '@/utils/ValidationSchema';

interface DeviceCreateModalProps {
	onSuccess: () => void;
}

const DeviceCreateModal: React.FC<DeviceCreateModalProps> = ({ onSuccess }) => {
	const { t } = useTranslation();
	const captcha = useCaptcha();

	const { data: respFotaDeviceModelsGet, loading: modelsLoading } = useHandleServer(['respFotaDeviceModelsGet'], basicFirmwareDeviceModelsGet);

	const [model, setModel] = useState(DEVICE_MODELS['ADMP50']);
	const {
		register,
		setValue,
		handleSubmit,
		formState: { errors },
	} = useForm<DeviceCreateRequest>({
		mode: 'onChange',
		defaultValues: {
			imei: '',
			name: null,
			phone: null,
			name_organization: null,
			request_configuration_on_connect: false,
			turnstile_token: '',
		},
	});

	const onSubmit = async (data: DeviceCreateRequest) => {
		if (!captcha.validate()) return;

		await basicDeviceCreate({
			...data,
			name: emptyToNull(data.name),
			phone: emptyToNull(data.phone),
			name_organization: emptyToNull(data.name_organization),
			model,
			turnstile_token: captcha.token,
		}).finally(captcha.reset);

		onSuccess();
	};

	if (modelsLoading || !respFotaDeviceModelsGet) {
		return (
			<div className="py-8">
				<Loading />
			</div>
		);
	}

	return (
		<form className="space-y-3">
			<LabeledField label="IMEI*">
				<GUInput type="number" {...register('imei', ValidationImeiSchema<DeviceCreateRequest, 'imei'>(t))} error={errors.imei?.message} />
			</LabeledField>

			<div className="flex flex-col">
				<label className="mb-1 text-sm">{t('label.model')}*</label>
				<GUISelect value={model} onChange={(e) => setModel(e.target.value)}>
					{respFotaDeviceModelsGet?.map((model, index) => (
						<option value={model} key={index}>
							{model}
						</option>
					))}
				</GUISelect>
			</div>

			<LabeledField label="label.device-name">
				<GUInput type="text" {...register('name')} error={errors.name?.message} />
			</LabeledField>

			<LabeledField label="label.phone">
				<GUInput
					type="tel"
					placeholder="+7xxxxxxxxxx"
					{...register('phone', ValidationPhoneSchema<DeviceCreateRequest, 'phone'>(t))}
					onChange={(e) => {
						let value = e.target.value;
						if (value === '+') {
							setValue('phone', '');
							return;
						}
						if (!value) {
							setValue('phone', '');
							return;
						}
						if (!value.startsWith('+')) {
							value = '+' + value;
						}
						setValue('phone', value, { shouldValidate: true });
					}}
					error={errors.phone?.message}
				/>
			</LabeledField>

			<LabeledField label="label.organization">
				<GUInput
					type="text"
					{...register('name_organization', {
						maxLength: {
							value: 100,
							message: t('message.validation-organization-name-length'),
						},
					})}
					error={errors.name_organization?.message}
				/>
			</LabeledField>

			{/* <div className="mt-2 flex items-center gap-3">
                <GUICheckbox checked={watch('request_configuration_on_connect')} onChange={(value) => setValue('request_configuration_on_connect', value)} />
                <p className="text-sm sm:text-base max-w-full sm:max-w-[30rem]">{t('message.request-device-configuration')}</p>
            </div> */}

			{/* Yandex SmartCaptcha */}
			<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} />

			<GUIButton type="submit" onClick={handleSubmit(onSubmit)}>
				{t('label.add')}
			</GUIButton>
		</form>
	);
};

export default DeviceCreateModal;
