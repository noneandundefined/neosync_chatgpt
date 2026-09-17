import i18next from 'i18next';

type SpamState = {
	count: number;
	firstAttempt: number;
	blockedUntil: number;
};

const spamMap = new Map<string, SpamState>();

type SpamOptions = {
	maxAttempts: number;
	windowMs: number;
	blockMs: number;
	getKey?: () => string;
	onBlocked?: (params: { message: string; seconds: number }) => void;
};

export class SpamGuardError extends Error {
	secondsLeft: number;

	constructor(secondsLeft: number, message: string) {
		super(message);
		this.secondsLeft = secondsLeft;
	}
}

export const runSpamGuardedAction = async <T>(actionKey: string, action: () => Promise<T>, options: SpamOptions) => {
	const now = Date.now();
	const key = options.getKey?.() ?? actionKey;

	const state = spamMap.get(key);

	if (!state) {
		spamMap.set(key, {
			count: 1,
			firstAttempt: now,
			blockedUntil: 0,
		});
	} else {
		if (state.blockedUntil > now) {
			const seconds = Math.ceil((state.blockedUntil - now) / 1000);

			const message = i18next.t('message.spam-too-many-requests', {
				seconds,
			});

			options.onBlocked?.({ message, seconds });

			throw new SpamGuardError(seconds, message);
		}

		if (now - state.firstAttempt > options.windowMs) {
			state.count = 1;
			state.firstAttempt = now;
		} else {
			state.count++;
		}

		if (state.count > options.maxAttempts) {
			state.blockedUntil = now + options.blockMs;

			const seconds = Math.ceil(options.blockMs / 1000);

			const message = i18next.t('message.spam-too-many-requests', {
				seconds,
			});

			options.onBlocked?.({ message, seconds });

			throw new SpamGuardError(seconds, message);
		}
	}

	return await action();
};
