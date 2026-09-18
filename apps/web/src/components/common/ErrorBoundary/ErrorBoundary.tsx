import { config } from '@/.config/config.client';
import i18next from 'i18next';
import React, { ReactNode } from 'react';

interface ErrorBoundaryProps {
	children: ReactNode;
}

interface ErrorBoundaryState {
	error?: Error;
	hasError: boolean;
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
	state = { hasError: false, error: new Error() };

	static getDerivedStateFromError(error: Error) {
		return { hasError: true, error };
	}

	render() {
		if (this.state.hasError) {
			return (
				<React.Fragment>
					<main className="flex flex-col items-center justify-center h-screen w-screen px-4 text-center">
						<h1 className="my-5 text-xl sm:text-3xl font-semibold">{i18next.t('message.something-went-wrong')}</h1>
					</main>
					{config.type.release === 'prod' ? (
						<p className="text-sm sm:text-lg text-gray-700 max-w-[90%] sm:max-w-[500px]">{i18next.t('message.please-refresh-or-try-later')}</p>
					) : (
						<p className="bg-black text-red-500 py-4 px-2 mt-4 text-xs sm:text-sm w-full sm:w-[80%] max-w-[800px] overflow-auto rounded">{this.state.error.stack}</p>
					)}
				</React.Fragment>
			);
		}

		return this.props.children;
	}
}

export default ErrorBoundary;
