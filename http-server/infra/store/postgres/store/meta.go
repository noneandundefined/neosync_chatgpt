package store

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"
)

type MetaStore struct {
	db *sql.DB
}

func (s *MetaStore) Update_MaintenanceIsActive(ctx context.Context, isActive bool) error {
	query := `
		UPDATE maintenance SET is_active = $1 WHERE id = 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, isActive)
	if err != nil {
		logger.Error("Update_MaintenanceIsActive req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *MetaStore) Get_MaintenanceIsActive(ctx context.Context) (*models.Maintenance, error) {
	query := `
		SELECT * FROM maintenance WHERE id = 1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	maintenance, err := pgqx.QueryRowContext[models.Maintenance](ctx, s.db, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_MaintenanceIsActive req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return maintenance, nil
}
