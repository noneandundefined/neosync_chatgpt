import Hex from 'crypto-js/enc-hex';
import SHA256 from 'crypto-js/sha256';

const wordArrayToUint8Array = (wordArray: any): Uint8Array => {
	const hex = wordArray.toString(Hex);
	const bytes = new Uint8Array(hex.length / 2);

	for (let i = 0; i < bytes.length; i++) {
		bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
	}

	return bytes;
};

const sha256 = (input: string): Uint8Array => {
	const hashWordArray = SHA256(input);
	return wordArrayToUint8Array(hashWordArray);
};

const hashMatchesDifficulty = async (input: string, difficulty: number): Promise<boolean> => {
	const hash = sha256(input);

	for (let i = 0; i < difficulty; i++) {
		const byteIdx = Math.floor(i / 2);
		const half = i % 2;
		const val = hash[byteIdx];

		if (half === 0 && val >> 4 !== 0) return false;
		if (half === 1 && (val & 0x0f) !== 0) return false;
	}

	return true;
};

export const PoWDDosDecision = async (challenge: string, difficulty: number): Promise<string> => {
	let nonce = 0;

	while (true) {
		const input = challenge + nonce;
		if (await hashMatchesDifficulty(input, difficulty)) {
			return nonce.toString();
		}

		nonce++;

		if (nonce % 1000 === 0) await new Promise((r) => setTimeout(r, 0));
	}
};
