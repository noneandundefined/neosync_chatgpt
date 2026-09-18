import { useState } from 'react';
import StatusIcon from '../common/StatusIcon';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { Loading } from '../../Loader/Loading';
import { formatLocalDate } from '@/utils/TimeUtils';
import Dropdown from '@/components/ui/Dropdown/Dropdown';
import { buildRoute, ROUTES } from '@/constants/constants';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicDeviceCommandCancel, basicDeviceCommandDelete, basicDeviceCommandHistory, DeviceCommand } from '@/rest/deviceCommandAPI';

import DotsHorizontal from '@/components/@icons/dots-horizontal';
import CalendarBlankOutline from '@/components/@icons/calendar-blank-outline';

interface CommandHistoryModalProps {
	closeModal: any;
}

interface CommandHistoryListProps {
	data: DeviceCommand[] | null;
	reload: any;
	onSuccess?: () => void;
}

const CommandHistoryModal: React.FC<CommandHistoryModalProps> = ({ closeModal }) => {
	const { data: respDeviceCommandHistory, loading, reload } = useHandleServer(['respDeviceCommandHistory'], () => basicDeviceCommandHistory());

	if (loading) {
		return <Loading />;
	}

	return <CommandHistoryList data={respDeviceCommandHistory} reload={reload} onSuccess={() => closeModal()} />;
};

export const CommandHistoryList: React.FC<CommandHistoryListProps> = ({ data, reload, onSuccess }) => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	const [openMenuIndex, setOpenMenuIndex] = useState<number | null>(null);

	if (!data) {
		return <p className="my-5 text-center text-[#ccc]">{t('message.command-history-empty')}</p>;
	}

	return (
		<div>
			{data.map((cmdHistory, index) => (
				<div key={index} className="relative flex items-center justify-between py-2 mr-2 border-b border-b-[#f3f3f3]">
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

					<div
						className="flex items-center cursor-pointer"
						onClick={(e) => {
							e.stopPropagation();
							setOpenMenuIndex(openMenuIndex === index ? null : index);
						}}
					>
						<DotsHorizontal fill="#49525f" />
					</div>

					<Dropdown open={openMenuIndex === index} close={() => setOpenMenuIndex(null)} className="absolute top-0 right-0 bg-white py-1 border border-[#e4e4e4] min-w-[11rem] z-[100]" stopPropagation>
						<p
							className="text-sm hover:bg-[#f7f7f7] px-3 py-[9px] md:py-[5px] cursor-pointer"
							onClick={() => {
								navigate(buildRoute(ROUTES.DEVICE_DETAILS_COMMAND, { id: cmdHistory.id }));
								if (onSuccess) {
									onSuccess();
								}
							}}
						>
							{t('label.details')}
						</p>
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
				</div>
			))}
		</div>
	);
};

export default CommandHistoryModal;
