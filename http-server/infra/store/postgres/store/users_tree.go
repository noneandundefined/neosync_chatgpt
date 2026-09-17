package store

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"time"
)

func scanUserTreePage(rows *sql.Rows) ([]models.UserTree, int, error) {
	defer rows.Close()

	users := make([]models.UserTree, 0)
	total := 0

	for rows.Next() {
		var user models.UserTree
		var rowTotal int
		if err := rows.Scan(&user.UserUUID, &user.ParentUUID, &user.RoleCode, &user.Email, &rowTotal); err != nil {
			return nil, 0, err
		}

		total = rowTotal
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserStore) Get_UsersForTree(ctx context.Context, roleCode, viewerUUID string, dealerUUID *string, search string, limit, offset int, all bool) ([]models.UserTree, int, error) {
	var (
		inner string
		args  []any
	)

	selectCols := `
		SELECT user_uuid, parent_uuid, code, email
		FROM view_users
	`

	treeExcludeRoles := ` AND code NOT IN ('` + constants.Role_SuperAdmin + `', '` + constants.Role_Support + `')`

	switch roleCode {
	case constants.Role_SuperAdmin:
		inner = selectCols + ` WHERE email != 'info@neomatica.ru' AND user_uuid != $1` + treeExcludeRoles
		args = []any{viewerUUID}

	case constants.Role_Support:
		inner = `
			WITH RECURSIVE subtree AS (
				SELECT user_uuid FROM user_cores WHERE user_uuid = $1
				UNION ALL
				SELECT user_cores.user_uuid
				FROM user_cores
				INNER JOIN subtree ON user_cores.parent_uuid = subtree.user_uuid
			)
			SELECT view_users.user_uuid, view_users.parent_uuid, view_users.code, view_users.email
			FROM view_users
			WHERE view_users.user_uuid IN (SELECT user_uuid FROM subtree)
				AND view_users.user_uuid != $1
				AND view_users.email != 'info@neomatica.ru'
				AND view_users.code NOT IN ('` + constants.Role_SuperAdmin + `', '` + constants.Role_Support + `')
		`
		args = []any{viewerUUID}

	case constants.Role_AdminL2:
		inner = selectCols + ` WHERE parent_uuid = $1 AND user_uuid != $1` + treeExcludeRoles
		args = []any{viewerUUID}

	case constants.Role_AdminL2Support, constants.Role_User:
		if dealerUUID == nil || *dealerUUID == "" {
			return []models.UserTree{}, 0, nil
		}

		inner = selectCols + `
			WHERE user_uuid != $2
				AND (
					user_uuid = $1
					OR (
						parent_uuid = $1
						AND code IN ($3, $4)
					)
				)
		`
		args = []any{*dealerUUID, viewerUUID, constants.Role_User, constants.Role_AdminL2Support}

	default:
		return []models.UserTree{}, 0, nil
	}

	searchParam := len(args) + 1
	query := `
		SELECT user_uuid, parent_uuid, code, email, COUNT(*) OVER() AS total
		FROM (` + inner + `) AS tree_users
		WHERE ` + userEmailSearch("", searchParam) + `
		ORDER BY email ASC
	`
	args = append(args, search)
	query, args = appendLimitOffset(query, args, all, limit, offset)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Get_UsersForTree req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, err
	}

	return scanUserTreePage(rows)
}
