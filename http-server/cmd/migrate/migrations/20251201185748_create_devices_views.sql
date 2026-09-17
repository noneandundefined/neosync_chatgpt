-- +goose Up
-- +goose StatementBegin
-- +View
CREATE OR REPLACE VIEW view_devices_full AS
SELECT
    devices.id,
    devices.created_at,
    devices.updated_at,
    devices.user_uuid,
    devices.imei,
    devices.status,
    devices.device_model,
    devices.activated,
    device_confs.password,
    device_confs.request_configuration_on_connect,
    syncs.firmware_version,
    syncs.firmware_version_upd,
    syncs.cfg_version,
    syncs.last_mod_time,
    syncs.cfg_hash,
    user_cores.email AS user_email,
    groups.name AS group_name,
    group_user_devices.group_id
FROM devices
LEFT JOIN device_confs ON devices.id = device_confs.device_id
LEFT JOIN syncs ON devices.id = syncs.device_id
LEFT JOIN user_cores ON devices.user_uuid = user_cores.user_uuid
LEFT JOIN group_user_devices ON devices.id = group_user_devices.device_id
LEFT JOIN groups ON group_user_devices.group_id = groups.id;
-- +View
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- -View
DROP VIEW IF EXISTS view_devices_full;
-- -View
-- +goose StatementEnd
