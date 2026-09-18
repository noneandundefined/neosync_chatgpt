package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"
)

type ConfigurationTemplateStore struct {
	db *sql.DB
}

func (s *ConfigurationTemplateStore) Create_ConfigurationTemplate(ctx context.Context, template *models.ConfigurationTemplate, cfgData []byte) (uint64, error) {
	query := `
		INSERT INTO configuration_templates (user_uuid, name, model, type_of_saving, cfg_data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var id uint64
	if err := s.db.QueryRowContext(
		ctx,
		query,
		template.UserUUID,
		template.Name,
		template.Model,
		template.TypeOfSaving,
		cfgData,
	).Scan(&id); err != nil {
		logger.Error("Create_ConfigurationTemplate req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	return id, nil
}

func (s *ConfigurationTemplateStore) Get_ConfigurationTemplatesByUserUuid(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.ConfigurationTemplate, int, int, error) {
	allowedColumnSort := map[string]string{
		"name":       "configuration_templates.name",
		"model":      "configuration_templates.model",
		"created_at": "configuration_templates.created_ut",
		"updated_at": "configuration_templates.updated_at",
	}

	orderBy := "ORDER BY configuration_templates.created_ut DESC"
	if col, ok := allowedColumnSort[columnSortKey]; ok {
		dir := "ASC"
		if columnSortDir == "desc" {
			dir = "DESC"
		}

		orderBy = fmt.Sprintf("ORDER BY %s %s", col, dir)
	}

	searchFilter := ""
	args := []any{userUuid}
	argIndex := 2

	if search != "" {
		searchFilter = fmt.Sprintf(" AND (configuration_templates.name ILIKE $%d OR configuration_templates.model ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM configuration_templates
		WHERE user_uuid = $1%s
	`, searchFilter)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		logger.Error("Get_ConfigurationTemplatesByUserUuid req={%s}: Failed count: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	var totalAll int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM configuration_templates WHERE user_uuid = $1`, userUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_ConfigurationTemplatesByUserUuid req={%s}: Failed total_all: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, created_ut, updated_at, user_uuid, name, model, type_of_saving
		FROM configuration_templates
		WHERE user_uuid = $1%s
		%s
		LIMIT $%d OFFSET $%d
	`, searchFilter, orderBy, argIndex, argIndex+1)

	args = append(args, limit, offset)

	templates, err := pgqx.QueryContext[models.ConfigurationTemplate](ctx, s.db, query, args...)
	if err != nil {
		logger.Error("Get_ConfigurationTemplatesByUserUuid req={%s}: Failed query: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return templates, total, totalAll, nil
}

func (s *ConfigurationTemplateStore) Get_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) (*models.ConfigurationTemplate, error) {
	query := `
		SELECT id, created_ut, updated_at, user_uuid, name, model, type_of_saving
		FROM configuration_templates
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	template, err := pgqx.QueryRowContext[models.ConfigurationTemplate](ctx, s.db, query, id, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ConfigurationTemplateByIdAndUserUuid req={%s}: Failed query: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return template, nil
}

func (s *ConfigurationTemplateStore) Get_ConfigurationTemplateCfgDataByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) ([]byte, error) {
	query := `
		SELECT cfg_data
		FROM configuration_templates
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var cfgData []byte
	if err := s.db.QueryRowContext(ctx, query, id, userUuid).Scan(&cfgData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ConfigurationTemplateCfgDataByIdAndUserUuid req={%s}: Failed query: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return cfgData, nil
}

func (s *ConfigurationTemplateStore) Delete_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) error {
	query := `
		DELETE FROM configuration_templates
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id, userUuid)
	if err != nil {
		logger.Error("Delete_ConfigurationTemplateByIdAndUserUuid req={%s}: Failed exec: %s", ctx.Value("XREQID").(string), err.Error())
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

func (s *ConfigurationTemplateStore) Update_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid, name, typeOfSaving string, cfgData []byte) error {
	query := `
		UPDATE configuration_templates
		SET name = $3, type_of_saving = $4, cfg_data = $5
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id, userUuid, name, typeOfSaving, cfgData)
	if err != nil {
		logger.Error("Update_ConfigurationTemplateByIdAndUserUuid req={%s}: Failed exec: %s", ctx.Value("XREQID").(string), err.Error())
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
