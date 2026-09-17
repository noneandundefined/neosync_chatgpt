import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { SIM_PRIORITY } from '@/constants/Sim.constant';
import CFGSelect from '@/components/ui/Select/CfgSelect';

const SimPriority = () => {
	const { t } = useTranslation();

	return (
		<>
			<div className="min-w-auto sm:min-w-[22rem]">
				<div className="flex flex-col gap-3">
					<p>{t('label.priority')} SIM:</p>
					<CFGSelect uid={UIDs.SIM_STATUS} array={SIM_PRIORITY} />
				</div>
			</div>
		</>
	);
};

export default SimPriority;
