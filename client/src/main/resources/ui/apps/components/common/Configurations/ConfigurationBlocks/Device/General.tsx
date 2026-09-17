import React from 'react';
import Modal from '@/components/Modal/Modal';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import OpenInNew from '@/components/@icons/open-in-new';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import { useModalContext } from '@/context/useModalContext';
import LabeledField from '@/components/ui/Form/LabeledField';
import CfgInputPassword from '@/components/ui/Input/CfgInputPassword';
import ModalDevicePassword from '@/components/Modal/ModalDevicePassword';
import { useOptionalConfigurationContext } from '@/context/useConfigurationContext';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { BEACON_MODE_ARRAY_BASE, BEACON_MODE_ARRAY_HYBRID } from '@/constants/BeaconMode.constant';

const General = () => {
	const { t } = useTranslation();
	const device = useOptionalConfigurationContext()?.device;
	const { support, isTemplate } = useConfigurationWrapperContext();

	const { open } = useModalContext();

	return (
		<div className="space-y-4 min-w-auto sm:min-w-[20rem]">
			{!isTemplate && device && (
				<>
					<LabeledField label="label.device-name">
						<CFGInput uid={UIDs.DEVICE_NAME} />
					</LabeledField>

					<LabeledField label="label.password">
						<CfgInputPassword uid={UIDs.AUTH_GLOBAL_PASS} className="w-full" disabled />

						<span
							className="inline w-[fit-content] border-b border-b-[transparent] hover:border-b-[#1d3c5d] cursor-pointer"
							onClick={() => {
								open(
									<Modal title={t('label.new-password')}>
										<ModalDevicePassword device={device} />
									</Modal>
								);
							}}
						>
							<div className="inline-flex gap-1 items-center">
								<p className="text-sm text-[#1d3c5d]">{t('label.change-password')}</p>
								<OpenInNew fill="#1d3c5d" size={13} className="mt-[2px]" />
							</div>
						</span>
					</LabeledField>
				</>
			)}

			<LabeledField label="label.device-mode">
				{/* Check mode ADMP50\Anothers */}
				{support.ModeHybrid ? (
					<CFGSelect uid={UIDs.DEVICE_MODE} array={BEACON_MODE_ARRAY_HYBRID} />
				) : (
					<CFGSelect uid={UIDs.DEVICE_FUNCTIONS} bit={BITs.DEVICE_FUNCTIONS.DEVICE_MODE_TRACKER} array={BEACON_MODE_ARRAY_BASE} />
				)}
			</LabeledField>

			{support.PowerBtn && (
				<React.Fragment>
					<div className="flex flex-row items-center justify-between">
						<label className="max-w-[60%]">{t('message.automatic-power-via-usb')}</label>
						<Tooltip title={t('message.tooltip-auto-power-on-via-usb')}>
							<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.AUTOMATIC_POWER_VIA_USB} vOff={0} vOn={1} />
						</Tooltip>
					</div>

					<div className="flex flex-row items-center justify-between">
						<label className="max-w-[60%]">{t('message.forbidding-power-power-button')}</label>
						<Tooltip title={t('message.tooltip-disable-power-button-shutdown')}>
							<CFGSwitch uid={UIDs.CUSTOM_MASK_DEFAULT_FALSE_1} bit={BITs.CUSTOM_MASK_DEFAULT_FALSE_1.FORBIDDING_POWER_POWER_BUTTON} vOff={0} vOn={1} />
						</Tooltip>
					</div>
				</React.Fragment>
			)}
		</div>
	);
};

export default General;
