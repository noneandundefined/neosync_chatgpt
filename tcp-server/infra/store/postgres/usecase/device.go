package usecase

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/util"
	"strings"
)

type DeviceUseCase struct {
	Device *adm.ADMDevice
	Store  store.Storage
	Db     *sql.DB
	Cache  *memory.Cache
}

func (u *DeviceUseCase) Create_Device() error {
	ctx := context.Background()

	if u.Device.Imei == nil {
		return errors.New("device IMEI is null")
	}

	imei := util.PrepareImei(*u.Device.Imei)

	device, err := u.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return err
	}

	if device != nil {
		/* If device is activated */
		if device.Activated != nil && *device.Activated {
			return nil
		}

		/* If device is not activated */
		return u.Activated_Device(ctx, device.ID)
	}

	var dId uint64
	errTx := store.WithTx(u.Db, ctx, func(tx *sql.Tx) error {
		/* Create device */
		deviceId, err := u.Store.Devices.Create_Device(ctx, tx, &models.Device{
			IMEI: imei,
		})
		if err != nil {
			return err
		}

		dId = deviceId

		/* Create deviceConf */
		if err := u.Store.Devices.Create_DeviceConf(ctx, tx, &models.DeviceConf{
			DeviceID:                      deviceId,
			Password:                      strings.TrimSpace(*u.Device.Pass),
			RequestConfigurationOnConnect: false,
		}); err != nil {
			return err
		}

		syncModel := &models.Sync{
			DeviceID:        dId,
			FirmwareVersion: u.Device.FirmwareVersion,
			CfgVersion:      u.Device.CfgVersion,
			LastModTime:     u.Device.LastModTime,
			CfgHash:         u.Device.CfgHash,
		}

		return u.Store.Syncs.Create_Sync(ctx, tx, syncModel)
	})

	if errTx != nil {
		return errTx
	}

	cfgModel := &models.Configuration{
		DeviceID: dId,
		CfgHash:  u.Device.CfgHash,
		CfgData:  nil,
	}

	return u.Store.Configurations.Create_Configuration(ctx, cfgModel)
}

func (u *DeviceUseCase) Activated_Device(ctx context.Context, dId uint64) error {
	if u.Device.Imei == nil {
		return errors.New("device IMEI is null")
	}

	imei := util.PrepareImei(*u.Device.Imei)

	if err := u.Store.Devices.Update_DeviceActivatedByImei(context.Background(), imei, true); err != nil {
		return err
	}

	errTx := store.WithTx(u.Db, ctx, func(tx *sql.Tx) error {
		syncModel := &models.Sync{
			DeviceID:        dId,
			FirmwareVersion: u.Device.FirmwareVersion,
			CfgVersion:      u.Device.CfgVersion,
			LastModTime:     u.Device.LastModTime,
			CfgHash:         u.Device.CfgHash,
		}

		if err := u.Store.Syncs.Update_SyncByTx(ctx, tx, syncModel); err != nil {
			return err
		}

		cfgModel := &models.Configuration{
			DeviceID: dId,
			CfgHash:  u.Device.CfgHash,
			CfgData:  nil,
		}

		if err := u.Store.Configurations.Update_ConfigurationByTx(ctx, tx, cfgModel); err != nil {
			return err
		}

		return nil
	})

	return errTx
}
