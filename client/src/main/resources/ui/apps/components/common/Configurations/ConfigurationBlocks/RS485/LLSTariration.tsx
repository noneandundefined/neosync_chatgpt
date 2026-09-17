import { useTranslation } from 'react-i18next';
import Colors from '@/constants/Color.constant';
import { LLS_COUNT } from '@/constants/DeviceModels.constant';
import GUICheckbox from '@/components/ui/Checkbox/GUICheckbox';
import { Dispatch, SetStateAction, useEffect, useRef, useState } from 'react';
import { basicConfigurationExportLLSTariration } from '@/rest/configurationAPI';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';
import { CACHEKEYs_TERMINAL_TARIRATION_SELECTED, CACHEKEYs_TERMINAL_TARIRATION_TABLES } from '@/constants/CacheKeys.constants';

import Plus from '@/components/@icons/plus';
import Close from '@/components/@icons/close';
import FolderOpen from '@/components/@icons/folder-open';
import ContentSave from '@/components/@icons/content-save';
import PlaylistRemove from '@/components/@icons/playlist-remove';

const LLSTariration = () => {
	const { t } = useTranslation();
	const { imei } = useConfigurationWrapperContext();

	const keyTarirationTables = CACHEKEYs_TERMINAL_TARIRATION_TABLES(imei);
	const keyTarirationSelected = CACHEKEYs_TERMINAL_TARIRATION_SELECTED(imei);

	const [selected, setSelected] = useState<number[]>([]);

	const [tables, setTables] = useState<{
		[sensorId: number]: { id: number; value: number; level: number }[];
	}>({});

	useEffect(() => {
		const saved = sessionStorage.getItem(keyTarirationTables);
		const savedSelected = sessionStorage.getItem(keyTarirationSelected);

		if (saved) setTables(JSON.parse(saved));
		if (savedSelected) setSelected(JSON.parse(savedSelected));
	}, []);

	useEffect(() => {
		sessionStorage.setItem(keyTarirationTables, JSON.stringify(tables));
	}, [tables]);

	useEffect(() => {
		sessionStorage.setItem(keyTarirationSelected, JSON.stringify(selected));
	}, [selected]);

	const fileInputRef = useRef<HTMLInputElement | null>(null);
	const handleOpenFile = () => {
		fileInputRef.current?.click();
	};

	const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
		const file = e.target.files?.[0];
		if (!file) return;

		const reader = new FileReader();

		reader.onload = () => {
			let text = reader.result as string;

			if (/[\uFFFD]/.test(text)) {
				const reader1251 = new FileReader();
				reader1251.onload = () => parseCsv(reader1251.result as string);
				reader1251.readAsText(file, 'windows-1251');
			} else {
				parseCsv(text);
			}
		};

		reader.readAsText(file, 'UTF-8');
	};

	const parseCsv = (csvText: string) => {
		const lines = csvText.split(/\r?\n/);

		const newTables: typeof tables = {};
		const newSelected: number[] = [];

		let currentSensor: number | null = null;

		for (let raw of lines) {
			const line = raw.trim();
			if (!line) continue;

			const sensorMatch = line.match(/Датчик\s+(\d+)/i);
			if (sensorMatch) {
				currentSensor = Number(sensorMatch[1]);
				newTables[currentSensor] = [];
				if (!newSelected.includes(currentSensor)) newSelected.push(currentSensor);
				continue;
			}

			if (line.startsWith('№') || currentSensor === null) continue;

			const parts = line.split(';').map((p) => p.trim());
			if (parts.length < LLS_COUNT) continue;

			if (!/^\d+$/.test(parts[0])) continue;

			const index = Number(parts[0]);
			const value = Number(parts[1].replace(',', '.'));
			const level = Number(parts[2].replace(',', '.'));

			newTables[currentSensor]?.push({
				id: Date.now() + index + currentSensor,
				value,
				level,
			});
		}

		setTables(newTables);
		setSelected(newSelected);
	};

	const addRow = () => {
		selected.forEach((sensorId) => {
			setTables((prev) => ({
				...prev,
				[sensorId]: [...(prev[sensorId] || []), { id: Date.now(), value: 0, level: 0 }],
			}));
		});
	};

	const toggleSensor = (index: number) => {
		setSelected((prev) => (prev.includes(index) ? prev.filter((x) => x !== index) : [...prev, index]));
	};

	return (
		<div className="w-full space-y-6">
			<p className="text-sm">{t('message.tariration-description')}</p>
			<p className="text-sm font-semibold">{t('message.tariration-note')}</p>

			<div className="flex flex-col sm:flex-row justify-between xl:justify-start items-center sm:items-start gap-5">
				<div className="flex flex-col xl:flex-row gap-5">
					<p className="text-[14px]">{t('message.tariration-select-sensors')}</p>

					<div className="flex flex-col items-center sm:items-start space-y-2 ">
						{Array.from({ length: LLS_COUNT }, (_, index) => (
							<div className="flex items-center gap-3" key={index}>
								<GUICheckbox checked={selected.includes(index)} onChange={() => toggleSensor(index)} />
								<p className="text-sm max-w-full sm:max-w-[30rem]">
									{t('label.sensor-short')} {index}
								</p>
							</div>
						))}
					</div>
				</div>

				<div className="flex flex-wrap items-start gap-3 ml-5">
					<div className="flex flex-col items-center justify-start max-w-[4rem]">
						<Plus size={40} fill="#008c25" className="p-[6px] border border-[#dedede] cursor-pointer bg-[#eee] hover:bg-[#d9d9d9]" onClick={addRow} />

						<p className="text-[12px] text-center mt-2">{t('message.tariration-add-values')}</p>
					</div>

					<div className="flex flex-col items-center justify-start max-w-[4.3rem]" onClick={handleOpenFile}>
						<FolderOpen size={40} fill="#000" className="p-[6px] border border-[#dedede] cursor-pointer bg-[#eee] hover:bg-[#d9d9d9]" />

						<p className="text-[12px] text-center mt-2">{t('message.tariration-open-table-file')}</p>
					</div>

					<input type="file" accept=".csv" ref={fileInputRef} onChange={handleFileUpload} className="hidden" />

					<div className="flex flex-col items-center justify-start max-w-[4.3rem]">
						<ContentSave
							size={40}
							fill="#000"
							className="p-[6px] border border-[#dedede] cursor-pointer bg-[#eee] hover:bg-[#d9d9d9]"
							onClick={() =>
								basicConfigurationExportLLSTariration(imei, {
									tariration: tables,
								})
							}
						/>

						<p className="text-[12px] text-center mt-2">{t('message.tariration-save-tables')}</p>
					</div>
				</div>
			</div>

			<div className="flex flex-wrap items-start gap-3">
				{selected.length > 0 ? (
					selected.sort((a, b) => a - b).map((index) => <LLSTarirationSensor sensorId={index} tables={tables} setTables={setTables} />)
				) : (
					<div className="w-full flex justify-center my-5">
						<p className="font-medium text-[#ccc]">{t('message.tariration-no-sensors-selected')}</p>
					</div>
				)}
			</div>
		</div>
	);
};

const LLSTarirationSensor = ({
	sensorId,
	tables,
	setTables,
}: {
	sensorId: number;
	tables: {
		[sensorId: number]: {
			id: number;
			value: number;
			level: number;
		}[];
	};
	setTables: Dispatch<
		SetStateAction<{
			[sensorId: number]: {
				id: number;
				value: number;
				level: number;
			}[];
		}>
	>;
}) => {
	const { t } = useTranslation();

	const rows = tables[sensorId] || [];

	const deleteRow = (sensorId: number, rowId: number) => {
		setTables((prev) => ({
			...prev,
			[sensorId]: prev[sensorId].filter((r) => r.id !== rowId),
		}));
	};

	const clearTable = (sensorId: number) => {
		setTables((prev) => ({ ...prev, [sensorId]: [] }));
	};

	const updateCell = (sensorId: number, rowId: number, field: 'value' | 'level', newValue: number) => {
		setTables((prev) => ({
			...prev,
			[sensorId]: prev[sensorId].map((r) => (r.id === rowId ? { ...r, [field]: newValue } : r)),
		}));
	};

	return (
		<div className="flex-1 min-w-[14rem] min-h-[16rem] max-h-[26rem] border border-[#6af] p-1 overflow-y-auto hide-scrollbar">
			<div className="flex items-center justify-between">
				<p className="text-[14px] font-medium text-center w-full">
					{t('message.tariration-sensor-table')} {sensorId}
				</p>
				<PlaylistRemove size={32} fill={Colors.color_error_base} className="h-[26px] w-[30px] border border-[#dedede] cursor-pointer bg-[#eee] hover:bg-[#d9d9d9]" onClick={() => clearTable(sensorId)} />
			</div>

			<table className="w-full border-collapse text-[13px]">
				<thead>
					<tr className="border-b border-gray-300">
						<th className="text-left font-normal px-1 py-1 w-[2rem]">№</th>
						<th className="text-left font-normal px-1 py-1">{t('label.value').toUpperCase()}</th>
						<th className="text-left font-normal px-1 py-1">{t('label.level').toUpperCase()}</th>
						<th className="w-[1rem]"></th>
					</tr>
				</thead>

				<tbody>
					{rows.map((row, i) => (
						<tr key={row.id} className="border-b border-gray-200">
							<td className="px-1 py-1">{i}</td>

							<td className="px-1 py-1">
								<input type="number" className="w-full border-0 hover:border-0 rounded-none px-1 text-sm" value={row.value} onChange={(e) => updateCell(sensorId, row.id, 'value', Number(e.target.value))} />
							</td>

							<td className="px-1 py-1 inline-block">
								<input type="number" className="w-full border-0 hover:border-0 rounded-none px-1 text-sm" value={row.level} onChange={(e) => updateCell(sensorId, row.id, 'level', Number(e.target.value))} />
							</td>

							<td>
								<div className="h-[20px] w-[20px] flex items-center justify-center border border-[#dedede] bg-[#eee] hover:bg-[#d9d9d9] cursor-pointer" onClick={() => deleteRow(sensorId, row.id)}>
									<Close size={14} fill={Colors.color_error_base} />
								</div>
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
};

export default LLSTariration;
