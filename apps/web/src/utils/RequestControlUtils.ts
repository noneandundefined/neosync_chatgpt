import i18next from 'i18next';
import { generateUUID } from './UuidUtils';

const inFlightActions = new Set<string>();
const cooldownUntil = new Map<string, number>();
const requestIdempotency = new Map<string, string>();

const defaultCooldownMs = 1500;

export type GuardedActionOptions = {
	cooldownMs?: number;
	getCooldownMessage?: (seconds: number) => string;
	onBlocked?: (message: string) => void;
};

export class ClientGuardError extends Error {
	code: 'ERR_CLIENT_INFLIGHT' | 'ERR_CLIENT_COOLDOWN' | 'ERR_CLIENT_CIRCUIT';

	constructor(code: 'ERR_CLIENT_INFLIGHT' | 'ERR_CLIENT_COOLDOWN' | 'ERR_CLIENT_CIRCUIT', message: string) {
		super(message);
		this.code = code;
	}
}

export const createIdempotencyKey = (actionKey: string) => {
	const cached = requestIdempotency.get(actionKey);
	if (cached) return cached;

	const generated = `${actionKey}:${generateUUID()}`;

	requestIdempotency.set(actionKey, generated);
	return generated;
};

export const clearIdempotencyKey = (actionKey: string) => {
	requestIdempotency.delete(actionKey);
};

export const runGuardedAction = async <T>(actionKey: string, action: () => Promise<T>, options: GuardedActionOptions = {}) => {
	const now = Date.now();
	const cooldownMs = options.cooldownMs ?? defaultCooldownMs;
	const blockedUntil = cooldownUntil.get(actionKey) ?? 0;

	if (inFlightActions.has(actionKey)) {
		throw new ClientGuardError('ERR_CLIENT_INFLIGHT', i18next.t('message.client-action-in-progress'));
	}

	if (now < blockedUntil) {
		const leftMs = blockedUntil - now;
		const leftSeconds = Math.max(1, Math.ceil(leftMs / 1000));
		const message = options.getCooldownMessage ? options.getCooldownMessage(leftSeconds) : i18next.t('message.client-action-please-wait-seconds', { seconds: leftSeconds });

		options.onBlocked?.(message);

		throw new ClientGuardError('ERR_CLIENT_COOLDOWN', message);
	}

	inFlightActions.add(actionKey);

	try {
		return await action();
	} finally {
		inFlightActions.delete(actionKey);
		cooldownUntil.set(actionKey, Date.now() + cooldownMs);
		clearIdempotencyKey(actionKey);
	}
};
