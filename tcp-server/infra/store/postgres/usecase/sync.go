package usecase

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/util"
	"time"
)

type SyncUseCase struct {
	Device *adm.ADMDevice
	Store  store.Storage
	Db     *sql.DB
	Cache  *memory.Cache
}

const (
	toDevice   = 0x01
	fromDevice = 0x02
	equal      = 0x00
)

func configurationLastModTime(d *models.Device_DeviceConf_Sync) uint32 {
	if d == nil {
		return 0
	}

	if ts := common.GetCfgLastModTime(d.CfgData); ts != 0 {
		return ts
	}

	return d.LastModTime
}

/* 0x01 - toDevice | 0x02 - fromDevice | 0x00 - equal */
func (u *SyncUseCase) Check_LastModTime(last_mod_time_db, last_mod_time_device uint32) uint8 {
	if last_mod_time_db > last_mod_time_device {
		return toDevice
	}

	if last_mod_time_device > last_mod_time_db {
		return fromDevice
	}

	return equal
}

func EffectiveConfigurationHash(d *models.Device_DeviceConf_Sync) uint32 {
	if d == nil {
		return 0
	}

	if d.CfgHash != 0 {
		return d.CfgHash
	}

	if len(d.CfgData) == 0 {
		return 0
	}

	return common.GetCfgHash(d.CfgData)
}

func IsConfigurationMissing(d *models.Device_DeviceConf_Sync) bool {
	if d == nil || len(d.CfgData) == 0 {
		return true
	}

	return EffectiveConfigurationHash(d) == 0
}

func (u *SyncUseCase) configurationSyncDirection(d *models.Device_DeviceConf_Sync) uint8 {
	if u.Device.CfgHash == 0 || u.Device.CfgHash == EffectiveConfigurationHash(d) {
		return equal
	}
	if d.ForceNeosyncConfigurationPriority {
		return toDevice
	}
	return u.Check_LastModTime(configurationLastModTime(d), u.Device.LastModTime)
}

func ShouldPreserveNeosyncConfiguration(d *models.Device_DeviceConf_Sync, deviceCfgHash uint32) bool {
	return d != nil && d.ForceNeosyncConfigurationPriority && !IsConfigurationMissing(d) && EffectiveConfigurationHash(d) != deviceCfgHash
}

func (u *SyncUseCase) Synchronization(session *memory.SessionMemory, d *models.Device_DeviceConf_Sync) error {
	if u.Device == nil || u.Device.Imei == nil {
		return errors.New("device or imei is not found")
	}

	if d == nil {
		return errors.New("device payload is not found")
	}

	imei := util.PrepareImei(*u.Device.Imei)

	if session == nil || imei == "" {
		return errors.New("session or imei is not found")
	}

	dbCfgHash := EffectiveConfigurationHash(d)

	if IsConfigurationMissing(d) {
		if u.Device.CfgHash != 0 {
			session.TransitData(imei, constants.ADM_RC_TYPE_GET_CFG, nil)
		}

		return nil
	}

	if u.Device.CfgHash == 0 {
		return nil
	}

	if u.Device.CfgHash == dbCfgHash {
		return nil
	}

	dbLastMod := configurationLastModTime(d)
	alreadyPushed := d.CfgPushedHash != 0 && d.CfgPushedHash == dbCfgHash

	switch u.configurationSyncDirection(d) {
	case toDevice: /* Записываем в терминал */
		if !util.HasDeviceOwner(d.UserUUID) {
			logger.Info("Synchronization imei={%s}: skip SET_CFG, device has no owner", imei)
			return nil
		}

		if alreadyPushed && (!d.ForceNeosyncConfigurationPriority || d.CfgSyncStatus != constants.CFG_SYNC_STATUS_CONFIRMED) {
			if d.CfgSyncStatus == constants.CFG_SYNC_STATUS_FAILED {
				logger.Info("Synchronization imei={%s}: already pushed cfg_hash={%d}, skip overwrite", imei, dbCfgHash)
				return nil
			}

			if !d.ForceNeosyncConfigurationPriority {
				logger.Info("Synchronization imei={%s}: already pushed cfg_hash={%d}, pull once instead of SET_CFG", imei, dbCfgHash)
				session.TransitData(imei, constants.ADM_RC_TYPE_GET_CFG, nil)
			}

			errMsg := constants.CFG_SYNC_ERROR_HASH_MISMATCH
			return u.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(context.Background(), d.ID, constants.CFG_SYNC_STATUS_FAILED, &errMsg)
		}

		if err := u.markConfigurationSyncPending(d.ID); err != nil {
			return err
		}

		logger.Info("Synchronization imei={%s}: Pushing configuration to device db_hash={%d} device_hash={%d} db_last_mod={%d} device_last_mod={%d} bytes={%d}", imei, dbCfgHash, u.Device.CfgHash, dbLastMod, u.Device.LastModTime, len(d.CfgData))

		if err := session.TransitData(imei, constants.ADM_RC_TYPE_SET_CFG, d.CfgData); err != nil {
			errMsg := err.Error()
			if statusErr := u.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(context.Background(), d.ID, constants.CFG_SYNC_STATUS_FAILED, &errMsg); statusErr != nil {
				logger.Error("Synchronization imei={%s}: failed to persist SET_CFG delivery error: %s", imei, statusErr.Error())
			}
			return err
		}

		session.MarkConfigurationPushed(d.ID, d.CfgData, dbCfgHash, imei)

		d.CfgPushedHash = dbCfgHash
		return nil

	case fromDevice: /* Записываем в БД */
		logger.Info("Synchronization imei={%s}: Pulling configuration from device db_hash={%d} device_hash={%d} db_last_mod={%d} device_last_mod={%d}", imei, dbCfgHash, u.Device.CfgHash, dbLastMod, u.Device.LastModTime)
		session.TransitData(imei, constants.ADM_RC_TYPE_GET_CFG, nil)
		return nil

	case equal: /* Не затираем трекер */
		logger.Info("Synchronization imei={%s}: last_mod equal, skip overwrite db_hash={%d} device_hash={%d} last_mod={%d}", imei, dbCfgHash, u.Device.CfgHash, dbLastMod)
		return nil
	}

	return nil
}

func (u *SyncUseCase) markConfigurationSyncPending(deviceID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, deviceID, constants.CFG_SYNC_STATUS_PENDING, nil)
}

func (u *SyncUseCase) ResolveConfigurationSyncStatus(d *models.Device_DeviceConf_Sync, deviceCfgHash uint32, verifyFailure bool) error {
	if d == nil {
		return errors.New("device payload is not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if deviceCfgHash == EffectiveConfigurationHash(d) {
		return u.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, d.ID, constants.CFG_SYNC_STATUS_CONFIRMED, nil)
	}

	if verifyFailure && d.CfgSyncStatus == constants.CFG_SYNC_STATUS_PENDING {
		errMsg := constants.CFG_SYNC_ERROR_HASH_MISMATCH
		return u.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, d.ID, constants.CFG_SYNC_STATUS_FAILED, &errMsg)
	}

	return nil
}
