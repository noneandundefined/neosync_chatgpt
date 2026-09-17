package usecase

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/store/postgres/store"
	"time"
)

type DeviceUseCase struct {
	store store.Storage
	db    *sql.DB
}

func (u *DeviceUseCase) Update_Configuration(ctx context.Context, deviceId uint64, cfg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := u.store.Configurations.Update_ConfigurationByDeviceId(ctx, deviceId, cfg); err != nil {
		return err
	}

	return nil
}
