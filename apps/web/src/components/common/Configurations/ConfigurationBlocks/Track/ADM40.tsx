import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const EMPTY_ADDR = '000000000000';
const ADM40_ADDR = '001B446789AB';

const ADM40 = () => {
	const { imei, section, manager } = useConfigurationWrapperContext();

	/** Values for ADM40 ADDR */
	const { value: adm40Addr, handleChange: adm40AddrChange, saveDraft: adm40AddrSave } = useConfigurationField(imei, section, manager, UIDs.ADM40_ADDR);

	/** Enabled ADM40 */
	const enabled = adm40Addr !== EMPTY_ADDR;

	return (
		<div className="space-y-5">
			<LabeledField label="message.use-adm40" className="flex items-center justify-between">
				<GUISwitch
					onChange={async (e) => {
						const checked = e.target.checked;

						if (!checked) {
							adm40AddrChange(EMPTY_ADDR);
							await adm40AddrSave(EMPTY_ADDR);
						} else {
							adm40AddrChange(ADM40_ADDR);
							await adm40AddrSave(ADM40_ADDR);
						}
					}}
					checked={enabled}
				/>
			</LabeledField>

			{enabled && (
				<LabeledField label="label.mac-address">
					<CFGInput uid={UIDs.ADM40_ADDR} />
				</LabeledField>
			)}
		</div>
	);
};

export default ADM40;
