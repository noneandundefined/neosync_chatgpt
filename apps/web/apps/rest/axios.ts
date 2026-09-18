import Qs from 'qs';
import i18next from 'i18next';
import { ROUTES } from '@/constants/constants';
import { PoWDDosDecision } from '@/utils/PowDDosUtils';
import { showServerErrorToast } from '@/utils/ToastUtils';
import { CACHEKEYs } from '@/constants/CacheKeys.constants';
import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { config as configClient } from '@/.config/config.client';

import { ClientGuardError, createIdempotencyKey } from '@/utils/RequestControlUtils';
import { trackApiOutcome, trackSearchOutcome } from '@/utils/ProductAnalytics';

interface AxiosRequestConfigWithRetry extends AxiosRequestConfig {
	_retry?: boolean;
	_powRetries?: number;
	_analyticsStartedAt?: number;
	skipErrorHandler?: boolean;
	actionKey?: string;
	skipIdempotency?: boolean;
}

/** Axios base url */
const axiosClient = axios.create({
	baseURL: configClient.type.release == 'dev' ? configClient.links.URL_BACKEND_DEV : configClient.type.release == 'prod' ? configClient.links.URL_BACKEND_PROD : configClient.links.URL_BACKEND_PREPROD,
	paramsSerializer: (params) => Qs.stringify(params, { arrayFormat: 'comma' }),
	withCredentials: true,
});

function clearAuthorization() {
	localStorage.setItem(CACHEKEYs.USER_ACCESS, 'false');
	Object.keys(localStorage)
		.filter((key) => key.startsWith('neosync:table'))
		.forEach((key) => localStorage.removeItem(key));
	if (window.location.pathname !== ROUTES.SIGNIN) window.location.replace(ROUTES.SIGNIN);
}

let authCheckInFlight: Promise<boolean> | null = null;

// One refresh at a time, including other tabs: a rotated token is single-use.
export function checkAuthSession(): Promise<boolean> {
	if (authCheckInFlight) return authCheckInFlight;
	const check = async () => {
		try {
			const response = await axios.get('/auth/check', {
				baseURL: axiosClient.defaults.baseURL,
				withCredentials: true,
			});
			return response.data.message.authenticated === true;
		} catch (error) {
			if (axios.isAxiosError(error) && error.response?.status === 401) {
				clearAuthorization();
				return false;
			}
			throw error;
		}
	};
	const run = () => (typeof navigator !== 'undefined' && navigator.locks ? navigator.locks.request('neosync-auth-refresh', check) : check());
	authCheckInFlight = run().finally(() => {
		authCheckInFlight = null;
	});
	return authCheckInFlight;
}

const capitalize = (str: string | undefined) => (str ? str.charAt(0).toUpperCase() + str.slice(1) : str);
const isMutationMethod = (method?: string) => ['post', 'put', 'patch', 'delete'].includes((method || '').toLowerCase());

const CIRCUIT_OPEN_MS = 15000;
const CIRCUIT_FAIL_THRESHOLD = 15;
const circuitState = {
	consecutiveFailures: 0,
	openUntil: 0,
	toastShownForOpenUntil: 0,
};

/** Axios helper for request */
axiosClient.interceptors.request.use(
	(config) => {
		const now = Date.now();
		(config as AxiosRequestConfigWithRetry)._analyticsStartedAt = now;

		if (isMutationMethod(config.method) && now < circuitState.openUntil) {
			if (circuitState.toastShownForOpenUntil !== circuitState.openUntil) {
				circuitState.toastShownForOpenUntil = circuitState.openUntil;

				showServerErrorToast(i18next.t('message.client-circuit-open'));
			}

			throw new ClientGuardError('ERR_CLIENT_CIRCUIT', i18next.t('message.client-circuit-open'));
		}

		/** Language header */
		const currentLanguage = i18next.language || 'en';
		config.headers['Accept-Language'] = currentLanguage;

		/** RequestId header */
		const requestID = localStorage.getItem(CACHEKEYs.NEOSYNC_X_REQ_ID);
		if (requestID) {
			config.headers['X-Request-ID'] = requestID;
		}

		/** Idempotency */
		const typedConfig = config as AxiosRequestConfigWithRetry;
		if (isMutationMethod(config.method) && !typedConfig.skipIdempotency) {
			const actionKey = typedConfig.actionKey || `${(config.method || 'post').toLowerCase()}:${config.url || 'unknown'}`;
			config.headers['X-Idempotency-Key'] = createIdempotencyKey(actionKey);
		}

		return config;
	},
	(error) => Promise.reject(error)
);

/** Axios helper for response */
axiosClient.interceptors.response.use(
	(response) => {
		const request = response.config as AxiosRequestConfigWithRetry;
		trackApiOutcome((request.method || 'get').toLowerCase(), request.url || '', true, Date.now() - (request._analyticsStartedAt || Date.now()), response.status);
		const message = response.data?.message;
		const resultCount = typeof message?.total === 'number' ? message.total : Array.isArray(message?.items) ? message.items.length : Array.isArray(message) ? message.length : -1;
		if (resultCount >= 0) trackSearchOutcome(request.url || '', resultCount);
		circuitState.consecutiveFailures = 0;

		const requestID = response.headers['x-request-id'];
		if (requestID) {
			localStorage.setItem(CACHEKEYs.NEOSYNC_X_REQ_ID, requestID);
		}

		return response;
	},
	async (error) => {
		if (error instanceof ClientGuardError) {
			return Promise.reject(error);
		}

		if (axios.isCancel(error) || error.code === 'ERR_CANCELED') {
			return Promise.reject(error);
		}

		/** Internet */
		if (error.message === 'Network Error' || error.code === 'ERR_NETWORK' || error.message.includes('Network request failed')) {
			showServerErrorToast(capitalize(i18next.t('message.internet-error')));
		}

		if (error instanceof AxiosError) {
			const failedRequest = error.config as AxiosRequestConfigWithRetry | undefined;
			if (failedRequest) {
				trackApiOutcome(
					(failedRequest.method || 'get').toLowerCase(),
					failedRequest.url || '',
					false,
					Date.now() - (failedRequest._analyticsStartedAt || Date.now()),
					error.response?.status,
					String(error.response?.data?.code || error.code || 'request_failed').slice(0, 128)
				);
			}
			if (error.response?.status === 429 || (error.response?.status && error.response.status >= 500)) {
				circuitState.consecutiveFailures += 1;

				if (circuitState.consecutiveFailures >= CIRCUIT_FAIL_THRESHOLD) {
					circuitState.openUntil = Date.now() + CIRCUIT_OPEN_MS;
				}
			} else {
				circuitState.consecutiveFailures = 0;
			}

			if (error?.response?.status == 429 && error.response.data?.message?.challenge && (failedRequest?._powRetries || 0) < 2) {
				const { challenge, difficulty } = error.response.data.message;

				const nonce = await PoWDDosDecision(challenge, difficulty);

				const config = {
					...error.config,
					_powRetries: (failedRequest?._powRetries || 0) + 1,
					headers: {
						...error.config?.headers,
						'Pow-Challenge': challenge,
						'Pow-Nonce': nonce,
					},
				};

				return axiosClient(config);
			}

			const originalRequest = error.config as AxiosRequestConfigWithRetry | undefined;
			if (error.response?.status === 401) {
				if (error.response.data?.message?.code === 'auth_refresh_required' && originalRequest && !originalRequest._retry) {
					originalRequest._retry = true;
					if (await checkAuthSession()) return axiosClient(originalRequest);
				} else {
					clearAuthorization();
				}
				return Promise.reject(error);
			}
			if (originalRequest?.skipErrorHandler) return Promise.reject(error);

			if (error.response && error.response.status) {
				showServerErrorToast(capitalize(error.response.data?.message) || capitalize(i18next.t('message.server-error')));
			} else if (error.request) {
				showServerErrorToast(capitalize(i18next.t('message.request-error')));
			} else {
				showServerErrorToast(capitalize(i18next.t('message.unknown-request-error')));
			}
		} else {
			showServerErrorToast(capitalize(i18next.t('message.try-again-error')));
		}

		return Promise.reject(error);
	}
);

export default axiosClient;
