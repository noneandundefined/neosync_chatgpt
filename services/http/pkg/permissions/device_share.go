package permissions

import (
	"context"
	"neomatica/neosync/infra/store/postgres/models"
)

type GroupDeviceShareStore interface {
	AggregatedDevicePermissions(ctx context.Context, memberUserUUID string, deviceID uint64) (canSendCommands, canReadConfig, canEditConfig bool, err error)
}

type ShareLevel int

const (
	ShareViewDevice ShareLevel = iota
	ShareReadConfig
	ShareEditConfig
	ShareSendCommand
)

/* CanAccessDevice — владелец, main role или участник шэринга группы с нужным уровнем прав */
func CanAccessDevice(ctx context.Context, g GroupDeviceShareStore, device models.UserOwner, deviceID uint64, user *models.UserAuth, roleCode string, level ShareLevel) bool {
	if IsMainRole(roleCode) {
		return true
	}

	if DeviceBelongsToUser(device, user) {
		return true
	}

	return deviceShareAccess(ctx, g, user, deviceID, level)
}

func DeviceAccess(ctx context.Context, g GroupDeviceShareStore, device models.UserOwner, deviceID uint64, user *models.UserAuth, roleCode string, level ShareLevel) bool {
	return CanAccessDevice(ctx, g, device, deviceID, user, roleCode, level)
}

func deviceShareAccess(ctx context.Context, g GroupDeviceShareStore, user *models.UserAuth, deviceID uint64, level ShareLevel) bool {
	if g == nil {
		return false
	}

	canSendCommands, canReadConfig, canEditConfig, err := g.AggregatedDevicePermissions(ctx, user.UserContact.UserUUID, deviceID)
	if err != nil || (!canSendCommands && !canReadConfig && !canEditConfig) {
		return false
	}

	switch level {
	case ShareViewDevice:
		return canSendCommands || canReadConfig || canEditConfig
	case ShareReadConfig:
		return canReadConfig || canEditConfig
	case ShareEditConfig:
		return canEditConfig
	case ShareSendCommand:
		return canSendCommands
	default:
		return false
	}
}
