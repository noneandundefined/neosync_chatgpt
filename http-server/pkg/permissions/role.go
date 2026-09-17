package permissions

import "neomatica/neosync/infra/constants"

func IsMainRole(roleCode string) bool {
	return roleCode == constants.Role_SuperAdmin || roleCode == constants.Role_Support
}
