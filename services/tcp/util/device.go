package util

func HasDeviceOwner(userUUID *string) bool {
	return userUUID != nil && *userUUID != ""
}
