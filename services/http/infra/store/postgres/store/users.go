package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/pgqx"
	"strings"
	"time"

	"github.com/lib/pq"
)

type UserStore struct {
	db *sql.DB
}

func userEmailSearch(alias string, paramIndex int) string {
	placeholder := fmt.Sprintf("$%d", paramIndex)
	column := "email"
	if alias != "" {
		column = alias + ".email"
	}

	return `(` + placeholder + ` = '' OR ` + column + ` ILIKE '%' || ` + placeholder + ` || '%')`
}

func (s *UserStore) Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error {
	query := `
		INSERT INTO user_cores (user_uuid, parent_uuid, email, password, refresh_token)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.ParentUUID, user.Email, user.Password, user.RefreshToken); err != nil {
		if strings.Contains(err.Error(), `user_cores_email_key`) {
			return httperr.Err_DuplicateEmail
		}

		logger.Error("Create_UserCore req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Create_UserContact(ctx context.Context, tx *sql.Tx, user *models.UserContact) error {
	query := `
		INSERT INTO user_contacts (user_uuid, name_organization, locality, phone, language)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.NameOrganization, user.Locality, user.Phone, user.Language); err != nil {
		logger.Error("Create_UserContact req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Create_UserRole(ctx context.Context, tx *sql.Tx, user *models.UserRole) error {
	query := `
		INSERT INTO user_roles (user_uuid, role_code)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.RoleCode); err != nil {
		logger.Error("Create_UserRole req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Create_UserAccess(ctx context.Context, tx *sql.Tx, user *models.UserAccess) error {
	query := `
		INSERT INTO user_accesses (
			user_uuid,
			access_treker_create,
			access_treker_edit,
			access_treker_delete,
			access_group_manage,
			access_configuration_read,
			access_configuration_apply,
			access_configuration_history,
			access_command_send,
			access_log_read
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.AccessTrekerCreate, user.AccessTrekerEdit, user.AccessTrekerDelete, user.AccessGroupManage, user.AccessConfigurationRead, user.AccessConfigurationApply, user.AccessConfigurationHistory, user.AccessCommandSend, user.AccessLogRead); err != nil {
		logger.Error("Create_UserAccess req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Get_ColumnTableUserCore(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'user_cores'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userCoreColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableUserCore req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return userCoreColumns, nil
}

func (s *UserStore) Get_ColumnTableUserRole(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'user_roles'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userRoleColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableUserRole req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return userRoleColumns, nil
}

func (s *UserStore) Get_ColumnTableUserContacts(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'user_contacts'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userContactColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableUserContacts req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return userContactColumns, nil
}

func (s *UserStore) Get_Users(ctx context.Context) ([]models.UserToOwner, error) {
	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_roles.role_code
		FROM user_cores
		LEFT JOIN user_roles ON user_cores.user_uuid = user_roles.user_uuid
		ORDER BY user_cores.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.UserToOwner](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Users req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Get_UserAuthByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error) {
	query := `
		SELECT * FROM view_user_auth
		WHERE user_uuid = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserAuth](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserAuthByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UsersWithParams(ctx context.Context, limit, offset int, search, role, columnSortKey, columnSortDir, userUuid string) ([]models.User, int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	allowedSortColumns := map[string]string{
		"email":             "email",
		"name_organization": "name_organization",
		"address":           "address",
		"phone":             "phone",
	}

	orderBy := "created_at DESC"

	if col, ok := allowedSortColumns[columnSortKey]; ok {
		orderBy = fmt.Sprintf("%s %s", col, columnSortDir)
	}

	query := fmt.Sprintf(`
		SELECT * FROM view_users
		WHERE email != 'info@neomatica.ru'
			AND user_uuid != $4
			AND ($3 = '' OR concat_ws(' ', email, name_organization, locality, phone, code) ILIKE '%%' || $3 || '%%')
			AND ($5 = '' OR code = $5)
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	users, err := pgqx.QueryContext[models.User](ctx, s.db, query, limit, offset, search, userUuid, role)
	if err != nil {
		logger.Error("Get_UsersWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryCount := `
		SELECT COUNT(*) FROM view_users
		WHERE email != 'info@neomatica.ru'
			AND ($1 = '' OR email ILIKE '%' || $1 || '%' OR phone ILIKE '%' || $1 || '%')
			AND ($2 = '' OR code = $2)
	`

	var total int
	if err := s.db.QueryRowContext(ctx, queryCount, search, role).Scan(&total); err != nil {
		logger.Error("Get_UsersWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryTotalAll := `
		SELECT COUNT(*) FROM view_users
		WHERE email != 'info@neomatica.ru' AND user_uuid != $1
	`

	var totalAll int
	if err := tx.QueryRowContext(ctx, queryTotalAll, userUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_UsersWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Get_UsersWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return users, total, totalAll, nil
}

func (s *UserStore) Get_UsersToOwner(ctx context.Context, role, userUuid string, parentUuid sql.NullString, search string, limit, offset int, all bool) ([]models.UserToOwner, int, error) {
	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_roles.role_code,
			COUNT(*) OVER() AS total
		FROM user_cores
		LEFT JOIN user_roles ON user_cores.user_uuid = user_roles.user_uuid
		WHERE
			(
				user_cores.user_uuid = $1
				OR (
					user_roles.role_code = 'USER'
					AND $2::varchar IS NOT NULL
					AND user_cores.parent_uuid = $2::varchar
				)
				OR (
					$2::varchar IS NULL
					AND user_roles.role_code = 'DEALER'
				)
			)
			AND ` + userEmailSearch("user_cores", 3) + `
		ORDER BY user_cores.email ASC
	`

	args := []interface{}{userUuid, parentUuid, search}
	query, args = appendLimitOffset(query, args, all, limit, offset)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Get_UsersToOwner req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]models.UserToOwner, 0)
	total := 0

	for rows.Next() {
		var user models.UserToOwner
		var rowTotal int
		if err := rows.Scan(&user.UserUUID, &user.Email, &user.RoleCode, &rowTotal); err != nil {
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

func (s *UserStore) Get_UsersByParentUuid(ctx context.Context, userUuid string) ([]models.User, error) {
	query := `
		SELECT * FROM view_users
		WHERE parent_uuid = $1 AND NOT (user_uuid = $1)
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.User](ctx, s.db, query, userUuid)
	if err != nil {
		logger.Error("Get_UsersByParentUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Get_UsersByParentUuidWithParams(ctx context.Context, limit, offset int, search, role, columnSortKey, columnSortDir, dealerUuid, currentUserUuid string) ([]models.User, int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	allowedSortColumns := map[string]string{
		"email":             "email",
		"name_organization": "name_organization",
		"address":           "address",
		"phone":             "phone",
	}

	orderBy := "created_at DESC"

	if col, ok := allowedSortColumns[columnSortKey]; ok {
		orderBy = fmt.Sprintf("%s %s", col, columnSortDir)
	}

	query := fmt.Sprintf(`
		SELECT * FROM view_users
		WHERE parent_uuid = $1
			AND user_uuid <> $2
			AND (
				$5 = ''
				OR concat_ws(' ', email, name_organization, locality, phone, code) ILIKE '%%' || $5 || '%%'
			)
			AND ($6 = '' OR code = $6)
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, orderBy)

	users, err := pgqx.QueryContext[models.User](ctx, s.db, query, dealerUuid, currentUserUuid, limit, offset, search, role)
	if err != nil {
		logger.Error("Get_UsersByParentUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryCount := `
		SELECT COUNT(*) FROM view_users
		WHERE parent_uuid = $1
			AND user_uuid <> $2
			AND (
				$3 = ''
				OR email ILIKE '%' || $3 || '%'
				OR phone ILIKE '%' || $3 || '%'
			)
			AND ($4 = '' OR code = $4)
	`

	var total int
	if err := s.db.QueryRowContext(ctx, queryCount, dealerUuid, currentUserUuid, search, role).Scan(&total); err != nil {
		logger.Error("Get_UsersByParentUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryTotalAll := `
		SELECT COUNT(*) FROM view_users
		WHERE parent_uuid = $1 AND user_uuid <> $2
	`

	var totalAll int
	if err := tx.QueryRowContext(ctx, queryTotalAll, dealerUuid, currentUserUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_UsersByParentUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Get_UsersByParentUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return users, total, totalAll, nil
}

func (s *UserStore) Get_UsersByParentUuidToOwner(ctx context.Context, userUuid string) ([]models.UserToOwner, error) {
	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_roles.role_code
		FROM user_cores
		LEFT JOIN user_roles ON user_cores.user_uuid = user_roles.user_uuid
		WHERE user_cores.parent_uuid = $1 AND NOT (user_cores.user_uuid = $1) AND user_roles.role_code = $2
		ORDER BY user_cores.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.UserToOwner](ctx, s.db, query, userUuid, constants.Role_AdminL2)
	if err != nil {
		logger.Error("Get_UsersByParentUuidToOwner req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Get_UserByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error) {
	query := `
		SELECT view_user_auth.*, user_cores.force_neosync_configuration_priority
		FROM view_user_auth
		JOIN user_cores ON user_cores.user_uuid = view_user_auth.user_uuid
		WHERE view_user_auth.user_uuid = $1 LIMIT 1
	`

	user, err := pgqx.QueryRowContext[models.UserAuth](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserLoginWithRoleCodeByUuid(ctx context.Context, userUuid string) (*models.UserLoginWithRoleCode, error) {
	query := `
		SELECT
			user_cores.email as login,
			user_roles.role_code
		FROM user_cores
		LEFT JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
		WHERE user_cores.user_uuid = $1
		LIMIT 1
	`

	user, err := pgqx.QueryRowContext[models.UserLoginWithRoleCode](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserMinimalByUuid req={%s} user={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), userUuid, err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserWithPasswordByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error) {
	query := `
		SELECT * FROM view_user_auth
		WHERE user_uuid = $1 LIMIT 1
	`

	user, err := pgqx.QueryRowContext[models.UserAuth](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserWithPasswordByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserCoreByUuid(ctx context.Context, userUuid string) (*models.UserCore, error) {
	query := `
		SELECT * FROM user_cores WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserCore](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserCoreByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserCoreByRefreshToken(ctx context.Context, token string) (*models.UserCore, error) {
	query := `
		SELECT * FROM user_cores WHERE refresh_token = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserCore](ctx, s.db, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserCoreByRefreshToken req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserAccessByUuid(ctx context.Context, userUuid string) (*models.UserAccess, error) {
	query := `
		SELECT * FROM user_accesses WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserAccess](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserAccessByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error) {
	query := `
		SELECT * FROM user_cores WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserCore](ctx, s.db, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserCoreByEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserRoleByUuid(ctx context.Context, userUuid string) (*models.UserRole, error) {
	query := `
		SELECT * FROM user_roles WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserRole](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserRoleByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserRolesByRoleCode(ctx context.Context, roleCode string) ([]models.UserRole, error) {
	query := `
		SELECT * FROM user_roles WHERE role_code = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.UserRole](ctx, s.db, query, roleCode)
	if err != nil {
		logger.Error("Get_UserRolesByRoleCode req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Update_UserCoreOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fields := []string{}
	values := []interface{}{}
	i := 1

	if user.Email != nil {
		fields = append(fields, fmt.Sprintf("email=$%d", i))
		values = append(values, *user.Email)
		i++
	}

	if user.Password != nil {
		fields = append(fields, fmt.Sprintf("password=$%d", i))
		values = append(values, *user.Password)
		i++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE user_cores SET %s WHERE user_uuid = $%d
	`, strings.Join(fields, ", "), i)

	values = append(values, userUuid)

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if strings.Contains(err.Error(), `user_cores_email_key`) {
			return httperr.Err_DuplicateEmail
		}

		logger.Error("Update_UserCoreOnesByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserContactOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fields := []string{}
	values := []interface{}{}
	i := 1

	if user.Phone != nil {
		fields = append(fields, fmt.Sprintf("phone=$%d", i))
		values = append(values, *user.Phone)
		i++
	}

	if user.NameOrganization != nil {
		fields = append(fields, fmt.Sprintf("name_organization=$%d", i))
		values = append(values, *user.NameOrganization)
		i++
	}

	if user.Locality != nil {
		fields = append(fields, fmt.Sprintf("locality=$%d", i))
		values = append(values, *user.Locality)
		i++
	}

	if user.Language != nil {
		fields = append(fields, fmt.Sprintf("language=$%d", i))
		values = append(values, *user.Language)
		i++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE user_contacts SET %s WHERE user_uuid = $%d
	`, strings.Join(fields, ", "), i)

	values = append(values, userUuid)

	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		logger.Error("Update_UserContactOnesByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserAccessOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fields := []string{}
	values := []interface{}{}
	i := 1

	if user.Accesses.AccessTrekerCreate != nil {
		fields = append(fields, fmt.Sprintf("access_treker_create=$%d", i))
		values = append(values, *user.Accesses.AccessTrekerCreate)
		i++
	}

	if user.Accesses.AccessTrekerEdit != nil {
		fields = append(fields, fmt.Sprintf("access_treker_edit=$%d", i))
		values = append(values, *user.Accesses.AccessTrekerEdit)
		i++
	}

	if user.Accesses.AccessTrekerDelete != nil {
		fields = append(fields, fmt.Sprintf("access_treker_delete=$%d", i))
		values = append(values, *user.Accesses.AccessTrekerDelete)
		i++
	}

	if user.Accesses.AccessGroupManage != nil {
		fields = append(fields, fmt.Sprintf("access_group_manage=$%d", i))
		values = append(values, *user.Accesses.AccessGroupManage)
		i++
	}

	if user.Accesses.AccessConfigurationRead != nil {
		fields = append(fields, fmt.Sprintf("access_configuration_read=$%d", i))
		values = append(values, *user.Accesses.AccessConfigurationRead)
		i++
	}

	if user.Accesses.AccessConfigurationApply != nil {
		fields = append(fields, fmt.Sprintf("access_configuration_apply=$%d", i))
		values = append(values, *user.Accesses.AccessConfigurationApply)
		i++
	}

	if user.Accesses.AccessConfigurationHistory != nil {
		fields = append(fields, fmt.Sprintf("access_configuration_history=$%d", i))
		values = append(values, *user.Accesses.AccessConfigurationHistory)
		i++
	}

	if user.Accesses.AccessCommandSend != nil {
		fields = append(fields, fmt.Sprintf("access_command_send=$%d", i))
		values = append(values, *user.Accesses.AccessCommandSend)
		i++
	}

	if user.Accesses.AccessLogRead != nil {
		fields = append(fields, fmt.Sprintf("access_log_read=$%d", i))
		values = append(values, *user.Accesses.AccessLogRead)
		i++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		`UPDATE user_accesses SET %s WHERE user_uuid = $%d`,
		strings.Join(fields, ", "),
		i,
	)

	values = append(values, userUuid)

	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		logger.Error("Update_UserAccessOnesByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserCoreRefreshToken(ctx context.Context, token, userUuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_cores SET refresh_token = $1 WHERE user_uuid = $2
	`

	if _, err := s.db.ExecContext(ctx, query, token, userUuid); err != nil {
		logger.Error("Update_UserCoreRefreshToken req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserPasswordByUuid(ctx context.Context, password, userUuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_cores SET password = $1 WHERE user_uuid = $2
	`

	if _, err := s.db.ExecContext(ctx, query, password, userUuid); err != nil {
		logger.Error("Update_UserPasswordByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_ForceNeosyncConfigurationPriority(ctx context.Context, userUuid string, enabled bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_cores SET force_neosync_configuration_priority = $1 WHERE user_uuid = $2
	`

	if _, err := s.db.ExecContext(ctx, query, enabled, userUuid); err != nil {
		logger.Error("Update_ForceNeosyncConfigurationPriority req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_CanViewChildGroupsByUuid(ctx context.Context, userUuid string, canViewChildGroups bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_cores SET can_view_child_groups = $1 WHERE user_uuid = $2
	`

	if _, err := s.db.ExecContext(ctx, query, canViewChildGroups, userUuid); err != nil {
		logger.Error("Update_CanViewChildGroupsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Delete_UsersByUuid(ctx context.Context, userUuids []string) error {
	if len(userUuids) == 0 {
		return httperr.Err_NotDeleted
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `SELECT delete_users_by_uuids($1)`, pq.Array(userUuids)); err != nil {
			logger.Error("Delete_UsersByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		return nil
	})
}

func (s *UserStore) Delete_UserByUuid(ctx context.Context, userUuid string) error {
	query := `
		SELECT delete_user_by_uuid($1)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	del, err := s.db.ExecContext(ctx, query, userUuid)
	if err != nil {
		logger.Error("Delete_UserByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	delRows, _ := del.RowsAffected()
	if delRows == 0 {
		return httperr.Err_NotDeleted
	}

	return nil
}
