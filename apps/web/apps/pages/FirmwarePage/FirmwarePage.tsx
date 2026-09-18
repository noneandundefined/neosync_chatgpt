import { useState } from 'react';
import PageLayout from '../PageLayout';
import FirmwaresList from '@/components/common/Firmwares/FirmwaresList';

const FirmwarePage = () => {
	const [fotaSysToken, setFotaSysToken] = useState<number>(0);

	return (
		<PageLayout>
			<FirmwaresList sysToken={fotaSysToken} setSysToken={setFotaSysToken} />
		</PageLayout>
	);
};

export default FirmwarePage;
