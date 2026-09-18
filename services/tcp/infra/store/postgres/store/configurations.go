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
)

type ConfigurationStore struct {
	db *sql.DB
}

func (s *ConfigurationStore) Create_Configuration(ctx context.Context, conf *models.Configuration) error {
	query := `
		INSERT INTO configurations (device_id, cfg_hash)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, conf.DeviceID, conf.CfgHash); err != nil {
		logger.Error("Create_Configuration device_id={%d}: Failed to execute sql: %s", conf.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *ConfigurationStore) Get_ConfigurationByDeviceId(ctx context.Context, deviceId uint64) (*models.Configuration, error) {
	query := `
		SELECT * FROM configurations WHERE device_id = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configuration, err := pgqx.QueryRowContext[models.Configuration](ctx, s.db, query, deviceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ConfigurationByDeviceId device_id={%d}: Failed to execute sql: %s", deviceId, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return configuration, nil
}

func (s *ConfigurationStore) Update_ConfigurationByDeviceId(ctx context.Context, conf *models.Configuration) error {
	query := `
		UPDATE configurations SET cfg_hash = $1, cfg_data = $2 WHERE device_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, conf.CfgHash, conf.CfgData, conf.DeviceID)
	if err != nil {
		logger.Error("Update_ConfigurationByDeviceId device_id={%d}: Failed to execute sql: %s", conf.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *ConfigurationStore) Update_Configuration(ctx context.Context, conf *models.Configuration) error {
	query := `
		UPDATE configurations SET cfg_hash = $1
		WHERE device_id = $2
	` /* UPDATE потому что на http сервере уже создается для не активированно трекера шаблон из нулей */

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, conf.CfgHash, conf.DeviceID)
	if err != nil {
		logger.Error("Update_Configuration device_id={%d}: Failed to execute sql: %s", conf.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *ConfigurationStore) Update_ConfigurationByTx(ctx context.Context, tx *sql.Tx, conf *models.Configuration) error {
	query := `
		UPDATE configurations SET cfg_hash = $1, cfg_data = $2 WHERE device_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, conf.CfgHash, conf.CfgData, conf.DeviceID)
	if err != nil {
		logger.Error("Update_ConfigurationByTx device_id={%d}: Failed to execute sql: %s", conf.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *ConfigurationStore) Update_ConfigurationSyncStatusByDeviceId(ctx context.Context, deviceID uint64, status string, syncError *string) error {
	query := `
		UPDATE configurations
		SET cfg_sync_status = $1, cfg_sync_error = $2
		WHERE device_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, status, syncError, deviceID)
	if err != nil {
		logger.Error("Update_ConfigurationSyncStatusByDeviceId device_id={%d}: Failed to execute sql: %s", deviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *ConfigurationStore) Update_ConfigurationPushedHashByDeviceId(ctx context.Context, deviceID uint64, cfgPushedHash uint32) error {
	query := `
		UPDATE configurations SET cfg_pushed_hash = $1 WHERE device_id = $2 AND cfg_hash = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, cfgPushedHash, deviceID)
	if err != nil {
		logger.Error("Update_ConfigurationPushedHashByDeviceId device_id={%d}: Failed to execute sql: %s", deviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}
