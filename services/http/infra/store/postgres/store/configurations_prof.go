package store

import (
	"context"
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/pgqx"
	"strings"
	"time"
)

func (s *ConfigurationStore) Get_ColumnTableConfigurationProf(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'configurations_prof'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	configurationColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableConfigurationProf req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return configurationColumns, nil
}

func (s *ConfigurationStore) Replace_ConfigurationProf(ctx context.Context, deviceID uint64, cfgHash uint32, prof *adm.ConfigurationProfValues) error {
	if prof == nil {
		return nil
	}

	unknownJSON, err := json.Marshal(prof.UnknownFields)
	if err != nil {
		return err
	}

	columns := []string{"device_id", "cfg_hash", "user_uuid", "unknown_fields"}
	placeholders := []string{"$1", "$2", "devices.user_uuid", "$3::jsonb"}
	args := []any{deviceID, cfgHash, string(unknownJSON)}
	argIdx := 4

	for _, schema := range constants.ConfigProfSchemas() {
		col := constants.ConfigProfColumnName(schema.Name)
		columns = append(columns, col)

		val, ok := prof.Values[col]
		if !ok {
			placeholders = append(placeholders, "NULL")
			continue
		}

		sqlType := constants.ConfigProfSQLType(schema.Type)
		switch sqlType {
		case "JSONB":
			raw, err := json.Marshal(val)
			if err != nil {
				return err
			}
			placeholders = append(placeholders, fmt.Sprintf("$%d::jsonb", argIdx))
			args = append(args, string(raw))
		default:
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
			args = append(args, val)
		}

		argIdx++
	}

	setParts := make([]string, 0, len(columns)-1)
	for _, col := range columns[1:] {
		setParts = append(setParts, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	query := fmt.Sprintf(`
		INSERT INTO configurations_prof (%s)
		SELECT
			%s
		FROM devices
		WHERE devices.id = $1
		ON CONFLICT (device_id) DO UPDATE SET %s
	`, strings.Join(columns, ", "), strings.Join(placeholders, ", "), strings.Join(setParts, ", "))

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		logger.Error("Replace_ConfigurationProf req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
