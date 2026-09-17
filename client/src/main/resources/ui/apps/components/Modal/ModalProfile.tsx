import { useForm } from 'react-hook-form';
import { useEffect, useState } from 'react';
import { GUInput } from '../ui/Input/GUInput';
import Tooltip from '@/components/ui/Tooltip';
import GUIButton from '../ui/Button/GUIButton';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import { hasPermission } from '@/constants/Roles.constant';
import { Loading } from '@/components/common/Loader/Loading';
import { useRole } from '@/context/RoleContext/useRoleContext';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { ACCESS_RULES, BASE_ACCESS_MASK } from '@/constants/Accesses.constant';
import { UserUpdateRequest } from '@/interface/user/userUpdateRequest.interface';
import { ValidationEmailSchema, ValidationPasswordSchema, ValidationPhoneSchema } from '@/utils/ValidationSchema';
import { basicUserGetMe, basicUserUpdateCvcg, basicUserUpdateMe, basicUserUpdateConfigurationPriority } from '@/rest/userAPI';

interface ModalProfileProps {
	onClose: () => void;
}

const ModalProfile: React.FC<ModalProfileProps> = ({ onClose }) => {
	const { role } = useRole();
	const { t } = useTranslation();

	const queryClient = useQueryClient();

	const { data: respUserGetMe, loading: userLoading } = useHandleServer(['respUserGetMe'], basicUserGetMe, { staleTime: 0 });

	const [tab, setTab] = useState<'main' | 'accesses'>('main');
	const [cvcg, setCvcg] = useState<boolean>(false);
	const [configurationPriority, setConfigurationPriority] = useState(false);
	const [prioritySaving, setPrioritySaving] = useState(false);

	const [accessesState, setAccessesState] = useState<typeof BASE_ACCESS_MASK>({ ...BASE_ACCESS_MASK });
	const {
		register,
		handleSubmit,
		setValue,
		formState: { errors },
	} = useForm<UserUpdateRequest>({
		mode: 'onChange',
		defaultValues: {
			email: '',
			password: '',
			phone: '',
		},
	});

	useEffect(() => {
		if (respUserGetMe) {
			setCvcg(respUserGetMe.can_view_child_groups ?? false);
			setConfigurationPriority(respUserGetMe.force_neosync_configuration_priority ?? false);
		}
	}, [respUserGetMe]);

	useEffect(() => {
		if (respUserGetMe) {
			let phone = respUserGetMe.user_contact.phone || '';
			if (phone && !phone.startsWith('+')) phone = '+' + phone;

			setValue('email', respUserGetMe.email, { shouldValidate: true });
			setValue('phone', phone, { shouldValidate: true });
			setValue('password', '', { shouldValidate: true });

			setAccessesState({
				access_treker_create: respUserGetMe.access_treker_create ?? false,
				access_treker_edit: respUserGetMe.access_treker_edit ?? false,
				access_treker_delete: respUserGetMe.access_treker_delete ?? false,
				access_group_manage: respUserGetMe.access_group_manage ?? false,
				access_configuration_read: respUserGetMe.access_configuration_read ?? false,
				access_configuration_apply: respUserGetMe.access_configuration_apply ?? false,
				access_configuration_history: respUserGetMe.access_configuration_history ?? false,
				access_command_send: respUserGetMe.access_command_send ?? false,
				access_log_read: respUserGetMe.access_log_read ?? false,
			});
		}
	}, [respUserGetMe]);

	const onSubmit = async (data: UserUpdateRequest) => {
		if (!respUserGetMe) return;

		const updatedFields: Partial<UserUpdateRequest> = {};

		if (data.email !== respUserGetMe.email) updatedFields.email = data.email;
		if (data.phone !== respUserGetMe.user_contact.phone) updatedFields.phone = data.phone;
		if (data.password) updatedFields.password = data.password;

		Object.keys(accessesState).forEach((key) => {
			const k = key as keyof typeof BASE_ACCESS_MASK;
			if (accessesState[k] !== respUserGetMe[k]) {
				if (!updatedFields.accesses) updatedFields.accesses = { ...accessesState };
			}
		});

		if (Object.keys(updatedFields).length === 0) {
			onClose();
			return;
		}

		await basicUserUpdateMe(updatedFields);
		await queryClient.invalidateQueries({ queryKey: ['respUserGetMe'] });

		onClose();
	};

	const saveConfigurationPriority = async (enabled: boolean) => {
		if (prioritySaving || !respUserGetMe) return;

		const previous = configurationPriority;

		setPrioritySaving(true);
		setConfigurationPriority(enabled);

		try {
			await queryClient.cancelQueries({ queryKey: ['respUserGetMe'] });
			await basicUserUpdateConfigurationPriority(enabled);

			setConfigurationPriority(enabled);

			void queryClient.invalidateQueries({ queryKey: ['respUserGetMe'], refetchType: 'none' });
		} catch {
			setConfigurationPriority(previous);
		} finally {
			setPrioritySaving(false);
		}
	};

	const grouped = ACCESS_RULES.reduce<Record<string, typeof ACCESS_RULES>>((acc, rule) => {
		if (!acc[rule.group]) acc[rule.group] = [];
		acc[rule.group].push(rule);
		return acc;
	}, {});

	const isGroupChecked = (group: string) => grouped[group].every((r) => !!accessesState[r.key as keyof typeof BASE_ACCESS_MASK]);

	if (userLoading || !respUserGetMe) {
		return (
			<div className="py-8">
				<Loading />
			</div>
		);
	}

	return (
		<div className="w-full max-w-[95%] sm:max-w-[600px] mx-auto mt-4 rounded">
			<div className="flex border-b mb-5">
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'main' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('main')}>
					{t('label.basic')}
				</p>
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'accesses' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('accesses')}>
					{t('message.access-rights')}
				</p>
			</div>

			{tab === 'main' && (
				<form className="flex flex-col space-y-3">
					<div className="flex-1">
						<label className="block text-sm mb-1">{t('label.login')}*</label>
						<GUInput type="email" {...register('email', ValidationEmailSchema<UserUpdateRequest, 'email'>(t))} error={errors.email?.message} />
					</div>
					{role !== null && hasPermission(role, 'user:change_password') && (
						<div className="flex-1">
							<label>{t('label.password')}</label>
							<GUInput type="text" {...register('password', ValidationPasswordSchema<UserUpdateRequest, 'password'>(t))} error={errors.password?.message} />
						</div>
					)}
					<div className="flex-1">
						<label>{t('label.phone')}</label>
						<GUInput type="tel" {...register('phone', ValidationPhoneSchema<UserUpdateRequest, 'phone'>(t))} placeholder="+7xxxxxxxxxx" error={errors.phone?.message} />
					</div>

					{role !== null && hasPermission(role, 'view:groups:child') && (
						<GUICheckbox
							label={t('message.show-groups-child-accounts')}
							checked={cvcg}
							onChange={async (checked) => {
								setCvcg(checked);
								await basicUserUpdateCvcg(checked);

								await queryClient.invalidateQueries({ queryKey: ['respUserGetMe'] });
							}}
						/>
					)}

					<fieldset disabled={prioritySaving} className="border-0 p-0 my-1" aria-busy={prioritySaving}>
						<div className="flex items-center justify-between gap-3">
							<p className='text-[#666] text-[14px]'>{t('message.force-neosync-configuration-priority')}</p>

							<Tooltip title={t('message.force-neosync-configuration-priority-description')}>
								<GUISwitch
									checked={configurationPriority}
									onChange={(event) => saveConfigurationPriority(event.target.checked)}
								/>
							</Tooltip>
						</div>
					</fieldset>

					<GUIButton type="submit" onClick={handleSubmit(onSubmit)}>
						{t('label.save')}
					</GUIButton>
				</form>
			)}

			{tab === 'accesses' && (
				<div className="p-4 space-y-4">
					{Object.entries(grouped).map(([group, rules]) => (
						<div key={group}>
							<div className="flex items-center gap-1">
								<GUICheckbox checked={isGroupChecked(group)} />

								<p>{group === 'main' ? t('label.mains') : group === 'configuration' ? t('label.configuration') : group === 'command' ? t('label.commands') : t('label.other')}</p>
							</div>

							<div className="ml-6 space-y-2 mt-2">
								{rules.map((rule) => (
									<div key={rule.key} className="flex items-center gap-1">
										<GUICheckbox checked={!!accessesState?.[rule.key as keyof typeof BASE_ACCESS_MASK]} />
										<p>{t(rule.label)}</p>
									</div>
								))}
							</div>
						</div>
					))}
				</div>
			)}
		</div>
	);
};

export default ModalProfile;
