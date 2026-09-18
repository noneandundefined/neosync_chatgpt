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
	"neomatica/neosync-tcp/pkg/protocol/common"
	"time"

	"github.com/lib/pq"
)

type CompanyStore struct {
	db *sql.DB
}

const companyTaskSelect = `
		SELECT
			now() AS database_now,
			devices.imei AS device_imei,
			company_tasks.id,
			company_tasks.created_at,
			company_tasks.updated_at,
			company_tasks.company_id,
			company_tasks.device_id,
			company_tasks.status,
			company_tasks.error_message,
			company_tasks.attempts,
			company_tasks.cfg_data,
			company_tasks.expires_at,
			company_tasks.last_sent_at,
			company_tasks.confirmation_deadline,
			company_tasks.next_attempt_at,
			company_tasks.attempt_token,
			company_tasks.lease_until,
			companies.created_at AS company_created_at,
			companies.status AS company_status,
			companies.ttl AS company_ttl
		FROM company_tasks
		INNER JOIN companies ON companies.id = company_tasks.company_id
		INNER JOIN devices ON devices.id = company_tasks.device_id`

func (s *CompanyStore) Get_ActiveCompanyTaskByDeviceId(ctx context.Context, deviceID uint64) (*models.CompanyTaskWithCompany, error) {
	query := companyTaskSelect + `
		WHERE company_tasks.device_id = $1
			AND company_tasks.status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
			AND companies.status NOT IN ('cancelled', 'failed', 'completed')
		ORDER BY company_tasks.created_at ASC, company_tasks.id ASC
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	task, err := pgqx.QueryRowContext[models.CompanyTaskWithCompany](ctx, s.db, query, deviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ActiveCompanyTaskByDeviceId device_id={%d}: Failed to execute sql: %s", deviceID, err.Error())
		return nil, tcperrors.HandleSQLError(ctx, err)
	}

	return task, nil
}

const refreshCompanyStatusQuery = `
		WITH stats AS (
			SELECT
				COUNT(*) AS total,
				COUNT(*) FILTER (WHERE status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')) AS active,
				COUNT(*) FILTER (WHERE status = 'completed') AS completed,
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
		logger.Error("Refresh_CompanyStatus company_id={%d}: Failed to execute sql: %s", companyID, err.Error())
		return tcperrors.HandleSQLError(ctx, err)
	}

	return tx.Commit()
}

/* Update_CompanyTaskDelivery применяет переход только к прочитанной версии задачи. */
func (s *CompanyStore) Update_CompanyTaskDelivery(ctx context.Context, task *models.CompanyTaskWithCompany, delivery *models.CompanyTaskDelivery) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Serialize task transitions within a campaign before recalculating its status.
	var companyID uint64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM companies WHERE id = $1 FOR UPDATE`, task.CompanyID).Scan(&companyID); err != nil {
		return false, err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE company_tasks SET status = $2, error_message = $3,
			attempt_token = $4, confirmation_deadline = $5, next_attempt_at = $6,
			lease_until = $7, last_sent_at = CASE WHEN $8 THEN now() ELSE last_sent_at END,
			attempts = COALESCE(attempts, 0) + CASE WHEN $9 THEN 1 ELSE 0 END
		WHERE id = $1 AND status = $10 AND updated_at = $11
			AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
			AND ($2 = 'expired' OR expires_at > now())
			AND (NOT $9 OR (status IN ('pending', 'queued') AND confirmation_deadline IS NULL
				AND lease_until IS NULL AND (next_attempt_at IS NULL OR next_attempt_at <= now())))
		`, task.ID, delivery.Status, delivery.ErrorMessage, delivery.AttemptToken, delivery.ConfirmationDeadline, delivery.NextAttemptAt, delivery.LeaseUntil, delivery.Sent, delivery.Claim, task.Status, task.UpdatedAt)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, nil
	}

	cfgHash := common.GetCfgHash(task.CfgData)
	if delivery.Claim {
		result, err := tx.ExecContext(ctx, `UPDATE configurations SET cfg_hash = $2, cfg_data = $3,
			cfg_sync_status = 'pending', cfg_sync_error = NULL WHERE device_id = $1`, task.DeviceID, cfgHash, task.CfgData)
		if err != nil {
			return false, err
		}

		count, err := result.RowsAffected()
		if err != nil {
			return false, err
		}

		if count == 0 {
			return false, errors.New("device configuration not found")
		}
	}

	if delivery.Status == constants.COMPANY_TASK_STATUS_COMPLETED {
		if _, err := tx.ExecContext(ctx, `UPDATE configurations SET cfg_hash = $2, cfg_data = $3,
			cfg_sync_status = 'confirmed', cfg_sync_error = NULL WHERE device_id = $1`, task.DeviceID, cfgHash, task.CfgData); err != nil {
			return false, err
		}
	}

	if delivery.Status == constants.COMPANY_TASK_STATUS_FAILED {
		if _, err := tx.ExecContext(ctx, `UPDATE configurations SET cfg_sync_status = 'failed', cfg_sync_error = $2
			WHERE device_id = $1 AND cfg_hash = $3`, task.DeviceID, delivery.ErrorMessage, cfgHash); err != nil {
			return false, err
		}
	}

	if _, err := tx.ExecContext(ctx, refreshCompanyStatusQuery, task.CompanyID); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func (s *CompanyStore) Update_CompanyTaskObservation(ctx context.Context, taskID uint64, cfgHash uint32) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `UPDATE company_tasks
		SET last_observed_cfg_hash = $2, last_observed_at = now()
		WHERE id = $1 AND status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')`, taskID, cfgHash)

	return err
}

func (s *CompanyStore) Get_DueCompanyTasks(ctx context.Context) ([]models.CompanyTaskWithCompany, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tasks, err := pgqx.QueryContext[models.CompanyTaskWithCompany](ctx, s.db, companyTaskSelect+`
		WHERE company_tasks.status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')
		AND (company_tasks.expires_at <= now() OR company_tasks.confirmation_deadline <= now() OR company_tasks.lease_until <= now()
			OR (company_tasks.status = 'queued' AND company_tasks.next_attempt_at <= now()))
		ORDER BY company_tasks.expires_at, company_tasks.id LIMIT 500`)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *CompanyStore) Get_CompanyTaskConnectedImeis(ctx context.Context, imeis []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `SELECT DISTINCT devices.imei FROM devices
		INNER JOIN company_tasks ON company_tasks.device_id = devices.id
		WHERE devices.imei = ANY($1)
			AND company_tasks.status IN ('pending', 'queued', 'sending', 'awaiting_confirmation')`, pq.Array(imeis))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var imei string

		if err := rows.Scan(&imei); err != nil {
			return nil, err
		}

		result = append(result, imei)
	}

	return result, rows.Err()
}
