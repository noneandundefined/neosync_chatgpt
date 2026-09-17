package pgqx

import (
	"context"
	"database/sql"
)

/* QueryContext перебор элементов в результате запроса к БД */
func QueryContext[T any](ctx context.Context, db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []T

	for rows.Next() {
		var item T

		fieldMap, err := collectFields(&item)
		if err != nil {
			return nil, err
		}

		scanArgs := make([]any, len(columns))
		for i, col := range columns {
			ptr, ok := fieldMap[col]
			if !ok {
				var dummy any
				ptr = &dummy
			}

			scanArgs[i] = ptr
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

/* QueryRowContext перебор элементов в результате запроса к БД */
func QueryRowContext[T any](ctx context.Context, db *sql.DB, query string, args ...any) (*T, error) {
	var item T

	row := db.QueryRowContext(ctx, query, args...)

	fieldMap, err := collectFields(&item)
	if err != nil {
		return nil, err
	}

	columns, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer columns.Close()

	columnNames, err := columns.Columns()
	if err != nil {
		return nil, err
	}

	scanArgs := make([]any, len(columnNames))
	for i, col := range columnNames {
		ptr, ok := fieldMap[col]
		if !ok {
			var dummy any
			ptr = &dummy
		}

		scanArgs[i] = ptr
	}

	if err := row.Scan(scanArgs...); err != nil {
		return nil, err
	}

	return &item, nil
}
