import React from 'react';
import Close from '@/components/@icons/close';
import { useTranslation } from 'react-i18next';
import LabeledField from '@/components/ui/Form/LabeledField';

interface SensorProps {
	index: number;
	sensor: string;
	typeLabel: string;
	onRemove: () => void;

	telemetry?: Record<string, any>;

	fields?: {
		label: string;
		value: any;
		unit: string;
	}[];
}

const Sensor: React.FC<SensorProps> = ({ index, sensor, typeLabel, onRemove, telemetry, fields = [] }) => {
	const { t } = useTranslation();

	const noData = t('label.no-data');

	const isMissing = (value: any) => value === noData || value === null || value === undefined || value === '';
	const hasIncompleteData = fields.some((field) => isMissing(field.value));
	const displayFields = hasIncompleteData ? fields.map((field) => ({ ...field, value: noData })) : fields;

	const lmtValue = hasIncompleteData ? null : (telemetry?.lmt ?? null);

	return (
		<div className="w-full md:w-auto rounded-[6px] border">
			<div className="bg-[#f1f2f3] p-2 rounded-tl-[6px] rounded-tr-[6px]">
				<div className="flex justify-between items-center min-w-[17rem]">
					<span className="flex items-center gap-2">
						<div className="flex items-center justify-center bg-[#395d95] min-w-[1.8rem] min-h-[1.8rem] rounded">
							<p className="text-white text-[14px] font-medium">{index}</p>
						</div>
						<p className="text-[14px]">{typeLabel ? t(typeLabel) : t('label.adm-unknown')}</p>
					</span>
					<div className="flex items-center gap-3">
						<span className="text-[13px]">{sensor}</span>
						<div onClick={onRemove} className="cursor-pointer rounded-[6px]">
							<Close size={19} fill="#d1242f" onClick={() => {}} />
						</div>
					</div>
				</div>
			</div>

			<div className="p-2 space-y-3 min-w-[20rem]">
				{displayFields.map((field, index) => (
					<LabeledField key={index} label={field.label} className="flex items-center gap-3">
						<div className="flex items-center gap-3">
							<input value={field.value} className="p-1" disabled />
							<span className="text-[13px]">{t(field.unit)}</span>
						</div>
					</LabeledField>
				))}

				{/* LMT */}
				<div className="flex justify-end">
					{lmtValue === null ? (
						<p className="font-medium text-[#d1242f]">N/A</p>
					) : (
						<p className="font-medium text-[#019100]">
							{lmtValue} {t('label.second-short')}.
						</p>
					)}
				</div>
			</div>
		</div>
	);
};

export default Sensor;
