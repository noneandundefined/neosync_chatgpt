import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { TRACK_ALTERNATIVE_MODE } from '@/constants/Track.constants';

const SubstitutionCoordinatesAltSource = () => {
	return (
		<div className="space-y-3">
			<LabeledField label="message.alternative-mode">
				<CFGSelect uid={UIDs.COORDINATE_SUBSTITUTION_ALTERNATIVE} array={TRACK_ALTERNATIVE_MODE} />
			</LabeledField>

			<LabeledField label="message.coordinate-validity-area-multiplier">
				<CFGInput uid={UIDs.MULTIPLIER_COORDINATE_VALIDITY_AREA} />
			</LabeledField>
		</div>
	);
};

export default SubstitutionCoordinatesAltSource;
