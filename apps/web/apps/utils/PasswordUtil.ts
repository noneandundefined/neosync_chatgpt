export function generatePassword(length = 6): string {
	const letters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
	const digits = '0123456789';
	const all = letters + digits;

	let password = '';

	password += letters.charAt(Math.floor(Math.random() * letters.length));
	password += digits.charAt(Math.floor(Math.random() * digits.length));

	for (let i = 2; i < length; i++) {
		password += all.charAt(Math.floor(Math.random() * all.length));
	}

	return shuffle(password.split('')).join('');
}

function shuffle<T>(arr: T[]): T[] {
	for (let i = arr.length - 1; i > 0; i--) {
		const j = Math.floor(Math.random() * (i + 1));
		[arr[i], arr[j]] = [arr[j], arr[i]];
	}

	return arr;
}
