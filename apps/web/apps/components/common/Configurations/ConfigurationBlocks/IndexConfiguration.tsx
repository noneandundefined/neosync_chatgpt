import React from 'react';

interface IndexConfigurationProps {
	title: string;
	content: React.ReactNode;
	maxWidth?: string;
	className?: string;
}

const IndexConfiguration: React.FC<IndexConfigurationProps> = ({ title, content, maxWidth = 'full', className = '' }) => {
	return (
		<div className={`flex-1 border border-[#d1d5db] my-3 ${className}`} style={{ maxWidth: maxWidth }}>
			<p className="font-medium text-[#49525f] text-[16px] bg-[#f1f2f3] p-1 px-3 min-w-[23.5vw] whitespace-nowrap">{title}</p>

			<div className="flex flex-col gap-4 my-2 p-3">
				<div className="w-full">{content}</div>
			</div>
		</div>
	);
};

export default IndexConfiguration;
