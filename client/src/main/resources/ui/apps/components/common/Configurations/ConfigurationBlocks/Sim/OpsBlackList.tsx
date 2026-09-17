import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface OpsBlackListProps {
	index: number;
}

const OpsBlackList: React.FC<OpsBlackListProps> = ({ index }) => {
	const { t } = useTranslation();
	const { manager } = useConfigurationWrapperContext();

	const opsBlackList = useConfigurationVal(manager, UIDs.OPS_BLACK_LIST) as number[];

	const sliceStart = index === 1 ? 0 : opsBlackList.length / 2;
	const sliceEnd = index === 1 ? opsBlackList.length / 2 : opsBlackList.length;

	const chunk = opsBlackList.slice(sliceStart, sliceEnd);

	return (
		<div className="space-y-1 border border-[#d1d5db]">
			<p className="uppercase text-sm bg-[#f1f2f3] p-1">{t('message.ops-black-list')}</p>

			<div className="flex flex-row gap-1 p-2">
				{chunk.map((_, idx) => (
					<CFGArrayItemInput key={sliceStart + idx} uid={UIDs.OPS_BLACK_LIST} index={sliceStart + idx} className="max-w-[5rem] max-h-[2rem]" />
				))}
			</div>
		</div>
	);
};

export default OpsBlackList;
