-- +goose Up
-- +goose StatementBegin
ALTER TABLE configuration_history
    ADD COLUMN IF NOT EXISTS cfg_sync_status VARCHAR(50) NULL,
    ADD COLUMN IF NOT EXISTS cfg_sync_error TEXT NULL,
    ADD COLUMN IF NOT EXISTS origin VARCHAR(32) NOT NULL DEFAULT 'unknown';

CREATE INDEX IF NOT EXISTS idx_configuration_history_device_hash
    ON configuration_history(device_id, cfg_hash, apply_at DESC);

-- Preserve the currently active configuration as the baseline version when it
-- is not already present in the history.
INSERT INTO configuration_history (
    apply_at,
    device_id,
    cfg_hash,
    cfg_data,
    cfg_sync_status,
    cfg_sync_error,
    origin
)
SELECT
    COALESCE(c.updated_at, timezone('UTC', now())),
    c.device_id,
    c.cfg_hash,
    c.cfg_data,
    c.cfg_sync_status,
    c.cfg_sync_error,
    'baseline'
FROM configurations c
WHERE c.cfg_data IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM configuration_history h
      WHERE h.device_id = c.device_id
        AND h.cfg_hash = c.cfg_hash
        AND h.cfg_data = c.cfg_data
  );

-- If the active binary was already present in the old rolling history, attach
-- the current synchronization outcome to that latest matching snapshot.
UPDATE configuration_history h
SET cfg_sync_status = c.cfg_sync_status,
    cfg_sync_error = c.cfg_sync_error,
    origin = CASE WHEN h.origin = 'unknown' THEN 'baseline' ELSE h.origin END
FROM configurations c
WHERE h.id = (
    SELECT h2.id
    FROM configuration_history h2
    WHERE h2.device_id = c.device_id
      AND h2.cfg_hash = c.cfg_hash
      AND h2.cfg_data = c.cfg_data
    ORDER BY h2.apply_at DESC, h2.id DESC
    LIMIT 1
);


DROP TRIGGER IF EXISTS trigger_configuration_history ON configurations;

CREATE OR REPLACE FUNCTION trg_configuration_history()
RETURNS TRIGGER AS $$
DECLARE
    configuration_changed BOOLEAN := FALSE;
BEGIN
    IF TG_OP = 'INSERT' THEN
        configuration_changed := NEW.cfg_data IS NOT NULL;
    ELSE
        configuration_changed :=
            NEW.cfg_data IS NOT NULL
            AND (
                NEW.cfg_hash IS DISTINCT FROM OLD.cfg_hash
                OR NEW.cfg_data IS DISTINCT FROM OLD.cfg_data
            );
    END IF;

    IF configuration_changed THEN
        -- Record every transition, including an intentional rollback to a
        -- configuration that appeared earlier in history. Consecutive
        -- duplicates are already excluded by configuration_changed.
        INSERT INTO configuration_history (
            apply_at,
            device_id,
            cfg_hash,
            cfg_data,
            cfg_sync_status,
            cfg_sync_error,
            origin
        )
        VALUES (
            timezone('UTC', now()),
            NEW.device_id,
            NEW.cfg_hash,
            NEW.cfg_data,
            NULL,
            NULL,
            'unknown'
        );

        -- Keep history bounded. 100 binary configurations are enough for
        -- support/debugging while preventing unbounded database growth.
        DELETE FROM configuration_history
        WHERE id IN (
            SELECT id
            FROM configuration_history
            WHERE device_id = NEW.device_id
            ORDER BY apply_at DESC, id DESC
            OFFSET 100
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_configuration_history
AFTER INSERT OR UPDATE ON configurations
FOR EACH ROW
EXECUTE FUNCTION trg_configuration_history();

COMMENT ON COLUMN configuration_history.cfg_sync_status IS 'Latest synchronization status recorded for this configuration version';
COMMENT ON COLUMN configuration_history.cfg_sync_error IS 'Last non-empty synchronization error observed for this configuration version';
COMMENT ON COLUMN configuration_history.origin IS 'Configuration origin: baseline, neosync, tracker or unknown';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_configuration_history ON configurations;
DROP FUNCTION IF EXISTS trg_configuration_history();

DROP INDEX IF EXISTS idx_configuration_history_device_hash;

ALTER TABLE configuration_history
    DROP COLUMN IF EXISTS origin,
    DROP COLUMN IF EXISTS cfg_sync_error,
    DROP COLUMN IF EXISTS cfg_sync_status;

DELETE FROM configuration_history
WHERE id IN (
    SELECT id
    FROM (
        SELECT
            id,
            ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY apply_at DESC, id DESC) AS row_num
        FROM configuration_history
    ) ranked
    WHERE ranked.row_num > 15
);

CREATE OR REPLACE FUNCTION trg_configuration_history()
RETURNS TRIGGER AS $$
DECLARE
    oldest_id INT;
    cnt INT;
BEGIN
    IF NEW.cfg_hash IS DISTINCT FROM OLD.cfg_hash OR NEW.cfg_data IS DISTINCT FROM OLD.cfg_data THEN
        SELECT COUNT(*) INTO cnt
        FROM configuration_history
        WHERE device_id = NEW.device_id;

        IF cnt >= 15 THEN
            SELECT id INTO oldest_id
            FROM configuration_history
            WHERE device_id = NEW.device_id
            ORDER BY apply_at ASC
            LIMIT 1;

            UPDATE configuration_history
            SET cfg_hash = NEW.cfg_hash,
                cfg_data = NEW.cfg_data,
                apply_at = now()
            WHERE id = oldest_id;
        ELSE
            INSERT INTO configuration_history (apply_at, device_id, cfg_hash, cfg_data)
            VALUES (now(), NEW.device_id, NEW.cfg_hash, NEW.cfg_data);
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_configuration_history
AFTER UPDATE ON configurations
FOR EACH ROW
EXECUTE FUNCTION trg_configuration_history();
-- +goose StatementEnd
