import { useParams } from 'react-router-dom';

export const useDecryptedImei = (): string | null => {
	const { imei } = useParams();
	if (!imei) return null;

	return imei;
};
