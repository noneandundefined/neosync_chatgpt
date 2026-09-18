import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import { useBitmask } from '@/hooks/Configuration/useBitmask';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { DYNAMIC_BLACK_BOX, getDbbChunks } from '@/constants/DynamicBlackBox.constant';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const TransmittedDataBlocks = () => {
	const { t } = useTranslation();

	const { model, manager, imei, section, isTemplate } = useConfigurationWrapperContext();

	const { value: maskSession, handleChange: handleBlackBoxStructureChange, saveDraft: saveBlackBoxStructureDraft } = useConfigurationField(imei, section, manager, UIDs.BLACK_BOX_STRUCTURE);
	if (maskSession == null) return;

	const chunks = getDbbChunks(model, isTemplate);

	const { isEnabled, toggle } = useBitmask(chunks.NAVIGATION_CHUNK, maskSession);

	const handleToggle = async (bit: number, isNavigation: boolean) => {
		if (isNavigation) return;

		let newMask = toggle(bit);

		const modelChunks = getDbbChunks(model, isTemplate);
		const LBS_BIT = modelChunks.LBS_CHUNK;

		if ('ANALOG_CHUNK' in modelChunks && bit === LBS_BIT && BitHelperUtils.checkBit(newMask, LBS_BIT)) {
			newMask = BitHelperUtils.setBit(newMask, modelChunks.ANALOG_CHUNK);
		}

		handleBlackBoxStructureChange(newMask);
		await saveBlackBoxStructureDraft(newMask);
	};

	return (
		<>
			<div>
				<p className="font-medium">{t('message.dynamic-black-box')}</p>
				<p className="max-w-[42rem] text-[13px] my-4">{t('message.dynamic-black-box-description')}</p>

				<div className="max-w-full">
					<div>
						{DYNAMIC_BLACK_BOX.filter((item) => item.state in chunks).map((item, index) => {
							const bit = chunks[item.state as keyof typeof chunks];
							const isNavigation = item.state === 'NAVIGATION_CHUNK';
							const isChecked = isNavigation || isEnabled(bit);

							return (
								<div key={index} className="flex justify-between items-center my-4">
									<p className="font-medium">{t(item.title)}</p>
									<GUISwitch checked={isChecked} disabled={isNavigation} onChange={() => handleToggle(bit, isNavigation)} />
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</>
	);
};

export default TransmittedDataBlocks;
