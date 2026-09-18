package usecase

import (
	"context"
	"database/sql"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
)

type UseCase struct {
	Configurations interface{}
	Devices        interface {
		Create_Device() error
		Activated_Device(ctx context.Context, dId uint64) error
	}
	Syncs interface {
		Synchronization(session *memory.SessionMemory, d *models.Device_DeviceConf_Sync) error
		ResolveConfigurationSyncStatus(d *models.Device_DeviceConf_Sync, deviceCfgHash uint32, verifyFailure bool) error
	}
	CompanyTasks interface {
		TryApplyPendingTask(ctx context.Context, session *memory.SessionMemory, d *models.Device_DeviceConf_Sync) (bool, error)
		ObserveDeviceConfiguration(ctx context.Context, deviceID uint64, cfgHash uint32) (bool, error)
	}
}

func NewUseCase(device *adm.ADMDevice, db *sql.DB, store store.Storage, cache *memory.Cache) UseCase {
	return UseCase{
		Configurations: &ConfigurationUseCase{device, store, db, cache},
		Devices:        &DeviceUseCase{device, store, db, cache},
		Syncs:          &SyncUseCase{device, store, db, cache},
		CompanyTasks:   &CompanyTaskUseCase{device, store, db},
	}
}
