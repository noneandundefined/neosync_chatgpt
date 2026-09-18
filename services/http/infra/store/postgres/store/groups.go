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
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/pkg/pgqx"
	"neomatica/neosync/types"
	"time"

	"github.com/lib/pq"
)

type GroupStore struct {
	db *sql.DB
}

func (s *GroupStore) Create_Group(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (user_uuid, name, description, can_edit_group, can_manage_devices, can_read_config, can_edit_config, can_send_commands) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, group.UserUuid, group.Name, group.Description, group.CanEditGroup, group.CanManageDevices, group.CanReadConfig, group.CanEditConfig, group.CanSendCommands).Scan(&group.ID); err != nil {
		logger.Error("Create_Group req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupStore) Get_ColumnTableGroup(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'groups'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	groupColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableGroup req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return groupColumns, nil
}

func (s *GroupStore) Get_Groups(ctx context.Context) ([]models.GroupName, error) {
	query := `
		SELECT id, name FROM groups
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	groups, err := pgqx.QueryContext[models.GroupName](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Groups req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return groups, nil
}

func (s *GroupStore) Get_GroupById(ctx context.Context, groupId uint64) (*models.Group, error) {
	query := `
		SELECT * FROM view_groups_with_objects
		WHERE view_groups_with_objects.id = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	group, err := pgqx.QueryRowContext[models.Group](ctx, s.db, query, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_GroupById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return group, nil
}

func (s *GroupStore) Get_GroupsByUuid(ctx context.Context, userUuid string) ([]models.Group, error) {
	query := `
		SELECT
			view_groups_with_objects.id,
			view_groups_with_objects.created_at,
			view_groups_with_objects.updated_at,
			view_groups_with_objects.user_uuid,
			view_groups_with_objects.name,
			view_groups_with_objects.description,
			view_groups_with_objects.can_edit_group,
			view_groups_with_objects.can_manage_devices,
			view_groups_with_objects.can_read_config,
			view_groups_with_objects.can_edit_config,
			view_groups_with_objects.can_send_commands,
			view_groups_with_objects.objects
		FROM view_groups_with_objects
		WHERE (
			view_groups_with_objects.user_uuid = $1
			OR EXISTS (
				SELECT 1 FROM group_members
				WHERE group_members.group_id = view_groups_with_objects.id AND group_members.member_user_uuid = $1
			)
		)
		ORDER BY view_groups_with_objects.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	groups, err := pgqx.QueryContext[models.Group](ctx, s.db, query, userUuid)
	if err != nil {
		logger.Error("Get_GroupsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return groups, nil
}

type groupRowWithTotal struct {
	models.Group
	TotalCount int `db:"total_count"`
}

func (s *GroupStore) Get_GroupsByUuidWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Group, int, int, error) {
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
		"name":        "view_groups_with_objects.name",
		"description": "view_groups_with_objects.description",
		"objects":     "view_groups_with_objects.objects",
		"created_at":  "view_groups_with_objects.created_at",
	}

	orderBy := "view_groups_with_objects.created_at DESC"

	if col, ok := allowedSortColumns[columnSortKey]; ok {
		orderBy = fmt.Sprintf("%s %s", col, columnSortDir)
	}

	query := fmt.Sprintf(`
		SELECT
			view_groups_with_objects.id,
			view_groups_with_objects.created_at,
			view_groups_with_objects.updated_at,
			view_groups_with_objects.user_uuid,
			view_groups_with_objects.name,
			view_groups_with_objects.description,
			view_groups_with_objects.can_edit_group,
			view_groups_with_objects.can_manage_devices,
			view_groups_with_objects.can_read_config,
			view_groups_with_objects.can_edit_config,
			view_groups_with_objects.can_send_commands,
			view_groups_with_objects.objects,
			COUNT(*) OVER() AS total_count
		FROM view_groups_with_objects
		WHERE
			(
				-- My groups
				view_groups_with_objects.user_uuid = $1

				-- Shared with me
				OR EXISTS (
					SELECT 1 FROM group_members
					WHERE group_members.group_id = view_groups_with_objects.id AND group_members.member_user_uuid = $1
				)

				-- Groups child users
				OR (
					EXISTS (
						SELECT 1
						FROM user_cores
						JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
						WHERE user_cores.user_uuid = view_groups_with_objects.user_uuid
							AND user_roles.role_code = 'DEALER_SUPPORT'
							AND user_cores.parent_uuid = $1
							AND (
								SELECT can_view_child_groups
								FROM user_cores
								WHERE user_uuid = $1
							) = true
					)
				)
			)
			AND (
	 			$2 = ''
	 			OR concat_ws(' ', view_groups_with_objects.name, view_groups_with_objects.description) ILIKE '%%' || $2 || '%%'
	 		)
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, orderBy)

	rows, err := pgqx.QueryContext[groupRowWithTotal](ctx, s.db, query, userUuid, search, limit, offset)
	if err != nil {
		logger.Error("Get_GroupsByUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if len(rows) == 0 {
		return []models.Group{}, 0, 0, nil
	}

	total := rows[0].TotalCount

	groups := make([]models.Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, r.Group)
	}

	queryTotalAll := `
		SELECT COUNT(*) FROM view_groups_with_objects
		WHERE
			(
				-- My groups
				view_groups_with_objects.user_uuid = $1

				-- Shared with me
				OR EXISTS (
					SELECT 1 FROM group_members
					WHERE group_members.group_id = view_groups_with_objects.id AND group_members.member_user_uuid = $1
				)

				-- Groups child users
				OR (
					EXISTS (
						SELECT 1
						FROM user_cores
						JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
						WHERE user_cores.user_uuid = view_groups_with_objects.user_uuid
							AND user_roles.role_code = 'DEALER_SUPPORT'
							AND user_cores.parent_uuid = $1
							AND (
								SELECT can_view_child_groups
								FROM user_cores
								WHERE user_uuid = $1
							) = true
					)
				)
			)
	`

	var totalAll int
	if err := tx.QueryRowContext(ctx, queryTotalAll, userUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_GroupsByUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Get_GroupsByUuidWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return groups, total, totalAll, nil
}

func (s *GroupStore) Get_GroupWithDevices(ctx context.Context, groupId uint64, authToken *types.AuthToken) (*models.GroupWithDevices, error) {
	userUuid := authToken.User.UserContact.UserUUID
	parentUserUuid := permissions.ParentUserUuid(authToken)

	queryGroup := `
		SELECT * FROM view_groups_with_objects
		WHERE view_groups_with_objects.id = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	group, err := pgqx.QueryRowContext[models.Group](ctx, s.db, queryGroup, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_GroupWithDevices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	result := &models.GroupWithDevices{
		Group: group,
	}

	queryAssigned := `
		SELECT id, imei
		FROM fn_get_devices_in_group($1, $2)
	`

	assigned, err := pgqx.QueryContext[models.GDevice](ctx, s.db, queryAssigned, groupId, userUuid)
	if err != nil {
		logger.Error("Get_GroupWithDevices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	result.Devices.Assigned = assigned

	queryAvailable := `
		SELECT id, imei
		FROM fn_get_devices_not_in_group($1, $2, $3, $4, $5)
	`

	roles := []string{constants.Role_SuperAdmin, constants.Role_Support}
	available, err := pgqx.QueryContext[models.GDevice](ctx, s.db, queryAvailable, groupId, parentUserUuid, userUuid, authToken.RoleCode, pq.Array(roles))
	if err != nil {
		logger.Error("Get_GroupByObjects req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	result.Devices.Available = available

	return result, nil
}

func (s *GroupStore) Get_ListDeviceIDsInGroup(ctx context.Context, groupID uint64) ([]uint64, error) {
	query := `
		SELECT device_id FROM group_user_devices WHERE group_id = $1 ORDER BY device_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	type row struct {
		DeviceID uint64 `db:"device_id"`
	}

	rows, err := pgqx.QueryContext[row](ctx, s.db, query, groupID)
	if err != nil {
		logger.Error("Get_ListDeviceIDsInGroup req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	out := make([]uint64, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].DeviceID)
	}

	return out, nil
}

func (s *GroupStore) Get_Assigned_And_Available_Devices(ctx context.Context, groupId uint64, authToken *types.AuthToken) ([]models.GDevice, []models.GDevice, error) {
	userUuid := authToken.User.UserContact.UserUUID
	parentUserUuid := permissions.ParentUserUuid(authToken)

	queryAssigned := `
		SELECT id, imei
		FROM fn_get_devices_in_group($1, $2)
	`

	assigned, err := pgqx.QueryContext[models.GDevice](ctx, s.db, queryAssigned, groupId, userUuid)
	if err != nil {
		logger.Error("Get_Assigned_And_Available_Devices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, nil, err
	}

	queryAvailable := `
		SELECT id, imei
		FROM fn_get_devices_not_in_group($1, $2, $3, $4, $5)
	`

	roles := []string{constants.Role_SuperAdmin, constants.Role_Support}
	available, err := pgqx.QueryContext[models.GDevice](ctx, s.db, queryAvailable, groupId, parentUserUuid, userUuid, authToken.RoleCode, pq.Array(roles))
	if err != nil {
		logger.Error("Get_Assigned_And_Available_Devices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, nil, err
	}

	return assigned, available, nil
}

func (s *GroupStore) Get_GroupNamesByUuid(ctx context.Context, userUuid string) ([]models.GroupName, error) {
	query := `
		SELECT
			id, name
		FROM groups
		WHERE (
			-- My groups
			user_uuid = $1

			-- Groups child users
			OR (
				EXISTS (
					SELECT 1
					FROM user_cores
					JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
					WHERE user_cores.user_uuid = groups.user_uuid
						AND user_roles.role_code = 'DEALER_SUPPORT'
						AND user_cores.parent_uuid = $1
						AND (
							SELECT can_view_child_groups
							FROM user_cores
							WHERE user_uuid = $1
						) = true
				)
			)
		)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	groups, err := pgqx.QueryContext[models.GroupName](ctx, s.db, query, userUuid)
	if err != nil {
		logger.Error("Get_GroupNamesByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return groups, nil
}

func (s *GroupStore) Update_AddDevicesToGroup(ctx context.Context, device *models.GroupDevices) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT fn_add_devices_to_group($1, $2, $3)
	`

	if _, err := s.db.ExecContext(ctx, query, device.GroupId, device.UserUuid, pq.Array(device.DevicesId)); err != nil {
		logger.Error("Update_AddDevicesToGroup req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupStore) Update_AddDevicesToGroupTx(ctx context.Context, tx *sql.Tx, device *models.GroupDevices) error {
	query := `
		SELECT fn_add_devices_to_group($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, device.GroupId, device.UserUuid, pq.Array(device.DevicesId)); err != nil {
		logger.Error("Update_AddDevicesToGroupTx req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupStore) Update_RemoveDevicesFromGroup(ctx context.Context, device *models.GroupDevices) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT fn_remove_devices_from_group($1, $2, $3)
	`

	if _, err := s.db.ExecContext(ctx, query, device.GroupId, device.UserUuid, pq.Array(device.DevicesId)); err != nil {
		logger.Error("Update_RemoveDevicesFromGroup req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupStore) Update_GroupById(ctx context.Context, groupId uint64, group *models.Group) error {
	query := `
		UPDATE groups SET name = $1, description = $2 WHERE id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, group.Name, group.Description, groupId)
	if err != nil {
		logger.Error("Update_GroupById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_GroupById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *GroupStore) Delete_GroupsById(ctx context.Context, groupIds []uint64) error {
	if len(groupIds) == 0 {
		return httperr.Err_NotDeleted
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `CREATE TEMP TABLE tmp_ids(group_id INTEGER PRIMARY KEY) ON COMMIT DROP`); err != nil {
			logger.Error("Delete_GroupsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		stmt, err := tx.PrepareContext(ctx, `INSERT INTO tmp_ids (group_id) VALUES ($1)`)
		if err != nil {
			logger.Error("Delete_GroupsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}
		defer stmt.Close()

		for _, id := range groupIds {
			if _, err := stmt.ExecContext(ctx, id); err != nil {
				logger.Error("Delete_GroupsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
				return err
			}
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM groups USING tmp_ids WHERE groups.id = tmp_ids.group_id`); err != nil {
			logger.Error("Delete_GroupsById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		return nil
	})
}

func (s *GroupStore) Delete_GroupById(ctx context.Context, groupId uint64) error {
	query := `
		DELETE FROM groups WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	del, err := s.db.ExecContext(ctx, query, groupId)
	if err != nil {
		logger.Error("Delete_GroupById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	delAffected, err := del.RowsAffected()
	if err != nil {
		logger.Error("Delete_GroupById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if delAffected == 0 {
		return httperr.Err_NotDeleted
	}

	return nil
}
