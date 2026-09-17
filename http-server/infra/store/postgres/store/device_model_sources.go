package store

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/pgqx"
	"time"
)

type DeviceModelSources struct {
	db *sql.DB
}

func (s *DeviceModelSources) Get_Sources(ctx context.Context) ([]models.DeviceModelSource, error) {
	query := `
		SELECT * FROM materialized_device_model_sources_device_model
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deviceModelSources, err := pgqx.QueryContext[models.DeviceModelSource](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Sources req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return deviceModelSources, nil
}

func (s *DeviceModelSources) Get_Models(ctx context.Context) ([]string, error) {
	query := `
		SELECT device_model FROM materialized_device_model_sources_device_model
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	results, err := pgqx.QueryContext[models.DeviceModelResult](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Models req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	models := make([]string, len(results))
	for i, r := range results {
		models[i] = r.DeviceModel
	}

	return models, nil
}

func (s *DeviceModelSources) Get_FirmwareUrl(ctx context.Context) ([]string, error) {
	query := `
		SELECT firmware_url FROM materialized_device_model_sources_device_model
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	results, err := pgqx.QueryContext[models.FirmwareUrlResult](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_FirmwareUrl req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	headUrls := make([]string, len(results))
	for i, r := range results {
		headUrls[i] = r.FirmwareUrl
	}

	return headUrls, nil
}
