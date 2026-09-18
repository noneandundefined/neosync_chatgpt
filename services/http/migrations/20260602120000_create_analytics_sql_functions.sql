-- +goose Up
-- +goose StatementBegin
-- Разбор JSONB-массивов (int[]) в байты с позицией (1-based), как fuel_sensor_address_list.
CREATE OR REPLACE FUNCTION analytic_jsonb_int_bytes(p_data JSONB)
RETURNS TABLE(byte_val INT, pos BIGINT)
LANGUAGE sql
IMMUTABLE
AS $$
	SELECT (elem.value)::INT, elem.ordinality::BIGINT
	FROM jsonb_array_elements(COALESCE(p_data, '[]'::JSONB)) WITH ORDINALITY AS elem(value, ordinality);
$$;

-- Слоты из массива байтов: bytes_per_slot байт на слот, не более max_slots слотов.
CREATE OR REPLACE FUNCTION analytic_byte_slots(
	p_data JSONB,
	p_bytes_per_slot INT DEFAULT 6,
	p_max_slots INT DEFAULT 8
)
RETURNS TABLE(slot INT, not_empty BOOLEAN)
LANGUAGE sql
IMMUTABLE
AS $$
	WITH raw AS (
		SELECT b.byte_val, b.pos
		FROM analytic_jsonb_int_bytes(p_data) AS b
		WHERE p_bytes_per_slot > 0
			AND p_max_slots > 0
			AND b.pos <= (p_bytes_per_slot * p_max_slots)::BIGINT
	),
	grouped AS (
		SELECT
			((raw.pos - 1) / p_bytes_per_slot)::INT AS slot,
			BOOL_OR(raw.byte_val <> 0) AS not_empty
		FROM raw
		GROUP BY ((raw.pos - 1) / p_bytes_per_slot)
	)
	SELECT grouped.slot, grouped.not_empty
	FROM grouped
	WHERE grouped.not_empty;
$$;

-- Типизированные слоты: type_mask — бит slot=1 → BLE, иначе RS-485 (как fuel_sensor_type).
CREATE OR REPLACE FUNCTION analytic_typed_byte_slots(
	p_address_list JSONB,
	p_type_mask INT,
	p_bytes_per_slot INT DEFAULT 6,
	p_max_slots INT DEFAULT 8
)
RETURNS TABLE(slot INT, not_empty BOOLEAN, is_ble BOOLEAN, is_rs485 BOOLEAN)
LANGUAGE sql
IMMUTABLE
AS $$
	SELECT
		s.slot,
		s.not_empty,
		((COALESCE(p_type_mask, 0) >> s.slot) & 1) = 1 AS is_ble,
		((COALESCE(p_type_mask, 0) >> s.slot) & 1) = 0 AS is_rs485
	FROM analytic_byte_slots(p_address_list, p_bytes_per_slot, p_max_slots) AS s;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS analytic_typed_byte_slots(JSONB, INT, INT, INT);
DROP FUNCTION IF EXISTS analytic_byte_slots(JSONB, INT, INT);
DROP FUNCTION IF EXISTS analytic_jsonb_int_bytes(JSONB);
-- +goose StatementEnd
