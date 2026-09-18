package store

import (
	"context"
	"database/sql"
	"errors"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"

	"time"

	"github.com/lib/pq"
)

type GroupMembersStore struct {
	db *sql.DB
}

/* AggregatedDevicePermissions — права участника на устройство из groups (шаблон группы) */
func (s *GroupMembersStore) AggregatedDevicePermissions(ctx context.Context, memberUserUUID string, deviceID uint64) (canSendCommands, canReadConfig, canEditConfig bool, err error) {
	query := `
		SELECT
			COALESCE(BOOL_OR(groups.can_send_commands), false),
			COALESCE(BOOL_OR(groups.can_read_config), false),
			COALESCE(BOOL_OR(groups.can_edit_config), false)
		FROM group_members
		INNER JOIN groups ON groups.id = group_members.group_id
		INNER JOIN group_user_devices ON group_user_devices.group_id = group_members.group_id AND group_user_devices.device_id = $2
		WHERE group_members.member_user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = s.db.QueryRowContext(ctx, query, memberUserUUID, deviceID).Scan(&canSendCommands, &canReadConfig, &canEditConfig)
	if errors.Is(err, sql.ErrNoRows) {
		return false, false, false, nil
	}

	if err != nil {
		logger.Error("AggregatedDevicePermissions req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return false, false, false, err
	}

	return canSendCommands, canReadConfig, canEditConfig, nil
}

func (s *GroupMembersStore) Get_ColumnTableGroupMember(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'group_members'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	groupMemberColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableGroupMember req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return groupMemberColumns, nil
}

func (s *GroupMembersStore) Get_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) (*models.GroupMember, error) {
	query := `
		SELECT created_at, group_id, member_user_uuid
		FROM group_members
		WHERE group_id = $1 AND member_user_uuid = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	group_member, err := pgqx.QueryRowContext[models.GroupMember](ctx, s.db, query, groupID, memberUserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_GroupMember req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return group_member, nil
}

func (s *GroupMembersStore) Get_GroupMembers(ctx context.Context, groupID uint64) ([]models.GroupMember, error) {
	query := `
		SELECT created_at, group_id, member_user_uuid
		FROM group_members
		WHERE group_id = $1
		ORDER BY created_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	group_members, err := pgqx.QueryContext[models.GroupMember](ctx, s.db, query, groupID)
	if err != nil {
		logger.Error("Get_GroupMembers req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return group_members, nil
}

/* Get_UserHasGroupAccess — владелец, участник шэринга или дилер с can_view_child_groups на владельца группы */
func (s *GroupMembersStore) Get_UserHasGroupAccess(ctx context.Context, groupID uint64, userUUID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM groups WHERE groups.id = $1 AND groups.user_uuid = $2
			UNION ALL
			SELECT 1 FROM group_members WHERE group_members.group_id = $1 AND group_members.member_user_uuid = $2
			UNION ALL
			SELECT 1
			FROM groups
			JOIN user_cores ON user_cores.user_uuid = groups.user_uuid
			JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
			WHERE groups.id = $1
				AND user_roles.role_code = 'DEALER_SUPPORT'
				AND user_cores.parent_uuid = $2
				AND (SELECT can_view_child_groups FROM user_cores WHERE user_uuid = $2) = true
		)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var ok bool
	if err := s.db.QueryRowContext(ctx, query, groupID, userUUID).Scan(&ok); err != nil {
		logger.Error("Get_UserHasGroupAccess req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return false, err
	}

	return ok, nil
}

func (s *GroupMembersStore) Update_GroupShareDefaults(ctx context.Context, groupID uint64, canEditGroup, canManageDevices, canReadConfig, canEditConfig, canSendCommands bool) error {
	query := `
		UPDATE groups
		SET can_edit_group = $1,
			can_manage_devices = $2,
			can_read_config = $3,
			can_edit_config = $4,
			can_send_commands = $5
		WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, canEditGroup, canManageDevices, canReadConfig, canEditConfig, canSendCommands, groupID)
	if err != nil {
		logger.Error("Update_GroupShareDefaults req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupMembersStore) Upsert_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) error {
	query := `
		INSERT INTO group_members (group_id, member_user_uuid) VALUES ($1, $2)
		ON CONFLICT (group_id, member_user_uuid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, groupID, memberUserUUID); err != nil {
		logger.Error("Upsert_GroupMember req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupMembersStore) Upsert_GroupMembers(ctx context.Context, groupID uint64, memberUUIDs []string) error {
	if len(memberUUIDs) == 0 {
		return nil
	}

	query := `
		INSERT INTO group_members (group_id, member_user_uuid)
		SELECT $1, unnest($2::varchar[])
		ON CONFLICT (group_id, member_user_uuid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, groupID, pq.Array(memberUUIDs)); err != nil {
		logger.Error("Upsert_GroupMembers req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *GroupMembersStore) Delete_GroupMembers(ctx context.Context, groupID uint64, memberUUIDs []string) error {
	if len(memberUUIDs) == 0 {
		return sql.ErrNoRows
	}

	query := `
		DELETE FROM group_members
		WHERE group_id = $1 AND member_user_uuid = ANY($2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, groupID, pq.Array(memberUUIDs))
	if err != nil {
		logger.Error("Delete_GroupMembers req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *GroupMembersStore) Delete_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) error {
	query := `
		DELETE FROM group_members
		WHERE group_id = $1 AND member_user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, groupID, memberUserUUID)
	if err != nil {
		logger.Error("Delete_GroupMember req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
