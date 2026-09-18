import { useForm } from 'react-hook-form';
import Plus from '@/components/@icons/plus';
import Minus from '@/components/@icons/minus';
import { useTranslation } from 'react-i18next';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUIButton from '@/components/ui/Button/GUIButton';
import { Loading } from '@/components/common/Loader/Loading';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { GUITextarea } from '@/components/ui/Textarea/GUITextarea';
import { useCallback, useEffect, useMemo, useState, type UIEvent } from 'react';
import { GroupUpdateRequest } from '@/interface/group/groupUpdateRequest.interface';
import { BASE_GROUP_RIGHTs_MASK, GROUP_RIGHTs_RULES } from '@/constants/Accesses.constant';
import { basicGroupUpdate, basicGroupWithDevicesGet, groupRightsStateToApiPayload } from '@/rest/groupAPI';

interface GroupEditModalProps {
	onSuccess: () => void;
	groupId: number;
	userUuid?: string;
}

type NameDescForm = Pick<GroupUpdateRequest, 'name' | 'description'>;

const PAGE_SIZE = 25;

const GroupEditModal: React.FC<GroupEditModalProps> = ({ onSuccess, groupId, userUuid }) => {
	const { t } = useTranslation();

	const [search, setSearch] = useState('');
	const [tab, setTab] = useState<'main' | 'objects' | 'right'>('main');

	const [availableVisible, setAvailableVisible] = useState(PAGE_SIZE);
	const [assignedVisible, setAssignedVisible] = useState(PAGE_SIZE);

	const [pendingAdd, setPendingAdd] = useState<Set<number>>(new Set());
	const [pendingRemove, setPendingRemove] = useState<Set<number>>(new Set());

	const fetchGroupWithDevicesGet = useCallback(() => basicGroupWithDevicesGet(groupId), [groupId]);
	const { data: resp, loading: groupLoading } = useHandleServer(['fetchGroupWithDevicesGet', groupId], fetchGroupWithDevicesGet, { staleTime: 0 });

	const isOwner = !!userUuid && !!resp?.group.user_uuid && resp.group.user_uuid === userUuid;

	const {
		register,
		watch,
		setValue,
		getValues,
		handleSubmit,
		formState: { errors },
	} = useForm<NameDescForm>({
		mode: 'onChange',
		defaultValues: {
			name: '',
			description: '',
		},
	});

	const [rightsState, setRightsState] = useState<typeof BASE_GROUP_RIGHTs_MASK>({ ...BASE_GROUP_RIGHTs_MASK });

	useEffect(() => {
		if (!resp) return;

		setValue('name', resp.group.name, {
			shouldValidate: true,
		});

		setValue('description', resp.group.description, {
			shouldValidate: true,
		});

		setPendingAdd(new Set());
		setPendingRemove(new Set());

		setRightsState({
			can_edit_group: resp.group.can_edit_group ?? false,
			can_manage_devices: resp.group.can_manage_devices ?? false,
			can_read_config: resp.group.can_read_config ?? false,
			can_edit_config: resp.group.can_edit_config ?? false,
			can_send_commands: resp.group.can_send_commands ?? false,
		});
	}, [resp, setValue]);

	const saveAll = async () => {
		const data = getValues();

		const payload: GroupUpdateRequest = {
			name: data.name,
			description: data.description ?? '',
			devices_add: Array.from(pendingAdd),
			devices_remove: Array.from(pendingRemove),
			...(isOwner
				? groupRightsStateToApiPayload(rightsState)
				: groupRightsStateToApiPayload({
						can_edit_group: resp?.group.can_edit_group ?? false,
						can_manage_devices: resp?.group.can_manage_devices ?? false,
						can_read_config: resp?.group.can_read_config ?? false,
						can_edit_config: resp?.group.can_edit_config ?? false,
						can_send_commands: resp?.group.can_send_commands ?? false,
					})),
		};

		await basicGroupUpdate(groupId, payload);

		setPendingAdd(new Set());
		setPendingRemove(new Set());

		onSuccess();
	};

	const togglePending = (id: number, add: boolean) => {
		if (add) {
			setPendingAdd((prev) => new Set(prev).add(id));
			setPendingRemove((prev) => {
				const s = new Set(prev);
				s.delete(id);
				return s;
			});
		} else {
			setPendingRemove((prev) => new Set(prev).add(id));
			setPendingAdd((prev) => {
				const s = new Set(prev);
				s.delete(id);
				return s;
			});
		}
	};

	const baseAssignedIds = useMemo(() => new Set(resp?.devices?.assigned?.map((d) => d.id) ?? []), [resp]);

	const virtualGroupDevices = useMemo(() => {
		const set = new Set(baseAssignedIds);

		pendingRemove.forEach((id) => set.delete(id));
		pendingAdd.forEach((id) => set.add(id));

		return set;
	}, [baseAssignedIds, pendingAdd, pendingRemove]);

	const devices = useMemo(() => [...(resp?.devices.assigned ?? []), ...(resp?.devices.available ?? [])], [resp]);

	const filterDevices = (assigned: boolean) => devices.filter((d) => virtualGroupDevices.has(d.id) === assigned && d.imei.toLowerCase().includes(search));

	const assigned = useMemo(() => filterDevices(true), [devices, virtualGroupDevices, search]);
	const available = useMemo(() => filterDevices(false), [devices, virtualGroupDevices, search]);

	useEffect(() => {
		setAvailableVisible(PAGE_SIZE);
		setAssignedVisible(PAGE_SIZE);
	}, [search]);

	const handleListScroll = (event: UIEvent<HTMLDivElement>, kind: 'available' | 'assigned') => {
		const el = event.currentTarget;
		if (el.scrollHeight - el.scrollTop - el.clientHeight > 48) return;

		if (kind === 'available') {
			setAvailableVisible((prev) => Math.min(prev + PAGE_SIZE, available.length));
			return;
		}

		setAssignedVisible((prev) => Math.min(prev + PAGE_SIZE, assigned.length));
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

	if (groupLoading || !resp) {
		return (
			<div className="py-8">
				<Loading />
			</div>
		);
	}

	return (
		<div className="w-full mx-auto mt-4 rounded">
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
				<form className="pt-4 space-y-4">
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

					<GUIButton type="button" onClick={handleSubmit(saveAll)}>
						{t('label.save')}
					</GUIButton>
				</form>
			)}

			{tab === 'objects' && (
				<div>
					<div className="py-3 flex flex-col sm:flex-row justify-between gap-4">
						<div className="min-w-[47%]">
							<input
								type="text"
								placeholder={t('label.search')}
								className="w-full md:min-w-[16rem] md:max-w-[35rem] bg-wgite border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm"
								value={search}
								onChange={(e) => setSearch(e.target.value)}
							/>

							<div className="mt-2 max-h-[200px] sm:max-h-[280px] overflow-y-auto" onScroll={(event) => handleListScroll(event, 'available')}>
								{available.length ? (
									available
										.slice(0, availableVisible)
										.map((obj) => <DeviceRow key={obj.id} obj={obj} label={t('label.add')} icon={<Plus size={19} fill="#49525f" />} onClick={() => togglePending(obj.id, true)} />)
								) : (
									<p className="text-center text-gray-500 py-3 text-sm">{t('label.no-devices')}</p>
								)}
							</div>
						</div>

						<div className="h-auto w-[1px] bg-[#ccc]" />

						<div className="min-w-[48%] max-h-[200px] sm:max-h-[280px] overflow-y-auto" onScroll={(event) => handleListScroll(event, 'assigned')}>
							{assigned.length ? (
								assigned
									.slice(0, assignedVisible)
									.map((obj) => <DeviceRow key={obj.id} obj={obj} action="remove" label={t('label.remove')} icon={<Minus size={19} fill="#49525f" />} onClick={() => togglePending(obj.id, false)} />)
							) : (
								<p className="text-center text-gray-500 py-3 text-sm">{t('label.no-devices')}</p>
							)}
						</div>
					</div>

					<GUIButton type="button" className="mt-3" onClick={handleSubmit(saveAll)}>
						{t('label.save')}
					</GUIButton>
				</div>
			)}

			{tab === 'right' && (
				<div>
					<div className={!isOwner ? 'pointer-events-none opacity-60' : undefined}>
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
					</div>

					{isOwner && (
						<GUIButton type="button" onClick={handleSubmit(saveAll)}>
							{t('label.save')}
						</GUIButton>
					)}
				</div>
			)}
		</div>
	);
};

const DeviceRow = ({ obj, action, label, icon, onClick }: any) => (
	<div className="group py-[0.6rem] sm:py-[0.4rem] px-[0.7rem] flex items-center justify-between cursor-pointer hover:bg-[#eee]">
		<div className="flex items-center gap-5">
			{action === 'remove' && <div className="w-[5px] h-[5px] rounded-full bg-[#49525f]" />}
			<p className="text-[14px] text-[#49525f]">{obj.imei}</p>
		</div>

		<div className="flex items-center sm:opacity-0 group-hover:opacity-100" onClick={onClick}>
			<p className="text-[#49525f] hover:text-[#000]">{label}</p>
			{icon}
		</div>
	</div>
);

export default GroupEditModal;
