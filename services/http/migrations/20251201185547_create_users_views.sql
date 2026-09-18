-- +goose Up
-- +goose StatementBegin
-- +View
CREATE OR REPLACE VIEW view_users AS
SELECT
    user_contacts.id,
    user_contacts.created_at,
    user_contacts.user_uuid,
    user_cores.parent_uuid,
    user_cores.can_view_child_groups,
    user_cores.email,
    parent_user.email AS parent_email,

    -- user_contacts.*
    user_contacts.name_organization,
    user_contacts.locality,
    user_contacts.language,
    user_contacts.phone,

    -- role.*
    roles.code AS code,
    roles.name_key AS name_key,
    roles.desc_key AS desc_key,

    -- accesses.*
    user_accesses.access_treker_create,
    user_accesses.access_treker_edit,
    user_accesses.access_treker_delete,
    user_accesses.access_group_manage,
    user_accesses.access_configuration_read,
    user_accesses.access_configuration_apply,
    user_accesses.access_configuration_history,
    user_accesses.access_command_send,
    user_accesses.access_log_read
FROM user_cores
LEFT JOIN user_contacts ON user_contacts.user_uuid = user_cores.user_uuid
LEFT JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
LEFT JOIN user_accesses ON user_accesses.user_uuid = user_cores.user_uuid
LEFT JOIN roles ON roles.code = user_roles.role_code
LEFT JOIN user_cores AS parent_user ON user_cores.parent_uuid = parent_user.user_uuid;

CREATE OR REPLACE VIEW view_user_auth AS
SELECT
    user_contacts.id,
    user_contacts.created_at,
    user_contacts.user_uuid,
    user_cores.parent_uuid,
    user_cores.can_view_child_groups,
    user_cores.email,
    parent_user.email AS parent_email,
    user_cores.password,
    user_cores.refresh_token,

    -- user_contacts.*
    user_contacts.name_organization,
    user_contacts.locality,
    user_contacts.language,
    user_contacts.phone,

    -- role.*
    roles.code AS code,
    roles.name_key AS name_key,
    roles.desc_key AS desc_key,

    -- accesses.*
    user_accesses.access_treker_create,
    user_accesses.access_treker_edit,
    user_accesses.access_treker_delete,
    user_accesses.access_group_manage,
    user_accesses.access_configuration_read,
    user_accesses.access_configuration_apply,
    user_accesses.access_configuration_history,
    user_accesses.access_command_send,
    user_accesses.access_log_read
FROM user_cores
LEFT JOIN user_contacts ON user_contacts.user_uuid = user_cores.user_uuid
LEFT JOIN user_roles ON user_roles.user_uuid = user_cores.user_uuid
LEFT JOIN user_accesses ON user_accesses.user_uuid = user_cores.user_uuid
LEFT JOIN roles ON roles.code = user_roles.role_code
LEFT JOIN user_cores AS parent_user ON user_cores.parent_uuid = parent_user.user_uuid;
-- +View
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- -View
DROP VIEW IF EXISTS view_users;
DROP VIEW IF EXISTS view_user_auth;
-- -View
-- +goose StatementEnd
