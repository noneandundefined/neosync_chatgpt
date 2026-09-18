package store

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/pgqx"
	"neomatica/neosync-tcp/pkg/protocol"
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

	_, err := tx.ExecContext(ctx, query, sync.DeviceID, sync.FirmwareVersion, sync.CfgVersion, sync.LastModTime, sync.CfgHash)
	if err != nil {
		logger.Error("Create_Sync device_id={%d}: Failed to execute sql: %s", sync.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *SyncStore) Get_SyncByDeviceId(ctx context.Context, deviceId uint64) (*models.Sync, error) {
	query := `
		SELECT * FROM syncs WHERE device_id = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	sync, err := pgqx.QueryRowContext[models.Sync](ctx, s.db, query, deviceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_SyncByDeviceId device_id={%d}: Failed to execute sql: %s", deviceId, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return sync, nil
}

func (s *SyncStore) Update_Sync(ctx context.Context, sync *models.Sync) error {
	query := `
		UPDATE syncs SET firmware_version = $1, cfg_version = $2, last_mod_time = $3, cfg_hash = $4
		WHERE device_id = $5
	` /* UPDATE потому что на http сервере уже создается для не активированно трекера шаблон из нулей */

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, sync.FirmwareVersion, sync.CfgVersion, sync.LastModTime, sync.CfgHash, sync.DeviceID)
	if err != nil {
		logger.Error("Update_Sync device_id={%d}: Failed to execute sql: %s", sync.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *SyncStore) Update_SyncByDeviceId(ctx context.Context, sync *models.Sync) error {
	query := `
		UPDATE syncs SET firmware_version = $1, cfg_version = $2, last_mod_time = $3, cfg_hash = $4
		WHERE device_id = $5
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, sync.FirmwareVersion, sync.CfgVersion, sync.LastModTime, sync.CfgHash, sync.DeviceID)
	if err != nil {
		logger.Error("Update_SyncByDeviceId device_id={%d}: Failed to execute sql: %s", sync.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *SyncStore) Update_SyncByTx(ctx context.Context, tx *sql.Tx, sync *models.Sync) error {
	query := `
		UPDATE syncs SET firmware_version = $1, cfg_version = $2, last_mod_time = $3, cfg_hash = $4
		WHERE device_id = $5
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, sync.FirmwareVersion, sync.CfgVersion, sync.LastModTime, sync.CfgHash, sync.DeviceID)
	if err != nil {
		logger.Error("Update_SyncByTx device_id={%d}: Failed to execute sql: %s", sync.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *SyncStore) Update_SyncTimeAndHashByDeviceID(ctx context.Context, protocol protocol.Protocol, deviceId uint64, pkt []byte) error {
	query := `
		UPDATE syncs SET last_mod_time = $1, cfg_hash = $2
		WHERE device_id = $3
		RETURNING id
	`

	/* Gets LastModTime and CfgHash */
	sync, err := protocol.Parse_SyncPacket(pkt)
	if err != nil {
		return tcperrors.HandleSQLError(ctx, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, sync.LastModTime, sync.CfgHash, deviceId); err != nil {
		logger.Error("Update_SyncTimeAndHashByDeviceID device_id={%d}: Failed to execute sql: %s", deviceId, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *SyncStore) Complete_FirmwareUpdateByDeviceId(ctx context.Context, deviceID uint64) error {
	query := `
		UPDATE syncs
		SET firmware_update_status = $1,
		    firmware_updated_at = NOW()
		WHERE device_id = $2 AND firmware_update_status = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, constants.FW_UPDATE_STATUS_IDLE, deviceID, constants.FW_UPDATE_STATUS_PENDING); err != nil {
		logger.Error("Complete_FirmwareUpdateByDeviceId device_id={%d}: Failed to execute sql: %s", deviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}
