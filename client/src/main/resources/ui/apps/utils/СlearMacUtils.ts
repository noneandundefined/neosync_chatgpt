export const clearMacBlock = (arr: number[], start: number) => {
	for (let i = 0; i < 6; i++) {
		arr[start + i] = 0;
	}
};
