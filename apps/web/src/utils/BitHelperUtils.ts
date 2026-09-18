export class BitHelperUtils {
	/** Установить бит (1) */
	static setBit(mask: number, bit: number): number {
		return mask | (1 << bit);
	}

	/** Сбросить бит (0) */
	static clearBit(mask: number, bit: number): number {
		return mask & ~(1 << bit);
	}

	/** Инвертировать бит */
	static invertBit(mask: number, bit: number): number {
		return mask ^ (1 << bit);
	}

	/** Изменить boolean бит */
	static toggleBit(mask: number, bit: number, enabled: boolean): number {
		return enabled ? this.setBit(mask, bit) : this.clearBit(mask, bit);
	}

	static toggleBits(mask: number, bits: number[], enabled: boolean): number {
		return bits.reduce((acc, bit) => this.toggleBit(acc, bit, enabled), mask);
	}

	/** Проверить бит (возвращает 0 или 1) */
	static checkBit(mask: number, bit: number): boolean {
		return (mask & (1 << bit)) !== 0;
	}

	/** Принудительно изменить бит */
	static changeBit(mask: number, bit: number, value: 0 | 1): number {
		return (~(1 << bit) & mask) | (value << bit);
	}

	/** Установить несколько бит (маской) */
	static setMask(mask: number, bit: number): number {
		return bit | mask;
	}

	/** Очистить несколько бит (маской) */
	static clearMask(mask: number, bit: number): number {
		return bit & ~mask;
	}
}
