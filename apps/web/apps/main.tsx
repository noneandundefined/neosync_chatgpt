import { Suspense, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { HelmetProvider } from 'react-helmet-async';
import { ErrorBoundary } from 'react-error-boundary';
import { BrowserRouter, useLocation } from 'react-router-dom';

import { ROUTES } from '@/constants/constants';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider, useQueryClient } from '@tanstack/react-query';

import Router from '@/router';
import Fallback from './components/Fallback/Fallback';

import { ToastContainer } from 'react-toastify';
// import { basicMetaTcpHealth } from './rest/metaAPI';
import TechnicalWork from './components/AppContainer/TechnicalWork';
import { RoleProvider } from './context/RoleContext/useRoleContext';

import ErrorFallback from './components/AppContainer/ErrorFallback';

import '@/utils/i18n';
import 'rc-slider/assets/index.css';
import { useAntiDebug } from './hooks/useAntiDebug';
import { ModalProvider } from './context/useModalContext';
import { TcpHealthContext } from './context/useTcpHealth';
// import { useHandleServer } from './hooks/Server/useHandleServer';
import { ModalMoveProvider } from './context/useModalMoveContext';
import CookieConsentBanner from './components/AppContainer/CookieConsentBanner';
import TerminalTemporarilyLimited from './components/AppContainer/TerminalTemporarilyLimited';
import ProductAnalyticsTracker from './components/AppContainer/ProductAnalyticsTracker';
// import NotificationProvider from './components/Notification/NotificationProvider';

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 1000 * 30,
			refetchOnWindowFocus: false,
			refetchOnReconnect: true,
			refetchOnMount: true,
		},
	},
});

const App = () => {
	useAntiDebug({
		redirectUrl: '/',
	});

	if (window.__APP_CONFIG__?.TECHNICAL_WORK === true) {
		return <TechnicalWork />;
	}

	return (
		<ErrorBoundary fallback={<ErrorFallback />}>
			<Suspense fallback={<Fallback />}>
				<Router />
			</Suspense>
		</ErrorBoundary>
	);
};

const Root = () => {
	const location = useLocation();
	const queryClient = useQueryClient();

	// const { data, loading } = useHandleServer(['respMetaTcpHealth'], basicMetaTcpHealth, { refetchInterval: 30000 });

	useEffect(() => {
		const isAuth =
			location.pathname.includes(ROUTES.SIGNIN) || location.pathname.includes(ROUTES.HOW_GET_ACCESS) || location.pathname.includes(ROUTES.RESET_PASSWORD_REQ) || location.pathname.includes(ROUTES.RESET_PASSWORD_NEW);

		if (isAuth) {
			queryClient.removeQueries({
				predicate: (query) => query.queryKey[0] !== 'respMetaTcpHealth',
			});
		}
	}, [location.pathname, queryClient]);

	const tcpUnavailable = false;

	return (
		<TcpHealthContext.Provider value={!tcpUnavailable}>
			<ProductAnalyticsTracker />
			{tcpUnavailable && <TerminalTemporarilyLimited />}

			<div className={tcpUnavailable ? 'pt-6' : 'pt-0'}>
				<ToastContainer position="top-center" autoClose={3850} hideProgressBar={false} newestOnTop={false} closeOnClick={true} closeButton={true} theme="light" limit={2} />
				<App />

				<CookieConsentBanner />
			</div>
		</TcpHealthContext.Provider>
	);
};

createRoot(document.getElementById('root')!).render(
	<BrowserRouter>
		<RoleProvider>
			<HelmetProvider>
				<QueryClientProvider client={queryClient}>
					<ModalProvider>
						<ModalMoveProvider>
							<Root />

							<ReactQueryDevtools initialIsOpen={true} />
						</ModalMoveProvider>
					</ModalProvider>
				</QueryClientProvider>
			</HelmetProvider>
		</RoleProvider>
	</BrowserRouter>
);
