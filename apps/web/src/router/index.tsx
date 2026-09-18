import { useEffect, useState } from 'react';
import { getAuthState } from '@/private-route';
import { ROUTES } from '@/constants/constants';
import { useTranslation } from 'react-i18next';
import config, { CustomRouteConfig } from './config';
import Fallback from '@/components/Fallback/Fallback';
import { Navigate, Route, Routes, useLocation } from 'react-router-dom';

const CustomRoute: React.FC<CustomRouteConfig> = ({ loginRequired = true, redirectIfLogged = true, component: Component, title }) => {
	const { t } = useTranslation();
	const location = useLocation();

	const [isLoading, setIsLoading] = useState<boolean>(true);
	const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);

	useEffect(() => {
		if (title) {
			document.title = `${t(title)} - NeoSync`;
		} else {
			document.title = `NeoSync`;
		}
	}, [title, t]);

	useEffect(() => {
		getAuthState().then((logged) => {
			setIsLoggedIn(logged);
			setIsLoading(false);
		});
	}, []);

	if (isLoading) {
		return <Fallback />;
	}

	if (loginRequired) {
		if (isLoggedIn) {
			return <Component />;
		} else {
			return <Navigate to={ROUTES.SIGNIN} state={{ from: location }} replace />;
		}
	} else {
		return isLoggedIn && redirectIfLogged ? <Navigate to="/" replace /> : <Component />;
	}
};

const AppRouter: React.FC = () => {
	return (
		<Routes>
			{config.map((route) => (
				<Route key={route.path} path={route.path} element={<CustomRoute {...route} />} />
			))}
		</Routes>
	);
};

export default AppRouter;
