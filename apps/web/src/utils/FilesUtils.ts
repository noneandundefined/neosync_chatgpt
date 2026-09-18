class FileHelper {
	private readonly BIN_MAX_SIZE_2_KB = 2 * 1024;
	private readonly BIN_MAX_SIZE_5_KB = 5 * 1024;
	private readonly BIN_MAX_SIZE_10_KB = 10 * 1024;
	private readonly BIN_MAX_SIZE_1_MB = 1 * 1024 * 1024;

	public fileMax2Kb(file: File): boolean {
		return file.size <= this.BIN_MAX_SIZE_2_KB;
	}

	public fileMax5Kb(file: File): boolean {
		return file.size <= this.BIN_MAX_SIZE_5_KB;
	}

	public fileMax10Kb(file: File): boolean {
		return file.size <= this.BIN_MAX_SIZE_10_KB;
	}

	public fileMax1Mb(file: File): boolean {
		return file.size <= this.BIN_MAX_SIZE_1_MB;
	}

	public isBin(file: File): Promise<boolean> {
		return new Promise((resolve) => {
			const reader = new FileReader();

			const blob = file.slice(0, 2048);

			reader.onload = () => {
				const result = reader.result as ArrayBuffer;
				const bytes = new Uint8Array(result);

				let nonPrintableCount = 0;
				const threshold = 0.115;

				for (const byte of bytes) {
					if (byte < 9 || (byte > 13 && byte < 32) || byte > 126) {
						nonPrintableCount++;
					}
				}

				const ratio = nonPrintableCount / bytes.length;
				resolve(ratio > threshold);
			};

			reader.onerror = () => resolve(false);

			reader.readAsArrayBuffer(blob);
		});
	}
}

export default new FileHelper();
