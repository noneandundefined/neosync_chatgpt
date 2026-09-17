import React from 'react';
import Header from '@/components/Header/Header';
import SideBar from '@/components/SideBar/SideBar';
import { useTcpHealth } from '@/context/useTcpHealth';

interface PageLayoutProps {
	children: React.ReactNode;
}

const PageLayout: React.FC<PageLayoutProps> = ({ children }) => {
	const tcpHealth = useTcpHealth();

	return (
		<div className={`flex flex-col ${tcpHealth ? 'h-screen' : 'h-[calc(100vh-24px)]'}`}>
			<Header show={true} />

			<div className="flex flex-1 min-h-0 w-full px-[0.4rem] lg:px-[1rem] xl:px-[2rem] gap-5">
				<SideBar />

				<main className="flex flex-col flex-1 mt-4 pb-3 min-h-0 min-w-0">{children}</main>
			</div>
		</div>
	);
};

export default PageLayout;
