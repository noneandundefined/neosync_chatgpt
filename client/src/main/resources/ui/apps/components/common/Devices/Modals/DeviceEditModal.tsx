import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { useCallback, useEffect } from 'react';
import { emptyToNull } from '@/utils/emptyToNullUtils';
import GUIButton from '@/components/ui/Button/GUIButton';
import GUISelect from '@/components/ui/Select/GUISelect';
import { GUInput } from '@/components/ui/Input/GUInput.tsx';
import LabeledField from '@/components/ui/Form/LabeledField';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { DEVICE_MODELS } from '@/constants/DeviceModels.constant';
import { basicFirmwareDeviceModelsGet } from '@/rest/firmwaresAPI';
import { basicDeviceGetByImei, basicDeviceUpdate } from '@/rest/deviceAPI';
import { ValidationImeiSchema, ValidationPhoneSchema } from '@/utils/ValidationSchema';
import { DeviceUpdateRequest } from '@/interface/device/deviceUpdateRequest.interface';

interface DeviceEditModalProps {
	imei: string;
	onSuccess: () => void;
}

const DeviceEditModal: React.FC<DeviceEditModalProps> = ({ imei, onSuccess }) => {
	const { t } = useTranslation();

	const fetchDevice = useCallback(() => basicDeviceGetByImei(imei), [imei]);
	const { data: respDeviceGet, loading: deviceLoading } = useHandleServer(['respDeviceGet', imei], fetchDevice, { staleTime: 0 });

	const { data: respFotaDeviceModelsGet, loading: modelsLoading } = useHandleServer(['respFotaDeviceModelsGet'], basicFirmwareDeviceModelsGet);

	const [model, setModel] = useState<string>('');
	const {
		register,
		setValue,
		reset,
		handleSubmit,
		formState: { errors },
	} = useForm<DeviceUpdateRequest>({
		mode: 'onChange',
		defaultValues: {
			imei: '',
			name: null,
			phone: null,
			name_organization: null,
			request_configuration_on_connect: false,
		},
	});

	useEffect(() => {
		if (!respDeviceGet) return;

		setModel(respDeviceGet.device_model ?? DEVICE_MODELS['ADMP50']);
		reset({
			imei: respDeviceGet.imei,
			model: respDeviceGet.device_model ?? '',
			name: respDeviceGet.name,
			phone: respDeviceGet.phone,
			name_organization: respDeviceGet.name_organization,
			request_configuration_on_connect: respDeviceGet.request_configuration_on_connect,
		});
	}, [respDeviceGet, reset]);

	const onSubmit = async (data: DeviceUpdateRequest) => {
		const updatedFields: Partial<DeviceUpdateRequest> = {};

		if (data.imei !== respDeviceGet?.imei) updatedFields.imei = data.imei;
		if (emptyToNull(data.name) !== emptyToNull(respDeviceGet?.name)) updatedFields.name = emptyToNull(data.name);
		if (emptyToNull(data.phone) !== emptyToNull(respDeviceGet?.phone)) updatedFields.phone = emptyToNull(data.phone);
		if (emptyToNull(data.name_organization) !== emptyToNull(respDeviceGet?.name_organization)) updatedFields.name_organization = emptyToNull(data.name_organization);
		if (model !== respDeviceGet?.device_model) updatedFields.model = model;
		if (data.request_configuration_on_connect !== respDeviceGet?.request_configuration_on_connect) updatedFields.request_configuration_on_connect = data.request_configuration_on_connect;

		if (Object.keys(updatedFields).length === 0) return;

		await basicDeviceUpdate(imei, updatedFields);
		onSuccess();
	};

	if (deviceLoading || modelsLoading || !respDeviceGet) {
		return (
			<div className="py-8">
				<Loading />
			</div>
		);
	}

	return (
		<form className="space-y-3">
			<LabeledField label="IMEI*">
				<GUInput type="number" {...register('imei', ValidationImeiSchema<DeviceUpdateRequest, 'imei'>(t))} error={errors.imei?.message} />
			</LabeledField>

			<div className="flex-1 mb-4">
				<label className="block mb-1 text-sm">{t('label.model')}*</label>
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
					{...register('phone', ValidationPhoneSchema<DeviceUpdateRequest, 'phone'>(t))}
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
                <GUICheckbox checked={watch('request_configuration_on_connect')} onChange={(value: any) => setValue('request_configuration_on_connect', value)} />
                <p className="text-sm sm:text-base max-w-full sm:max-w-[30rem]">{t('message.request-device-configuration')}</p>
            </div> */}

			<GUIButton type="submit" onClick={handleSubmit(onSubmit)}>
				{t('label.save')}
			</GUIButton>
		</form>
	);
};

export default DeviceEditModal;
