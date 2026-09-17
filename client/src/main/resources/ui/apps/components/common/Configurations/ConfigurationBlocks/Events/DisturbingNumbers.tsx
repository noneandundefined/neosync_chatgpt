import { UIDs } from '@/constants/UID.constant';
import CfgPhonesInput from '@/components/ui/Input/CfgPhones';

const DisturbingNumbers = () => {
	return <CfgPhonesInput uid={UIDs.ALARM_PHONE} placeholder="+7xxxxxxxxxx" />;
};

export default DisturbingNumbers;
