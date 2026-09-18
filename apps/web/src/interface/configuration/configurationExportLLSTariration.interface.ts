interface TarirationRow {
	id: number;
	value: number;
	level: number;
}

export interface ConfigurationExportLLSTariration {
	tariration: Record<number, TarirationRow[]>;
}
