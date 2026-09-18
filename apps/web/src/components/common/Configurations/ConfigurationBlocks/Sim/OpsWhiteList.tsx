import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationVal } from '@/hooks/Configuration/useConfigurationVal';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface OpsWhiteListProps {
	index: number; // 1 OR 2
}

const OpsWhiteList: React.FC<OpsWhiteListProps> = ({ index }) => {
	const { t } = useTranslation();
	const { manager } = useConfigurationWrapperContext();

	const opsWhiteList = useConfigurationVal(manager, UIDs.OPS_WHITE_LIST_2) as number[];

	const sliceStart = index === 1 ? 0 : opsWhiteList.length / 2;
	const sliceEnd = index === 1 ? opsWhiteList.length / 2 : opsWhiteList.length;

	const chunk = opsWhiteList.slice(sliceStart, sliceEnd);

	return (
		<div className="space-y-1 border border-[#d1d5db]">
			<p className="uppercase text-sm bg-[#f1f2f3] p-1">{t('message.ops-white-list')}</p>

			<div className="flex flex-wrap sm:flex-nowrap justify-start sm:justify-between gap-1 p-2">
				{chunk.map((_, idx) => (
					<CFGArrayItemInput key={sliceStart + idx} uid={UIDs.OPS_WHITE_LIST_2} index={sliceStart + idx} className="max-w-[3rem] sm:max-w-full max-h-[2rem]" />
				))}
			</div>
		</div>
	);
};

export default OpsWhiteList;
