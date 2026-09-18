-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_members (
    created_at TIMESTAMPTZ NOT NULL DEFAULT (timezone('UTC', now())),
    group_id BIGINT NOT NULL,
    member_user_uuid VARCHAR(255) NOT NULL,

    PRIMARY KEY (group_id, member_user_uuid),
    CONSTRAINT fk_group_members_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_group_members_user FOREIGN KEY (member_user_uuid) REFERENCES user_cores(user_uuid) ON DELETE CASCADE
);

-- +Indexes
CREATE INDEX IF NOT EXISTS idx_group_members_member_uuid ON group_members(member_user_uuid);
-- +Indexes

-- +Comments
COMMENT ON TABLE group_members IS 'Users who are members of shared groups';
COMMENT ON COLUMN group_members.created_at IS 'Timestamp when the member was added to the group';
COMMENT ON COLUMN group_members.group_id IS 'Foreign key to the group';
COMMENT ON COLUMN group_members.member_user_uuid IS 'Uuid of the member user';
-- +Comments
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS group_members;
-- +goose StatementEnd
