package store

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"
	"unicode/utf8"
)

type AnalyticStore struct {
	db *sql.DB
}

func normalizeSQLValue(v interface{}) interface{} {
	if b, ok := v.([]byte); ok {
		return string(b)
	}

	return v
}

func (s *AnalyticStore) Create_AnalyticsRecord(ctx context.Context, analytic *models.AnalyticRecord) error {
	query := `
		INSERT INTO analytic_records (user_uuid, device_id, cfg_hash, rmq_published, device_online, error_reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, analytic.UserUUID, analytic.DeviceID, analytic.CfgHash, analytic.RMQPublished, analytic.DeviceOnline, analytic.ErrorReason)
	if err != nil {
		logger.Error("Create_AnalyticsRecord req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *AnalyticStore) Create_AnalyticsUser(ctx context.Context, analytic *models.AnalyticsUser) error {
	query := `
		INSERT INTO analytics_users (user_uuid, is_online)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, analytic.UserUUID, analytic.IsOnline)
	if err != nil {
		logger.Error("Create_AnalyticsUser req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *AnalyticStore) Get_ColumnTableAnalyticsUsers(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'analytics_users'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	analyticsUsersColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableAnalyticsUsers req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return analyticsUsersColumns, nil
}

func (s *AnalyticStore) Get_ColumnTableAnalyticRecords(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'analytic_records'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	analyticRecordColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableAnalyticRecord req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return analyticRecordColumns, nil
}

func (s *AnalyticStore) Get_ColumnTableProductAnalyticsEvents(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'product_analytics_events' AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
}

func (s *AnalyticStore) Exec_SelectQuery(ctx context.Context, query string, args ...any) (*models.AnalyticQueryResult, error) {
	// q := strings.TrimSpace(strings.ToUpper(query))
	// if !strings.HasPrefix(q, "SELECT") {
	// 	return nil, errors.New("only SELECT queries allowed")
	// }

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		const snippetMax = 1800
		snippet := query
		if utf8.RuneCountInString(snippet) > snippetMax {
			r := []rune(snippet)
			snippet = string(r[:snippetMax]) + "…"
		}

		logger.Error("Exec_SelectQuery req={%s}: Failed to exec sql: %s | query_snippet={%s}", ctx.Value("XREQID").(string), err.Error(), snippet)
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		logger.Error("Exec_SelectQuery req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	result := &models.AnalyticQueryResult{
		Columns: cols,
		Rows:    make([][]interface{}, 0),
	}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			logger.Error("Exec_SelectQuery req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		for i := range values {
			values[i] = normalizeSQLValue(values[i])
		}

		result.Rows = append(result.Rows, values)
	}

	return result, nil
}
