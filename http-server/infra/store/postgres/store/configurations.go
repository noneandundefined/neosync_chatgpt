package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/adm/admparser"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/pgqx"
	"strings"
	"time"
)

type ConfigurationStore struct {
	db *sql.DB
}

func (s *ConfigurationStore) Create_Configuration(ctx context.Context, tx *sql.Tx, cfg *models.Configuration) error {
	query := `
		INSERT INTO configurations (device_id, cfg_hash)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, cfg.DeviceID, cfg.CfgHash)
	if err != nil {
		logger.Error("Create_Configuration req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *ConfigurationStore) Create_ConfigurationsBatch(ctx context.Context, tx *sql.Tx, cfgs []*models.Configuration) error {
	vStrings := make([]string, 0, len(cfgs))
	vArgs := make([]interface{}, 0, len(cfgs)*2)

	for idx, c := range cfgs {
		vStrings = append(vStrings, fmt.Sprintf("($%d, $%d)", idx*2+1, idx*2+2))
		vArgs = append(vArgs, c.DeviceID, c.CfgHash)
	}

	query := `
		INSERT INTO configurations (device_id, cfg_hash)
		VALUES` + strings.Join(vStrings, ",") + `
	`

	_, err := tx.ExecContext(ctx, query, vArgs...)
	if err != nil {
		logger.Error("Create_ConfigurationsBatch req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *ConfigurationStore) Get_ColumnTableConfiguration(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'configurations'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configurationColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableConfiguration req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configurationColumns, nil
}

func (s *ConfigurationStore) Get_ConfigurationByDeviceId(ctx context.Context, id uint64) (*models.Configuration, error) {
	query := `
		SELECT
		    id,
		    device_id,
		    cfg_hash,
		    cfg_data,
		    cfg_sync_status,
		    cfg_sync_error
		FROM configurations
		WHERE device_id = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configuration, err := pgqx.QueryRowContext[models.Configuration](ctx, s.db, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ConfigurationByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configuration, nil
}

func (s *ConfigurationStore) Get_ReferenceConfigurationByModelAndUserUuid(ctx context.Context, model, userUuid string) (*models.Configuration, error) {
	query := `
		SELECT
		    configurations.id,
		    configurations.device_id,
		    configurations.cfg_hash,
		    configurations.cfg_data,
		    configurations.cfg_sync_status,
		    configurations.cfg_sync_error
		FROM configurations
		INNER JOIN devices ON devices.id = configurations.device_id
		WHERE devices.device_model = $1
		  AND (devices.user_uuid = $2 OR devices.owner_uuid = $2)
		  AND devices.activated = true
		  AND configurations.cfg_data IS NOT NULL
		ORDER BY configurations.updated_at DESC
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configuration, err := pgqx.QueryRowContext[models.Configuration](ctx, s.db, query, model, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ReferenceConfigurationByModelAndUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configuration, nil
}

func (s *ConfigurationStore) Get_ConfigurationHistoriesByDeviceId(ctx context.Context, id uint64) ([]models.ConfigurationHistory, error) {
	query := `
		SELECT id, apply_at, device_id, cfg_hash
		FROM configuration_history
		WHERE device_id = $1
		ORDER BY apply_at DESC
		LIMIT 155
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configurationHistories, err := pgqx.QueryContext[models.ConfigurationHistory](ctx, s.db, query, id)
	if err != nil {
		logger.Error("Get_ConfigurationHistoriesByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configurationHistories, nil
}

func (s *ConfigurationStore) Get_ConfigurationHistoryByCfgHash(ctx context.Context, cfgHash uint32) (*models.ConfigurationHistory, error) {
	query := `
		SELECT *
		FROM configuration_history
		WHERE cfg_hash = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configurationHistory, err := pgqx.QueryRowContext[models.ConfigurationHistory](ctx, s.db, query, cfgHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ConfigurationHistoryByCfgHash req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configurationHistory, nil
}

func (s *ConfigurationStore) Update_ConfigurationByDeviceId(ctx context.Context, id uint64, cfg []byte) error {
	query := `
		UPDATE configurations SET cfg_hash = $1, cfg_data = $2 WHERE device_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, admparser.GetCfgHash(cfg), cfg, id)
	if err != nil {
		logger.Error("Update_ConfigurationByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_ConfigurationByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
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
		logger.Error("Update_ConfigurationSyncStatusByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
