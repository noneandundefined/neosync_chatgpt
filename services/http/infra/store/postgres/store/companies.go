package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"

	"github.com/lib/pq"
)

type CompanyStore struct {
	db *sql.DB
}

func (s *CompanyStore) Create_Company(ctx context.Context, tx *sql.Tx, company *models.Company) (uint64, error) {
	query := `
		INSERT INTO companies (user_uuid, name, status, ttl, existing_task_action, launch_mode, incompatible_action, configuration_source, configuration_source_data, source_metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var id uint64
	if err := tx.QueryRowContext(
		ctx,
		query,
		company.UserUUID,
		company.Name,
		company.Status,
		company.TTL,
		company.ExistingTaskAction,
		company.LaunchMode,
		company.IncompatibleAction,
		company.ConfigurationSource,
		company.ConfigurationSourceData,
		string(company.SourceMetadata),
	).Scan(&id); err != nil {
		logger.Error("Create_Company req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	return id, nil
}

func (s *CompanyStore) Create_CompanyTask(ctx context.Context, companyID, deviceID uint64, cfgData []byte) error {
	query := `
		INSERT INTO company_tasks (company_id, device_id, status, cfg_data, expires_at)
		SELECT $1, $2, 'pending', $3, created_at + ttl * INTERVAL '24 hours' FROM companies WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, companyID, deviceID, cfgData); err != nil {
		logger.Error("Create_CompanyTask company_id={%d} device_id={%d}: Failed to exec sql: %s", companyID, deviceID, err.Error())
		return err
	}

	return nil
}

func (s *CompanyStore) Get_CompanyForProcessing(ctx context.Context, companyID uint64, userUUID string) (*models.Company, error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			user_uuid,
			name,
			status,
			ttl,
			existing_task_action,
			launch_mode,
			incompatible_action,
			configuration_source,
			configuration_source_data
		FROM companies
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var company models.Company
	err := s.db.QueryRowContext(ctx, query, companyID, userUUID).Scan(
		&company.ID,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.UserUUID,
		&company.Name,
		&company.Status,
		&company.TTL,
		&company.ExistingTaskAction,
		&company.LaunchMode,
		&company.IncompatibleAction,
		&company.ConfigurationSource,
		&company.ConfigurationSourceData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_CompanyForProcessing company_id={%d}: Failed to exec sql: %s", companyID, err.Error())
		return nil, err
	}

	return &company, nil
}

func (s *CompanyStore) Update_CompanyStatus(ctx context.Context, companyID uint64, status string) error {
	query := `
		UPDATE companies
		SET status = $2
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, companyID, status); err != nil {
		logger.Error("Update_CompanyStatus company_id={%d}: Failed to exec sql: %s", companyID, err.Error())
		return err
	}

	return nil
}

func (s *CompanyStore) Cancel_PendingCompanyTasksByDeviceIds(ctx context.Context, tx *sql.Tx, deviceIDs []uint64) error {
	if len(deviceIDs) == 0 {
		return nil
	}

	query := `
		UPDATE company_tasks
		SET status = 'cancelled'
		WHERE device_id = ANY($1) AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, `SELECT id FROM companies WHERE id IN (
		SELECT company_id FROM company_tasks WHERE device_id = ANY($1)
		AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
	) ORDER BY id FOR UPDATE`, pq.Array(deviceIDs))
	if err != nil {
		return err
	}

	companyIDs := make([]uint64, 0)

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		companyIDs = append(companyIDs, id)
	}

	rowsErr := rows.Err()
	rows.Close()

	if rowsErr != nil {
		return rowsErr
	}

	if _, err := tx.ExecContext(ctx, query, pq.Array(deviceIDs)); err != nil {
		logger.Error("Cancel_PendingCompanyTasksByDeviceIds req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	for _, companyID := range companyIDs {
		if _, err := tx.ExecContext(ctx, refreshCompanyStatusQuery, companyID); err != nil {
			return err
		}
	}

	return nil
}

func (s *CompanyStore) Get_PendingCompanyTaskDeviceIds(ctx context.Context, deviceIDs []uint64) ([]uint64, error) {
	if len(deviceIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT DISTINCT device_id
		FROM company_tasks
		WHERE device_id = ANY($1) AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, pq.Array(deviceIDs))
	if err != nil {
		logger.Error("Get_PendingCompanyTaskDeviceIds req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []uint64
	for rows.Next() {
		var deviceID uint64
		if err := rows.Scan(&deviceID); err != nil {
			return nil, err
		}

		result = append(result, deviceID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CompanyStore) Get_CompaniesByUserUuid(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Company, int, int, error) {
	allowedColumnSort := map[string]string{
		"name":       "companies.name",
		"status":     "companies.status",
		"created_at": "companies.created_at",
	}

	orderBy := "ORDER BY companies.created_at DESC"
	if col, ok := allowedColumnSort[columnSortKey]; ok {
		dir := "ASC"
		if columnSortDir == "desc" {
			dir = "DESC"
		}

		orderBy = fmt.Sprintf("ORDER BY %s %s", col, dir)
	}

	baseQuery := `
		FROM companies
		LEFT JOIN company_tasks ON company_tasks.company_id = companies.id
		WHERE companies.user_uuid = $1
			AND ($2 = '' OR companies.name ILIKE '%' || $2 || '%')
		GROUP BY companies.id
	`

	countFilteredQuery := `
		SELECT COUNT(*) FROM (
			SELECT companies.id
			` + baseQuery + `
		) AS filtered_companies
	`

	countAllQuery := `
		SELECT COUNT(*) FROM companies WHERE user_uuid = $1
	`

	dataQuery := `
		SELECT
			companies.id,
			companies.created_at,
			companies.updated_at,
			companies.user_uuid,
			companies.name,
			companies.status,
			companies.ttl,
			companies.existing_task_action,
			companies.launch_mode,
			companies.incompatible_action,
			companies.configuration_source,
			COUNT(company_tasks.id)::int AS tasks_count
		` + baseQuery + `
		` + orderBy + `
		LIMIT $3 OFFSET $4
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var totalFiltered int
	if err := s.db.QueryRowContext(ctx, countFilteredQuery, userUuid, search).Scan(&totalFiltered); err != nil {
		logger.Error("Get_CompaniesByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	var totalAll int
	if err := s.db.QueryRowContext(ctx, countAllQuery, userUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_CompaniesByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	companies, err := pgqx.QueryContext[models.Company](ctx, s.db, dataQuery, userUuid, search, limit, offset)
	if err != nil {
		logger.Error("Get_CompaniesByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return companies, totalFiltered, totalAll, nil
}

func (s *CompanyStore) Retry_FailedCompanyTasksByCompanyId(ctx context.Context, companyID uint64) (int64, error) {
	query := `
		UPDATE company_tasks
		SET status = 'pending', error_message = NULL, attempts = 0,
			attempt_token = NULL, lease_until = NULL, confirmation_deadline = NULL,
			next_attempt_at = NULL, last_sent_at = NULL, updated_at = timezone('UTC', now())
		WHERE company_id = $1 AND status = 'failed' AND expires_at > now()
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var lockedID uint64

	if err := tx.QueryRowContext(ctx, `SELECT id FROM companies WHERE id = $1 FOR UPDATE`, companyID).Scan(&lockedID); err != nil {
		return 0, err
	}

	result, err := tx.ExecContext(ctx, query, companyID)
	if err != nil {
		logger.Error("Retry_FailedCompanyTasksByCompanyId company_id={%d}: Failed to exec sql: %s", companyID, err.Error())
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, refreshCompanyStatusQuery, companyID); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return affected, nil
}

func (s *CompanyStore) Cancel_PendingCompanyTasksByCompanyId(ctx context.Context, companyID uint64) (int64, error) {
	query := `
		UPDATE company_tasks
		SET status = 'cancelled', updated_at = timezone('UTC', now())
		WHERE company_id = $1 AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var lockedID uint64

	if err := tx.QueryRowContext(ctx, `SELECT id FROM companies WHERE id = $1 FOR UPDATE`, companyID).Scan(&lockedID); err != nil {
		return 0, err
	}

	result, err := tx.ExecContext(ctx, query, companyID)
	if err != nil {
		logger.Error("Cancel_PendingCompanyTasksByCompanyId company_id={%d}: Failed to exec sql: %s", companyID, err.Error())
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, refreshCompanyStatusQuery, companyID); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return affected, nil
}

const refreshCompanyStatusQuery = `
		WITH stats AS (
			SELECT
				COUNT(*) AS total,
				COUNT(*) FILTER (WHERE status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')) AS active,
				COUNT(*) FILTER (WHERE status = 'completed') AS completed,
				COUNT(*) FILTER (WHERE status = 'failed') AS failed,
				COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled
			FROM company_tasks
			WHERE company_id = $1
		)
		UPDATE companies
		SET
			status = CASE
				WHEN (SELECT active FROM stats) > 0 THEN 'running'
				WHEN (SELECT total FROM stats) > 0 AND (SELECT completed FROM stats) = (SELECT total FROM stats) THEN 'completed'
				WHEN (SELECT total FROM stats) > 0 AND (SELECT cancelled FROM stats) = (SELECT total FROM stats) THEN 'cancelled'
				WHEN (SELECT total FROM stats) > 0 THEN 'failed'
				ELSE status
			END,
			updated_at = timezone('UTC', now())
		WHERE id = $1
	`

func (s *CompanyStore) Refresh_CompanyStatus(ctx context.Context, companyID uint64) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var lockedID uint64

	if err := tx.QueryRowContext(ctx, `SELECT id FROM companies WHERE id = $1 FOR UPDATE`, companyID).Scan(&lockedID); err != nil {
		return err
	}
	
	if _, err := tx.ExecContext(ctx, refreshCompanyStatusQuery, companyID); err != nil {
		logger.Error("Refresh_CompanyStatus company_id={%d}: Failed to exec sql: %s", companyID, err.Error())
		return err
	}

	return tx.Commit()
}

func (s *CompanyStore) Get_CompanyByIdAndUserUuid(ctx context.Context, companyID uint64, userUuid string) (*models.CompanyWithTasks, error) {
	query := `
		SELECT
			companies.id,
			companies.created_at,
			companies.updated_at,
			companies.user_uuid,
			companies.name,
			companies.status,
			companies.ttl,
			companies.existing_task_action,
			companies.launch_mode,
			companies.incompatible_action,
			companies.configuration_source,
			COALESCE(
				companies.source_metadata, '{}'::jsonb
			) AS source_metadata,
			companies.configuration_source_data,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', company_tasks.id,
						'created_at', company_tasks.created_at,
						'updated_at', company_tasks.updated_at,
						'company_id', company_tasks.company_id,
						'device_id', company_tasks.device_id,
						'device_imei', devices.imei,
						'device_model', devices.device_model,
						'status', company_tasks.status,
						'error_message', company_tasks.error_message,
						'attempts', company_tasks.attempts,
						'last_sent_at', company_tasks.last_sent_at,
						'confirmation_deadline', company_tasks.confirmation_deadline,
						'next_attempt_at', company_tasks.next_attempt_at,
						'expires_at', company_tasks.expires_at
					)
				) FILTER (WHERE company_tasks.id IS NOT NULL),
				'[]'
			) AS tasks
		FROM companies
		LEFT JOIN company_tasks ON company_tasks.company_id = companies.id
		LEFT JOIN devices ON devices.id = company_tasks.device_id
		WHERE companies.id = $1 AND companies.user_uuid = $2
		GROUP BY companies.id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result models.CompanyWithTasks
	var tasksJSON []byte

	err := s.db.QueryRowContext(ctx, query, companyID, userUuid).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.UserUUID,
		&result.Name,
		&result.Status,
		&result.TTL,
		&result.ExistingTaskAction,
		&result.LaunchMode,
		&result.IncompatibleAction,
		&result.ConfigurationSource,
		&result.SourceMetadata,
		&result.ConfigurationSnapshot,
		&tasksJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_CompanyByIdAndUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if err := json.Unmarshal(tasksJSON, &result.Tasks); err != nil {
		return nil, err
	}

	result.TasksCount = len(result.Tasks)

	return &result, nil
}
