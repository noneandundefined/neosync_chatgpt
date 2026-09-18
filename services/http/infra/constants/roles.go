package constants

const (
	Role_SuperAdmin     = "SUPERADMIN"
	Role_Support        = "SUPPORT"
	Role_AdminL2        = "DEALER"
	Role_AdminL2Support = "DEALER_SUPPORT"
	Role_User           = "USER"
)

var RoleToSubordinateRole = map[string]string{
	Role_Support:        Role_AdminL2,
	Role_AdminL2:        Role_AdminL2Support,
	Role_AdminL2Support: Role_User,
}

func CanCreateRole(creatorRoleID, targetRoleID string) (string, bool) {
	switch creatorRoleID {
	case Role_SuperAdmin:
		return Role_Support, targetRoleID == Role_Support
	case Role_Support:
		return Role_AdminL2, targetRoleID == Role_AdminL2
	case Role_AdminL2:
		return Role_AdminL2Support, targetRoleID == Role_AdminL2Support
	case Role_AdminL2Support:
		return Role_User, targetRoleID == Role_User
	default:
		return "NONE", false
	}
}
