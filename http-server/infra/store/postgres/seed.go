package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"neomatica/neosync/infra/constants"

	"github.com/google/uuid"
)

func InsertSuperAdminIfNotExist(ctx context.Context, db *sql.DB) error {
	var count int

	queryCheckSuperAdmin := `
		SELECT COUNT(*) FROM user_roles WHERE role_code = $1
	`

	row := db.QueryRowContext(ctx, queryCheckSuperAdmin, constants.Role_SuperAdmin)
	if err := row.Scan(&count); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error select super_admin")
	}

	if count > 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	uuid := uuid.NewString()

	queryInsertUserCore := `
		INSERT INTO user_cores (user_uuid, email, password)
		VALUES ($1, $2, $3)
	`

	_, errUserCore := tx.ExecContext(ctx, queryInsertUserCore, uuid, "info@neomatica.ru", "LM3bncrGcV88Vz3BDTiOrpnuExIToL7GzW0Ju/9/")
	if errUserCore != nil {
		return errUserCore
	}

	queryInsertUserContact := `
		INSERT INTO user_contacts (user_uuid, name_organization, locality, phone, language)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, errUserContact := tx.ExecContext(ctx, queryInsertUserContact, uuid, "ООО \"Неоматика\"", "Пермь", "78007753404", "ru")
	if errUserContact != nil {
		return errUserContact
	}

	queryInsertUserRole := `
		INSERT INTO user_roles (user_uuid, role_code)
		VALUES ($1, $2)
	`

	_, errUserRole := tx.ExecContext(ctx, queryInsertUserRole, uuid, constants.Role_SuperAdmin)
	if errUserRole != nil {
		return errUserRole
	}

	queryInsertUserAccess := `
		INSERT INTO user_accesses (user_uuid, access_treker_create, access_treker_edit, access_treker_delete, access_group_manage, access_configuration_read, access_configuration_apply, access_configuration_history, access_command_send, access_log_read)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, errUserAccess := tx.ExecContext(ctx, queryInsertUserAccess, uuid, true, true, true, true, true, true, true, true, true)
	if errUserAccess != nil {
		return errUserAccess
	}

	return tx.Commit()
}
