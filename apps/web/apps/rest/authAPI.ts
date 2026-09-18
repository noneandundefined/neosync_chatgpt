import i18next from 'i18next';
import axiosClient, { checkAuthSession } from './axios';
import { toast } from 'react-toastify';
import { ROUTES } from '@/constants/constants';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import { runSpamGuardedAction } from '@/utils/SpamGuardUtils';
import { SignInRequestWithToken } from '@/interface/auth/signInRequest.interface';

const apiPath = '/auth';

export interface SignInResponse {
	login: string;
	role_code: string;
}

/**
 * Вход пользователя
 */
export const basicAuthSignIn = async (payload: SignInRequestWithToken): Promise<SignInResponse> => {
	return await runSpamGuardedAction(
		`auth:signin:${payload.login}:${payload.password}`,
		async () => {
			const response = await axiosClient.post(`${apiPath}/signin`, payload);
			localStorage.setItem(CACHEKEYs.USER_ACCESS, 'true');

			toast.success(i18next.t('message.auth-success'));

			return response.data.message;
		},
		{
			maxAttempts: 6,
			windowMs: 10_000,
			blockMs: 60_000,
			onBlocked: ({ message }) => {
				showServerErrorToast(message);
			},
		}
	);
};

/**
 * Проверка аунтетифицированного пользователя
 */
export const basicAuthCheck = async (): Promise<boolean> => {
	return checkAuthSession();
};

/**
 * Выход пользователя
 */
export const basicAuthSignOut = async (toastShow: boolean) => {
	const response = await axiosClient.post(`${apiPath}/signout`);
	localStorage.setItem(CACHEKEYs.USER_ACCESS, 'false');

	if (toastShow) {
		toast.success(response.data.message);
	}

	/** Clear local values caches */
	Object.keys(localStorage)
		.filter((key) => key.startsWith('neosync:table'))
		.forEach((key) => localStorage.removeItem(key));

	window.location.replace(ROUTES.SIGNIN);
};

/**
 * Запрос на сброс пароля
 */
export const basicRequestPasswordReset = async (email: string) => {
	await runSpamGuardedAction(
		`auth:password-reset-req:${email}`,
		async () => {
			const response = await axiosClient.post(`${apiPath}/password/reset/request?email=${email}`);
			toast.success(response.data.message);
		},
		{
			maxAttempts: 8,
			windowMs: 10_000,
			blockMs: 60_000,
			onBlocked: ({ message }) => {
				showServerErrorToast(message);
			},
		}
	);
};

/**
 * Сброс пароля
 */
export const basicPasswordReset = async (uuid: string, exp: string, sig: string, password: string, turnstile_token: string) => {
	await runSpamGuardedAction(
		`auth:password-confirm:${uuid}:${sig}:${password}`,
		async () => {
			const response = await axiosClient.post(`${apiPath}/password/reset/confirm?uuid=${uuid}&exp=${exp}&sig=${sig}`, {
				password: password,
				turnstile_token: turnstile_token,
			});

			showServerErrorToast(response.data.message);
		},
		{
			maxAttempts: 6,
			windowMs: 10_000,
			blockMs: 60_000,
			onBlocked: ({ message }) => {
				toast.error(message);
			},
		}
	);
};
