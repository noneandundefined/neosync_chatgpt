import i18next from 'i18next';
import axiosClient from './axios';
import { toast } from 'react-toastify';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { ACTIONKEYs } from '@/constants/CacheKeys.constants';
import { runGuardedAction } from '@/utils/RequestControlUtils';
import { UserUpdateRequest } from '@/interface/user/userUpdateRequest.interface';
import { UserCreateRequest } from '@/interface/auth/userCreateRequest.interface';
import { ColumnSortState } from '@/components/common/Table/GenericTable/GenericTable';

const apiPath = '/users';

export interface UserCoreResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	parent_uuid?: string | null;
	can_view_child_groups: boolean;
	email: string;
	password: string;
	refresh_token?: string | null;
}

export interface UserContactResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	name_organization: string;
	locality: string;
	phone: string;
	language: string;
}

export interface UserRoleResponse {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
	role_code: string;
}

interface UserPermissions {
	access_treker_create?: boolean;
	access_treker_edit?: boolean;
	access_treker_delete?: boolean;
	access_group_manage?: boolean;
	access_configuration_read?: boolean;
	access_configuration_apply?: boolean;
	access_configuration_history?: boolean;
	access_command_send?: boolean;
	access_log_read?: boolean;
}

export interface UserAccessResponse extends UserPermissions {
	id: number;
	created_at: string;
	updated_at: string;
	user_uuid: string;
}

export interface UserLoginWithRoleCodeResponse {
	login: string;
	role_code: string;
}

export interface Role {
	id: number;
	role_code: string;
	description?: string;
}

export interface UserResponse extends UserPermissions {
	user_contact: UserContactResponse;
	parent_uuid?: string | null;
	can_view_child_groups: boolean;
	email: string;
	parent_email?: string | null;
	role: Role;
}

export interface UserAuthResponse extends UserPermissions {
	force_neosync_configuration_priority: boolean;
	user_contact: UserContactResponse;
	parent_uuid?: string | null;
	can_view_child_groups: boolean;
	email: string;
	parent_email?: string | null;
	password: string;
	refresh_token?: string | null;
	role: Role;
}

export interface UserToOwnerResponse {
	user_uuid: string;
	email: string;
	role_code: string;
}

export interface UsersWPResponse {
	items: UserResponse[];
	page: number;
	limit: number;
	total: number;
	total_all: number;
	total_pages: number;
}

export interface UserTreeResponse {
	user_uuid: string;
	parent_uuid?: string | null;
	role_code: string;
	email: string;
}

export interface UsersPageResponse<T> {
	items: T[];
	page: number;
	limit: number;
	total: number;
	has_more: boolean;
}

/**
 * Получение данных пользователя
 */
export const basicUserGet = async (uuid: string, signal?: AbortSignal): Promise<UserAuthResponse> => {
	const response = await axiosClient.get(`${apiPath}/${uuid}`, { signal });
	return response.data.message;
};

/**
 * Обновление данных пользователя
 */
export const basicUserUpdate = async (uuid: string, payload: UserUpdateRequest) => {
	const response = await axiosClient.patch(`${apiPath}/${uuid}`, payload);
	toast.success(response.data.message);
};

/**
 * Получение профиля
 */
export const basicUserGetMe = async (signal?: AbortSignal): Promise<UserAuthResponse> => {
	const response = await axiosClient.get(`${apiPath}/me`, { signal });
	return response.data.message;
};

/**
 * Получение доступов профиля
 */
export const basicUserGetMeAccesses = async (signal?: AbortSignal): Promise<UserAccessResponse> => {
	const response = await axiosClient.get(`${apiPath}/me/accesses`, { signal });
	return response.data.message;
};

/**
 * Получение минимальной информации профиля
 */
export const basicUserGetMeLoginWithRoleCode = async (signal?: AbortSignal): Promise<UserLoginWithRoleCodeResponse> => {
	const response = await axiosClient.get(`${apiPath}/me/login-with-role`, { signal });
	return response.data.message;
};

/**
 * Обновления профиля
 */
export const basicUserUpdateMe = async (payload: UserUpdateRequest) => {
	const response = await axiosClient.patch(`${apiPath}/me`, payload);
	toast.success(response.data.message);
};

/**
 * Рекурсивное дерево пользователей (по иерархии роли)
 */
export const basicUsersGetTree = async (params?: { search?: string; page?: number; limit?: number; all?: boolean; signal?: AbortSignal }): Promise<UsersPageResponse<UserTreeResponse>> => {
	const response = await axiosClient.get(`${apiPath}/tree`, {
		signal: params?.signal,
		params: {
			search: params?.search || undefined,
			page: params?.all ? undefined : params?.page,
			limit: params?.all ? undefined : params?.limit,
			all: params?.all ? true : undefined,
		},
	});
	return response.data.message;
};

/**
 * Получение пользователей с пагинацией
 */
export const basicUsersGetWithParams = async (page: number, limit: number, search: string = '', role: string = '', columnSort: ColumnSortState, signal?: AbortSignal): Promise<UsersWPResponse> => {
	const response = await axiosClient.get(`${apiPath}?page=${page}&limit=${limit}&search=${encodeURIComponent(search)}&role=${role}&columnSortKey=${columnSort.key}&columnSortDir=${columnSort.direction}`, { signal });
	return response.data.message;
};

/**
 * Получение ВСЕХ пользователей
 */
export const basicUsersGetFull = async (): Promise<UserToOwnerResponse[]> => {
	const response = await axiosClient.get(`${apiPath}/full-list`);
	return response.data.message;
};

/**
 * Получение пользователей для transfer
 */
export const basicUsersToOwnerGet = async (params?: { search?: string; page?: number; limit?: number; all?: boolean; signal?: AbortSignal }): Promise<UsersPageResponse<UserToOwnerResponse>> => {
	const response = await axiosClient.get(`${apiPath}/owners`, {
		signal: params?.signal,
		params: {
			search: params?.search || undefined,
			page: params?.all ? undefined : params?.page,
			limit: params?.all ? undefined : params?.limit,
			all: params?.all ? true : undefined,
		},
	});
	return response.data.message;
};

/**
 * Создание пользователя
 */
export const basicUserCreate = async (payload: UserCreateRequest) => {
	const response = await axiosClient.post(`${apiPath}/create`, payload);
	toast.success(response.data.message);
};

/**
 * Обновление CVCG у дилера
 */
export const basicUserUpdateCvcg = async (cvcg: boolean): Promise<void> => {
	await axiosClient.patch(`${apiPath}/me/cvcg`, { can_view_child_groups: cvcg });
};

export const basicUserUpdateConfigurationPriority = async (enabled: boolean): Promise<void> => {
	await axiosClient.patch(`${apiPath}/me/configuration-priority`, { force_neosync_configuration_priority: enabled });
};

/**
 * Удаление пользователя
 */
export const basicUserDelete = async (uuid: string) => {
	const response = await axiosClient.delete(`${apiPath}/${uuid}`);
	toast.success(response.data.message);
};

/**
 * Удаление пользователей
 */
export const basicUserMassiveDelete = async (uuids: string[]) => {
	await runGuardedAction(
		ACTIONKEYs.USER_DELETE_MASSIVE,
		async () => {
			const response = await axiosClient.post(`${apiPath}/delete`, { uuids }, { actionKey: ACTIONKEYs.USER_DELETE_MASSIVE } as any);
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
 * Отправка на почту данные для входа пользователя
 */
export const basicUserSendEmail = async (uuid: string, turnstileToken: string, language: string) => {
	const response = await axiosClient.post(`${apiPath}/${uuid}/send/email`, {
		turnstile_token: turnstileToken,
		language,
	});

	toast.success(response.data.message);
};
