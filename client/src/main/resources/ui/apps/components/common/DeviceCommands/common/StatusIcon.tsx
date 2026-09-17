import Check from '@/components/@icons/check';
import Close from '@/components/@icons/close';
import Minus from '@/components/@icons/minus';
import Tooltip from '@/components/ui/Tooltip';

interface StatusIconProps {
	status: string;
	t: any;
}

const StatusIcon: React.FC<StatusIconProps> = ({ status, t }) => {
	switch (status) {
		case 'completed':
			return (
				<Tooltip title={t('label.status-success')}>
					<Check size={15} fill="#029b00" />
				</Tooltip>
			);

		case 'executionerror':
			return (
				<Tooltip title={t('label.status-failed')}>
					<Close size={15} fill="#e60c00" />
				</Tooltip>
			);

		case 'notcompleted':
			return (
				<Tooltip title={t('label.notcompleted')}>
					<div className="rotate-45">
						<Minus size={15} fill="#49525f" />
					</div>
				</Tooltip>
			);

		case 'inprogress':
			return (
				<Tooltip title={t('label.status-in-progress')}>
					<div className="relative flex justify-center items-center">
						<div className="w-[8px] h-[8px] rounded-full bg-[#9a6700]" />
						<div className="absolute w-[8px] h-[8px] border-b-2 border-[#9a6700] rounded-full p-[7px] animate-spin"></div>
					</div>
				</Tooltip>
			);

		default:
			return (
				<Tooltip title={t('label.status-pending')}>
					<div className="w-[8px] h-[8px] rounded-full bg-[#9a6700]" />
				</Tooltip>
			);
	}
};

export default StatusIcon;
