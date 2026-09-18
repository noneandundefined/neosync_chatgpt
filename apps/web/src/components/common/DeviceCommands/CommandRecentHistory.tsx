import React, { useState } from 'react';
import { Loading } from '../Loader/Loading';
import StatusIcon from './common/StatusIcon';
import { useTranslation } from 'react-i18next';
import { formatLocalDate } from '@/utils/TimeUtils';
import Dropdown from '@/components/ui/Dropdown/Dropdown';
import { formatCommandResponse } from '@/utils/CommandResponseUtils';
import { basicDeviceCommandCancel, basicDeviceCommandDelete, DeviceCommandWithExecutions } from '@/rest/deviceCommandAPI';

import ChevronUp from '@/components/@icons/chevron-up';
import ChevronDown from '@/components/@icons/chevron-down';
import DotsHorizontal from '@/components/@icons/dots-horizontal';
import CalendarBlankOutline from '@/components/@icons/calendar-blank-outline';

interface CommandRecentHistoryProps {
	data: DeviceCommandWithExecutions[] | null;
	loading: boolean;
	reload: () => void;
	openIndex?: number | null;
	onOpenIndexChange?: (index: number | null) => void;
}

const CommandRecentHistory: React.FC<CommandRecentHistoryProps> = ({ data, loading, reload, openIndex: controlledOpenIndex, onOpenIndexChange }) => {
	const { t } = useTranslation();

	const [internalOpenIndex, setInternalOpenIndex] = useState<number | null>(null);
	const [openMenuIndex, setOpenMenuIndex] = useState<number | null>(null);

	const openIndex = controlledOpenIndex !== undefined ? controlledOpenIndex : internalOpenIndex;
	const setOpenIndex = onOpenIndexChange ?? setInternalOpenIndex;

	if (loading) {
		return (
			<div>
				<Loading />
			</div>
		);
	}

	if (!data) {
		return null;
	}

	return (
		<div className="flex-1 !mt-5 space-y-3">
			<div className="relative flex">
				<p className="text-center font-medium text-[17px] text-[#49525f]">{t('message.recent-commands')}</p>
			</div>

			{data.slice(0, 7).map((cmdHistory, index) => (
				<div className="relative" key={cmdHistory.id}>
					<div
						className="flex items-center justify-between mr-2 my-2 cursor-pointer"
						onClick={(e) => {
							e.stopPropagation();
							setOpenIndex(openIndex === index ? null : index);
						}}
					>
						<div className="flex items-center gap-2">
							<div className="w-[1rem] h-[1rem] flex items-center justify-center">
								<StatusIcon status={cmdHistory.status} t={t} />
							</div>

							<p className="font-medium text-sm uppercase truncate max-w-[130px] md:max-w-[350px]">{cmdHistory.command}</p>

							<div className="flex items-center gap-1">
								<CalendarBlankOutline fill="#49525f" size={17} />
								<p className="text-[#49525f] text-[13px]">{formatLocalDate(cmdHistory.created_at)}</p>
							</div>
						</div>

						<div className="flex items-center gap-3">
							<div className="flex items-center">{openIndex === index ? <ChevronUp fill="#49525f" size={21} /> : <ChevronDown fill="#49525f" size={21} />}</div>
							<div
								className="flex items-center cursor-pointer"
								onClick={(e) => {
									e.stopPropagation();
									setOpenMenuIndex(openMenuIndex === index ? null : index);
								}}
							>
								<DotsHorizontal fill="#49525f" />
							</div>
						</div>
					</div>

					<Dropdown open={openMenuIndex === index} close={() => setOpenMenuIndex(null)} className="absolute top-0 right-0 bg-white py-1 border border-[#e4e4e4] min-w-[11rem] z-[100]" stopPropagation>
						{(cmdHistory.status === 'pending' || cmdHistory.status === 'inprogress') && (
							<p
								className="text-sm hover:bg-[#f7f7f7] px-3 py-[9px] md:py-[5px] cursor-pointer"
								onClick={async () => {
									await basicDeviceCommandCancel(cmdHistory.id);
									setOpenMenuIndex(null);
									reload();
								}}
							>
								{t('message.cancel-execution')}
							</p>
						)}
						<p
							className="text-sm text-[#d1242f] hover:text-white hover:bg-[#d1242f] px-3 py-[9px] md:py-[5px] cursor-pointer"
							onClick={async () => {
								await basicDeviceCommandDelete(cmdHistory.id);
								setOpenMenuIndex(null);
								reload();
							}}
						>
							{t('message.delete-command')}
						</p>
					</Dropdown>

					{openIndex === index && (
						<div className="max-h-[50vw] md:max-h-[15vw] overflow-y-auto">
							{cmdHistory.executions.map((device, executionIndex) => (
								<div className="space-y-3 m-2 bg-white border border-[#f3f3f3] p-2 rounded" key={executionIndex}>
									<div className="flex items-center gap-2">
										<div className="w-[1rem] h-[1rem] flex items-center justify-center">
											<StatusIcon status={device.status} t={t} />
										</div>

										<p className="font-medium">{device.imei}</p>
									</div>

									<div className="flex items-start justify-between gap-3 bg-[#f9f9f9] p-2 rounded">
										<p className="text-sm text-[#666]">{formatCommandResponse(device.response, t)}</p>
									</div>
								</div>
							))}
						</div>
					)}

					<div className="w-full h-[1px] bg-[#f3f3f3]" />
				</div>
			))}
		</div>
	);
};

export default CommandRecentHistory;
