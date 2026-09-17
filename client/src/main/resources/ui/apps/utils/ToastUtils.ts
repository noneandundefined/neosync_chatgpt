import { Id, toast } from 'react-toastify';

const MAX_ACTIVE_SERVER_TOASTS = 2;
const DEDUPE_WINDOW_MS = 5000;
const BURST_WINDOW_MS = 10000;
const BURST_LIMIT = 2;

const activeServerToastIds = new Set<Id>();
const lastServerMessageAt = new Map<string, number>();
const serverToastHistory: number[] = [];

const normalizeMessage = (message: string) => message.trim().toLowerCase();

const cleanupHistory = (now: number) => {
	while (serverToastHistory.length > 0 && now - serverToastHistory[0] > BURST_WINDOW_MS) {
		serverToastHistory.shift();
	}
};

const canShowServerToast = (message: string) => {
	const now = Date.now();
	const normalized = normalizeMessage(message);
	const lastShownAt = lastServerMessageAt.get(normalized);

	cleanupHistory(now);

	if (activeServerToastIds.size >= MAX_ACTIVE_SERVER_TOASTS) {
		return false;
	}

	if (lastShownAt && now - lastShownAt < DEDUPE_WINDOW_MS) {
		return false;
	}

	if (serverToastHistory.length >= BURST_LIMIT) {
		return false;
	}

	lastServerMessageAt.set(normalized, now);
	serverToastHistory.push(now);
	return true;
};

export const showServerErrorToast = (message?: string) => {
	if (!message || !canShowServerToast(message)) {
		toast.clearWaitingQueue();
		return false;
	}

	const toastId = toast.error(message, {
		onClose: () => {
			activeServerToastIds.delete(toastId);
		},
	});

	activeServerToastIds.add(toastId);
	toast.clearWaitingQueue();

	return true;
};

export const showServerSuccessToast = (message?: string) => {
	if (!message || !canShowServerToast(message)) {
		toast.clearWaitingQueue();
		return false;
	}

	const toastId = toast.success(message, {
		onClose: () => {
			activeServerToastIds.delete(toastId);
		},
	});

	activeServerToastIds.add(toastId);
	toast.clearWaitingQueue();

	return true;
};
