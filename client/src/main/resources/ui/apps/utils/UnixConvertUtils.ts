export const convertUnixTimestampToPrettyFormat = (unixTimestamp: number): string => {
	const date = new Date(unixTimestamp * 1000);

	const dates = (value: number) => value.toString().padStart(2, '0');

	const day = dates(date.getDate());
	const month = dates(date.getMonth() + 1);
	const year = date.getFullYear();

	const hours = dates(date.getHours());
	const minutes = dates(date.getMinutes());
	const seconds = dates(date.getSeconds());

	return `${day}.${month}.${year} ${hours}:${minutes}:${seconds}`;
};

export const convertNumberToPrettyHexString = (int64: number): string => {
	if (int64 < 0) {
		int64 += Math.pow(2, 32);
	}

	return int64.toString(16).toUpperCase().padStart(8, '0');
};
