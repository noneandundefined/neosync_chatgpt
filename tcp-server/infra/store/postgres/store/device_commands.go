package store

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/pkg/pgqx"
	"time"
)

type DeviceCommandStore struct {
	db *sql.DB
}

func (s *DeviceCommandStore) Get_DeviceCommandsByImeiAndSendMode(ctx context.Context, imei, sendMode string) ([]models.DeviceCommandWithExecutions, error) {
	query := `
		SELECT
			device_commands.id,
			device_command_executions.id AS execution_id,
			device_commands.send_mode,
			device_commands.command,
			device_command_executions.device_imei,
			device_command_executions.status
		FROM device_commands
		JOIN device_command_executions ON device_command_executions.command_id = device_commands.id
		WHERE device_command_executions.device_imei = $1
			AND device_commands.send_mode = $2
			AND device_command_executions.status IN ('inprogress','pending')
		ORDER BY device_commands.created_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	commands, err := pgqx.QueryContext[models.DeviceCommandWithExecutions](ctx, s.db, query, imei, sendMode)
	if err != nil {
		logger.Error("Get_DeviceCommandsByImeiAndSendMode imei={%s} send_mode={%s}: Failed to exec sql: %s", imei, sendMode, err.Error())
		return nil, err
	}

	return commands, nil
}

func (s *DeviceCommandStore) Update_DeviceCommandStatusProgressByCommand(ctx context.Context, imei, command string) error {
	query := `
		WITH picked AS (
			SELECT
				device_command_executions.id,
				device_command_executions.command_id
			FROM device_command_executions
			JOIN device_commands
				ON device_commands.id = device_command_executions.command_id
			WHERE device_command_executions.device_imei = $1
				AND device_commands.command = $2
				AND device_command_executions.status = 'pending'
			ORDER BY
				device_command_executions.status DESC,
				device_commands.created_at ASC
			LIMIT 1
		),
		updated AS (
			UPDATE device_command_executions
			SET status = 'inprogress'
			FROM picked
			WHERE device_command_executions.id = picked.id
			RETURNING picked.command_id
		)
		UPDATE device_commands
		SET status = 'inprogress'
		FROM updated
		WHERE device_commands.id = updated.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei, command)
	if err != nil {
		logger.Error("Update_DeviceCommandStatusProgressByCommand imei={%s}: Failed to exec sql: %s", imei, err.Error())
		return err
	}

	return nil
}

func (s *DeviceCommandStore) Claim_DeviceCommandExecution(ctx context.Context, imei, command string) (uint64, uint64, error) {
	query := `
		WITH picked AS (
			SELECT
				device_command_executions.id AS execution_id,
				device_command_executions.command_id
			FROM device_command_executions
			JOIN device_commands
				ON device_commands.id = device_command_executions.command_id
			WHERE device_command_executions.device_imei = $1
				AND device_commands.command = $2
				AND device_command_executions.status = 'pending'
			ORDER BY device_commands.created_at ASC
			LIMIT 1
		),
		updated AS (
			UPDATE device_command_executions
			SET status = 'inprogress'
			FROM picked
			WHERE device_command_executions.id = picked.execution_id
			RETURNING picked.execution_id, picked.command_id
		)
		UPDATE device_commands
		SET status = 'inprogress'
		FROM updated
		WHERE device_commands.id = updated.command_id
		RETURNING updated.execution_id, updated.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var executionID, commandID uint64
	err := s.db.QueryRowContext(ctx, query, imei, command).Scan(&executionID, &commandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, err
		}

		logger.Error("Claim_DeviceCommandExecution imei={%s}: Failed to exec sql: %s", imei, err.Error())
		return 0, 0, err
	}

	return executionID, commandID, nil
}

func (s *DeviceCommandStore) Update_DeviceCommandExecutionByID(ctx context.Context, executionID uint64, status, response string) error {
	query := `
		WITH updated AS (
			UPDATE device_command_executions
			SET status = $1, response = $2
			WHERE id = $3 AND status = 'inprogress'
			RETURNING command_id
		)
		UPDATE device_commands
		SET status = $1
		FROM updated
		WHERE device_commands.id = updated.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, status, nullIfEmpty(response), executionID)
	if err != nil {
		logger.Error("Update_DeviceCommandExecutionByID execution_id={%d}: Failed to exec sql: %s", executionID, err.Error())
		return err
	}

	return nil
}

func (s *DeviceCommandStore) Update_DeviceCommandExecutionTimeoutByID(ctx context.Context, executionID uint64) error {
	var commandID uint64
	var sendMode string

	querySelect := `
		SELECT device_commands.id, device_commands.send_mode
		FROM device_command_executions
		JOIN device_commands ON device_commands.id = device_command_executions.command_id
		WHERE device_command_executions.id = $1
			AND device_command_executions.status = 'inprogress'
	`

	queryUpdateExecution := `
		UPDATE device_command_executions
		SET status = $1, response = $2
		WHERE id = $3 AND status = 'inprogress'
	`

	queryUpdateCommand := `
		UPDATE device_commands SET status = $1 WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, querySelect, executionID).Scan(&commandID, &sendMode); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			logger.Error("Update_DeviceCommandExecutionTimeoutByID execution_id={%d}: Failed to exec sql: %s", executionID, err.Error())
			return err
		}

		status := constants.CMD_STATUS_EXECUTIONERROR
		response := "message.device-not-connected"

		if sendMode == constants.CMD_SEND_MODE_ON_CONNECT {
			status = constants.CMD_STATUS_PENDING
			response = ""
		}

		if _, err := tx.ExecContext(ctx, queryUpdateExecution, status, nullIfEmpty(response), executionID); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, queryUpdateCommand, status, commandID); err != nil {
			return err
		}

		return nil
	})
}

func (s *DeviceCommandStore) Update_DeviceCommandExecutionInProgressByID(ctx context.Context, executionID uint64) error {
	query := `
		WITH updated AS (
			UPDATE device_command_executions
			SET status = 'inprogress'
			WHERE id = $1 AND status = 'pending'
			RETURNING command_id
		)
		UPDATE device_commands
		SET status = 'inprogress'
		FROM updated
		WHERE device_commands.id = updated.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, executionID)
	if err != nil {
		logger.Error("Update_DeviceCommandExecutionInProgressByID execution_id={%d}: Failed to exec sql: %s", executionID, err.Error())
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *DeviceCommandStore) Update_DeviceCommandMarkExecution(ctx context.Context, imei, status, command, response string) error {
	var commandExecutionID, commandID uint64

	querySelect := `
		SELECT
			device_command_executions.id,
			device_commands.id AS command_id
        FROM device_command_executions
        JOIN device_commands ON device_commands.id = device_command_executions.command_id
        WHERE device_command_executions.device_imei = $1
        	AND device_commands.command = $2
        	AND device_command_executions.status = 'inprogress'
        ORDER BY device_commands.created_at ASC
        LIMIT 1
	`

	queryUpdateExecution := `
		UPDATE device_command_executions
		SET status = $1, response = $2
		WHERE id = $3
	`

	queryUpdateCommand := `
		UPDATE device_commands SET status = $1 WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, querySelect, imei, command).Scan(&commandExecutionID, &commandID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			logger.Error("Update_DeviceCommandMarkExecution imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		if _, err := tx.ExecContext(ctx, queryUpdateExecution, status, response, commandExecutionID); err != nil {
			logger.Error("Update_DeviceCommandMarkExecution imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		if _, err := tx.ExecContext(ctx, queryUpdateCommand, status, commandID); err != nil {
			logger.Error("Update_DeviceCommandMarkExecution imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		return nil
	})
}

func (s *DeviceCommandStore) Update_DeviceCommandsMarkExecutionBatch(ctx context.Context, imei, status, response string) error {
	query := `
		WITH upd AS (
			UPDATE device_command_executions
			SET status = $2, response = $3
			FROM device_commands
			WHERE device_commands.id = device_command_executions.command_id
				AND device_command_executions.device_imei = $1
				AND device_command_executions.status IN ('pending', 'inprogress')
				AND device_commands.send_mode != $4
			RETURNING device_command_executions.command_id
		)
		UPDATE device_commands
		SET status = $2
		WHERE id IN (SELECT DISTINCT command_id FROM upd)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei, status, response, constants.CMD_SEND_MODE_ON_CONNECT)
	if err != nil {
		logger.Error("Update_DeviceCommandsMarkExecutionBatch imei={%s}: Failed to exec sql: %s", imei, err.Error())
		return err
	}

	return nil
}

func (s *DeviceCommandStore) Update_DeviceCommandsResetOnConnectOnDisconnect(ctx context.Context, imei string) error {
	query := `
		WITH reset AS (
			UPDATE device_command_executions
			SET status = 'pending', response = NULL
			FROM device_commands
			WHERE device_commands.id = device_command_executions.command_id
				AND device_command_executions.device_imei = $1
				AND device_commands.send_mode = $2
				AND device_command_executions.status = 'inprogress'
			RETURNING device_command_executions.command_id
		)
		UPDATE device_commands
		SET status = 'pending'
		WHERE id IN (SELECT DISTINCT command_id FROM reset)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei, constants.CMD_SEND_MODE_ON_CONNECT)
	if err != nil {
		logger.Error("Update_DeviceCommandsResetOnConnectOnDisconnect imei={%s}: Failed to exec sql: %s", imei, err.Error())
		return err
	}

	return nil
}

func (s *DeviceCommandStore) Update_DeviceCommandHandleTimeout(ctx context.Context, imei, command string) error {
	var commandExecutionID, commandID uint64
	var sendMode string

	querySelect := `
		SELECT
			device_command_executions.id,
			device_commands.id AS command_id,
			device_commands.send_mode
		FROM device_command_executions
		JOIN device_commands ON device_commands.id = device_command_executions.command_id
		WHERE device_command_executions.device_imei = $1
			AND device_commands.command = $2
			AND device_command_executions.status = 'inprogress'
		ORDER BY device_commands.created_at ASC
		LIMIT 1
	`

	queryUpdateExecution := `
		UPDATE device_command_executions
		SET status = $1, response = $2
		WHERE id = $3
	`

	queryUpdateCommand := `
		UPDATE device_commands SET status = $1 WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, querySelect, imei, command).Scan(&commandExecutionID, &commandID, &sendMode); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			logger.Error("Update_DeviceCommandHandleTimeout imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		status := constants.CMD_STATUS_EXECUTIONERROR
		response := "message.device-not-connected"

		if sendMode == constants.CMD_SEND_MODE_ON_CONNECT {
			status = constants.CMD_STATUS_PENDING
			response = ""
		}

		if _, err := tx.ExecContext(ctx, queryUpdateExecution, status, nullIfEmpty(response), commandExecutionID); err != nil {
			logger.Error("Update_DeviceCommandHandleTimeout imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		if _, err := tx.ExecContext(ctx, queryUpdateCommand, status, commandID); err != nil {
			logger.Error("Update_DeviceCommandHandleTimeout imei={%s}: Failed to exec sql: %s", imei, err.Error())
			return err
		}

		return nil
	})
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}

	return value
}

func (s *DeviceCommandStore) Update_DeviceCommandClearProgressStatus(ctx context.Context, imei string) error {
	query := `
		WITH timedout AS (
			UPDATE device_command_executions
			SET status = 'pending'
			WHERE device_imei = $1
			AND status = 'inprogress'
			AND updated_at < now() - interval '30 seconds'
			RETURNING command_id
		)
		UPDATE device_commands
		SET status = 'pending'
		FROM timedout
		WHERE device_commands.id = timedout.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei)
	if err != nil {
		logger.Error("Update_DeviceCommandClearProgressStatus imei={%s}: Failed to exec sql: %s", imei, err.Error())
		return err
	}

	return nil
}
