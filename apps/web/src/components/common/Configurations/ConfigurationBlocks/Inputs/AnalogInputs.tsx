import Slider from 'rc-slider';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import { useEffect, useRef, useState } from 'react';
import IndexConfiguration from '../IndexConfiguration';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import CFGArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const AnalogInputs = () => {
	const { t } = useTranslation();
	const { support } = useConfigurationWrapperContext();

	return (
		<div className="flex flex-col flex-wrap">
			{Array.from({ length: support.InputCount }, (_, index) => (
				<IndexConfiguration title={`${t('label.input-short')} ${index}`} content={<Input index={index} />} key={index} />
			))}
		</div>
	);
};

const Input: React.FC<{ index: number }> = ({ index }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	const { value: ainFalseValue = [], handleChange: handleChangeFalse, saveDraft: saveFalseDraft } = useConfigurationField(imei, section, manager, UIDs.AIN_FALSE_HIGH);

	const { value: ainTrueValue = [], handleChange: handleChangeTrue, saveDraft: saveTrueDraft } = useConfigurationField(imei, section, manager, UIDs.AIN_TRUE_LOW);

	const [falseVal, setFalseVal] = useState<number>(ainFalseValue[index] || 0);
	const [trueVal, setTrueVal] = useState<number>(ainTrueValue[index] || 0);

	useEffect(() => {
		setFalseVal(ainFalseValue[index] || 0);
	}, [ainFalseValue, index]);

	useEffect(() => {
		setTrueVal(ainTrueValue[index] || 0);
	}, [ainTrueValue, index]);

	const [enabled, setEnabled] = useState<boolean>(false);
	const userToggledRef = useRef(false);
	const hasAny = (ainFalseValue[index] ?? 0) !== 0 || (ainTrueValue[index] ?? 0) !== 0;

	useEffect(() => {
		if (!userToggledRef.current) {
			setEnabled(hasAny);
		}
	}, [hasAny]);

	const effectiveOnFnFalse = userToggledRef.current ? (enabled ? 4000 : 0) : undefined;
	const effectiveOnFnTrue = userToggledRef.current ? (enabled ? 7000 : 0) : undefined;

	const styles = {
		disabledStyle: {
			opacity: enabled ? 1 : 0.5,
			pointerEvents: (enabled ? 'auto' : 'none') as React.CSSProperties['pointerEvents'],
		},
	};

	const handleSliderFalse = async (val: number[]) => {
		const newVal = val[1];
		const trueLow = ainTrueValue[index] ?? 0;

		if (trueLow !== 0 && newVal >= trueLow) return;

		if (ainFalseValue[index] === newVal) return;

		const newArr = [...ainFalseValue];
		newArr[index] = newVal;

		handleChangeFalse(newArr);
		await saveFalseDraft(newArr);
	};

	const handleSliderTrue = async (val: number[]) => {
		const newVal = val[0];
		const falseHigh = ainFalseValue[index] ?? 0;

		if (falseHigh !== 0 && newVal <= falseHigh) return;

		if (ainTrueValue[index] === newVal) return;

		const newArr = [...ainTrueValue];
		newArr[index] = newVal;

		handleChangeTrue(newArr);
		await saveTrueDraft(newArr);
	};

	return (
		<div className="w-full sm:min-w-[10rem] space-y-4">
			<div className="flex flex-col sm:flex-row sm:items-center gap-3" style={styles.disabledStyle}>
				<p className="font-normal text-[14px]">{t('label.voltage')}:</p>
				<input type="text" value={0} className="sm:max-w-[5rem]" readOnly />
			</div>

			<div className="flex justify-between sm:justify-start items-center gap-3">
				<p className="font-normal text-[14px]">{t('message.use-as-discrete')}:</p>
				<GUISwitch
					checked={enabled}
					onChange={() => {
						const next = !enabled;
						userToggledRef.current = true;
						setEnabled(next);
					}}
				/>
			</div>

			<div className="flex flex-col xl:flex-row xl:items-center gap-3" style={styles.disabledStyle}>
				<p className="font-normal text-[14px] whitespace-nowrap">{t('message.logical-zero')}:</p>
				<div className="flex flex-col xl:flex-row xl:items-center gap-8 w-full">
					<div className="flex items-center gap-2">
						<input type="text" value={0} className="xl:max-w-[5rem]" style={{ flex: 3 }} readOnly />
						<span>-</span>
						<CFGArrayItemInput uid={UIDs.AIN_FALSE_HIGH} index={index} value={falseVal} onFn={effectiveOnFnFalse} className="xl:max-w-[5rem]" />
						<label>{t('label.mV')}</label>
					</div>

					<Slider
						range
						min={0}
						max={60000}
						className="flex-1 w-full"
						value={[0, falseVal]}
						onChange={(val) => {
							if (Array.isArray(val)) setFalseVal(val[1]);
						}}
						onAfterChange={(val) => {
							if (Array.isArray(val)) handleSliderFalse(val as number[]);
						}}
						trackStyle={[{ backgroundColor: 'blue' }]}
						handleStyle={[{ borderColor: 'blue' }, { borderColor: 'blue' }]}
					/>
				</div>
			</div>

			<div className="flex flex-col xl:flex-row xl:items-center gap-3" style={styles.disabledStyle}>
				<p className="font-normal text-[14px] whitespace-nowrap">{t('message.logical-one')}:</p>
				<div className="flex flex-col xl:flex-row xl:items-center gap-8 w-full">
					<div className="flex items-center gap-2">
						<CFGArrayItemInput uid={UIDs.AIN_TRUE_LOW} index={index} value={trueVal} onFn={effectiveOnFnTrue} className="xl:max-w-[5rem]" />
						<span>-</span>
						<input type="text" className="xl:max-w-[5rem]" value={!enabled ? 0 : 60000} style={{ flex: 3 }} readOnly />
						<label>{t('label.mV')}</label>
					</div>

					<Slider
						range
						min={0}
						max={60000}
						className="flex-1 w-full"
						value={[trueVal, 60000]}
						onChange={(val) => {
							if (Array.isArray(val)) setTrueVal(val[0]);
						}}
						onAfterChange={(val) => {
							if (Array.isArray(val)) handleSliderTrue(val as number[]);
						}}
						trackStyle={[{ backgroundColor: 'green' }]}
						handleStyle={[{ borderColor: 'green' }, { borderColor: 'green' }]}
					/>
				</div>
			</div>
		</div>
	);
};

export default AnalogInputs;
