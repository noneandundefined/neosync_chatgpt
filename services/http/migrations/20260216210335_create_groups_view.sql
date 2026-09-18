-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW view_groups_with_objects AS
SELECT
	groups.id,
	groups.created_at,
	groups.updated_at,
	groups.user_uuid,
	groups.name,
	groups.description,
	groups.can_edit_group,
	groups.can_manage_devices,
	groups.can_read_config,
	groups.can_edit_config,
	groups.can_send_commands,
	COUNT(group_user_devices.device_id) AS objects
FROM groups
LEFT JOIN group_user_devices ON group_user_devices.group_id = groups.id
GROUP BY
	groups.id,
	groups.created_at,
	groups.updated_at,
	groups.user_uuid,
	groups.name,
	groups.description,
	groups.can_edit_group,
	groups.can_manage_devices,
	groups.can_read_config,
	groups.can_edit_config,
	groups.can_send_commands;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS view_groups_with_objects;
-- +goose StatementEnd
