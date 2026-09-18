import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import LabeledField from '@/components/ui/Form/LabeledField';

const FilteringEmissions = () => {
	const { t } = useTranslation();

	return (
		<div className="space-y-3">
			<LabeledField label="message.discarded-points-during-ejection">
				<CFGInput uid={UIDs.DISCARDED_POINTS_DURING_EJECTION} />
			</LabeledField>

			<LabeledField label="message.distance-between-points">
				<div className="flex items-center gap-2">
					<CFGInput uid={UIDs.DISTANCE_BETWEEN_POINTS} />
					<label>{t('label.m')}</label>
				</div>
			</LabeledField>
		</div>
	);
};

export default FilteringEmissions;
