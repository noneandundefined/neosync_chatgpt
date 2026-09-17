package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/pgqx"
	"strings"
	"time"

	"github.com/lib/pq"
)

type DeviceStore struct {
	db *sql.DB
}

func (s *DeviceStore) Create_Device(ctx context.Context, tx *sql.Tx, device *models.Device) (uint64, error) {
	var id uint64

	query := `
		INSERT INTO devices (user_uuid, owner_uuid, name, phone, name_organization, imei, status, device_model)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, device.UserUUID, device.OwnerUUID, device.Name, device.Phone, device.NameOrganization, device.IMEI, device.Status, device.DeviceModel).Scan(&id)
	if err != nil {
		logger.Error("Create_Device req={%s} imei={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), device.IMEI, err.Error())
		return 0, err
	}

	return id, nil
}

func (s *DeviceStore) Create_DevicesBatch(ctx context.Context, tx *sql.Tx, devices []*models.Device) ([]*models.Device, error) {
	devicesInsert := make([]*models.Device, 0, len(devices))
	vStrings := make([]string, 0, len(devices))
	vArgs := make([]interface{}, 0, len(devices)*5)

	for idx, d := range devices {
		vStrings = append(vStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", idx*5+1, idx*5+2, idx*5+3, idx*5+4, idx*5+5))
		vArgs = append(vArgs, d.UserUUID, d.OwnerUUID, d.IMEI, d.Status, d.DeviceModel)
	}

	query := `
		INSERT INTO devices (user_uuid, owner_uuid, imei, status, device_model)
		VALUES ` + strings.Join(vStrings, ",") + `
		RETURNING id, imei
	`

	rows, err := tx.QueryContext(ctx, query, vArgs...)
	if err != nil {
		logger.Error("Create_DevicesBatch req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device models.Device

		if err := rows.Scan(
			&device.ID,
			&device.IMEI,
		); err != nil {
			return nil, err
		}

		devicesInsert = append(devicesInsert, &device)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Create_DevicesBatch req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return devicesInsert, nil
}

func (s *DeviceStore) Create_DeviceConf(ctx context.Context, tx *sql.Tx, device *models.DeviceConf) error {
	query := `
		INSERT INTO device_confs (device_id, request_configuration_on_connect)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, device.DeviceID, device.RequestConfigurationOnConnect)
	if err != nil {
		logger.Error("Create_DeviceConf req={%s} device_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), device.DeviceID, err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Create_DeviceConfBatch(ctx context.Context, tx *sql.Tx, devices []*models.DeviceConf) error {
	vStrings := make([]string, 0, len(devices))
	vArgs := make([]interface{}, 0, len(devices)*2)

	for idx, d := range devices {
		vStrings = append(vStrings, fmt.Sprintf("($%d, $%d)", idx*2+1, idx*2+2))
		vArgs = append(vArgs, d.DeviceID, d.RequestConfigurationOnConnect)
	}

	query := `
		INSERT INTO device_confs (device_id, request_configuration_on_connect)
		VALUES` + strings.Join(vStrings, ",") + `
	`

	_, err := tx.ExecContext(ctx, query, vArgs...)
	if err != nil {
		logger.Error("Create_DeviceConfBatch req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Get_Devices(ctx context.Context) ([]models.DeviceShort, error) {
	query := `
		SELECT
			devices.imei
		FROM devices
		ORDER BY devices.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	devices, err := pgqx.QueryContext[models.DeviceShort](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Devices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return devices, nil
}

func (s *DeviceStore) Get_ColumnTableDevice(ctx context.Context) ([]models.DatabaseSchemaColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_name = 'devices'
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deviceColumns, err := pgqx.QueryContext[models.DatabaseSchemaColumn](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_ColumnTableDevice req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return deviceColumns, nil
}

func (s *DeviceStore) Get_DeviceStatusByImei(ctx context.Context, imei string) (*models.DeviceStatus, error) {
	query := `
		SELECT
			updated_at, status, device_extended_model
		FROM devices
		WHERE imei = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deviceStatus, err := pgqx.QueryRowContext[models.DeviceStatus](ctx, s.db, query, imei)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceStatusByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return deviceStatus, nil
}

func (s *DeviceStore) Get_DeviceDetailedStatusByImeis(ctx context.Context, imeis []string) ([]models.DeviceDetailedStatus, error) {
	query := `
		SELECT
			devices.id,
			devices.imei,
			devices.status,
			devices.device_model,
			configurations.cfg_sync_status AS configuration_sync_status,
			syncs.firmware_version
		FROM devices
		LEFT JOIN configurations ON configurations.device_id = devices.id
		LEFT JOIN syncs ON syncs.device_id = devices.id
		WHERE devices.imei = ANY($1)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deviceStatus, err := pgqx.QueryContext[models.DeviceDetailedStatus](ctx, s.db, query, pq.Array(imeis))
	if err != nil {
		logger.Error("Get_DeviceDetailedStatusByImeis req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return deviceStatus, nil
}

/* Return - devices, totalFiltered, totalAll */
func (s *DeviceStore) Get_DevicesWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Device_DeviceConf_Sync, int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	allowedColumnSort := map[string]string{
		"imei":              "devices.imei",
		"device_model":      "devices.device_model",
		"name":              "devices.name",
		"phone":             "devices.phone",
		"name_organization": "devices.name_organization",
		"email":             "user_cores.email",
		"firmware_version":  "syncs.firmware_version",
		"last_mod_time":     "syncs.last_mod_time",
	}

	orderBy := "ORDER BY (devices.status = true) DESC"

	if col, ok := allowedColumnSort[columnSortKey]; ok {
		dir := "ASC"
		if columnSortDir == "desc" {
			dir = "DESC"
		}

		orderBy += fmt.Sprintf(", %s %s", col, dir)
	}

	baseCTE := `
		WITH filtered AS (
			SELECT devices.id
			FROM devices
			LEFT JOIN device_confs ON devices.id = device_confs.device_id
			LEFT JOIN syncs ON devices.id = syncs.device_id
			LEFT JOIN user_cores ON
				(devices.owner_uuid IS NOT NULL AND devices.owner_uuid = user_cores.user_uuid)
				OR
				(devices.owner_uuid IS NULL AND devices.user_uuid = user_cores.user_uuid)
			WHERE (
				$1 = ''
				OR concat_ws(' ',
					devices.name,
					devices.phone,
					devices.name_organization,
					devices.imei,
					devices.device_model,
					user_cores.email,
					syncs.last_mod_time,
					syncs.firmware_version
				) ILIKE '%' || $1 || '%'
			)
		)
	`

	query := baseCTE + fmt.Sprintf(`
		SELECT
			devices.*,
			device_confs.request_configuration_on_connect,
			configurations.updated_at AS configuration_updated_at,
			configurations.cfg_sync_status AS configuration_sync_status,
			configurations.cfg_sync_error AS configuration_sync_error,
			syncs.firmware_version,
			syncs.firmware_version_upd,
			syncs.firmware_updated_at,
			COALESCE(syncs.firmware_update_status, 'idle') AS firmware_update_status,
			syncs.cfg_version,
			syncs.last_mod_time,
			syncs.cfg_hash,
			user_cores.email AS user_email
		FROM devices
		JOIN filtered ON filtered.id = devices.id
		LEFT JOIN device_confs ON devices.id = device_confs.device_id
		LEFT JOIN configurations ON devices.id = configurations.device_id
		LEFT JOIN syncs ON devices.id = syncs.device_id
		LEFT JOIN user_cores ON
			(devices.owner_uuid IS NOT NULL AND devices.owner_uuid = user_cores.user_uuid)
			OR
			(devices.owner_uuid IS NULL AND devices.user_uuid = user_cores.user_uuid)
		%s
		LIMIT $2 OFFSET $3
	`, orderBy)

	devices, err := pgqx.QueryContext[models.Device_DeviceConf_Sync](ctx, s.db, query, search, limit, offset)
	if err != nil {
		logger.Error("Get_DevicesWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryCount := baseCTE + `
		SELECT COUNT(*) FROM filtered
	`

	var total int
	if err := tx.QueryRowContext(ctx, queryCount, search).Scan(&total); err != nil {
		logger.Error("Get_DevicesWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryTotalAll := `
		SELECT COUNT(*) FROM devices
	`

	var totalAll int
	if err := tx.QueryRowContext(ctx, queryTotalAll).Scan(&totalAll); err != nil {
		logger.Error("Get_DevicesWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Get_DevicesWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return devices, total, totalAll, nil
}

/* Return - devices, totalFiltered, totalAll */
func (s *DeviceStore) Get_DevicesByUuidsWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, roleCode string, userUuid, ownerUuid *string) ([]models.Device_DeviceConf_Sync, int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	allowedColumnSort := map[string]string{
		"imei":              "devices.imei",
		"device_model":      "devices.device_model",
		"name":              "devices.name",
		"phone":             "devices.phone",
		"name_organization": "devices.name_organization",
		"email":             "user_cores.email",
		"firmware_version":  "syncs.firmware_version",
		"last_mod_time":     "syncs.last_mod_time",
	}

	orderBy := "ORDER BY (devices.status = true) DESC"

	if col, ok := allowedColumnSort[columnSortKey]; ok {
		dir := "ASC"
		if columnSortDir == "desc" {
			dir = "DESC"
		}

		orderBy += fmt.Sprintf(", %s %s", col, dir)
	}

	userJoinOn := "devices.user_uuid = user_cores.user_uuid"
	if roleCode == constants.Role_AdminL2 || roleCode == constants.Role_AdminL2Support {
		userJoinOn = `(
				(devices.owner_uuid IS NOT NULL AND devices.owner_uuid = user_cores.user_uuid)
				OR
				(devices.owner_uuid IS NULL AND devices.user_uuid = user_cores.user_uuid)
			)`
	}

	baseCTE := `
		WITH filtered AS (
			SELECT devices.id
			FROM devices
			LEFT JOIN device_confs ON devices.id = device_confs.device_id
			LEFT JOIN syncs ON devices.id = syncs.device_id
			LEFT JOIN user_cores ON ` + userJoinOn + `
			WHERE
				-- DEALER/DEALER_SUPPORT
				($1::varchar IS NULL OR devices.user_uuid = $1)
				-- USER
				AND ($2::varchar IS NULL OR devices.owner_uuid = $2)
				AND (
					$3 = ''
					OR concat_ws(' ',
						devices.name,
						devices.phone,
						devices.name_organization,
						devices.imei,
						devices.device_model,
						user_cores.email,
						syncs.last_mod_time,
						syncs.firmware_version
					) ILIKE '%' || $3 || '%'
				)
		)
	`

	query := baseCTE + fmt.Sprintf(`
		SELECT
			devices.*,
			device_confs.request_configuration_on_connect,
			configurations.updated_at AS configuration_updated_at,
			configurations.cfg_sync_status AS configuration_sync_status,
			configurations.cfg_sync_error AS configuration_sync_error,
			syncs.firmware_version,
			syncs.firmware_version_upd,
			syncs.firmware_updated_at,
			COALESCE(syncs.firmware_update_status, 'idle') AS firmware_update_status,
			syncs.cfg_version,
			syncs.last_mod_time,
			syncs.cfg_hash,
			user_cores.email AS user_email
		FROM devices
		JOIN filtered ON filtered.id = devices.id
		LEFT JOIN device_confs ON devices.id = device_confs.device_id
		LEFT JOIN configurations ON devices.id = configurations.device_id
		LEFT JOIN syncs ON devices.id = syncs.device_id
		LEFT JOIN user_cores ON `+userJoinOn+`
		%s
		LIMIT $4 OFFSET $5
	`, orderBy)

	devices, err := pgqx.QueryContext[models.Device_DeviceConf_Sync](ctx, s.db, query, userUuid, ownerUuid, search, limit, offset)
	if err != nil {
		logger.Error("Get_DevicesByUuidsWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryCount := baseCTE + `
		SELECT COUNT(*) FROM filtered
	`

	var total int
	if err := tx.QueryRowContext(ctx, queryCount, userUuid, ownerUuid, search).Scan(&total); err != nil {
		logger.Error("Get_DevicesByUuidsWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	queryTotalAll := `
		SELECT COUNT(*) FROM devices
		WHERE
			-- DEALER/DEALER_SUPPORT
			($1::varchar IS NULL OR devices.user_uuid = $1)
			-- USER
			AND ($2::varchar IS NULL OR devices.owner_uuid = $2)
	`

	var totalAll int
	if err := tx.QueryRowContext(ctx, queryTotalAll, userUuid, ownerUuid).Scan(&totalAll); err != nil {
		logger.Error("Get_DevicesByUuidsWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Get_DevicesByUuidsWithParams req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, 0, err
	}

	return devices, total, totalAll, nil
}

func (s *DeviceStore) Get_DeviceConfByDeviceId(ctx context.Context, deviceId uint64) (*models.DeviceConf, error) {
	query := `
		SELECT * FROM device_confs
		WHERE device_id = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.DeviceConf](ctx, s.db, query, deviceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceConfByDeviceId req={%s} device_id={%d}: Failed to exec sql: %s", ctx.Value("XREQID").(string), deviceId, err.Error())
		return nil, err
	}

	return device, err
}

func (s *DeviceStore) Get_DevicesByImeis(ctx context.Context, imeis []string) ([]models.Device, error) {
	query := `
		SELECT * FROM devices WHERE imei = ANY($1)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	devices, err := pgqx.QueryContext[models.Device](ctx, s.db, query, pq.Array(imeis))
	if err != nil {
		logger.Error("Get_DevicesByImeis req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return devices, nil
}

func (s *DeviceStore) Get_DeviceByImei(ctx context.Context, imei string) (*models.Device, error) {
	query := `
		SELECT * FROM devices WHERE imei = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.Device](ctx, s.db, query, imei)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return device, nil
}

func (s *DeviceStore) Get_DeviceAndConfByImei(ctx context.Context, imei string) (*models.Device_DeviceConf_Sync, error) {
	query := `
		SELECT
			devices.*,
			device_confs.request_configuration_on_connect,
			configurations.updated_at AS configuration_updated_at,
			configurations.cfg_sync_status AS configuration_sync_status,
			configurations.cfg_sync_error AS configuration_sync_error,
			syncs.firmware_version,
			syncs.firmware_version_upd,
			syncs.firmware_updated_at,
			COALESCE(syncs.firmware_update_status, 'idle') AS firmware_update_status,
			syncs.cfg_version,
			syncs.last_mod_time,
			syncs.cfg_hash,
			user_cores.email AS user_email
		FROM devices
		LEFT JOIN device_confs ON devices.id = device_confs.device_id
		LEFT JOIN configurations ON devices.id = configurations.device_id
		LEFT JOIN syncs ON devices.id = syncs.device_id
		LEFT JOIN user_cores ON devices.user_uuid = user_cores.user_uuid
		WHERE devices.imei = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	device, err := pgqx.QueryRowContext[models.Device_DeviceConf_Sync](ctx, s.db, query, imei)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_DeviceAndConfByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return device, nil
}

func (s *DeviceStore) Update_DeviceStatusByImei(ctx context.Context, imei string, status bool) error {
	query := `
		UPDATE devices SET status = $1 WHERE imei = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, status, imei)
	if err != nil {
		logger.Error("Update_DeviceStatusByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_DeviceStatusByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *DeviceStore) Update_DeviceOwnerUuidByImei(ctx context.Context, imei string, ownerUuid *string) error {
	var ownerUuidParam sql.NullString
	if ownerUuid != nil && *ownerUuid != "" {
		ownerUuidParam = sql.NullString{String: *ownerUuid, Valid: true}
	} else {
		ownerUuidParam = sql.NullString{Valid: false}
	}

	query := `
		UPDATE devices
		SET owner_uuid = $1::varchar
		WHERE imei = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, ownerUuidParam, imei)
	if err != nil {
		logger.Error("Update_DeviceOwnerUuidByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DeviceUuidByImei(ctx context.Context, imei string, userUuid, ownerUuid *string) error {
	toNull := func(v *string) sql.NullString {
		if v != nil && *v != "" {
			return sql.NullString{String: *v, Valid: true}
		}

		return sql.NullString{Valid: false}
	}

	query := `
		UPDATE devices SET user_uuid = $1, owner_uuid = $2 WHERE imei = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, toNull(userUuid), toNull(ownerUuid), imei)
	if err != nil {
		logger.Error("Update_DeviceUuidByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DeviceUuidByDeviceId(ctx context.Context, deviceId uint64, userUuid *string) error {
	query := `
		UPDATE devices SET user_uuid = $1, owner_uuid = NULL WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUuid, deviceId)
	if err != nil {
		logger.Error("Update_DeviceUuidByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DeviceUuidByDeviceIdTx(ctx context.Context, tx *sql.Tx, deviceId uint64, userUuid, ownerUuid *string) error {
	query := `
		UPDATE devices SET user_uuid = $1, owner_uuid = $2 WHERE id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userUuid, ownerUuid, deviceId)
	if err != nil {
		logger.Error("Update_DeviceUuidByDeviceIdTx req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DeviceByImei(ctx context.Context, tx *sql.Tx, imei string, device *models.UpdateDevice) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fields := []string{}
	values := []interface{}{}
	i := 1

	if device.Imei != nil {
		fields = append(fields, fmt.Sprintf("imei=$%d", i))
		values = append(values, *device.Imei)
		i++
	}

	if device.DeviceModel != nil {
		fields = append(fields, fmt.Sprintf("device_model=$%d", i))
		values = append(values, *device.DeviceModel)
		i++
	}

	if device.Name != nil {
		fields = append(fields, fmt.Sprintf("name=$%d", i))
		values = append(values, *device.Name)
		i++
	}

	if device.Phone != nil {
		fields = append(fields, fmt.Sprintf("phone=$%d", i))
		values = append(values, *device.Phone)
		i++
	}

	if device.NameOrganization != nil {
		fields = append(fields, fmt.Sprintf("name_organization=$%d", i))
		values = append(values, *device.NameOrganization)
		i++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE devices SET %s WHERE imei = $%d
	`, strings.Join(fields, ", "), i)

	values = append(values, imei)

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		logger.Error("Update_DeviceByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DevicesUuidByImeiToNull(ctx context.Context, imeis []string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `SELECT update_devices_uuid_by_imeis_to_null($1)`, pq.Array(imeis)); err != nil {
			logger.Error("Update_DevicesUuidByImeiToNull req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		return nil
	})
}

func (s *DeviceStore) Update_DeviceUnlinkFromAccountByImei(ctx context.Context, imei, userUuid string) error {
	query := `
		UPDATE devices
		SET user_uuid = NULL, owner_uuid = NULL
		WHERE imei = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei, userUuid)
	if err != nil {
		logger.Error("Update_DeviceUnlinkFromAccountByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DevicesUnlinkFromAccountByImei(ctx context.Context, imeis []string, userUuid string) error {
	if len(imeis) == 0 {
		return httperr.Err_NotDeleted
	}

	query := `
		UPDATE devices
		SET user_uuid = NULL, owner_uuid = NULL
		WHERE imei = ANY($1) AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, pq.Array(imeis), userUuid)
	if err != nil {
		logger.Error("Update_DevicesUnlinkFromAccountByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Update_DeviceConfByDeviceId(ctx context.Context, tx *sql.Tx, deviceId uint64, device *models.UpdateDevice) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fields := []string{}
	values := []interface{}{}
	i := 1

	if device.RequestConfigurationOnConnect != nil {
		fields = append(fields, fmt.Sprintf("request_configuration_on_connect=$%d", i))
		values = append(values, *device.RequestConfigurationOnConnect)
		i++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE device_confs SET %s WHERE device_id = $%d
	`, strings.Join(fields, ", "), i)

	values = append(values, deviceId)

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		logger.Error("Update_DeviceConfByDeviceId req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Delete_DevicesByImei(ctx context.Context, imeis []string) error {
	if len(imeis) == 0 {
		return httperr.Err_NotDeleted
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `SELECT delete_devices_by_imeis($1)`, pq.Array(imeis)); err != nil {
			logger.Error("Delete_DevicesByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		return nil
	})
}

func (s *DeviceStore) Delete_DevicesByImeiAndUserUuid(ctx context.Context, imeis []string, userUuid string) error {
	if len(imeis) == 0 {
		return httperr.Err_NotDeleted
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `SELECT delete_devices_by_imeis_and_user_uuid($1, $2)`, pq.Array(imeis), userUuid); err != nil {
			logger.Error("Delete_DevicesByImeiAndUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		return nil
	})
}

func (s *DeviceStore) Delete_DeviceByImeiAndUserUuid(ctx context.Context, imei, userUuid string) error {
	query := `
		DELETE FROM devices
		WHERE imei = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei, userUuid)
	if err != nil {
		logger.Error("Delete_DeviceByImeiAndUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *DeviceStore) Delete_DeviceByImei(ctx context.Context, imei string) error {
	query := `
		DELETE FROM devices WHERE imei = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, imei)
	if err != nil {
		logger.Error("Delete_DeviceByImei req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
