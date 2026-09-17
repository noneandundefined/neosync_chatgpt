package permissions

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/types"
)

func ParentUserUuid(authToken *types.AuthToken) string {
	if authToken == nil {
		return ""
	}

	/* User: user see parent devices */
	targetUuid := authToken.User.UserContact.UserUUID

	if authToken.User.Role.Code == constants.Role_AdminL2Support || authToken.User.Role.Code == constants.Role_User {
		targetUuid = *authToken.User.ParentUUID
	}

	return targetUuid
}
