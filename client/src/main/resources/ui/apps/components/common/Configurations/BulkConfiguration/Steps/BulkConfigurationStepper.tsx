type BulkConfigurationStepperProps = {
	number: number;
	isLast?: boolean;
};

const BulkConfigurationStepper: React.FC<BulkConfigurationStepperProps> = ({ number, isLast }) => {
	return (
		<div className="hidden md:flex flex-col items-center self-stretch w-8">
			<div className="h-7 w-7 rounded-full bg-[#49525f] text-white flex items-center justify-center font-medium shrink-0">{number}</div>

			{!isLast && <div className="flex-1 w-px bg-gray-300 mt-2" />}
		</div>
	);
};

export default BulkConfigurationStepper;
