import { useMemo, useState } from 'react';
import OpsWhiteList from './OpsWhiteList';
import OpsBlackList from './OpsBlackList';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { arrayParseData } from '@/utils/ArrayUtils';
import IndexConfiguration from '../IndexConfiguration';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { detectOperatorId, NET_SIM, OPERATORS } from '@/constants/Sim.constant';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import ChevronUpCircleOutline from '@/components/@icons/chevron-up-circle-outline';
import { DEVICE_MODELS, getStaticSupport } from '@/constants/DeviceModels.constant';
import ChevronDownCircleOutline from '@/components/@icons/chevron-down-circle-outline';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface SimsProps {
	enabledApn: boolean;
}

export const Sims: React.FC<SimsProps> = ({ enabledApn }) => {
	const { t } = useTranslation();
	const { model, isTemplate } = useConfigurationWrapperContext();

	const count = enabledApn ? 2 : getStaticSupport((model ?? DEVICE_MODELS.ADM333V2) as keyof typeof DEVICE_MODELS, isTemplate).SimCount;

	return Array.from({ length: count }, (_, index) => <IndexConfiguration title={`${enabledApn ? t('label.apn').toUpperCase() : t('label.sim')}${index}`} content={<Sim index={index} enabledApn={enabledApn} />} key={index} />);
};

interface SimProps {
	index: number;
	enabledApn: boolean;
}

const Sim: React.FC<SimProps> = ({ index, enabledApn }) => {
	const { t } = useTranslation();
	const { imei, section, manager, support, isTemplate } = useConfigurationWrapperContext();

	const apnField = useConfigurationField(imei, section, manager, UIDs.APN_NAME);
	const loginField = useConfigurationField(imei, section, manager, UIDs.APN_USER);
	const passwordField = useConfigurationField(imei, section, manager, UIDs.APN_PASS);

	const apn = String(arrayParseData(apnField.value)[index] ?? '');
	const login = String(arrayParseData(loginField.value)[index] ?? '');
	const password = String(arrayParseData(passwordField.value)[index] ?? '');
	const operatorValue = useMemo(() => detectOperatorId(apn, login, password), [apn, login, password]);

	const [openOpsLists, setOpenOpsLists] = useState<boolean>(false);

	return (
		<div className="w-full sm:min-w-[22rem] space-y-4">
			<LabeledField label="message.operator-presets">
				<CFGSelect uid={UIDs.APN_USER} index={index} array={OPERATORS} presetType="operator" _value={operatorValue} writeFromExtra="login" />
			</LabeledField>

			<LabeledField label="APN">
				<CFGArrayItemInput uid={UIDs.APN_NAME} index={index} />
			</LabeledField>

			<LabeledField label="label.only-login">
				<CFGArrayItemInput uid={UIDs.APN_USER} index={index} />
			</LabeledField>

			<LabeledField label="label.password">
				<CFGArrayItemInput uid={UIDs.APN_PASS} index={index} />
			</LabeledField>

			<LabeledField label="PIN">
				<CFGArrayItemInput uid={UIDs.SIM_PIN} index={index} placeholder="PIN" />
			</LabeledField>

			{support.SimNet && (isTemplate || support.ModeHybrid) && (
				<LabeledField label="">
					<CFGSelect uid={UIDs.NET_MODE} index={index} array={NET_SIM} />
				</LabeledField>
			)}

			{enabledApn && (
				<LabeledField label="label.address">
					<CFGArrayItemInput uid={UIDs.SERVER_HOST} index={index} readOnly={true} />
				</LabeledField>
			)}

			{enabledApn && (
				<LabeledField label="label.port">
					<CFGArrayItemInput uid={UIDs.SERVER_PORT} index={index} readOnly={true} />
				</LabeledField>
			)}

			{/* ops lists */}
			<div className="border border-[#d1d5db]">
				<div className="flex flex-row items-center gap-1 bg-[#f1f2f3] p-1 cursor-pointer" onClick={() => setOpenOpsLists(!openOpsLists)}>
					{openOpsLists ? <ChevronUpCircleOutline fill="#333" size={19} /> : <ChevronDownCircleOutline fill="#333" size={19} />}

					<p className="text-[#333]">{t('message.ops-title-lists')}</p>
				</div>

				{openOpsLists && (
					<div className="p-2 space-y-3">
						<OpsWhiteList index={index + 1} />
						<OpsBlackList index={index + 1} />
					</div>
				)}
			</div>
		</div>
	);
};

export default Sim;
