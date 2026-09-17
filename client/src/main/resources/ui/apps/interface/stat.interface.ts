export interface StatInfo {
	sim: string;
	server0: string;
	server1: string;
	lte: string;
	gsmTime: string;
}

export const parseStatInfo = (statHex: string | number | null): StatInfo | null => {
	if (!statHex) return null;

	const stat = typeof statHex === 'string' ? parseInt(statHex, 16) : statHex;

	return {
		sim: stat & 0x0002 ? 'SIM1' : 'SIM0',
		server0: stat & 0x0004 ? 'Отключен' : 'Пакет с данными успешно отправлен',
		server1: stat & 0x0004 ? 'Отключен' : 'Пакет с данными успешно отправлен',
		lte: stat & 0x0800 ? 'Да' : 'Нет',
		gsmTime: stat & 0x4000 ? 'Синхронизировано' : 'Не синхронизировано',
	};
};
