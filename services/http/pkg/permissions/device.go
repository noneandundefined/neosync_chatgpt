package permissions

import "neomatica/neosync/infra/store/postgres/models"

func HasDeviceAccount(userUUID *string) bool {
	return userUUID != nil && *userUUID != ""
}

func DeviceBelongsToUser(deviceOwner models.UserOwner, user *models.UserAuth) bool {
	if deviceOwner.GetUserUUID() == nil {
		return false
	}

	if *deviceOwner.GetUserUUID() == user.UserContact.UserUUID {
		return true
	}

	if user.ParentUUID != nil && *deviceOwner.GetUserUUID() == *user.ParentUUID {
		return true
	}

	return false
}
