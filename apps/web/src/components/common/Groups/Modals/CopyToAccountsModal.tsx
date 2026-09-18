import useCaptcha from '@/hooks/useCaptcha';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import GUIButton from '@/components/ui/Button/GUIButton';
import { useModalContext } from '@/context/useModalContext';
import { ROLEs_AND_NAME } from '@/constants/Roles.constant';
import { Loading, SpinnerLoading } from '../../Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicUsersGetTree, UserTreeResponse } from '@/rest/userAPI';
import SmartCaptchaWidget from '@/components/Security/SmartCaptchaWidget';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { basicGroupMembersGet, basicGroupMembersRemove, basicGroupMembersUpsert } from '@/rest/groupAPI';

const PAGE_SIZE = 25;

interface CopyToAccountsModalProps {
	groupId: number;
	onSuccess?: () => void;
}

const CopyToAccountsModal: React.FC<CopyToAccountsModalProps> = ({ groupId, onSuccess }) => {
	const captcha = useCaptcha();

	const { t } = useTranslation();
	const { close } = useModalContext();
	const queryClient = useQueryClient();

	const [page, setPage] = useState(0);
	const [search, setSearch] = useState('');
	const [hasMore, setHasMore] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [usersLoading, setUsersLoading] = useState(true);
	const [debouncedSearch, setDebouncedSearch] = useState('');
	const [users, setUsers] = useState<UserTreeResponse[]>([]);
	const [selectedUuids, setSelectedUuids] = useState<Set<string>>(() => new Set());
	const [selectedInfo, setSelectedInfo] = useState<Record<string, UserTreeResponse>>({});
	const [membersQueryEnabled, setMembersQueryEnabled] = useState(false);

	useEffect(() => {
		queryClient.removeQueries({ queryKey: ['groupMembers', groupId] });
		setSelectedUuids(new Set());
		setSelectedInfo({});
		setMembersQueryEnabled(true);
	}, [groupId, queryClient]);

	const { data: respGroupMembers, loading: membersLoading } = useHandleServer(['groupMembers', groupId], () => basicGroupMembersGet(groupId), {
		enabled: membersQueryEnabled,
		staleTime: 0,
		gcTime: 0,
	});

	const memberUuidSet = useMemo(() => new Set((respGroupMembers ?? []).map((m) => m.member_user_uuid)), [respGroupMembers]);

	const loadingMoreRef = useRef(false);
	const searchRef = useRef(debouncedSearch);

	searchRef.current = debouncedSearch;

	useEffect(() => {
		const timeout = setTimeout(() => setDebouncedSearch(search.trim()), 300);
		return () => clearTimeout(timeout);
	}, [search]);

	useEffect(() => {
		if (!membersQueryEnabled || membersLoading || !respGroupMembers) return;

		setSelectedUuids(new Set(respGroupMembers.map((m) => m.member_user_uuid)));
	}, [membersQueryEnabled, membersLoading, respGroupMembers]);

	const loadPage = useCallback(async (nextPage: number, append: boolean) => {
		if (loadingMoreRef.current) return;

		loadingMoreRef.current = true;
		const searchUsed = searchRef.current;
		setUsersLoading(true);

		try {
			const response = await basicUsersGetTree({
				search: searchUsed,
				page: nextPage,
				limit: PAGE_SIZE,
			});

			if (searchRef.current !== searchUsed) return;

			setUsers((prev) => (append ? [...prev, ...response.items] : response.items));
			setPage(response.page);
			setHasMore(response.has_more);
			setSelectedInfo((prev) => {
				const next = { ...prev };
				response.items.forEach((user) => {
					next[user.user_uuid] = user;
				});
				return next;
			});
		} finally {
			loadingMoreRef.current = false;
			setUsersLoading(false);
		}
	}, []);

	useEffect(() => {
		setUsers([]);
		setPage(0);
		setHasMore(true);
		void loadPage(1, false);
	}, [debouncedSearch, loadPage]);

	const toggleUser = (user: UserTreeResponse) => {
		setSelectedInfo((prev) => ({ ...prev, [user.user_uuid]: user }));
		setSelectedUuids((prev) => {
			const next = new Set(prev);

			if (next.has(user.user_uuid)) {
				next.delete(user.user_uuid);
			} else {
				next.add(user.user_uuid);
			}

			return next;
		});
	};

	const selectedUsers = useMemo(() => [...selectedUuids].map((uuid) => selectedInfo[uuid]).filter(Boolean), [selectedUuids, selectedInfo]);

	const selectedToGrant = useMemo(() => [...selectedUuids].filter((uuid) => !memberUuidSet.has(uuid)), [selectedUuids, memberUuidSet]);

	const selectedToRevoke = useMemo(() => [...memberUuidSet].filter((uuid) => !selectedUuids.has(uuid)), [selectedUuids, memberUuidSet]);

	const hasChanges = selectedToGrant.length > 0 || selectedToRevoke.length > 0;
	const loading = (usersLoading && users.length === 0) || membersLoading;

	const handleListScroll = (event: React.UIEvent<HTMLDivElement>) => {
		const el = event.currentTarget;
		if (!hasMore || usersLoading) return;
		if (el.scrollHeight - el.scrollTop - el.clientHeight > 48) return;

		void loadPage(page + 1, true);
	};

	const handleConfirm = async () => {
		if (!hasChanges || !captcha.validate()) return;

		setSubmitting(true);
		try {
			const payloadBase = { turnstile_token: captcha.token };

			if (selectedToGrant.length > 0) {
				await basicGroupMembersUpsert(groupId, {
					...payloadBase,
					member_uuids: selectedToGrant,
				});
			}

			if (selectedToRevoke.length > 0) {
				await basicGroupMembersRemove(groupId, {
					...payloadBase,
					member_uuids: selectedToRevoke,
				});
			}

			setSelectedUuids(new Set());
			queryClient.removeQueries({ queryKey: ['groupMembers', groupId] });
			onSuccess?.();
			close();
		} finally {
			captcha.reset();
			setSubmitting(false);
		}
	};

	const renderUserRow = (user: UserTreeResponse) => {
		const isSelected = selectedUuids.has(user.user_uuid);

		return (
			<div
				className={`flex mr-2 items-center justify-between px-5 py-1 rounded-lg border ${isSelected ? 'hover:bg-[#ffd1ce] hover:border-[#de1303]' : 'hover:bg-[#E8EDF5] hover:border-[#7d92b4]'} ${isSelected && 'border-[#7d92b4] bg-[#E8EDF5]'} transition cursor-pointer`}
				key={user.user_uuid}
				onClick={() => toggleUser(user)}
			>
				<div className="flex items-center gap-3">
					<div className="border border-[#395d95] p-[2px] rounded-full">
						<div className="w-[2rem] h-[2rem] rounded-[50%] flex items-center justify-center" style={{ background: '#395d95' }}>
							<p className="text-[#fff] text-[15px]">{user.email.charAt(0).toUpperCase()}</p>
						</div>
					</div>

					<div>
						<p>{user.email}</p>
					</div>
				</div>

				<p className="text-sm text-[#CCCCCC]">{user.role_code ? t(`label.${ROLEs_AND_NAME[user.role_code]}`) : t('label.role-not-found')}</p>
			</div>
		);
	};

	return (
		<div>
			{selectedUuids.size > 0 && (
				<div className="mb-3 p-2 rounded border">
					<div className="flex items-center justify-between mb-4">
						<p className="text-sm text-[#333] font-medium">
							{t('label.selected')}: {selectedUuids.size}
						</p>
					</div>

					<div className="flex flex-wrap gap-2 max-h-[80px] overflow-y-auto">
						{selectedUsers.map((user) => (
							<div key={user.user_uuid} className="flex items-center gap-1 px-4 py-1 rounded-full text-[12px] border border-[#49525f]">
								<p className="text-[#49525f]">{user.email}</p>
							</div>
						))}
					</div>
				</div>
			)}

			<input
				type="text"
				placeholder={t('label.search')}
				className="w-full mb-2 md:min-w-[16rem] bg-wgite border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm"
				value={search}
				onChange={(e) => setSearch(e.target.value)}
			/>

			<div className="space-y-2 max-h-[220px] overflow-y-auto" onScroll={handleListScroll}>
				{loading ? <Loading /> : users.length > 0 ? users.map(renderUserRow) : <p className="my-3 text-center text-[#999]">{t('label.no-users')}</p>}

				{usersLoading && users.length > 0 && <SpinnerLoading />}
			</div>

			<SmartCaptchaWidget ref={captcha.widgetRef} onVerify={captcha.onVerify} className="mt-6 text-left" />

			<GUIButton className="mt-5" onClick={handleConfirm} disabled={!hasChanges || submitting}>
				{t('label.confirm')}
			</GUIButton>
		</div>
	);
};

export default CopyToAccountsModal;
