import CryptoJS from 'crypto-js';

const salt = import.meta.env.VITE_ENCRYPTION_SALT!;
const keySize = import.meta.env.VITE_ENCRYPTION_KEY_SIZE!;
const iterations = import.meta.env.VITE_ENCRYPTION_ITERATION!;
const iv = CryptoJS.enc.Hex.parse(import.meta.env.VITE_ENCRYPTION_IV!);

const getKey = () => {
	return CryptoJS.PBKDF2(salt, CryptoJS.enc.Utf8.parse(salt), {
		keySize: keySize,
		iterations: iterations,
	});
};

export const encrypt = (text: string): string => {
	const key = getKey();
	const encoded = CryptoJS.AES.encrypt(text, key, {
		iv: iv,
		mode: CryptoJS.mode.CBC,
		padding: CryptoJS.pad.Pkcs7,
	}).toString();

	return encoded.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
};

export const decrypt = (text: string): string => {
	const key = getKey();
	let base64 = text.replace(/-/g, '+').replace(/_/g, '/');
	while (base64.length % 4 !== 0) base64 += '=';

	const decoded = CryptoJS.AES.decrypt(base64, key, {
		iv: iv,
		mode: CryptoJS.mode.CBC,
		padding: CryptoJS.pad.Pkcs7,
	});
	return decoded.toString(CryptoJS.enc.Utf8);
};
