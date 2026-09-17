import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import GUIButton from '@/components/ui/Button/GUIButton';
import { LANGUAGES } from '@/constants/Language.constant';
import { basicUserGet, basicUserUpdate } from '@/rest/userAPI';
import React, { useCallback, useEffect, useState } from 'react';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { ACCESS_RULES, BASE_ACCESS_MASK } from '@/constants/Accesses.constant';
import { UserUpdateRequest } from '@/interface/user/userUpdateRequest.interface';
import { ValidationEmailSchema, ValidationPasswordRequiredSchema, ValidationPhoneSchema } from '@/utils/ValidationSchema';
import GUISelect from '@/components/ui/Select/GUISelect';

interface UserEditModalProps {
	uuid: string;
	onSuccess: () => void;
}

const UserEditModal: React.FC<UserEditModalProps> = ({ uuid, onSuccess }) => {
	const { t } = useTranslation();

	const fetchUserGet = useCallback(() => basicUserGet(uuid), [uuid]);
	const { data: respUserGet, loading: userLoading } = useHandleServer(['respUserGet', uuid], fetchUserGet, { staleTime: 0 });

	const [tab, setTab] = useState<'main' | 'accesses'>('main');
	const [accessesState, setAccessesState] = useState<typeof BASE_ACCESS_MASK>({ ...BASE_ACCESS_MASK });

	const [language, setLanguage] = useState<string>('');

	const {
		register,
		setValue,
		handleSubmit,
		reset,
		formState: { errors },
	} = useForm<UserUpdateRequest>({
		mode: 'onChange',
		defaultValues: {
			email: '',
			password: '',
			phone: '',
			name_organization: '',
			locality: '',
		},
	});

	useEffect(() => {
		if (!respUserGet) return;

		setLanguage(respUserGet.user_contact.language);

		reset({
			email: respUserGet.email,
			password: respUserGet.password,
			phone: respUserGet.user_contact.phone,
			name_organization: respUserGet.user_contact.name_organization,
			locality: respUserGet.user_contact.locality,
			language: respUserGet.user_contact.language,
		});

		setAccessesState({
			access_treker_create: respUserGet.access_treker_create ?? false,
			access_treker_edit: respUserGet.access_treker_edit ?? false,
			access_treker_delete: respUserGet.access_treker_delete ?? false,
			access_group_manage: respUserGet.access_group_manage ?? false,
			access_configuration_read: respUserGet.access_configuration_read ?? false,
			access_configuration_apply: respUserGet.access_configuration_apply ?? false,
			access_configuration_history: respUserGet.access_configuration_history ?? false,
			access_command_send: respUserGet.access_command_send ?? false,
			access_log_read: respUserGet.access_log_read ?? false,
		});
	}, [respUserGet, reset]);

	const onSubmit = async (data: UserUpdateRequest) => {
		const updatedFields: Partial<UserUpdateRequest> = {};

		if (data.email !== respUserGet?.email) updatedFields.email = data.email;
		if (data.password !== respUserGet?.password) updatedFields.password = data.password;
		if (data.phone !== respUserGet?.user_contact.phone) updatedFields.phone = data.phone;
		if (data.name_organization !== respUserGet?.user_contact.name_organization) updatedFields.name_organization = data.name_organization;
		if (data.locality !== respUserGet?.user_contact.locality) updatedFields.locality = data.locality;
		if (language !== respUserGet?.user_contact.language) updatedFields.language = language;

		updatedFields.accesses = accessesState;

		await basicUserUpdate(uuid, updatedFields);
		onSuccess();
	};

	// CHECKBOX
	const grouped = ACCESS_RULES.reduce<Record<string, typeof ACCESS_RULES>>((acc, rule) => {
		if (!acc[rule.group]) acc[rule.group] = [];
		acc[rule.group].push(rule);
		return acc;
	}, {});

	const toggleOne = (key: keyof typeof BASE_ACCESS_MASK, value: boolean) => {
		setAccessesState((prev) => ({ ...prev, [key]: value }));
	};

	const toggleGroup = (group: string, value: boolean) => {
		const updated = { ...accessesState };
		grouped[group].forEach((r) => {
			updated[r.key as keyof typeof BASE_ACCESS_MASK] = value;
		});

		setAccessesState(updated);
	};

	const isGroupChecked = (group: string) => grouped[group].every((r) => !!accessesState[r.key as keyof typeof BASE_ACCESS_MASK]);

	if (userLoading || !respUserGet) {
		return (
			<div className="py-8">
				<Loading />
			</div>
		);
	}

	return (
		<div className="w-full sm:max-w-[600px] mx-auto mt-4 rounded">
			<div className="flex border-b mb-5">
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'main' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('main')}>
					{t('label.basic')}
				</p>
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'accesses' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('accesses')}>
					{t('message.access-rights')}
				</p>
			</div>

			<form className="flex flex-col space-y-3">
				{tab === 'main' && (
					<React.Fragment>
						<div className="flex-1">
							<label className="block text-sm mb-1">{t('label.login')}*</label>
							<GUInput type="email" {...register('email', ValidationEmailSchema<UserUpdateRequest, 'email'>(t))} error={errors.email?.message} />
						</div>
						<div className="flex-1">
							<label>{t('label.password')}</label>
							<GUInput type="text" {...register('password', ValidationPasswordRequiredSchema<UserUpdateRequest, 'password'>(t))} error={errors.password?.message} />
						</div>
						<div className="flex-1">
							<label className="block text-sm mb-1">{t('label.language')}</label>
							<GUISelect value={language} onChange={(e) => setLanguage(e.target.value)}>
								{LANGUAGES.map((lang, index) => (
									<option value={lang.id} key={index}>
										<p>{lang.value}</p>
									</option>
								))}
							</GUISelect>
						</div>
						<div className="flex-1">
							<label className="block text-sm mb-1">{t('label.name-organization')}</label>
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
						</div>
						<div className="flex-1">
							<label>{t('label.locality')}</label>
							<GUInput type="text" {...register('locality')} error={errors.locality?.message} />
						</div>
						<div className="flex-1">
							<label>{t('label.phone')}</label>
							<GUInput
								type="tel"
								placeholder="+7xxxxxxxxxx"
								{...register('phone', ValidationPhoneSchema<UserUpdateRequest, 'phone'>(t))}
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
						</div>

						<GUIButton type="submit" onClick={handleSubmit(onSubmit)}>
							{t('label.save')}
						</GUIButton>
					</React.Fragment>
				)}

				{tab === 'accesses' && (
					<div className="space-y-4">
						{Object.entries(grouped).map(([group, rules]) => (
							<div key={group}>
								<div className="flex items-center gap-1">
									<GUICheckbox checked={isGroupChecked(group)} onChange={(v) => toggleGroup(group, v)} />

									<p>{group === 'main' ? t('label.mains') : group === 'configuration' ? t('label.configuration') : group === 'command' ? t('label.commands') : t('label.other')}</p>
								</div>

								<div className="ml-6 space-y-2 mt-2">
									{rules.map((rule) => (
										<div key={rule.key} className="flex items-center gap-1">
											<GUICheckbox checked={!!accessesState[rule.key as keyof typeof BASE_ACCESS_MASK]} onChange={(v) => toggleOne(rule.key as keyof typeof BASE_ACCESS_MASK, v)} />
											<p>{t(rule.label)}</p>
										</div>
									))}
								</div>
							</div>
						))}

						<GUIButton type="submit" onClick={handleSubmit(onSubmit)}>
							{t('label.save')}
						</GUIButton>
					</div>
				)}
			</form>
		</div>
	);
};

export default UserEditModal;
