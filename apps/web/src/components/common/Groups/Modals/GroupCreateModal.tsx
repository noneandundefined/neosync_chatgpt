import { useState } from 'react';
import { toast } from 'react-toastify';
import { useForm } from 'react-hook-form';
import useCaptcha from '@/hooks/useCaptcha';
import { useTranslation } from 'react-i18next';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { GUITextarea } from '@/components/ui/Textarea/GUITextarea';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';
import { basicGroupCreate, groupRightsStateToApiPayload } from '@/rest/groupAPI';
import { GroupCreateRequest } from '@/interface/group/groupCrerateRequest.interface';
import { BASE_GROUP_RIGHTs_MASK, GROUP_RIGHTs_RULES } from '@/constants/Accesses.constant';

interface GroupCreateModalProps {
	onSuccess: () => void;
}

type MainForm = Pick<GroupCreateRequest, 'name' | 'description' | 'turnstile_token'>;

const GroupCreateModal: React.FC<GroupCreateModalProps> = ({ onSuccess }) => {
	const { t } = useTranslation();
	const captcha = useCaptcha();

	const [tab, setTab] = useState<'main' | 'objects' | 'right'>('main');
	const [rightsState, setRightsState] = useState<typeof BASE_GROUP_RIGHTs_MASK>({ ...BASE_GROUP_RIGHTs_MASK });

	const {
		register,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<MainForm>({
		mode: 'onChange',
		defaultValues: {
			name: '',
			description: '',
			turnstile_token: '',
		},
	});

	const submitCreate = async (data: MainForm) => {
		if (!captcha.validate()) {
			toast.error(t('message.validation-required-field'));
			return;
		}

		const payload: GroupCreateRequest = {
			...data,
			turnstile_token: captcha.token,
			...groupRightsStateToApiPayload(rightsState),
		};

		await basicGroupCreate(payload).finally(captcha.reset);

		onSuccess();
	};

	const grouped = GROUP_RIGHTs_RULES.reduce<Record<string, typeof GROUP_RIGHTs_RULES>>((acc, rule) => {
		if (!acc[rule.group]) acc[rule.group] = [];
		acc[rule.group].push(rule);
		return acc;
	}, {});

	const toggleOne = (key: keyof typeof BASE_GROUP_RIGHTs_MASK, value: boolean) => {
		setRightsState((prev) => ({ ...prev, [key]: value }));
	};

	const toggleGroup = (group: string, value: boolean) => {
		const updated = { ...rightsState };
		grouped[group].forEach((r) => {
			updated[r.key as keyof typeof BASE_GROUP_RIGHTs_MASK] = value;
		});

		setRightsState(updated);
	};

	const isGroupChecked = (group: string) => grouped[group].every((r) => !!rightsState[r.key as keyof typeof BASE_GROUP_RIGHTs_MASK]);

	return (
		<div className="w-full sm:max-w-[600px] mx-auto mt-4 rounded">
			<div className="flex border-b">
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'main' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('main')}>
					{t('label.basic')}
				</p>
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'objects' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('objects')}>
					{t('label.selecting-objects')}
				</p>
				<p className={`flex-1 text-center py-2 text-sm sm:text-base bg-transparent cursor-pointer ${tab === 'right' ? 'border-b-2 border-[#afb5c0]' : ''}`} onClick={() => setTab('right')}>
					{t('message.access-rights')}
				</p>
			</div>

			{tab === 'main' && (
				<div className="py-2 space-y-4">
					<div>
						<label className="block text-sm">{t('label.title')}*</label>
						<GUInput
							type="text"
							{...register('name', {
								required: t('message.validation-required-field'),
								maxLength: {
									value: 30,
									message: t('message.validation-group-name-length'),
								},
							})}
							error={errors.name?.message}
						/>
						<div className="text-right text-xs">{watch('name')?.length} / 30</div>
					</div>
					<div>
						<label className="block text-sm">{t('label.description')}</label>
						<GUITextarea
							{...register('description', {
								maxLength: {
									value: 300,
									message: t('message.validation-group-description-length'),
								},
							})}
							error={errors.description?.message}
						/>
						<div className="text-right text-xs">{watch('description')?.length} / 300</div>
					</div>
				</div>
			)}

			{tab === 'objects' && (
				<div className="p-4 gap-4">
					<p className="text-sm text-center text-[#999] flex justify-center my-5">{t('message.error-started-create-group')}</p>
				</div>
			)}

			{tab === 'right' && (
				<div className="space-y-4 my-5">
					{Object.entries(grouped).map(([group, rules]) => (
						<div key={group}>
							<div className="flex items-center gap-1">
								<GUICheckbox checked={isGroupChecked(group)} onChange={(v) => toggleGroup(group, v)} />
								<p>{group === 'main' ? t('label.mains') : group === 'device' ? t('label.device') : t('label.other')}</p>
							</div>
							<div className="ml-6 space-y-2 mt-2">
								{rules.map((rule) => (
									<div key={rule.key} className="flex items-center gap-1">
										<GUICheckbox checked={!!rightsState[rule.key as keyof typeof BASE_GROUP_RIGHTs_MASK]} onChange={(v) => toggleOne(rule.key as keyof typeof BASE_GROUP_RIGHTs_MASK, v)} />
										<p>{t(rule.label)}</p>
									</div>
								))}
							</div>
						</div>
					))}
				</div>
			)}

			<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} className="mt-6" />

			<GUIButton type="button" className="mt-4" onClick={handleSubmit(submitCreate)}>
				{t('label.create')}
			</GUIButton>
		</div>
	);
};

export default GroupCreateModal;
