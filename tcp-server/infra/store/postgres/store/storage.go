package store

import (
	"context"
	"database/sql"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/pkg/protocol"
)

type Storage struct {
	Analytics interface {
		Create_AnalyticsConfiguration(ctx context.Context, analytic *models.AnalyticsConfiguration) error

		Mark_AnalyticsConfigurationAppliedByDeviceHash(ctx context.Context, userUUID string, deviceID uint64, cfgHash uint32) error
	}
	Devices interface {
		Create_Device(ctx context.Context, tx *sql.Tx, device *models.Device) (uint64, error)
		Create_DeviceConf(ctx context.Context, tx *sql.Tx, device *models.DeviceConf) error

		Get_DeviceFullByImei(ctx context.Context, imei string) (*models.Device_DeviceConf_Sync, error)
		Get_DeviceByImei(ctx context.Context, imei string) (*models.Device, error)
		Get_DeviceRequestConfigurationOnConnectByImei(ctx context.Context, imei string) (*models.RequestConfigurationOnConnect, error)

		Update_DeviceStatusNotActive(ctx context.Context, connectedImeis []string) error
		Update_DeviceActivatedByImei(ctx context.Context, imei string, activated bool) error
		Update_DeviceStatusByImei(ctx context.Context, imei string, status bool) error
		Update_DeviceModelsByImei(ctx context.Context, model *string, extendedModel *string, imei string) error
		Update_DevicePasswordByDeviceId(ctx context.Context, password string, deviceId uint64) error
	}
	DeviceCommands interface {
		Get_DeviceCommandsByImeiAndSendMode(ctx context.Context, imei, sendMode string) ([]models.DeviceCommandWithExecutions, error)

		Update_DeviceCommandMarkExecution(ctx context.Context, imei, status, command, response string) error
		Claim_DeviceCommandExecution(ctx context.Context, imei, command string) (executionID, commandID uint64, err error)
		Update_DeviceCommandExecutionInProgressByID(ctx context.Context, executionID uint64) error
		Update_DeviceCommandExecutionByID(ctx context.Context, executionID uint64, status, response string) error
		Update_DeviceCommandExecutionTimeoutByID(ctx context.Context, executionID uint64) error
		Update_DeviceCommandsMarkExecutionBatch(ctx context.Context, imei, status, response string) error
		Update_DeviceCommandsResetOnConnectOnDisconnect(ctx context.Context, imei string) error
		Update_DeviceCommandHandleTimeout(ctx context.Context, imei, command string) error
		Update_DeviceCommandStatusProgressByCommand(ctx context.Context, imei, command string) error
		Update_DeviceCommandClearProgressStatus(ctx context.Context, imei string) error
	}
	Configurations interface {
		Create_Configuration(ctx context.Context, conf *models.Configuration) error

		Get_ConfigurationByDeviceId(ctx context.Context, deviceId uint64) (*models.Configuration, error)

		Update_Configuration(ctx context.Context, conf *models.Configuration) error
		Update_ConfigurationByDeviceId(ctx context.Context, conf *models.Configuration) error
		Update_ConfigurationByTx(ctx context.Context, tx *sql.Tx, conf *models.Configuration) error
		Update_ConfigurationSyncStatusByDeviceId(ctx context.Context, deviceID uint64, status string, syncError *string) error
		Update_ConfigurationPushedHashByDeviceId(ctx context.Context, deviceID uint64, cfgPushedHash uint32) error
	}
	Syncs interface {
		Create_Sync(ctx context.Context, tx *sql.Tx, sync *models.Sync) error

		Get_SyncByDeviceId(ctx context.Context, deviceId uint64) (*models.Sync, error)

		Update_SyncByDeviceId(ctx context.Context, sync *models.Sync) error
		Update_SyncByTx(ctx context.Context, tx *sql.Tx, sync *models.Sync) error
		Update_SyncTimeAndHashByDeviceID(ctx context.Context, protocol protocol.Protocol, deviceId uint64, pkt []byte) error
		Update_Sync(ctx context.Context, sync *models.Sync) error

		Complete_FirmwareUpdateByDeviceId(ctx context.Context, deviceID uint64) error
	}
	Companies interface {
		Get_ActiveCompanyTaskByDeviceId(ctx context.Context, deviceID uint64) (*models.CompanyTaskWithCompany, error)
		Get_DueCompanyTasks(ctx context.Context) ([]models.CompanyTaskWithCompany, error)
		Get_CompanyTaskConnectedImeis(ctx context.Context, imeis []string) ([]string, error)

		Update_CompanyTaskDelivery(ctx context.Context, task *models.CompanyTaskWithCompany, delivery *models.CompanyTaskDelivery) (bool, error)
		Update_CompanyTaskObservation(ctx context.Context, taskID uint64, cfgHash uint32) error

		Refresh_CompanyStatus(ctx context.Context, companyID uint64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Analytics:      &AnalyticStore{db},
		Devices:        &DeviceStore{db},
		DeviceCommands: &DeviceCommandStore{db},
		Configurations: &ConfigurationStore{db},
		Syncs:          &SyncStore{db},
		Companies:      &CompanyStore{db},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
