import i18next from 'i18next';
import axiosClient from './axios';
import { toast } from 'react-toastify';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { runGuardedAction } from '@/utils/RequestControlUtils';
import { BASE_GROUP_RIGHTs_MASK } from '@/constants/Accesses.constant';
import { ACTIONKEYs, CACHEKEYs } from '@/constants/CacheKeys.constants';
import { GroupUpdateRequest } from '@/interface/group/groupUpdateRequest.interface';
import { GroupCreateRequest } from '@/interface/group/groupCrerateRequest.interface';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';
import { GroupMembersBatchRequest } from '@/interface/group/groupMembersBatchRequest.interface';

const apiPath = '/groups';

export const groupRightsStateToApiPayload = (state: typeof BASE_GROUP_RIGHTs_MASK) => ({
	can_edit_group: state.can_edit_group,
	can_manage_devices: state.can_manage_devices,
	can_read_config: state.can_read_config,
	can_edit_config: state.can_edit_config,
	can_send_commands: state.can_send_commands,
});

export interface GroupResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	name: string;
	description: string;
	objects: number;
	can_edit_group?: boolean;
	can_manage_devices?: boolean;
	can_read_config?: boolean;
	can_edit_config?: boolean;
	can_send_commands?: boolean;
}

export interface GDeviceResponse {
	id: number;
	imei: string;
}

export interface GroupWithDevicesResponse {
	group: GroupResponse;
	devices: {
		assigned: GDeviceResponse[];
		available: GDeviceResponse[];
	};
}

export interface GroupMemberResponse {
	created_at: string;
	group_id: number;
	member_user_uuid: string;
}

export interface GroupsWPResponse {
	items: GroupResponse[];
	page: number;
	limit: number;
	total: number;
	total_all: number;
	total_pages: number;
}

export interface GroupFieldsNameResponse {
	id: number;
	name: string;
}

/**
 * Создания группы
 */
export const basicGroupCreate = async (payload: GroupCreateRequest) => {
	const response = await axiosClient.post(`${apiPath}`, payload);

	sessionStorage.removeItem(CACHEKEYs.GROUP_KEY);
	toast.success(response.data.message);
};

/**
 * Вывод групп с пагинацией
 */
export const basicGroupsGetWithParams = async (page: number, limit: number, search: string = '', columnSort: ColumnSortState, signal?: AbortSignal): Promise<GroupsWPResponse> => {
	const response = await axiosClient.get(`${apiPath}?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}&columnSortKey=${columnSort.key}&columnSortDir=${columnSort.direction}`, { signal });
	return response.data.message;
};

/**
 * Вывод ВСЕХ групп
 */
export const basicGroupsGetFull = async (): Promise<GroupFieldsNameResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/full-list`);
	return response.data.message;
};

/**
 * Вывод информации о группе с обьектами
 */
export const basicGroupWithDevicesGet = async (id: number, signal?: AbortSignal): Promise<GroupWithDevicesResponse> => {
	const response = await axiosClient.get(`${apiPath}/${id}`, {
		signal,
	});

	return response.data.message;
};

/**
 * Удаление группы
 */
export const basicGroupDelete = async (id: number) => {
	const response = await axiosClient.delete(`${apiPath}/${id}`);

	sessionStorage.removeItem(CACHEKEYs.GROUP_KEY);
	toast.success(response.data.message);
};

/**
 * Удаление групп
 */
export const basicGroupMassiveDelete = async (ids: number[]) => {
	await runGuardedAction(
		ACTIONKEYs.GROUP_DELETE_MASSIVE,
		async () => {
			const response = await axiosClient.post(`${apiPath}/delete`, { ids }, { actionKey: ACTIONKEYs.GROUP_DELETE_MASSIVE } as any);
			sessionStorage.removeItem(CACHEKEYs.GROUP_KEY);
			toast.success(response.data.message);
		},
		{
			cooldownMs: 2000,
			getCooldownMessage: (seconds) => i18next.t('message.client-action-please-wait-seconds', { seconds }),
			onBlocked: (msg) => showServerErrorToast(msg),
		}
	);
};

/**
 * Обновление данных группы
 */
export const basicGroupUpdate = async (id: number, payload: GroupUpdateRequest) => {
	const response = await axiosClient.put(`${apiPath}/${id}`, payload);
	sessionStorage.removeItem(CACHEKEYs.GROUP_KEY);
	toast.success(response.data.message);
};

/**
 * Получение fields name group
 */
export const basicGroupFieldsName = async (): Promise<GroupFieldsNameResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/fields/name`);
	return response.data.message;
};

/**
 * Участники группы
 */
export const basicGroupMembersGet = async (groupId: number, signal?: AbortSignal): Promise<GroupMemberResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/${groupId}/members`, { signal });
	return response.data.message ?? [];
};

/**
 * Массовое предоставление доступа к группе
 */
export const basicGroupMembersUpsert = async (groupId: number, payload: GroupMembersBatchRequest) => {
	const response = await axiosClient.post(`${apiPath}/${groupId}/members/upsert`, payload);
	toast.success(response.data.message);
};

/**
 * Массовый отзыв доступа к группе
 */
export const basicGroupMembersRemove = async (groupId: number, payload: GroupMembersBatchRequest) => {
	const response = await axiosClient.post(`${apiPath}/${groupId}/members/remove`, payload);
	toast.success(response.data.message);
};
