package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"

	"github.com/lib/pq"
)

type DeviceCommandStore struct {
	db *sql.DB
}

func (s *DeviceCommandStore) Create_DeviceCommand(ctx context.Context, tx *sql.Tx, command *models.DeviceCommand) (uint64, error) {
	var id uint64

	query := `
		INSERT INTO device_commands (execute_until, session_id, user_uuid, send_mode, command)
		VALUES (NOW() + interval '5 minutes', $1, $2, $3, $4) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	/* Get count commands */
	var count int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM device_commands WHERE user_uuid = $1", command.UserUUID).Scan(&count); err != nil {
		logger.Error("Create_DeviceCommand req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	if count < 15 {
		if err := tx.QueryRowContext(ctx, query, command.SessionID, command.UserUUID, command.SendMode, command.Command).Scan(&id); err != nil {
			logger.Error("Create_DeviceCommand req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return 0, err
		}
	} else {
		/* Если уже 15 команд, перезаписываем самую старую команду */
		if err := tx.QueryRowContext(ctx, `
			UPDATE device_commands
			SET
				created_at = NOW(),
				execute_until = NOW() + interval '5 minutes',
				session_id = $1,
				send_mode = $2,
				status = 'pending',
				command = $3
			WHERE id = (
				SELECT id FROM device_commands
				WHERE user_uuid = $4
				ORDER BY created_at ASC
				LIMIT 1
			)
			RETURNING id
		`, command.SessionID, command.SendMode, command.Command, command.UserUUID).Scan(&id); err != nil {
			logger.Error("Create_DeviceCommand req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return 0, err
		}
	}

	return id, nil
}

func (s *DeviceCommandStore) Create_Batch_DeviceCommandExecution(ctx context.Context, tx *sql.Tx, commandId uint64, imeis []string) error {
	queryInsert := `
		INSERT INTO device_command_executions (command_id, device_imei) VALUES ($1, $2)
	`

	queryDelete := `
		DELETE FROM device_command_executions WHERE command_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, queryDelete, commandId); err != nil {
		logger.Error("Create_Batch_DeviceCommandExecution req={%s} commandId={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), commandId, err.Error())
		return err
	}

	if len(imeis) == 0 {
		return nil
	}

	stmt, err := tx.PrepareContext(ctx, queryInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, imei := range imeis {
		_, err := stmt.ExecContext(ctx, commandId, imei)
		if err != nil {
			logger.Error("Create_Batch_DeviceCommandExecution req={%s} imei={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), imei, err.Error())
			return err
		}
	}

	return nil
}

func (s *DeviceCommandStore) Get_DeviceCommandWithExecutionsById(ctx context.Context, userUuid string, commandId uint64) (*models.DeviceCommandWithExecutions, error) {
	var result models.DeviceCommandWithExecutions
	var executionsJSON []byte

	query := `
		SELECT
			device_commands.id,
			device_commands.created_at,
			device_commands.updated_at,
			device_commands.execute_until,
			device_commands.session_id,
			device_commands.user_uuid,
			device_commands.send_mode,
			device_commands.status,
			device_commands.command,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'imei', device_command_executions.device_imei,
						'status', device_command_executions.status,
						'response', device_command_executions.response
					)
				) FILTER (WHERE device_command_executions.id IS NOT NULL),
				'[]'
			) AS executions
		FROM device_commands
		LEFT JOIN device_command_executions ON device_command_executions.command_id = device_commands.id
		WHERE device_commands.id = $1 AND device_commands.user_uuid = $2
		GROUP BY
			device_commands.id,
			device_commands.created_at,
			device_commands.updated_at,
			device_commands.user_uuid,
			device_commands.send_mode,
			device_commands.status,
			device_commands.command
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, commandId, userUuid).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.ExecuteUntil,
		&result.SessionID,
		&result.UserUUID,
		&result.SendMode,
		&result.Status,
		&result.Command,
		&executionsJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		logger.Error("Get_DeviceCommandWithExecutionsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if err := json.Unmarshal(executionsJSON, &result.Executions); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *DeviceCommandStore) Get_DeviceCommandWithExecutionsByUuid(ctx context.Context, userUuid string) ([]models.DeviceCommand, error) {
	query := `
		SELECT * FROM device_commands
		WHERE user_uuid = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	commands, err := pgqx.QueryContext[models.DeviceCommand](ctx, s.db, query, userUuid)
	if err != nil {
		logger.Error("Get_DeviceCommandWithExecutionsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return commands, nil
}

func (s *DeviceCommandStore) Get_DeviceCommandWithExecutionsByUuidAndSessId(ctx context.Context, userUuid, sessId string) ([]models.DeviceCommandWithExecutions, error) {
	var result []models.DeviceCommandWithExecutions

	query := `
		SELECT
			device_commands.id,
			device_commands.created_at,
			device_commands.updated_at,
			device_commands.execute_until,
			device_commands.session_id,
			device_commands.user_uuid,
			device_commands.send_mode,
			device_commands.status,
			device_commands.command,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'imei', device_command_executions.device_imei,
						'status', device_command_executions.status,
						'response', device_command_executions.response
					)
				) FILTER (WHERE device_command_executions.id IS NOT NULL),
				'[]'
			) AS executions
		FROM device_commands
		LEFT JOIN device_command_executions ON device_command_executions.command_id = device_commands.id
		WHERE device_commands.session_id = $1 AND device_commands.user_uuid = $2
		GROUP BY
			device_commands.id,
			device_commands.created_at,
			device_commands.updated_at,
			device_commands.user_uuid,
			device_commands.send_mode,
			device_commands.status,
			device_commands.command
		ORDER BY device_commands.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, sessId, userUuid)
	if err != nil {
		logger.Error("Get_DeviceCommandWithExecutionsByUuidAndSessId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cmd models.DeviceCommandWithExecutions
		var executionsJSON []byte

		err := rows.Scan(
			&cmd.ID,
			&cmd.CreatedAt,
			&cmd.UpdatedAt,
			&cmd.ExecuteUntil,
			&cmd.SessionID,
			&cmd.UserUUID,
			&cmd.SendMode,
			&cmd.Status,
			&cmd.Command,
			&executionsJSON,
		)
		if err != nil {
			logger.Error("Get_DeviceCommandWithExecutionsByUuidAndSessId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		if err := json.Unmarshal(executionsJSON, &cmd.Executions); err != nil {
			logger.Error("Get_DeviceCommandWithExecutionsByUuidAndSessId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		result = append(result, cmd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DeviceCommandStore) Update_DeviceCommandStatusProgressByImeiAndCommandId(ctx context.Context, imeis []string, commandId uint64) error {
	query := `
		WITH updated AS (
			UPDATE device_command_executions
			SET status = 'inprogress'
			WHERE device_imei = ANY($1)
			AND command_id = $2
			AND NOT EXISTS (
				SELECT 1
				FROM device_command_executions AS device_command_executions_2
				JOIN device_commands AS device_commands_2
					ON device_commands_2.id = device_command_executions_2.command_id
				WHERE device_command_executions_2.device_imei = ANY($1)
					AND device_command_executions_2.status = 'inprogress'
			)
			RETURNING command_id
		)
		UPDATE device_commands
		SET status = 'inprogress'
		FROM updated
		WHERE device_commands.id = updated.command_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, pq.Array(imeis), commandId); err != nil {
		logger.Error("Update_DeviceCommandStatusProgressByImeiAndCommandId req={%s} command_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), commandId, err.Error())
		return err
	}

	return nil
}

func (s *DeviceCommandStore) Update_DeviceCommandCancelByCommandId(ctx context.Context, userUuid string, commandId uint64) error {
	var imeis []string

	querySelect := `
		SELECT device_imei
		FROM device_command_executions
		WHERE command_id = $1 AND status IN ('pending','inprogress')
	`

	queryUpdExecutions := `
		UPDATE device_command_executions
			SET status = 'notcompleted'
		FROM device_commands
		WHERE device_commands.id = device_command_executions.command_id
			AND device_commands.user_uuid = $1
			AND device_command_executions.device_imei = ANY($2)
			AND device_commands.id = $3
	`

	queryUpdCommand := `
		UPDATE device_commands
			SET status = 'notcompleted'
		WHERE user_uuid = $1 AND id = $2
			AND status IN ('pending','inprogress')
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		/* Get IMEIs from command */
		rows, err := tx.QueryContext(ctx, querySelect, commandId)
		if err != nil {
			logger.Error("Update_DeviceCommandCancelByCommandId req={%s} command_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), commandId, err.Error())
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var imei string
			if err = rows.Scan(&imei); err != nil {
				return err
			}

			imeis = append(imeis, imei)
		}

		if _, err := tx.ExecContext(ctx, queryUpdExecutions, userUuid, pq.Array(imeis), commandId); err != nil {
			logger.Error("Update_DeviceCommandCancelByCommandId req={%s} command_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), commandId, err.Error())
			return err
		}

		if _, err := tx.ExecContext(ctx, queryUpdCommand, userUuid, commandId); err != nil {
			logger.Error("Update_DeviceCommandCancelByCommandId req={%s} command_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), commandId, err.Error())
			return err
		}

		return nil
	})
}

func (s *DeviceCommandStore) Delete_DeviceCommandByUuidAndId(ctx context.Context, userUuid string, command_id uint64) error {
	query := `
		DELETE FROM device_commands WHERE user_uuid = $1 and id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid, command_id); err != nil {
		logger.Error("Delete_DeviceCommandByUuidAndId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
