import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { BITs } from '@/constants/Bits.constant';
import { arrayParseData } from '@/utils/ArrayUtils';
import { useEffect, useMemo, useState } from 'react';
import CFGInput from '@/components/ui/Input/CfgInput';
import CFGSelect from '@/components/ui/Select/CfgSelect';
import GUISelect from '@/components/ui/Select/GUISelect';
import CFGSwitch from '@/components/ui/Checkbox/CFGSwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { DEVICE_MODELS, getStaticSupport } from '@/constants/DeviceModels.constant';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { STATIC_MODE, TRACK_ENTRY_STATUS_STATIC_MODE, TRACK_STATIC_MODE_2, TRACK_STATIC_MODE_BASE, TRACK_STATIC_MODE_BY_SPEED } from '@/constants/Track.constants';

const AIN_DISABLED = 255;

const isIgnitionMode = (mode: number) => mode === STATIC_MODE.AIN_NORMAL || mode === STATIC_MODE.AIN_INVERSE;

const parseUint = (raw: string | number | string[] | number[] | null | undefined, fallback = 0) => {
	const parsed = arrayParseData(raw);
	const candidate = parsed[0] ?? raw;

	if (candidate === '' || candidate === null || candidate === undefined) {
		return fallback;
	}

	const num = Number(candidate);

	return Number.isFinite(num) ? num : fallback;
};

const FreezingCoordinates = () => {
	const { t } = useTranslation();
	const { model, imei, support, section, manager, isTemplate } = useConfigurationWrapperContext();
	const staticSupport = getStaticSupport((model ?? DEVICE_MODELS.ADM333V2) as keyof typeof DEVICE_MODELS, isTemplate);

	const [mainMode, setMainMode] = useState(0);
	const [entryStatus, setEntryStatus] = useState(0);

	/** STATIC_MODE */
	const { value: staticMode, handleChange: handleStaticModeChange, saveDraft: saveStaticModeDraft } = useConfigurationField(imei, section, manager, UIDs.STATIC_MODE);
	/** STATIC_AIN_NUMBER — номер аналогового входа для INSTATIC X,Y; 255 = функция выключена */
	const { value: staticAinNumber, handleChange: handleAinChange } = useConfigurationField(imei, section, manager, UIDs.STATIC_AIN_NUMBER);

	const staticModeValue = useMemo(() => parseUint(staticMode, 0), [staticMode]);
	const ainValue = useMemo(() => parseUint(staticAinNumber, AIN_DISABLED), [staticAinNumber]);

	const ensureIgnitionInput = () => {
		if (ainValue === AIN_DISABLED) {
			handleAinChange(0);
		}
	};

	const clearIgnitionInput = () => {
		if (ainValue !== AIN_DISABLED) {
			handleAinChange(AIN_DISABLED);
		}
	};

	useEffect(() => {
		if (staticModeValue === STATIC_MODE.AIN_NORMAL) {
			setMainMode(STATIC_MODE.AIN_NORMAL);
			setEntryStatus(STATIC_MODE.AIN_NORMAL);
		} else if (staticModeValue === STATIC_MODE.AIN_INVERSE) {
			setMainMode(STATIC_MODE.AIN_NORMAL);
			setEntryStatus(STATIC_MODE.AIN_INVERSE);
		} else {
			setMainMode(staticModeValue);
			setEntryStatus(STATIC_MODE.OFF);
		}
	}, [staticModeValue]);

	useEffect(() => {
		if (!manager?.getField(UIDs.STATIC_AIN_NUMBER)) {
			return;
		}

		if (isIgnitionMode(staticModeValue) && ainValue === AIN_DISABLED) {
			handleAinChange(0);
		}
	}, [manager, staticModeValue, ainValue, handleAinChange]);

	const saveStatic = async (newMain: number, newEntry: number) => {
		let value = 0;

		if (newMain === STATIC_MODE.AIN_NORMAL) {
			value = newEntry === STATIC_MODE.AIN_INVERSE ? STATIC_MODE.AIN_INVERSE : STATIC_MODE.AIN_NORMAL;
		} else {
			value = newMain;
		}

		handleStaticModeChange(value);
		await saveStaticModeDraft(value);

		if (isIgnitionMode(value)) {
			ensureIgnitionInput();
		} else {
			clearIgnitionInput();
		}
	};

	const onMainChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
		const val = parseInt(e.target.value, 10);

		setMainMode(val);
		await saveStatic(val, entryStatus);
	};

	const onEntryStatusChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
		const val = parseInt(e.target.value, 10);

		setEntryStatus(val);
		await saveStatic(mainMode, val);
	};

	const onSpeedModeChange = (selectedId: string) => {
		const mode = Number(selectedId);

		if (mode === STATIC_MODE.AIN_NORMAL) {
			ensureIgnitionInput();
		} else {
			clearIgnitionInput();
		}
	};

	const track_static_array = isTemplate || support.ModeHybrid ? TRACK_STATIC_MODE_2 : TRACK_STATIC_MODE_BASE;
	const showIgnitionInput = isIgnitionMode(staticModeValue) || mainMode === STATIC_MODE.AIN_NORMAL;

	return (
		<>
			{isTemplate || support.StaticFreezingCoordinates ? (
				<div className="min-w-[20vw] space-y-4">
					<LabeledField label="message.coordinate-freeze-mode">
						<CFGSelect uid={UIDs.STATIC_MODE} array={TRACK_STATIC_MODE_BY_SPEED} onChange={onSpeedModeChange} />
					</LabeledField>

					<LabeledField label="message.coordinate-freeze-at-invalid" className="flex items-center gap-3">
						<CFGSwitch uid={UIDs.DEVICE_FUNCTION_2} bit={BITs.DEVICE_FUNCTION_2.COORDINATE_FREEZE_AT_INVALID} vOff={0} vOn={1} />
					</LabeledField>

					<LabeledField label="message.input-number" disabled={!showIgnitionInput}>
						<Tooltip title={t('message.track-index-of-input')}>
							<CFGInput uid={UIDs.STATIC_AIN_NUMBER} disabled={!showIgnitionInput} min={0} max={1} />
						</Tooltip>
					</LabeledField>

					<LabeledField label="label.speed" disabled={staticModeValue !== STATIC_MODE.PROGRAM}>
						<Tooltip position="bottom" title={t('message.bigger-value-field-triggers-device-speed')}>
							<div className="flex items-center gap-2">
								<CFGInput uid={UIDs.STATIC_PROGRAM_SPEED} disabled={staticModeValue !== STATIC_MODE.PROGRAM} />
								<label>{t('label.km-h')}</label>
							</div>
						</Tooltip>
					</LabeledField>
				</div>
			) : (
				<div className="min-w-[20vw] space-y-4">
					<LabeledField label="message.coordinate-freeze-mode">
						<Tooltip title={t('message.tooltip-freeze-coordinates-on-parking')}>
							<GUISelect value={mainMode} onChange={(e) => onMainChange(e)}>
								{track_static_array.map((track, index) => (
									<option value={track.id} key={index}>
										{t(track.value)}
									</option>
								))}
							</GUISelect>
						</Tooltip>
					</LabeledField>

					{(isTemplate || staticSupport.StaticFreezingCoordinates) && (
						<LabeledField label="message.coordinate-freeze-at-invalid" className="flex items-center gap-3">
							<CFGSwitch uid={UIDs.DEVICE_FUNCTION_2} bit={BITs.DEVICE_FUNCTION_2.COORDINATE_FREEZE_AT_INVALID} vOff={0} vOn={1} />
						</LabeledField>
					)}

					{mainMode === STATIC_MODE.AIN_NORMAL && (
						<>
							<LabeledField label="message.input-number">
								<Tooltip title={t('message.track-index-of-input')}>
									<CFGInput uid={UIDs.STATIC_AIN_NUMBER} min={0} max={1} />
								</Tooltip>
							</LabeledField>

							<LabeledField label="message.input-state-for-coordinate-freeze">
								<GUISelect value={entryStatus} onChange={(e) => onEntryStatusChange(e)}>
									{TRACK_ENTRY_STATUS_STATIC_MODE.map((track, index) => (
										<option value={track.id} key={index}>
											{t(track.value)}
										</option>
									))}
								</GUISelect>
							</LabeledField>
						</>
					)}
				</div>
			)}
		</>
	);
};

export default FreezingCoordinates;
