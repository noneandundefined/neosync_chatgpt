package store

import (
	"context"
	"database/sql"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"time"
)

const devicesStateGroupAccessSQL = `
	groups.user_uuid = $1
	OR (
		EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_members.group_id = groups.id
				AND group_members.member_user_uuid = $1
		)
	)
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
`

func deviceStateSearch(paramIndex int) string {
	placeholder := fmt.Sprintf("$%d", paramIndex)
	return `(` + placeholder + ` = '' OR concat_ws(' ', devices.imei, devices.name, devices.name_organization) ILIKE '%' || ` + placeholder + ` || '%')`
}

func appendLimitOffset(query string, args []interface{}, all bool, limit, offset int) (string, []interface{}) {
	if all {
		return query, args
	}

	n := len(args)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", n+1, n+2)
	return query, append(args, limit, offset)
}

func scanDevicesStateGroups(rows *sql.Rows) ([]models.DevicesForSendCommand, error) {
	defer rows.Close()

	list := make([]models.DevicesForSendCommand, 0)
	for rows.Next() {
		var group models.DevicesForSendCommand
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.DeviceCount); err != nil {
			return nil, err
		}

		list = append(list, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func scanDeviceStatePage(rows *sql.Rows) ([]models.DeviceShort, int, error) {
	defer rows.Close()

	devices := make([]models.DeviceShort, 0)
	total := 0

	for rows.Next() {
		var device models.DeviceShort
		var rowTotal int
		if err := rows.Scan(&device.IMEI, &device.Status, &device.Activated, &rowTotal); err != nil {
			return nil, 0, err
		}

		total = rowTotal
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (s *DeviceCommandStore) Get_DevicesStateByUserUuid(ctx context.Context, authToken *types.AuthToken, search string) ([]models.DevicesForSendCommand, error) {
	userUuid := authToken.User.UserContact.UserUUID
	parentUserUuid := permissions.ParentUserUuid(authToken)

	var ownerUuid *string
	if authToken.RoleCode == constants.Role_User {
		ownerUuid = &authToken.User.UserContact.UserUUID
	}

	searchClause := deviceStateSearch(4)

	query := `
		WITH user_group_devices AS (
			SELECT
				groups.id AS group_id,
				groups.name,
				devices.id AS device_id
			FROM groups
			JOIN group_user_devices ON group_user_devices.group_id = groups.id
			JOIN devices ON devices.id = group_user_devices.device_id
			WHERE (` + devicesStateGroupAccessSQL + `)
				AND ` + searchClause + `
		),
		all_user_devices AS (
			SELECT devices.id AS device_id
			FROM devices
			WHERE devices.user_uuid = $2
				AND ($3::varchar IS NULL OR devices.owner_uuid = $3)
				AND ` + searchClause + `
		),
		devices_in_groups AS (
			SELECT group_user_devices.device_id
			FROM group_user_devices
			JOIN groups ON group_user_devices.group_id = groups.id
			WHERE ` + devicesStateGroupAccessSQL + `
		)
		SELECT group_id, name, COUNT(*)::int AS device_count
		FROM user_group_devices
		GROUP BY group_id, name
		UNION ALL
		SELECT 0 AS group_id, NULL AS name, COUNT(*)::int AS device_count
		FROM all_user_devices
		WHERE NOT EXISTS (
			SELECT 1 FROM devices_in_groups
			WHERE devices_in_groups.device_id = all_user_devices.device_id
		)
		HAVING COUNT(*) > 0
		ORDER BY group_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userUuid, parentUserUuid, ownerUuid, search)
	if err != nil {
		logger.Error("Get_DevicesStateByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return scanDevicesStateGroups(rows)
}

func (s *DeviceCommandStore) Get_DevicesState(ctx context.Context, userUuid string, search string) ([]models.DevicesForSendCommand, error) {
	searchClause := deviceStateSearch(2)

	query := `
		WITH user_group_devices AS (
			SELECT
				groups.id AS group_id,
				groups.name,
				devices.id AS device_id
			FROM groups
			JOIN group_user_devices ON group_user_devices.group_id = groups.id
			JOIN devices ON devices.id = group_user_devices.device_id
			WHERE groups.user_uuid = $1
				AND ` + searchClause + `
		),
		all_user_devices AS (
			SELECT devices.id AS device_id
			FROM devices
			WHERE ` + searchClause + `
		),
		devices_in_groups AS (
			SELECT device_id FROM group_user_devices
			JOIN groups ON group_user_devices.group_id = groups.id
			WHERE groups.user_uuid = $1
		)
		SELECT group_id, name, COUNT(*)::int AS device_count
		FROM user_group_devices
		GROUP BY group_id, name
		UNION ALL
		SELECT 0 AS group_id, NULL AS name, COUNT(*)::int AS device_count
		FROM all_user_devices
		WHERE NOT EXISTS (
			SELECT 1 FROM devices_in_groups
			WHERE devices_in_groups.device_id = all_user_devices.device_id
		)
		HAVING COUNT(*) > 0
		ORDER BY group_id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userUuid, search)
	if err != nil {
		logger.Error("Get_DevicesState req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return scanDevicesStateGroups(rows)
}

func (s *DeviceCommandStore) Get_DevicesStateDevicesByUserUuid(ctx context.Context, authToken *types.AuthToken, groupID uint64, search string, limit, offset int, all bool) ([]models.DeviceShort, int, error) {
	userUuid := authToken.User.UserContact.UserUUID
	parentUserUuid := permissions.ParentUserUuid(authToken)

	var ownerUuid *string
	if authToken.RoleCode == constants.Role_User {
		ownerUuid = &authToken.User.UserContact.UserUUID
	}

	var (
		query string
		args  []interface{}
	)

	if groupID == 0 {
		searchClause := deviceStateSearch(4)
		query = `
			SELECT
				devices.imei,
				devices.status,
				devices.activated,
				COUNT(*) OVER() AS total
			FROM devices
			WHERE devices.user_uuid = $2
				AND ($3::varchar IS NULL OR devices.owner_uuid = $3)
				AND ` + searchClause + `
				AND NOT EXISTS (
					SELECT 1
					FROM group_user_devices
					JOIN groups ON group_user_devices.group_id = groups.id
					WHERE group_user_devices.device_id = devices.id
						AND (` + devicesStateGroupAccessSQL + `)
				)
			ORDER BY devices.imei
		`
		args = []interface{}{userUuid, parentUserUuid, ownerUuid, search}
	} else {
		searchClause := deviceStateSearch(3)
		query = `
			SELECT
				devices.imei,
				devices.status,
				devices.activated,
				COUNT(*) OVER() AS total
			FROM groups
			JOIN group_user_devices ON group_user_devices.group_id = groups.id
			JOIN devices ON devices.id = group_user_devices.device_id
			WHERE groups.id = $2
				AND (` + devicesStateGroupAccessSQL + `)
				AND ` + searchClause + `
			ORDER BY devices.imei
		`
		args = []interface{}{userUuid, groupID, search}
	}

	query, args = appendLimitOffset(query, args, all, limit, offset)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Get_DevicesStateDevicesByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, err
	}

	return scanDeviceStatePage(rows)
}

func (s *DeviceCommandStore) Get_DevicesStateDevices(ctx context.Context, userUuid string, groupID uint64, search string, limit, offset int, all bool) ([]models.DeviceShort, int, error) {
	var (
		query string
		args  []interface{}
	)

	if groupID == 0 {
		searchClause := deviceStateSearch(2)
		query = `
			SELECT
				devices.imei,
				devices.status,
				devices.activated,
				COUNT(*) OVER() AS total
			FROM devices
			WHERE ` + searchClause + `
				AND NOT EXISTS (
					SELECT 1
					FROM group_user_devices
					JOIN groups ON group_user_devices.group_id = groups.id
					WHERE group_user_devices.device_id = devices.id
						AND groups.user_uuid = $1
				)
			ORDER BY devices.imei
		`
		args = []interface{}{userUuid, search}
	} else {
		searchClause := deviceStateSearch(3)
		query = `
			SELECT
				devices.imei,
				devices.status,
				devices.activated,
				COUNT(*) OVER() AS total
			FROM groups
			JOIN group_user_devices ON group_user_devices.group_id = groups.id
			JOIN devices ON devices.id = group_user_devices.device_id
			WHERE groups.user_uuid = $1
				AND groups.id = $2
				AND ` + searchClause + `
			ORDER BY devices.imei
		`
		args = []interface{}{userUuid, groupID, search}
	}

	query, args = appendLimitOffset(query, args, all, limit, offset)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Get_DevicesStateDevices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, err
	}

	return scanDeviceStatePage(rows)
}
