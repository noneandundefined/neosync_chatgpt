import { useTranslation } from 'react-i18next';
import { Loading, SpinnerLoading } from '../Loader/Loading';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { DeviceShort, DevicesForSendCommand } from '@/rest/deviceCommandAPI';
import { basicDevicesState, basicDevicesStateDevices } from '@/rest/deviceAPI';
import React, { Dispatch, SetStateAction, useCallback, useEffect, useRef, useState } from 'react';

import ChevronUp from '@/components/@icons/chevron-up';
import ChevronDown from '@/components/@icons/chevron-down';

const PAGE_SIZE = 25;

interface CommandDevicesProps {
	selectedDevice: DeviceShort[];
	setSelectedDevice: Dispatch<SetStateAction<DeviceShort[]>>;
}

interface GroupDevicesState {
	items: DeviceShort[];
	page: number;
	hasMore: boolean;
	loading: boolean;
	total: number;
}

const CommandDevices: React.FC<CommandDevicesProps> = ({ selectedDevice, setSelectedDevice }) => {
	const { t } = useTranslation();

	const [search, setSearch] = useState<string>('');
	const [debouncedSearch, setDebouncedSearch] = useState<string>('');
	const [openGroups, setOpenGroups] = useState<number[]>([]);
	const [groupDevices, setGroupDevices] = useState<Record<number, GroupDevicesState>>({});

	const loadingMoreRef = useRef<Record<number, boolean>>({});
	const searchRef = useRef(debouncedSearch);
	searchRef.current = debouncedSearch;

	useEffect(() => {
		const timeout = setTimeout(() => setDebouncedSearch(search.trim()), 300);
		return () => clearTimeout(timeout);
	}, [search]);

	const { data: respDevicesState, loading: respDevicesStateLoading } = useHandleServer(['respDevicesState', debouncedSearch], () => basicDevicesState(debouncedSearch));

	useEffect(() => {
		setGroupDevices({});
		setOpenGroups([]);
		loadingMoreRef.current = {};
	}, [debouncedSearch]);

	const loadGroupPage = useCallback(async (groupId: number, page: number, append: boolean) => {
		if (loadingMoreRef.current[groupId]) return;

		loadingMoreRef.current[groupId] = true;
		const searchUsed = searchRef.current;

		setGroupDevices((prev) => ({
			...prev,
			[groupId]: {
				items: append ? (prev[groupId]?.items ?? []) : [],
				page: prev[groupId]?.page ?? 0,
				hasMore: prev[groupId]?.hasMore ?? true,
				loading: true,
				total: prev[groupId]?.total ?? 0,
			},
		}));

		try {
			const response = await basicDevicesStateDevices({
				groupId,
				search: searchUsed,
				page,
				limit: PAGE_SIZE,
			});

			if (searchRef.current !== searchUsed) return;

			setGroupDevices((prev) => ({
				...prev,
				[groupId]: {
					items: append ? [...(prev[groupId]?.items ?? []), ...response.items] : response.items,
					page: response.page,
					hasMore: response.has_more,
					loading: false,
					total: response.total,
				},
			}));
		} catch {
			setGroupDevices((prev) => ({
				...prev,
				[groupId]: {
					items: prev[groupId]?.items ?? [],
					page: prev[groupId]?.page ?? 0,
					hasMore: prev[groupId]?.hasMore ?? true,
					loading: false,
					total: prev[groupId]?.total ?? 0,
				},
			}));
		} finally {
			loadingMoreRef.current[groupId] = false;
		}
	}, []);

	const toggleDevice = (device: DeviceShort, groupId: number) => {
		setSelectedDevice((prev) => {
			const exists = prev.some((d) => d.imei === device.imei);

			if (exists) {
				return prev.filter((d) => d.imei !== device.imei);
			}

			return [...prev, { ...device, group_id: groupId }];
		});
	};

	const toggleGroupDevice = async (group: DevicesForSendCommand) => {
		const selectedInGroup = selectedDevice.filter((d) => d.group_id === group.group_id).length;
		const isAllSelected = group.device_count > 0 && selectedInGroup === group.device_count;

		if (isAllSelected) {
			setSelectedDevice((prev) => prev.filter((d) => d.group_id !== group.group_id));
			return;
		}

		let devices = groupDevices[group.group_id]?.items ?? [];

		if (devices.length < group.device_count) {
			const response = await basicDevicesStateDevices({
				groupId: group.group_id,
				search: debouncedSearch,
				all: true,
			});

			devices = response.items;

			setGroupDevices((prev) => ({
				...prev,
				[group.group_id]: {
					items: devices,
					page: 1,
					hasMore: false,
					loading: false,
					total: response.total,
				},
			}));
		}

		setSelectedDevice((prev) => {
			const merged = [...prev];

			devices.forEach((device) => {
				if (!merged.some((d) => d.imei === device.imei)) {
					merged.push({ ...device, group_id: group.group_id });
				}
			});

			return merged;
		});
	};

	const toggleGroupOpen = (groupId: number) => {
		const willOpen = !openGroups.includes(groupId);

		setOpenGroups((prev) => (willOpen ? [...prev, groupId] : prev.filter((id) => id !== groupId)));

		if (willOpen && !groupDevices[groupId]) {
			void loadGroupPage(groupId, 1, false);
		}
	};

	const handleListScroll = (event: React.UIEvent<HTMLDivElement>) => {
		const el = event.currentTarget;
		if (el.scrollHeight - el.scrollTop - el.clientHeight > 64) return;

		openGroups.forEach((groupId) => {
			const state = groupDevices[groupId];
			if (!state || state.loading || !state.hasMore) return;

			void loadGroupPage(groupId, state.page + 1, true);
		});
	};

	const filteredDevicesState = respDevicesState ?? [];

	return (
		<div className="min-w-[20rem] max-md:min-h-[50vh] md:h-full md:min-h-0 flex flex-col bg-white border border-[#e4e4e4] py-4 rounded">
			<div className="px-5 shrink-0">
				<h3 className="font-medium text-[17px] text-[#49525f]">{t('label.devices')}</h3>
			</div>

			<div className="px-5 mt-3 shrink-0">
				<input type="text" placeholder={t('label.search')} className="w-full bg-wgite border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm" value={search} onChange={(e) => setSearch(e.target.value)} />
			</div>

			<div className="flex-1 min-h-0 overflow-y-auto mt-3" onScroll={handleListScroll}>
				{respDevicesStateLoading ? (
					<Loading />
				) : (
					filteredDevicesState.map((state) => {
						const isOpen = openGroups.includes(state.group_id);
						const devices = groupDevices[state.group_id];
						const selectedInGroup = selectedDevice.filter((d) => d.group_id === state.group_id).length;
						const allSelected = state.device_count > 0 && selectedInGroup === state.device_count;

						return (
							<div className="border-t border-[#e4e4e4]" key={state.group_id}>
								<div className="flex cursor-pointer justify-between py-3 px-5" onClick={() => toggleGroupOpen(state.group_id)}>
									<div className="flex gap-1" onClick={(e) => e.stopPropagation()}>
										<GUICheckbox
											checked={allSelected}
											onChange={() => toggleGroupDevice(state)}
											label={
												<div className="!flex !items-center !gap-1">
													<h4 className="!text-[16px] !text-black">{state.name ? state.name : t('message.general-group')}</h4>
													<span className="!text-gray !text-[13px]">({state.device_count})</span>
												</div>
											}
										/>
									</div>

									{isOpen ? <ChevronUp fill="#49525f" /> : <ChevronDown fill="#49525f" />}
								</div>

								{isOpen && (
									<div className="space-y-1">
										{!devices || (devices.loading && devices.items.length === 0) ? (
											<Loading />
										) : (
											devices.items.map((device) => {
												const isChecked = selectedDevice.some((d) => d.imei === device.imei);

												return (
													<div className="flex items-center justify-between px-8 gap-2 py-[2px] hover:bg-[#f3f3f3]" key={device.imei}>
														<GUICheckbox size="15px" checked={isChecked} onChange={() => toggleDevice(device, state.group_id)} label={<p className="!text-[14px]">{device.imei}</p>} />
													</div>
												);
											})
										)}

										{devices?.loading && devices.items.length > 0 && <SpinnerLoading />}
									</div>
								)}
							</div>
						);
					})
				)}
			</div>
		</div>
	);
};

export default CommandDevices;
