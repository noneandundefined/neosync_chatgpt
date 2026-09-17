import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { TRACK_NAVTIMESYNC } from '@/constants/Track.constants';

const NavTimeSync = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-3">
			<LabeledField label="label.mode">
				<CFGSelect uid={UIDs.NAVTIMESYNC} array={TRACK_NAVTIMESYNC} />
			</LabeledField>

			<LabeledField label="message.navtimesync-ntp-nmea">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.NAVTIMESYNC_MAX_TIME_DIFF_SEC} />
					<label>{t('label.s')}</label>
				</div>
			</LabeledField>
		</div>
	);
};

export default NavTimeSync;
