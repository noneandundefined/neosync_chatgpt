import i18next from 'i18next';
import { ROUTES } from '@/constants/constants';
import { useLocation } from 'react-router-dom';
import { basicUserGetMeLoginWithRoleCode } from '@/rest/userAPI';
import { isRoleCode, RoleCode } from '@/constants/Roles.constant';
import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from 'react';

interface RoleContextType {
	login: string | null;
	role: RoleCode | null;
	setLogin: (login: string) => void;
	setRole: (role: RoleCode) => void;
}

const RoleContext = createContext<RoleContextType | undefined>(undefined);

export const useRole = (): RoleContextType => {
	const context = useContext(RoleContext);
	if (!context) {
		throw new Error(i18next.t('message.get-role-error'));
	}

	return context;
};

interface RoleProviderProps {
	children: ReactNode;
}

export const RoleProvider = ({ children }: RoleProviderProps) => {
	const [role, setRole] = useState<RoleCode | null>(null);
	const [login, setLogin] = useState<string | null>(null);

	const location = useLocation();

	const resetUser = useCallback(() => {
		setRole(null);
		setLogin(null);
	}, []);

	useEffect(() => {
		const isAuth =
			location.pathname.includes(ROUTES.SIGNIN) || location.pathname.includes(ROUTES.HOW_GET_ACCESS) || location.pathname.includes(ROUTES.RESET_PASSWORD_REQ) || location.pathname.includes(ROUTES.RESET_PASSWORD_NEW);

		if (isAuth) {
			resetUser();
			return;
		}

		const handle = async () => {
			const response = await basicUserGetMeLoginWithRoleCode();

			if (response && isRoleCode(response.role_code)) {
				setRole(response.role_code);
				setLogin(response.login);
			} else {
				setRole(null);
				setLogin(null);
			}
		};

		handle();
	}, [location.pathname, resetUser]);

	return <RoleContext.Provider value={{ login, role, setLogin, setRole }}>{children}</RoleContext.Provider>;
};
