-- +goose Up
-- +goose StatementBegin
-- +Function
CREATE OR REPLACE FUNCTION update_devices_uuid_by_imeis_to_null(p_imeis VARCHAR[])
RETURNS VOID AS $$
BEGIN
    UPDATE devices SET owner_uuid = NULL WHERE imei = ANY(p_imeis);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_devices_by_imeis(p_imeis VARCHAR[])
RETURNS VOID AS $$
BEGIN
    DELETE FROM devices WHERE imei = ANY(p_imeis);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_devices_by_imeis_and_user_uuid(p_imeis VARCHAR[], p_user_uuid VARCHAR)
RETURNS VOID AS $$
BEGIN
    DELETE FROM devices
    WHERE imei = ANY(p_imeis)
      AND user_uuid = p_user_uuid;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_devices_for_user(p_user VARCHAR)
RETURNS TABLE (
    id               BIGINT,
    created_at       timestamptz,
    updated_at       timestamptz,
    user_uuid        VARCHAR,
    imei             VARCHAR,
    status           BOOLEAN,
    device_model     VARCHAR,
    activated        BOOLEAN,
    password         VARCHAR,
    request_configuration_on_connect BOOLEAN,
    firmware_version SMALLINT,
    firmware_version_upd SMALLINT,
    cfg_version      SMALLINT,
    last_mod_time    BIGINT,
    cfg_hash         BIGINT,
    user_email       VARCHAR,
    group_name       VARCHAR,
    group_id         BIGINT
)
AS $$
BEGIN
    RETURN QUERY
        SELECT
            view_devices_full.id,
            view_devices_full.created_at,
            view_devices_full.updated_at,
            view_devices_full.user_uuid,
            view_devices_full.imei,
            view_devices_full.status,
            view_devices_full.device_model,
            view_devices_full.activated,
            view_devices_full.password,
            view_devices_full.request_configuration_on_connect,
            view_devices_full.firmware_version::SMALLINT,
            view_devices_full.firmware_version_upd::SMALLINT,
            view_devices_full.cfg_version::SMALLINT,
            view_devices_full.last_mod_time,
            view_devices_full.cfg_hash,
            view_devices_full.user_email,
            view_devices_full.group_name,
            view_devices_full.group_id
        FROM view_devices_full
        WHERE
            (
                EXISTS (
                    SELECT 1 FROM user_roles
                    WHERE user_roles.user_uuid = p_user AND user_roles.role_code IN ('SUPERADMIN', 'SUPPORT')
                )
                    OR view_devices_full.user_uuid = p_user
                    OR view_devices_full.group_id IN (
                        SELECT group_user_devices.group_id FROM group_user_devices
                        WHERE group_user_devices.user_uuid = p_user
                    )
            );
END;
$$ LANGUAGE plpgsql STABLE;
-- +Function
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- -Function
DROP FUNCTION IF EXISTS update_devices_uuid_by_imeis_to_null(VARCHAR[]);
DROP FUNCTION IF EXISTS delete_devices_by_imeis(VARCHAR[]);
DROP FUNCTION IF EXISTS delete_devices_by_imeis_and_user_uuid(VARCHAR[], VARCHAR);
DROP FUNCTION IF EXISTS fn_devices_for_user(VARCHAR);
-- -Function
-- +goose StatementEnd
