-- +goose Up
-- +goose StatementBegin
-- BLE-адреса (MacAddressArray) и SN (SNSArray) хранятся как JSONB-массивы строк,
-- а не int-байтов. analytic_jsonb_int_bytes должен безопасно читать и числа, и строки.
CREATE OR REPLACE FUNCTION analytic_jsonb_int_bytes(p_data JSONB)
RETURNS TABLE(byte_val INT, pos BIGINT)
LANGUAGE sql
IMMUTABLE
AS $$
	SELECT
		CASE jsonb_typeof(elem.value)
			WHEN 'number' THEN (elem.value)::INT
			WHEN 'string' THEN
				CASE
					WHEN (elem.value #>> '{}') ~ '^-?[0-9]+$' THEN (elem.value #>> '{}')::INT
					ELSE 0
				END
			ELSE 0
		END,
		elem.ordinality::BIGINT
	FROM jsonb_array_elements(COALESCE(p_data, '[]'::JSONB)) WITH ORDINALITY AS elem(value, ordinality);
$$;

CREATE OR REPLACE VIEW analytics_configuration_prof AS
WITH sensor_slot_counts AS (
	SELECT
		configurations_prof.device_id,
		configurations_prof.user_uuid,
		configurations_prof.updated_at,
		configurations_prof.server_host,
		configurations_prof.remote_server_host_address,
		configurations_prof.traffic_protocol,
		configurations_prof.points_period_move,
		configurations_prof.device_mode,
		configurations_prof.static_mode,
		(
			COALESCE((
				SELECT COUNT(*)::INT
				FROM analytic_byte_slots(configurations_prof.fuel_sensor_address_list, 6, 8)
			), 0)
			+ COALESCE((
				SELECT COUNT(*)::INT
				FROM jsonb_array_elements_text(COALESCE(configurations_prof.ble_adm_sensor_address_list, '[]'::jsonb)) AS ble_adm_sensor_mac
				WHERE NULLIF(TRIM(ble_adm_sensor_mac), '') IS NOT NULL
					AND REPLACE(UPPER(TRIM(ble_adm_sensor_mac)), ':', '') <> '000000000000'
			), 0)
			+ COALESCE((
				SELECT COUNT(*)::INT
				FROM jsonb_array_elements_text(COALESCE(configurations_prof.ble_fuel_sensor_address_list, '[]'::jsonb)) AS ble_fuel_sensor_mac
				WHERE NULLIF(TRIM(ble_fuel_sensor_mac), '') IS NOT NULL
					AND REPLACE(UPPER(TRIM(ble_fuel_sensor_mac)), ':', '') <> '000000000000'
			), 0)
			+ COALESCE((
				SELECT COUNT(*)::INT
				FROM analytic_byte_slots(configurations_prof.lls_fuel_sensors_addr, 6, 8)
			), 0)
			+ COALESCE((
				SELECT COUNT(*)::INT
				FROM jsonb_array_elements_text(COALESCE(configurations_prof.ow_temp_sn, '[]'::jsonb)) AS ow_temp_sn_value
				WHERE NULLIF(TRIM(ow_temp_sn_value), '') IS NOT NULL
			), 0)
		) AS sensor_count
	FROM configurations_prof
)
SELECT
	sensor_slot_counts.device_id,
	sensor_slot_counts.user_uuid,
	sensor_slot_counts.updated_at,
	sensor_slot_counts.server_host,
	sensor_slot_counts.remote_server_host_address,
	sensor_slot_counts.traffic_protocol,
	sensor_slot_counts.points_period_move,
	sensor_slot_counts.device_mode,
	sensor_slot_counts.static_mode,
	sensor_slot_counts.sensor_count,
	CASE
		WHEN sensor_slot_counts.server_host IS NOT NULL
			AND NULLIF(sensor_slot_counts.remote_server_host_address, '') IS NOT NULL
			THEN 'main_and_reserve'
		WHEN sensor_slot_counts.server_host IS NOT NULL THEN 'one_server'
		ELSE 'other'
	END AS server_schema_bucket,
	CASE
		WHEN sensor_slot_counts.sensor_count = 0 THEN 'no_sensors'
		ELSE 'one_or_more'
	END AS sensors_bucket_exclusive,
	CASE
		WHEN sensor_slot_counts.traffic_protocol::TEXT ILIKE '%tcp%' THEN 'tcp'
		WHEN sensor_slot_counts.traffic_protocol::TEXT ILIKE '%udp%' THEN 'udp'
		ELSE 'other'
	END AS protocol_bucket,
	CASE
		WHEN COALESCE(sensor_slot_counts.points_period_move, 0) = 0 THEN 'period_0'
		WHEN sensor_slot_counts.points_period_move <= 60 THEN 'from_30_to_60'
		WHEN sensor_slot_counts.points_period_move <= 300 THEN 'from_60_to_300'
		ELSE 'over_300'
	END AS period_bucket
FROM sensor_slot_counts;

COMMENT ON VIEW analytics_configuration_prof IS
	'Analytics buckets over configurations_prof: sensor count, server schema, protocol, transmission period (points_period_move = PERIOD in motion).';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS analytics_configuration_prof;

CREATE OR REPLACE FUNCTION analytic_jsonb_int_bytes(p_data JSONB)
RETURNS TABLE(byte_val INT, pos BIGINT)
LANGUAGE sql
IMMUTABLE
AS $$
	SELECT (elem.value)::INT, elem.ordinality::BIGINT
	FROM jsonb_array_elements(COALESCE(p_data, '[]'::JSONB)) WITH ORDINALITY AS elem(value, ordinality);
$$;
-- +goose StatementEnd
