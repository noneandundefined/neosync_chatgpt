import React from 'react';
import { useForm } from 'react-hook-form';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ROLES } from '@/constants/Roles.constant';
import { GUInput } from '@/components/ui/Input/GUInput';
import { generatePassword } from '@/utils/PasswordUtil';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import GUIButton from '@/components/ui/Button/GUIButton';
import { LANGUAGES } from '@/constants/Language.constant';
import { useRole } from '@/context/RoleContext/useRoleContext';
import { ACCESS_RULES, BASE_ACCESS_MASK } from '@/constants/Accesses.constant';
import { UserCreateRequest } from '@/interface/auth/userCreateRequest.interface';
import GUISelect from '@/components/ui/Select/GUISelect';
import { basicUserCreate } from '@/rest/userAPI';
import { ValidationEmailSchema, ValidationPasswordRequiredSchema, ValidationPhoneSchema } from '@/utils/ValidationSchema';

interface UserCreateModalProps {
	onSuccess: () => void;
}

const UserCreateModal: React.FC<UserCreateModalProps> = ({ onSuccess }) => {
	const { role } = useRole();
	const { t } = useTranslation();

	const [tab, setTab] = useState<'main' | 'accesses'>('main');
	const [accessesState, setAccessesState] = useState<typeof BASE_ACCESS_MASK>({ ...BASE_ACCESS_MASK });

	const [language, setLanguage] = useState<string>('ru');

	const {
		register,
		handleSubmit,
		setValue,
		formState: { errors },
	} = useForm<UserCreateRequest>({
		mode: 'onChange',
		defaultValues: {
			email: '',
			password: generatePassword(10),
			phone: '',
			locality: '',
			name_organization: '',
			role_code: '',
		},
	});

	useEffect(() => {
		const updatedRole = role == ROLES.SUPERADMIN ? ROLES.SUPPORT : role == ROLES.SUPPORT ? ROLES.DEALER : role == ROLES.DEALER ? ROLES.DEALER_SUPPORT : ROLES.USER;
		setValue('role_code', updatedRole);
	}, [role]);

	const onSubmit = async (data: UserCreateRequest) => {
		const payload = {
			...data,
			language,
			accesses: accessesState,
		};

		await basicUserCreate(payload);

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
							<GUInput type="email" {...register('email', ValidationEmailSchema<UserCreateRequest, 'email'>(t))} error={errors.email?.message} />
						</div>
						<div className="flex-1">
							<label>{t('label.password')}*</label>
							<GUInput type="text" {...register('password', ValidationPasswordRequiredSchema<UserCreateRequest, 'password'>(t))} error={errors.password?.message} />
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
								{...register('phone', ValidationPhoneSchema<UserCreateRequest, 'phone'>(t))}
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
							{t('label.add')}
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
							{t('label.add')}
						</GUIButton>
					</div>
				)}
			</form>
		</div>
	);
};

export default UserCreateModal;
