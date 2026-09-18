-- +goose Up
-- +goose StatementBegin
-- +Function
CREATE OR REPLACE FUNCTION delete_users_by_uuids(p_user_uuids VARCHAR[])
RETURNS VOID AS $$
BEGIN
    DELETE FROM user_cores WHERE user_uuid = ANY(p_user_uuids);
    DELETE FROM user_contacts WHERE user_contacts.user_uuid = ANY(p_user_uuids);
    DELETE FROM user_roles WHERE user_roles.user_uuid = ANY(p_user_uuids);
    DELETE FROM user_accesses WHERE user_accesses.user_uuid = ANY(p_user_uuids);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_user_by_uuid(p_user_uuid VARCHAR)
RETURNS VOID AS $$
BEGIN
    DELETE FROM user_cores WHERE user_uuid = p_user_uuid;
    DELETE FROM user_contacts WHERE user_uuid = p_user_uuid;
    DELETE FROM user_roles WHERE user_uuid = p_user_uuid;
    DELETE FROM user_accesses WHERE user_uuid = p_user_uuid;
END;
$$ LANGUAGE plpgsql;
-- +Function
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS delete_user_by_uuid(VARCHAR);
DROP FUNCTION IF EXISTS delete_users_by_uuids(VARCHAR[]);
-- +goose StatementEnd
