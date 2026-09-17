import React from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import { BitHelperUtils } from '@/utils/BitHelperUtils';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import GUISelect from '@/components/ui/Select/GUISelect';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { TelemetryResponse } from '@/interface/configuration/configurationTelemetryResponse.interface';
import { ADM20_MODE, ADM20_OUTPUT_ACTIVATION_MODE, ADM20_OUTPUT_STATE } from '@/constants/ADM20.constant';

interface ADM20ReaderProps {
	telemetry: TelemetryResponse | null;
}

const ADM20Reader: React.FC<ADM20ReaderProps> = ({ telemetry }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	/** Telemetry */
	const getValueOrNoData = (v: any) => (v === 0 ? 0 : (v ?? '-'));
	const cardUid = telemetry ? getValueOrNoData(telemetry?.adm20info?.card_uid) : '-';
	const tagId = telemetry ? getValueOrNoData(telemetry?.adm20info?.tag_id) : '-';
	const tagSn = telemetry ? getValueOrNoData(telemetry?.adm20info?.tag_sn) : '-';
	const tagRssi = telemetry ? getValueOrNoData(telemetry?.adm20info?.tag_rssi) : '-';
	const tagAdc = telemetry ? getValueOrNoData(telemetry?.adm20info?.tag_adc) : '-';
	const tagPeriod = telemetry ? getValueOrNoData(telemetry?.adm20info?.tag_period) : '-';
	/** Telemetry */

	/** Field. Mode */
	const { value: adm20Mode = 0, handleChange: handleAdm20ModeChange, saveDraft: saveAdm20ModeDraft } = useConfigurationField(imei, section, manager, UIDs.ADM20_MODE);
	/** Field. Adm20IdMode */
	const { value: adm20IdMode, handleChange: handleAdm20IdModeChange, saveDraft: saveAdm20IdModeDraft } = useConfigurationField(imei, section, manager, UIDs.ADM20_ID_MODE);
	/** Field. Address */
	const { value: adm20Addr, handleChange: handleAdm20AddrChange, saveDraft: saveAdm20AddrDraft } = useConfigurationField(imei, section, manager, UIDs.ADM20_ADDR);

	const cardDetectionChecked = BitHelperUtils.checkBit(adm20Mode, BITs.ADM20_MODE.ADM20_PERIODIC_RFID) || BitHelperUtils.checkBit(adm20Mode, BITs.ADM20_MODE.ADM20_GUARANTEED_RFID);
	const tagDetectionChecked = BitHelperUtils.checkBit(adm20Mode, 2);

	const handleBitCheckboxChange = async (checked: boolean, bits: number[]) => {
		let newMask = adm20Mode;
		bits.forEach((bit) => {
			if (checked) newMask |= 1 << bit;
			else newMask &= ~(1 << bit);
		});

		handleAdm20ModeChange(newMask);
		await saveAdm20ModeDraft(newMask);
	};

	const handleSelectModeChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
		const selected = Number(e.target.value);
		let newMask = adm20Mode;

		newMask &= ~((1 << 0) | (1 << 1));

		if (selected === 0) newMask |= 1 << 0;
		if (selected === 1) newMask |= 1 << 1;

		handleAdm20ModeChange(newMask);
		await saveAdm20ModeDraft(newMask);
	};

	const isAddrEnabled = adm20Addr !== 255;

	const handleAddrCheckboxChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
		const newValue = e.target.checked ? 0 : 255;

		handleAdm20AddrChange(newValue);
		await saveAdm20AddrDraft(newValue);
	};

	const handleCheckboxChange = async (checked: boolean, handleChangeFn: (v: number) => void, saveFn: (v: number) => Promise<void>) => {
		const value = checked ? 1 : 0;

		handleChangeFn(value);
		await saveFn(value);
	};

	const currentSelectMode = BitHelperUtils.checkBit(adm20Mode, 1) ? 1 : 0;

	return (
		<div className="flex flex-col gap-4 w-full">
			<div className="flex flex-col sm:flex-row sm:items-center gap-[2rem]">
				<div className="flex items-center justify-between sm:justify-start gap-4">
					<p className="font-medium">{t('label.adm20-reader')}</p>
					<GUISwitch checked={isAddrEnabled} onChange={handleAddrCheckboxChange} />
				</div>

				{isAddrEnabled && (
					<LabeledField label="label.address" className="flex items-center flex-row gap-3">
						<Tooltip title={t('message.tooltip-rs485-reader-address')}>
							<CFGInput uid={UIDs.ADM20_ADDR} disabled={!isAddrEnabled} />
						</Tooltip>
					</LabeledField>
				)}
			</div>

			{isAddrEnabled && (
				<React.Fragment>
					<div className="flex flex-col md:flex-row items-stretch gap-4 border border-[#d9d9d9] p-3">
						<div className="space-y-6 border border-[#d9d9d9] flex-1 p-3">
							<LabeledField label="message.adm20-card-detection" className="flex items-center flex-row gap-3">
								<GUISwitch checked={cardDetectionChecked} onChange={(e) => handleBitCheckboxChange(e.target.checked, [0, 1])} />
							</LabeledField>

							{cardDetectionChecked && (
								<LabeledField label="label.mode" className="flex flex-col md:items-center md:flex-row gap-x-5">
									<GUISelect value={currentSelectMode} onChange={handleSelectModeChange}>
										{ADM20_MODE.map((m) => (
											<option key={m.id} value={m.id}>
												{t(m.value)}
											</option>
										))}
									</GUISelect>
								</LabeledField>
							)}
						</div>

						<div className="space-y-6 border border-[#d9d9d9] p-3">
							<LabeledField label="message.adm20-tag-detection" className="flex items-center flex-row gap-3">
								<GUISwitch checked={tagDetectionChecked} onChange={(e) => handleBitCheckboxChange(e.target.checked, [2])} />
							</LabeledField>

							{tagDetectionChecked && (
								<LabeledField label="message.adm20-transmit-sn-uid" argv={':'} className="flex items-center flex-row gap-3">
									<GUISwitch checked={adm20IdMode} onChange={(e) => handleCheckboxChange(e.target.checked, handleAdm20IdModeChange, saveAdm20IdModeDraft)} />
								</LabeledField>
							)}
						</div>
					</div>

					<div className="flex flex-col space-y-5 border border-[#d9d9d9] p-3">
						<LabeledField label="message.adm20-output-activation-mode" className="flex flex-col sm:flex-row sm:items-center gap-x-3">
							<CFGSelect uid={UIDs.ADM20_OUTPUT_STATE_MASK} array={ADM20_OUTPUT_ACTIVATION_MODE} className="w-full" />
						</LabeledField>

						<LabeledField label="message.adm20-output-state" className="flex flex-col sm:flex-row sm:items-center gap-x-3">
							<CFGSelect uid={UIDs.ADM20_ALARM_STATE} array={ADM20_OUTPUT_STATE} className="w-full" />
						</LabeledField>
					</div>

					<div className="flex flex-col md:flex-row items-stretch gap-4">
						<div className="w-full flex flex-col space-y-5 border border-[#d9d9d9] p-3">
							<p className="uppercase text-[13px] font-medium">{t('message.adm20-applied-rfid-card')}</p>

							<LabeledField label="UID" className="flex items-center flex-row gap-3">
								<input type="text" className="mini" value={cardUid} disabled />
							</LabeledField>
						</div>

						<div className="w-full border border-[#d9d9d9] p-3">
							<p className="uppercase text-[13px] font-medium mb-3">{t('message.adm20-detected-tag')}</p>

							<LabeledField label="message.adm20-detected-tag-sn" className="flex items-center justify-between flex-row gap-3">
								<input type="text" className="mini" value={tagSn} disabled />
							</LabeledField>

							<LabeledField label="message.adm20-detected-tag-id" className="flex items-center justify-between flex-row gap-3">
								<input type="text" className="mini" value={tagId} disabled />
							</LabeledField>

							<LabeledField label="message.adm20-detected-tag-rssi" className="flex items-center justify-between flex-row gap-3">
								<input type="text" className="mini" value={tagRssi} disabled />
							</LabeledField>

							<LabeledField label="message.adm20-detected-tag-battery-voltage" className="flex items-center justify-between flex-row gap-3">
								<input type="text" className="mini" value={tagAdc} disabled />
							</LabeledField>

							<LabeledField label="message.adm20-detected-tag-broadcast-period" className="flex items-center justify-between flex-row gap-3">
								<input type="text" className="mini" value={tagPeriod} disabled />
							</LabeledField>
						</div>
					</div>
				</React.Fragment>
			)}
		</div>
	);
};

export default ADM20Reader;
