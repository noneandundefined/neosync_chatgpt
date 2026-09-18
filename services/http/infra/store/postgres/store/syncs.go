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

type SyncStore struct {
	db *sql.DB
}

func (s *SyncStore) Create_Sync(ctx context.Context, tx *sql.Tx, sync *models.Sync) error {
	query := `
		INSERT INTO syncs (device_id, firmware_version, cfg_version, last_mod_time, cfg_hash)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, sync.DeviceID, sync.FirmwareVersion, sync.CfgVersion, sync.LastModTime, sync.CfgHash); err != nil {
		logger.Error("Create_Sync req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *SyncStore) Create_SyncsBatch(ctx context.Context, tx *sql.Tx, syncs []*models.Sync) error {
	vStrings := make([]string, 0, len(syncs))
	vArgs := make([]interface{}, 0, len(syncs)*5)

	for idx, s := range syncs {
		vStrings = append(vStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", idx*5+1, idx*5+2, idx*5+3, idx*5+4, idx*5+5))
		vArgs = append(vArgs, s.DeviceID, s.FirmwareVersion, s.CfgVersion, s.LastModTime, s.CfgHash)
	}

	query := `
		INSERT INTO syncs (device_id, firmware_version, cfg_version, last_mod_time, cfg_hash)
		VALUES` + strings.Join(vStrings, ",") + `
	`

	if _, err := tx.ExecContext(ctx, query, vArgs...); err != nil {
		logger.Error("Create_SyncsBatch req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *SyncStore) Get_ColumnTableSync(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'syncs'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	syncColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableSync req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return syncColumns, nil
}

func (s *SyncStore) Get_SyncByDeviceId(ctx context.Context, id uint64) (*models.Sync, error) {
	query := `
		SELECT * FROM syncs WHERE device_id = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	sync, err := pgqx.QueryRowContext[models.Sync](ctx, s.db, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_SyncByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return sync, nil
}

func (s *SyncStore) Update_SyncByDeviceId(ctx context.Context, tx *sql.Tx, id uint64, cfg []byte) error {
	query := `
		UPDATE syncs SET cfg_hash = $1, last_mod_time = $2 WHERE device_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upd, err := tx.ExecContext(ctx, query, admparser.GetCfgHash(cfg), admparser.GetCfgLastModTime(cfg), id)
	if err != nil {
		logger.Error("Update_SyncByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_SyncByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *SyncStore) Update_SyncFirmwareVersionUpd(ctx context.Context, firmware uint16, model string) error {
	query := `
		UPDATE syncs SET firmware_version_upd = $1
		FROM devices
		WHERE syncs.device_id = devices.id
		AND devices.device_model = $2
		AND syncs.firmware_version_upd IS DISTINCT FROM $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, firmware, model)
	if err != nil {
		logger.Error("Update_SyncFirmwareVersionUpd req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	_, err = upd.RowsAffected()
	if err != nil {
		logger.Error("Update_SyncFirmwareVersionUpd req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *SyncStore) Update_FirmwareUpdateStatusByDeviceId(ctx context.Context, deviceID uint64, status string) error {
	query := `
		UPDATE syncs SET firmware_update_status = $1 WHERE device_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, status, deviceID); err != nil {
		logger.Error("Update_FirmwareUpdateStatusByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
