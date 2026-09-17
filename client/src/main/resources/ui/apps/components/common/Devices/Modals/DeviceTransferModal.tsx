import useCaptcha from '@/hooks/useCaptcha';
import { useTranslation } from 'react-i18next';
import GUIButton from '@/components/ui/Button/GUIButton';
import { Loading, SpinnerLoading } from '../../Loader/Loading';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';

import { useRole } from '@/context/RoleContext/useRoleContext';
import { basicUsersToOwnerGet, UserToOwnerResponse } from '@/rest/userAPI';
import { basicDeviceTransfer, Device_DeviceConf_Sync } from '@/rest/deviceAPI';
import { hasPermission, ROLEs_AND_NAME, ROLES } from '@/constants/Roles.constant';
import { DeviceTransferRequest } from '@/interface/device/deviceTransferRequest.interface';

const PAGE_SIZE = 25;

interface DeviceTransferModalProps {
	onSuccess: () => void;
	onTransfer?: (user: { user_uuid: string; email: string } | null) => void;
	device?: Device_DeviceConf_Sync;
	imeis?: string[];
}

const DeviceTransferModal: React.FC<DeviceTransferModalProps> = ({ onSuccess, onTransfer, device, imeis }) => {
	const { role } = useRole();
	const { t } = useTranslation();

	const captcha = useCaptcha();
	const canViewUsers = Boolean(role && hasPermission(role, 'view:users'));

	const [page, setPage] = useState(0);
	const [hasMore, setHasMore] = useState(true);
	const [search, setSearch] = useState<string>('');
	const [usersLoading, setUsersLoading] = useState(false);
	const [users, setUsers] = useState<UserToOwnerResponse[]>([]);
	const [debouncedSearch, setDebouncedSearch] = useState<string>('');
	const [selectedUserUuid, setSelectedUserUuid] = useState<string>('');
	const [selectedUser, setSelectedUser] = useState<UserToOwnerResponse | null>(null);

	const loadingMoreRef = useRef(false);
	const searchRef = useRef(debouncedSearch);

	searchRef.current = debouncedSearch;

	const isDealerTransfer = role === ROLES.DEALER || role === ROLES.DEALER_SUPPORT;

	const initialUserUuid = device && (!imeis || imeis.length <= 1) ? device.owner_uuid || device.user_uuid || '' : '';
	const hasOwnerChanged = selectedUserUuid !== initialUserUuid;

	useEffect(() => {
		setSelectedUserUuid(initialUserUuid);
		setSelectedUser(null);
	}, [initialUserUuid]);

	useEffect(() => {
		const timeout = setTimeout(() => setDebouncedSearch(search.trim()), 300);
		return () => clearTimeout(timeout);
	}, [search]);

	const loadPage = useCallback(
		async (nextPage: number, append: boolean) => {
			if (!canViewUsers || loadingMoreRef.current) return;

			loadingMoreRef.current = true;
			const searchUsed = searchRef.current;
			setUsersLoading(true);

			try {
				const response = await basicUsersToOwnerGet({
					search: searchUsed,
					page: nextPage,
					limit: PAGE_SIZE,
				});

				if (searchRef.current !== searchUsed) return;

				setUsers((prev) => (append ? [...prev, ...response.items] : response.items));
				setPage(response.page);
				setHasMore(response.has_more);
				setSelectedUser((prev) => response.items.find((user) => user.user_uuid === (prev?.user_uuid || initialUserUuid)) ?? prev);
			} finally {
				loadingMoreRef.current = false;
				setUsersLoading(false);
			}
		},
		[canViewUsers, initialUserUuid]
	);

	useEffect(() => {
		if (!canViewUsers) return;

		setUsers([]);
		setPage(0);
		setHasMore(true);
		void loadPage(1, false);
	}, [canViewUsers, debouncedSearch, loadPage]);

	const handleListScroll = (event: React.UIEvent<HTMLDivElement>) => {
		const el = event.currentTarget;
		if (!hasMore || usersLoading) return;
		if (el.scrollHeight - el.scrollTop - el.clientHeight > 48) return;

		void loadPage(page + 1, true);
	};

	const handleTransfer = async () => {
		if (!hasOwnerChanged) {
			onSuccess();
			return;
		}

		if (!captcha.validate()) return;

		const transferRequest: DeviceTransferRequest = {
			imeis: imeis && imeis.length > 0 ? imeis : device ? [device.imei] : [],
			uuid: selectedUserUuid,
			turnstile_token: captcha.token,
		};

		await basicDeviceTransfer(transferRequest).finally(captcha.reset);

		onTransfer?.(selectedUser);
		onSuccess();
	};

	return (
		<div>
			<input
				type="text"
				placeholder={t('label.search')}
				className="w-full mb-2 md:min-w-[16rem] bg-wgite border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm"
				value={search}
				onChange={(e) => setSearch(e.target.value)}
			/>

			<div className="space-y-2 max-h-[280px] overflow-y-auto" onScroll={handleListScroll}>
				{usersLoading && users.length === 0 ? (
					<Loading />
				) : users.length > 0 ? (
					users.map((user) => {
						const isSelected = selectedUserUuid === user.user_uuid;

						return (
							<div
								className={`flex mr-2 items-center justify-between px-5 py-1 rounded-lg border ${isSelected && !isDealerTransfer ? 'hover:bg-[#ffd1ce] hover:border-[#de1303]' : 'hover:bg-[#E8EDF5] hover:border-[#7d92b4]'} ${isSelected && 'border-[#7d92b4] bg-[#E8EDF5]'} transition cursor-pointer`}
								key={user.user_uuid}
								onClick={() => {
									if (user.user_uuid === selectedUserUuid) {
										if (isDealerTransfer) {
											return;
										}

										setSelectedUserUuid('');
										setSelectedUser(null);
										return;
									}

									setSelectedUserUuid(user.user_uuid);
									setSelectedUser(user);
								}}
							>
								<div className="flex items-center gap-3">
									<div className="border border-[#395d95] p-[2px] rounded-full">
										<div className="w-[2rem] h-[2rem] rounded-[50%] flex items-center justify-center" style={{ background: '#395d95' }}>
											<p className="text-[#fff] text-[15px]">{user.email.charAt(0).toUpperCase()}</p>
										</div>
									</div>

									<p>{user.email}</p>
								</div>

								<p className="text-sm text-[#CCCCCC]">{user.role_code ? t(`label.${ROLEs_AND_NAME[user.role_code]}`) : t('label.role-not-found')}</p>
							</div>
						);
					})
				) : (
					<p className="my-3 text-center text-[#999]">{t('label.no-users')}</p>
				)}

				{usersLoading && users.length > 0 && <SpinnerLoading />}
			</div>

			<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} className="mt-6 text-left" />

			<GUIButton className="mt-5" onClick={handleTransfer}>
				{t('label.confirm')}
			</GUIButton>
		</div>
	);
};

export default DeviceTransferModal;
