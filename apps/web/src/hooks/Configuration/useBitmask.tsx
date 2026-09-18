import { BitHelperUtils } from '@/utils/BitHelperUtils';
import { useEffect, useState } from 'react';

export const useBitmask = (initialMask: number = 0, maskSession?: string | number) => {
	const [mask, setMask] = useState<number>(initialMask);

	useEffect(() => {
		if (maskSession !== undefined) {
			setMask(Number(maskSession));
		} else {
			setMask(initialMask);
		}
	}, [maskSession, initialMask]);

	const isEnabled = (bit: number) => BitHelperUtils.checkBit(mask, bit);

	const toggle = (bit: number) => {
		const newMask = BitHelperUtils.invertBit(mask, bit);
		setMask(newMask);
		return newMask;
	};

	const setBit = (bit: number) => setMask(BitHelperUtils.setBit(mask, bit));
	const clearBit = (bit: number) => setMask(BitHelperUtils.clearBit(mask, bit));

	return { mask, setMask, isEnabled, toggle, setBit, clearBit };
};
