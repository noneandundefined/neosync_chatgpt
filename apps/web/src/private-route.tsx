import { CACHEKEYs } from './constants/CacheKeys.constants';
import { basicAuthCheck } from './rest/authAPI';

export async function getAuthState() {
	const cached = localStorage.getItem(CACHEKEYs.USER_ACCESS);
	if (cached !== null) {
		return cached === 'true';
	}

	try {
		const logged = await basicAuthCheck();
		localStorage.setItem(CACHEKEYs.USER_ACCESS, logged ? 'true' : 'false');
		return logged;
	} catch {
		return false;
	}
}
