export const ROLES = {
	SUPERADMIN: 'SUPERADMIN',
	SUPPORT: 'SUPPORT',
	DEALER: 'DEALER', // ADMINL2
	DEALER_SUPPORT: 'DEALER_SUPPORT',
	USER: 'USER',
} as const;

export const ROLES_ARRAY = ['SUPPORT', 'DEALER', 'DEALER_SUPPORT', 'USER'];

export function getAvailableRoles(currentRole: RoleCode | null) {
	if (currentRole === ROLES.SUPERADMIN) {
		return ROLES_ARRAY;
	}

	const index = ROLES_ARRAY.indexOf(currentRole ?? '');
	if (index === -1) return [];

	return ROLES_ARRAY.slice(index);
}

export type RoleCode = (typeof ROLES)[keyof typeof ROLES];

export const ROLEs_AND_NAME: Record<string, string> = {
	SUPERADMIN: 'role-super-admin',
	SUPPORT: 'role-support',
	DEALER: 'role-dealer',
	DEALER_SUPPORT: 'role-dealer-support',
	USER: 'role-user',
};

const ROLE_PERMISSIONS: Record<RoleCode, readonly string[]> = {
	SUPERADMIN: ['view:users', 'transfer:users', 'create:device', 'import:device', 'user:change_password', 'view:analytics'], // SUPERADMIN
	SUPPORT: ['view:users', 'transfer:users', 'create:device', 'import:device', 'view:analytics'], // SUPPORT
	DEALER: ['view:users', 'transfer:users', 'create:device', 'import:device', 'view:groups:child'], // DEALER
	DEALER_SUPPORT: ['view:users', 'transfer:users', 'create:device', 'import:device'], // DEALER_SUPPORT
	USER: ['create:device', 'import:device'], // USER
} as const;

export function isRoleCode(value: string): value is RoleCode {
	const ALLOWED_ROLEs: RoleCode[] = ['SUPERADMIN', 'SUPPORT', 'DEALER', 'DEALER_SUPPORT', 'USER'];
	return ALLOWED_ROLEs.includes(value as RoleCode);
}

export function hasPermission(roleId: RoleCode, permission: string): boolean {
	return (ROLE_PERMISSIONS[roleId] ?? []).includes(permission) ?? false;
}
