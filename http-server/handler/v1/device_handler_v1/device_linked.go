package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
)

func deviceIsOnCurrentAccount(exist models.Device, currentUserUuid, accountUserUuid, roleCode string) bool {
	if exist.UserUUID == nil {
		return false
	}

	if roleCode == constants.Role_User {
		return exist.OwnerUUID != nil && *exist.OwnerUUID == currentUserUuid
	}

	if *exist.UserUUID == currentUserUuid || *exist.UserUUID == accountUserUuid {
		return true
	}

	return exist.OwnerUUID != nil && *exist.OwnerUUID == currentUserUuid
}

func linkedDeviceError(tr locale.Translator, exist models.Device, currentUserUuid, accountUserUuid, roleCode string) string {
	if exist.UserUUID == nil {
		return ""
	}

	if deviceIsOnCurrentAccount(exist, currentUserUuid, accountUserUuid, roleCode) {
		return tr.TErr("device-already-in-your-account")
	}

	return tr.TErr("device-already-linked")
}
