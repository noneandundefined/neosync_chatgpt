package usecase

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/store/postgres/store"
)

type UseCase struct {
	Device interface {
		Update_Configuration(ctx context.Context, deviceId uint64, cfg []byte) error
	}
}

func NewUseCase(db *sql.DB, store store.Storage) UseCase {
	return UseCase{
		Device: &DeviceUseCase{store, db},
	}
}
