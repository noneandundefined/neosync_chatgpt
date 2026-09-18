-- +goose Up
-- +goose StatementBegin
-- +Function
-- Функция — получить устройства группы
CREATE OR REPLACE FUNCTION fn_get_devices_in_group(group_id_p BIGINT, user_uuid_p VARCHAR)
RETURNS TABLE (id BIGINT, imei VARCHAR) AS $$
BEGIN
	RETURN QUERY
	SELECT DISTINCT ON (devices.imei)
		devices.id,
		devices.imei
	FROM devices
	JOIN group_user_devices ON group_user_devices.device_id = devices.id
	JOIN groups ON group_user_devices.group_id = groups.id
	WHERE group_user_devices.group_id = group_id_p
		AND (
			groups.user_uuid = user_uuid_p
			OR EXISTS (
				SELECT 1
				FROM group_members
				WHERE group_members.group_id = groups.id
					AND group_members.member_user_uuid = user_uuid_p
			)
			OR (
				EXISTS (
					SELECT 1
					FROM user_cores
					WHERE user_cores.user_uuid = groups.user_uuid
						AND user_cores.parent_uuid = user_uuid_p
						AND (
							SELECT can_view_child_groups
							FROM user_cores
							WHERE user_uuid = user_uuid_p
						) = true
				)
			)
		)
	ORDER BY devices.imei, devices.id;
END;
$$ LANGUAGE plpgsql;

-- Функция — устройства НЕ в группе
CREATE OR REPLACE FUNCTION fn_get_devices_not_in_group(
    group_id_p BIGINT, parent_user_uuid_p VARCHAR,
    user_uuid_p VARCHAR, role_code_p VARCHAR,
    allowed_roles_p TEXT[]
)
RETURNS TABLE (id BIGINT, imei VARCHAR) AS $$
BEGIN
    RETURN QUERY
    SELECT devices.id, devices.imei
    FROM devices
    WHERE NOT EXISTS (
        SELECT 1
        FROM group_user_devices
        WHERE group_user_devices.device_id = devices.id
            AND group_user_devices.group_id = group_id_p
    )
    AND (
        -- SUPERADMIN/SUPPORT
        role_code_p = ANY(allowed_roles_p)

        -- USER
        OR (role_code_p = 'USER' AND devices.owner_uuid = user_uuid_p)

        -- DEALER/DEALER_SUPPORT
        OR (role_code_p <> 'USER' AND devices.user_uuid = parent_user_uuid_p)
    );
END;
$$ LANGUAGE plpgsql;

-- Функция — добавить устройства в группу
CREATE OR REPLACE FUNCTION fn_add_devices_to_group(group_id_p BIGINT, user_uuid_p VARCHAR, device_ids_p BIGINT[])
RETURNS VOID AS $$
BEGIN
	INSERT INTO group_user_devices (group_id, device_id)
	SELECT group_id_p, devices.id
	FROM devices
	WHERE devices.id = ANY(device_ids_p)
		AND EXISTS (
			SELECT 1
			FROM groups
			WHERE groups.id = group_id_p
				AND (
					groups.user_uuid = user_uuid_p
					OR (SELECT can_view_child_groups FROM user_cores WHERE user_uuid = user_uuid_p) = true
					OR (
						groups.can_manage_devices = true
						AND EXISTS (
							SELECT 1
							FROM group_members
							WHERE group_members.group_id = groups.id
								AND group_members.member_user_uuid = user_uuid_p
						)
					)
				)
		)
	ON CONFLICT (group_id, device_id) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Функция — удалить устройства из группы
CREATE OR REPLACE FUNCTION fn_remove_devices_from_group(group_id_p BIGINT, user_uuid_p VARCHAR, device_ids_p BIGINT[])
RETURNS VOID AS $$
BEGIN
	DELETE FROM group_user_devices
	USING groups g
	WHERE group_user_devices.group_id = groups.id
		AND group_user_devices.group_id = group_id_p
		AND group_user_devices.device_id = ANY(device_ids_p)
		AND (
			groups.user_uuid = user_uuid_p
			OR (SELECT can_view_child_groups FROM user_cores WHERE user_uuid = user_uuid_p) = true
			OR (
				groups.can_manage_devices = true
				AND EXISTS (
					SELECT 1
					FROM group_members
					WHERE group_members.group_id = groups.id
						AND group_members.member_user_uuid = user_uuid_p
				)
			)
		);
END;
$$ LANGUAGE plpgsql;
-- +Function
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS fn_get_devices_in_group;
DROP FUNCTION IF EXISTS fn_get_devices_not_in_group;
DROP FUNCTION IF EXISTS fn_add_devices_to_group;
DROP FUNCTION IF EXISTS fn_remove_devices_from_group;
-- +goose StatementEnd
