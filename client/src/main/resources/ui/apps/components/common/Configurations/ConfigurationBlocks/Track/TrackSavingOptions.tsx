import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import CFGInput from '@/components/ui/Input/CfgInput';
import GUIRange from '@/components/ui/Range/GUIRange';
import GUISwitch from '@/components/ui/Checkbox/GUISwitch';
import LabeledField from '@/components/ui/Form/LabeledField';
import { TRACK_RANGE_PRESETS } from '@/constants/Track.constants';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import CfgArrayItemInput from '@/components/ui/Input/CfgArrayItemInput';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

interface TrackSavingOptionsProps {
	expertSettings: boolean;
	setExpertSettings: Dispatch<SetStateAction<boolean>>;
}

const TrackSavingOptions: React.FC<TrackSavingOptionsProps> = ({ expertSettings, setExpertSettings }) => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	const { value: course, handleChange: handleCourseChange, saveDraft: saveCourseDraft } = useConfigurationField(imei, section, manager, UIDs.POINTS_TRACK_COURSE);
	const { value: crosstrack, handleChange: handleCrosstrackChange, saveDraft: saveCrosstrackDraft } = useConfigurationField(imei, section, manager, UIDs.POINTS_TRACK_CROSSTRACK);

	const computeCurrentRange = () => {
		for (const [rangeStr, preset] of Object.entries(TRACK_RANGE_PRESETS)) {
			const range = parseInt(rangeStr, 10);
			if (JSON.stringify(preset.course) === JSON.stringify(course) && JSON.stringify(preset.crosstrack) === JSON.stringify(crosstrack)) return range;
		}

		return 0;
	};

	const [range, setRange] = useState<number>(computeCurrentRange);

	useEffect(() => {
		setRange(computeCurrentRange());
	}, [course, crosstrack]);

	const handleTrackDetailChange = async (value: number) => {
		setRange(value);

		const preset = TRACK_RANGE_PRESETS[value as keyof typeof TRACK_RANGE_PRESETS];
		if (!preset) return;

		handleCourseChange(preset.course);
		await saveCourseDraft(preset.course);

		handleCrosstrackChange(preset.crosstrack);
		await saveCrosstrackDraft(preset.crosstrack);
	};

	return (
		<>
			<div className="w-full sm:min-w-[20vw] space-y-4">
				<p className="font-medium">{t('message.track-details')}</p>

				<LabeledField label="message.track-detail-level" className="flex flex-col lg:flex-row lg:items-center lg:gap-5">
					<div className="w-full sm:min-w-[33rem]">
						<GUIRange max={4} value={range} onChange={setRange} onChangeComplete={handleTrackDetailChange} />
					</div>
				</LabeledField>

				<div className="flex items-center gap-3 my-3">
					<GUISwitch checked={expertSettings} onChange={(e) => setExpertSettings(e.target.checked)} />
					<p className="font-medium">{t('message.expert-settings')}</p>
				</div>

				{expertSettings && (
					<div className="space-y-6">
						<div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
							<LabeledField label="message.minimum-speed">
								<Tooltip title={t('message.tooltip-speed-threshold')}>
									<div className="flex items-center gap-2">
										<CFGInput uid={UIDs.POINTS_TRACK_SPEED_STOP} className="mini min-w-full" />
										<label>{t('label.km-h')}</label>
									</div>
								</Tooltip>
							</LabeledField>

							<LabeledField label="label.distance">
								<Tooltip title={t('message.tooltip-distance-threshold')}>
									<div className="flex items-center gap-2">
										<CFGInput uid={UIDs.POINTS_TRACK_DISTANCE} className="mini min-w-full" />
										<label>{t('label.m')}</label>
									</div>
								</Tooltip>
							</LabeledField>

							<LabeledField label="message.speed-change">
								<Tooltip title={t('message.tooltip-acceleration-threshold')}>
									<div className="flex items-center gap-2">
										<CFGInput uid={UIDs.POINTS_TRACK_ACCELERATION} className="mini min-w-full" />
										<label>{t('label.km-h-s')}</label>
									</div>
								</Tooltip>
							</LabeledField>
						</div>

						<div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
							{['slow-motion', 'medium-motion', 'fast-motion'].map((motion, motionIndex) => (
								<div key={motionIndex} className="space-y-3">
									<p className="font-medium my-2">{t(`message.${motion}`)}</p>

									<LabeledField label="label.angle">
										<div className="flex items-center gap-3">
											<CfgArrayItemInput uid={UIDs.POINTS_TRACK_COURSE} index={motionIndex} className="mini min-w-full" />
											<label>{t('label.deg')}</label>
										</div>
									</LabeledField>

									<LabeledField label="label.deviation">
										<div className="flex items-center gap-3">
											<CfgArrayItemInput uid={UIDs.POINTS_TRACK_CROSSTRACK} index={motionIndex} className="mini min-w-full" />
											<label>{t('label.deg')}</label>
										</div>
									</LabeledField>
								</div>
							))}
						</div>
					</div>
				)}
			</div>
		</>
	);
};

export default TrackSavingOptions;
