import { UIDs } from '@/constants/UID.constant';
import CfgPhonesInput from '@/components/ui/Input/CfgPhones';

const AuthorizedPhones = () => {
	return (
		<div className="min-w-auto sm:min-w-[23rem]">
			<CfgPhonesInput uid={UIDs.AUTH_PHONE} placeholder="+7xxxxxxxxxx" />
		</div>
	);
};

export default AuthorizedPhones;
