package store

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/pgqx"
	"time"

	"github.com/lib/pq"
)

type DeviceStore struct {
	db *sql.DB
}

func (s *DeviceStore) Create_Device(ctx context.Context, tx *sql.Tx, device *models.Device) (uint64, error) {
	var id uint64

	query := `
		INSERT INTO devices (imei, activated)
		VALUES ($1, $2)
		ON CONFLICT (imei) DO UPDATE SET activated = EXCLUDED.activated
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, device.IMEI, true).Scan(&id)
	if err != nil {
		logger.Error("Create_Device imei={%s}: Failed to execute sql: %s", device.IMEI, err.Error())
		return 0, tcperrors.HandleSQLError(ctx, err)
	}

	return id, nil
}

func (s *DeviceStore) Create_DeviceConf(ctx context.Context, tx *sql.Tx, device *models.DeviceConf) error {
	query := `
		INSERT INTO device_confs (device_id, password, request_configuration_on_connect)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, device.DeviceID, device.Password, device.RequestConfigurationOnConnect)
	if err != nil {
		logger.Error("Create_DeviceConf device_id={%d}: Failed to execute sql: %s", device.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *DeviceStore) Get_DeviceFullByImei(ctx context.Context, imei string) (*models.Device_DeviceConf_Sync, error) {
	query := `
		SELECT
			devices.*,
			COALESCE(user_cores.force_neosync_configuration_priority, false) AS force_neosync_configuration_priority,
			device_confs.request_configuration_on_connect,
			COALESCE(configurations.updated_at, NOW()) as cfg_updated_at,
			COALESCE(configurations.cfg_hash, 0) as cfg_hash,
			COALESCE(configurations.cfg_data, ''::bytea) as cfg_data,
			COALESCE(configurations.cfg_sync_status, 'idle') as cfg_sync_status,
			configurations.cfg_sync_error as cfg_sync_error,
			COALESCE(configurations.cfg_pushed_hash, 0) as cfg_pushed_hash,
			syncs.firmware_version,
			syncs.cfg_version,
			syncs.last_mod_time,
			syncs.cfg_hash as sync_cfg_hash
		FROM devices
		LEFT JOIN user_cores ON devices.user_uuid = user_cores.user_uuid
		LEFT JOIN device_confs ON devices.id = device_confs.device_id
		LEFT JOIN configurations ON devices.id = configurations.device_id
		LEFT JOIN syncs ON devices.id = syncs.device_id
		WHERE devices.imei = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.Device_DeviceConf_Sync](ctx, s.db, query, imei)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceFullByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return device, nil
}

func (s *DeviceStore) Get_DeviceByImei(ctx context.Context, imei string) (*models.Device, error) {
	query := `
		SELECT * FROM devices WHERE imei = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.Device](ctx, s.db, query, imei)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return device, nil
}

func (s *DeviceStore) Get_DeviceRequestConfigurationOnConnectByImei(ctx context.Context, imei string) (*models.RequestConfigurationOnConnect, error) {
	query := `
		SELECT
			devices.imei,
			device_confs.request_configuration_on_connect
		FROM devices
		INNER JOIN device_confs ON devices.id = device_confs.device_id
		WHERE devices.imei = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.RequestConfigurationOnConnect](ctx, s.db, query, imei)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceRequestConfigurationOnConnectByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return device, nil
}

func (s *DeviceStore) Update_DeviceStatusNotActive(ctx context.Context, connectedImeis []string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, `
		UPDATE devices SET status = false
		WHERE status = true AND NOT (imei = ANY($1))
	`, pq.Array(connectedImeis)); err != nil {
		logger.Error("Update_DeviceStatusNotActive: Failed to execute sql: %s", err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	if len(connectedImeis) == 0 {
		return nil
	}

	if _, err := s.db.ExecContext(ctx, `
		UPDATE devices SET status = true
		WHERE status = false AND imei = ANY($1)
	`, pq.Array(connectedImeis)); err != nil {
		logger.Error("Update_DeviceStatusNotActive: Failed to execute sql: %s", err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *DeviceStore) Update_DeviceActivatedByImei(ctx context.Context, imei string, activated bool) error {
	query := `
		UPDATE devices SET activated = $1 WHERE imei = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, activated, imei)
	if err != nil {
		logger.Error("Update_DeviceActivatedByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *DeviceStore) Update_DeviceStatusByImei(ctx context.Context, imei string, status bool) error {
	query := `
		UPDATE devices SET status = $1 WHERE imei = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, status, imei)
	if err != nil {
		logger.Error("Update_DeviceStatusByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *DeviceStore) Update_DeviceModelsByImei(ctx context.Context, model *string, extendedModel *string, imei string) error {
	query := `
		UPDATE devices
		SET
			device_model = COALESCE($1::text, device_model),
			device_extended_model = COALESCE($2::text, device_extended_model)
		WHERE imei = $3
			AND (
				($1::text IS NOT NULL AND device_model IS DISTINCT FROM $1::text)
				OR
				($2::text IS NOT NULL AND device_extended_model IS DISTINCT FROM $2::text)
			)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var modelArg any
	if model != nil {
		modelArg = *model
	}

	var extendedArg any
	if extendedModel != nil {
		extendedArg = *extendedModel
	}

	_, err := s.db.ExecContext(ctx, query, modelArg, extendedArg, imei)
	if err != nil {
		logger.Error("Update_DeviceModelsByImei imei={%s}: Failed to execute sql: %s", imei, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *DeviceStore) Update_DevicePasswordByDeviceId(ctx context.Context, password string, deviceId uint64) error {
	query := `
		UPDATE device_confs SET password = $1 WHERE device_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, password, deviceId)
	if err != nil {
		logger.Error("Update_DevicePasswordByDeviceId device_id={%d}: Failed to execute sql: %s", deviceId, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}
