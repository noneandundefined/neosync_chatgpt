-- +goose Up
-- +goose StatementBegin
-- +View
CREATE MATERIALIZED VIEW materialized_device_model_sources_device_model AS
SELECT *
FROM device_model_sources
WITH DATA ;

-- REFRESH MATERIALIZED VIEW materialized_device_model_sources_device_model;
-- +View
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- -View
DROP MATERIALIZED VIEW IF EXISTS materialized_device_model_sources_device_model;
-- -View
-- +goose StatementEnd
