package store

import (
	"context"
	"database/sql"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/tcperrors"
	"time"
)

type AnalyticStore struct {
	db *sql.DB
}

func (s *AnalyticStore) Create_AnalyticsConfiguration(ctx context.Context, analytic *models.AnalyticsConfiguration) error {
	query := `
		INSERT INTO analytic_records (user_uuid, device_id, cfg_hash, rmq_published, device_online, error_reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, analytic.UserUUID, analytic.DeviceID, analytic.CfgHash, analytic.RMQPublished, analytic.DeviceOnline, analytic.ErrorReason)
	if err != nil {
		logger.Error("Create_AnalyticsConfiguration device_id={%d}: Failed to execute sql: %s", analytic.DeviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}

func (s *AnalyticStore) Mark_AnalyticsConfigurationAppliedByDeviceHash(ctx context.Context, userUUID string, deviceID uint64, cfgHash uint32) error {
	query := `
		WITH updated AS (
			UPDATE analytic_records
			SET rmq_published = true, device_online = true, error_reason = NULL
			WHERE id = (
				SELECT id
				FROM analytic_records
				WHERE device_id = $1 AND cfg_hash = $2
				ORDER BY created_at DESC
				LIMIT 1
			)
			RETURNING id
		)
		INSERT INTO analytic_records (user_uuid, device_id, cfg_hash, rmq_published, device_online, error_reason)
		SELECT $3, $1, $2, true, true, NULL
		WHERE NOT EXISTS (SELECT 1 FROM updated)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, deviceID, cfgHash, userUUID)
	if err != nil {
		logger.Error("Mark_AnalyticsConfigurationAppliedByDeviceHash device_id={%d}: Failed to execute sql: %s", deviceID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return nil
}
