import SHA256 from 'crypto-js/sha256';
import { enc } from 'crypto-js';

export const getHashHex = (data: string) => {
	const json = JSON.stringify(data);
	return SHA256(json).toString(enc.Hex);
};
