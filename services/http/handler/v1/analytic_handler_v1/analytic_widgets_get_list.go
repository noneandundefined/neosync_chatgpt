package analytic_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение виджетов для аналитики */

func (h *Handler) GetAnalyticWidgetsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	tr := middleware.TranslatorFromContext(r.Context())

	widgets := []WidgetResponse{
		{
			ID:       1,
			Section:  "devices",
			Type:     "stat",
			Title:    tr.T("analytic-widget-active-trackers-title"),
			Subtitle: tr.T("analytic-widget-active-trackers-subtitle"),
			Icon:     "car",
			Query: `
				SELECT COUNT(*)::numeric AS value
				FROM devices
				WHERE status = true
			`,
		},
		{
			ID:      2,
			Section: "overview",
			Type:    "stat",
			Title:   tr.T("analytic-widget-active-users-title"),
			Icon:    "account-multiple",
			Query: `
				SELECT COUNT(DISTINCT user_uuid)::numeric AS value
				FROM analytics_users
				WHERE user_uuid = user_uuid
			`,
		},
		{
			ID:      3,
			Section: "configurations",
			Type:    "stat",
			Title:   tr.T("analytic-widget-config-created-title"),
			Icon:    "poll",
			Query: `
				SELECT COUNT(*)::numeric AS value
				FROM analytic_records
				WHERE user_uuid = user_uuid
					AND device_id = device_id
			`,
		},
		{
			ID:      4,
			Section: "configurations",
			Type:    "stat",
			Title:   tr.T("analytic-widget-config-sent-title"),
			Icon:    "configuration-send",
			Query: `
				SELECT COUNT(*)::numeric AS value
				FROM analytic_records
				WHERE rmq_published = true
					AND user_uuid = user_uuid
					AND device_id = device_id
			`,
		},
		{
			ID:      5,
			Section: "errors",
			Type:    "stat",
			Title:   tr.T("analytic-widget-config-send-errors-title"),
			Icon:    "error",
			Query: `
				SELECT COUNT(*)::numeric AS value
				FROM analytic_records
				WHERE (error_reason IS NOT NULL OR device_online = false)
					AND user_uuid = user_uuid
					AND device_id = device_id
			`,
		},
		{
			ID:      6,
			Section: "configurations",
			Type:    "stat",
			Title:   tr.T("analytic-widget-config-success-apply-title"),
			Icon:    "success",
			Query: `
				SELECT COUNT(*)::numeric AS value
				FROM analytic_records
				WHERE rmq_published = true AND error_reason IS NULL AND device_online = true
					AND user_uuid = user_uuid
					AND device_id = device_id
			`,
		},
		{
			ID:       7,
			Section:  "configurations",
			Type:     "bar",
			Title:    tr.T("analytic-widget-popular-configs-title"),
			Subtitle: tr.T("analytic-widget-popular-configs-subtitle"),
			Query: `
				WITH bounds AS (
					SELECT {{from}}::timestamptz AS from_ts, {{to}}::timestamptz AS to_ts
				),
				group_devices AS (
					SELECT device_id
					FROM group_user_devices
					WHERE {{group_id}} IS NOT NULL AND group_id = {{group_id}}
				),
				base AS (
					SELECT *
					FROM analytics_configuration_prof
					WHERE updated_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)
						AND ({{user_uuid}} IS NULL OR user_uuid = {{user_uuid}})
						AND ({{group_id}} IS NULL OR device_id IN (SELECT device_id FROM group_devices))
				),
				unioned AS (
					SELECT 'server_schema'::text AS category, server_schema_bucket::text AS bucket, COUNT(*)::BIGINT AS value
					FROM base
					GROUP BY server_schema_bucket
					UNION ALL
					SELECT 'sensors'::text, 'no_sensors'::text, COUNT(*)::BIGINT
					FROM base
					WHERE sensor_count = 0
					UNION ALL
					SELECT 'sensors'::text, 'one_or_more'::text, COUNT(*)::BIGINT
					FROM base
					WHERE sensor_count >= 1
					UNION ALL
					SELECT 'sensors'::text, 'two_or_more'::text, COUNT(*)::BIGINT
					FROM base
					WHERE sensor_count >= 2
					UNION ALL
					SELECT 'protocol'::text, protocol_bucket::text, COUNT(*)::BIGINT AS value
					FROM base
					GROUP BY protocol_bucket
					UNION ALL
					SELECT 'period'::text, period_bucket::text, COUNT(*)::BIGINT AS value
					FROM base
					GROUP BY period_bucket
				)
				SELECT
					category,
					bucket,
					value,
					CASE
						WHEN SUM(value) OVER (PARTITION BY category) = 0 THEN 0
						ELSE ROUND((value * 100.0) / SUM(value) OVER (PARTITION BY category), 0)::BIGINT
					END AS percent
				FROM unioned
				ORDER BY category, value DESC
			`,
		},
		{
			ID:      8,
			Section: "configurations",
			Type:    "line_chart",
			Title:   tr.T("analytic-widget-config-trend-title"),
			Query: `
				SELECT
					DATE(created_at) AS date,
					COUNT(*) AS created,
					COUNT(*) FILTER (WHERE device_online = true AND rmq_published = true) AS sent
				FROM analytic_records
				WHERE user_uuid = user_uuid
					AND device_id = device_id
				GROUP BY DATE(created_at)
				ORDER BY DATE(created_at)
			`,
		},
		{
			ID:      9,
			Section: "configurations",
			Type:    "gauge_error",
			Title:   tr.T("analytic-widget-config-success-fail-title"),
			Query: `
				SELECT
					CASE
						WHEN (error_reason IS NOT NULL OR device_online = false) THEN 'failed'
						WHEN rmq_published = true THEN 'success'
						ELSE 'ignored'
					END AS status,
					COUNT(*) AS value
				FROM analytic_records
				WHERE user_uuid = user_uuid
					AND device_id = device_id
				GROUP BY status
			`,
		},
		{
			ID:      10,
			Section: "errors",
			Type:    "table",
			Title:   tr.T("analytic-widget-config-error-reasons-title"),
			Query: `
				WITH bounds AS (
					SELECT {{from}}::timestamptz AS from_ts, {{to}}::timestamptz AS to_ts
				),
				group_devices AS (
					SELECT device_id
					FROM group_user_devices
					WHERE {{group_id}} IS NOT NULL AND group_id = {{group_id}}
				),
				base AS (
					SELECT
						analytic_records.created_at,
						COALESCE(analytic_records.error_reason, 'tracker offline') AS error_reason,
						devices.imei AS device_imei,
						COALESCE(NULLIF(user_cores.email, ''), analytic_records.user_uuid) AS account_name
					FROM analytic_records
					LEFT JOIN devices ON devices.id = analytic_records.device_id
					LEFT JOIN user_cores ON user_cores.user_uuid = analytic_records.user_uuid
					WHERE (analytic_records.error_reason IS NOT NULL OR analytic_records.device_online = false)
						AND analytic_records.created_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)
						AND ({{user_uuid}} IS NULL OR analytic_records.user_uuid = {{user_uuid}})
						AND ({{group_id}} IS NULL OR analytic_records.device_id IN (SELECT device_id FROM group_devices))
				),
				top_reasons AS (
					SELECT
						'reason'::text AS row_type,
						error_reason AS reason,
						COUNT(*)::BIGINT AS reason_count,
						NULL::timestamptz AS operation_time,
						NULL::text AS account_name,
						NULL::text AS device_imei,
						NULL::text AS operation_reason
					FROM base
					GROUP BY error_reason
					ORDER BY reason_count DESC
					LIMIT 5
				),
				recent_ops AS (
					SELECT
						'recent'::text AS row_type,
						NULL::text AS reason,
						NULL::BIGINT AS reason_count,
						created_at AS operation_time,
						account_name,
						device_imei,
						error_reason AS operation_reason
					FROM base
					ORDER BY created_at DESC
					LIMIT 5
				)
				SELECT * FROM top_reasons
				UNION ALL
				SELECT * FROM recent_ops
			`,
		},
		{
			ID:      11,
			Section: "errors",
			Type:    "gauge_top_reasons",
			Title:   tr.T("analytic-widget-config-success-fail-title"),
			Query: `
				WITH bounds AS (
					SELECT {{from}}::timestamptz AS from_ts, {{to}}::timestamptz AS to_ts
				),
				group_devices AS (
					SELECT device_id
					FROM group_user_devices
					WHERE {{group_id}} IS NOT NULL AND group_id = {{group_id}}
				),
				base AS (
					SELECT error_reason, device_online
					FROM analytic_records
					WHERE (error_reason IS NOT NULL OR device_online = false)
						AND created_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)
						AND ({{user_uuid}} IS NULL OR user_uuid = {{user_uuid}})
						AND ({{group_id}} IS NULL OR device_id IN (SELECT device_id FROM group_devices))
				)
				SELECT reason, reason_count
				FROM (
					SELECT error_reason AS reason, COUNT(*)::BIGINT AS reason_count
					FROM base
					WHERE error_reason IS NOT NULL
					GROUP BY error_reason
					UNION ALL
					SELECT 'tracker offline' AS reason, COUNT(*)::BIGINT AS reason_count
					FROM base
					WHERE device_online = false
					GROUP BY 1
				) t
				ORDER BY reason_count DESC
			`,
		},
		{
			ID: 101, Section: "overview", Type: "stat", Title: "Сессии",
			Query: `SELECT COUNT(DISTINCT session_id)::numeric AS value FROM product_analytics_events WHERE session_id = session_id`,
		},
		{
			ID: 102, Section: "overview", Type: "stat", Title: "Полезные действия",
			Query: `SELECT COUNT(*)::numeric AS value FROM product_analytics_events WHERE success = true AND category IN ('device', 'configuration', 'command', 'firmware', 'group') AND user_uuid = user_uuid`,
		},
		{
			ID: 103, Section: "overview", Type: "stat", Title: "Пользователи с ошибками",
			Query: `SELECT COUNT(DISTINCT user_uuid)::numeric AS value FROM product_analytics_events WHERE success = false AND user_uuid = user_uuid`,
		},
		{
			ID: 104, Section: "overview", Type: "report_line", Title: "Активность сервиса",
			Query: `SELECT DATE(occurred_at) AS date, COUNT(DISTINCT session_id) AS sessions, COUNT(DISTINCT user_uuid) AS users, COUNT(*) FILTER (WHERE success = true) AS successful_actions, COUNT(*) FILTER (WHERE success = false) AS errors FROM product_analytics_events WHERE user_uuid = user_uuid GROUP BY DATE(occurred_at) ORDER BY DATE(occurred_at)`,
		},
		{
			ID: 201, Section: "funnels", Type: "funnel", Title: "Настройка устройства",
			Query: `SELECT step, value FROM (VALUES (1, 'configuration_opened', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'page_view' AND path LIKE '/configurations/devices/%')), (2, 'draft_saved', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'configuration_draft_save_success')), (3, 'apply_started', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name IN ('configuration_apply_success', 'configuration_apply_failed'))), (4, 'configuration_sent', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'configuration_apply_success'))) f(sort, step, value) ORDER BY sort`,
		},
		{
			ID: 202, Section: "funnels", Type: "funnel", Title: "Подключение нового устройства",
			Query: `SELECT step, value FROM (VALUES (1, 'device_creation_started', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name IN ('device_create_success', 'device_create_failed'))), (2, 'device_created', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'device_create_success')), (3, 'configuration_opened', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'page_view' AND path LIKE '/configurations/devices/%')), (4, 'configuration_applied', (SELECT COUNT(DISTINCT session_id) FROM product_analytics_events WHERE event_name = 'configuration_apply_success'))) f(sort, step, value) ORDER BY sort`,
		},
		{
			ID: 203, Section: "funnels", Type: "report_table", Title: "Время до полезного результата",
			Query: `WITH session_steps AS (SELECT session_id, MIN(occurred_at) FILTER (WHERE event_name IN ('session_started', 'session_authenticated_success', 'login_success')) AS started_at, MIN(occurred_at) FILTER (WHERE event_name IN ('device_create_success', 'configuration_apply_success', 'command_send_success')) AS value_at FROM product_analytics_events GROUP BY session_id) SELECT COUNT(*) FILTER (WHERE value_at IS NOT NULL) AS activated_sessions, ROUND(AVG(EXTRACT(EPOCH FROM (value_at - started_at))) FILTER (WHERE value_at >= started_at)) AS avg_seconds, ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (value_at - started_at))) FILTER (WHERE value_at >= started_at)) AS median_seconds FROM session_steps`,
		},
		{
			ID: 301, Section: "behavior", Type: "stat", Title: "Просмотры страниц",
			Query: `SELECT COUNT(*)::numeric AS value FROM product_analytics_events WHERE event_name = 'page_view' AND user_uuid = user_uuid`,
		},
		{
			ID: 302, Section: "behavior", Type: "stat", Title: "Среднее время на странице, сек.",
			Query: `SELECT COALESCE(ROUND(AVG(duration_ms) / 1000), 0)::numeric AS value FROM product_analytics_events WHERE event_name = 'page_duration' AND duration_ms BETWEEN 0 AND 3600000 AND user_uuid = user_uuid`,
		},
		{
			ID: 303, Section: "behavior", Type: "report_table", Title: "Популярные страницы",
			Query: `SELECT path AS page, COUNT(*) AS views, COUNT(DISTINCT user_uuid) AS users, ROUND(AVG(duration_ms) FILTER (WHERE event_name = 'page_duration') / 1000) AS avg_seconds FROM product_analytics_events WHERE event_name IN ('page_view', 'page_duration') AND user_uuid = user_uuid GROUP BY path ORDER BY views DESC NULLS LAST LIMIT 20`,
		},
		{
			ID: 304, Section: "behavior", Type: "report_table", Title: "Использование элементов интерфейса",
			Query: `SELECT COALESCE(properties->>'control', 'unknown') AS control, path, COUNT(*) AS clicks, COUNT(DISTINCT user_uuid) AS users FROM product_analytics_events WHERE event_name = 'control_clicked' AND user_uuid = user_uuid GROUP BY control, path ORDER BY clicks DESC LIMIT 20`,
		},
		{
			ID: 305, Section: "behavior", Type: "report_table", Title: "Поиск и пустые результаты",
			Query: `SELECT properties->>'api_path' AS api_path, COUNT(*) AS searches, COUNT(*) FILTER (WHERE event_name = 'search_zero_results') AS zero_results, ROUND(100.0 * COUNT(*) FILTER (WHERE event_name = 'search_zero_results') / NULLIF(COUNT(*), 0), 1) AS zero_percent FROM product_analytics_events WHERE category = 'search' AND user_uuid = user_uuid GROUP BY api_path ORDER BY searches DESC LIMIT 20`,
		},
		{
			ID: 306, Section: "behavior", Type: "report_table", Title: "Активность по ролям",
			Query: `SELECT COALESCE(role_code, 'anonymous') AS role, COUNT(DISTINCT user_uuid) AS users, COUNT(DISTINCT session_id) AS sessions, COUNT(*) AS events FROM product_analytics_events GROUP BY role_code ORDER BY events DESC`,
		},
		{
			ID: 401, Section: "errors", Type: "stat", Title: "Ошибки интерфейса и API",
			Query: `SELECT COUNT(*)::numeric AS value FROM product_analytics_events WHERE success = false AND user_uuid = user_uuid`,
		},
		{
			ID: 402, Section: "errors", Type: "report_table", Title: "Частые ошибки пользователей",
			Query: `SELECT event_name AS event, COALESCE(error_code, 'unknown') AS error, path, COUNT(*) AS occurrences, COUNT(DISTINCT user_uuid) AS users FROM product_analytics_events WHERE success = false AND user_uuid = user_uuid GROUP BY event_name, error_code, path ORDER BY occurrences DESC LIMIT 25`,
		},
		{
			ID: 403, Section: "errors", Type: "report_table", Title: "Поля с ошибками валидации",
			Query: `SELECT properties->>'field' AS field, properties->>'input_type' AS input_type, error_code, path, COUNT(*) AS occurrences FROM product_analytics_events WHERE event_name = 'form_validation_failed' AND user_uuid = user_uuid GROUP BY field, input_type, error_code, path ORDER BY occurrences DESC LIMIT 25`,
		},
		{
			ID: 404, Section: "errors", Type: "report_table", Title: "Восстановление после ошибок",
			Query: `WITH failures AS (SELECT session_id, event_name, occurred_at FROM product_analytics_events WHERE success = false), recovered AS (SELECT DISTINCT failures.session_id, failures.event_name FROM failures WHERE EXISTS (SELECT 1 FROM product_analytics_events WHERE product_analytics_events.session_id = failures.session_id AND product_analytics_events.success = true AND product_analytics_events.occurred_at > failures.occurred_at)) SELECT event_name AS failed_event, COUNT(*) AS failures, COUNT(*) FILTER (WHERE recovered.session_id IS NOT NULL) AS recovered, ROUND(100.0 * COUNT(*) FILTER (WHERE recovered.session_id IS NOT NULL) / NULLIF(COUNT(*), 0), 1) AS recovery_percent FROM failures LEFT JOIN recovered USING (session_id, event_name) GROUP BY event_name ORDER BY failures DESC LIMIT 20`,
		},
		{
			ID: 501, Section: "retention", Type: "report_table", Title: "Возвращаемость пользователей",
			Query: `WITH first_seen AS (SELECT user_uuid, MIN(DATE(occurred_at)) AS cohort FROM product_analytics_events GROUP BY user_uuid), activity AS (SELECT DISTINCT user_uuid, DATE(occurred_at) AS day FROM product_analytics_events) SELECT first_seen.cohort, COUNT(*) AS new_users, COUNT(*) FILTER (WHERE EXISTS (SELECT 1 FROM activity WHERE activity.user_uuid = first_seen.user_uuid AND activity.day = first_seen.cohort + 1)) AS day_1, COUNT(*) FILTER (WHERE EXISTS (SELECT 1 FROM activity WHERE activity.user_uuid = first_seen.user_uuid AND activity.day BETWEEN first_seen.cohort + 7 AND first_seen.cohort + 13)) AS day_7, COUNT(*) FILTER (WHERE EXISTS (SELECT 1 FROM activity WHERE activity.user_uuid = first_seen.user_uuid AND activity.day BETWEEN first_seen.cohort + 30 AND first_seen.cohort + 59)) AS day_30 FROM first_seen GROUP BY first_seen.cohort ORDER BY first_seen.cohort DESC LIMIT 20`,
		},
		{
			ID: 502, Section: "retention", Type: "report_line", Title: "Активные пользователи",
			Query: `SELECT DATE(occurred_at) AS date, COUNT(DISTINCT user_uuid) AS daily_users, COUNT(DISTINCT session_id) AS sessions FROM product_analytics_events WHERE user_uuid = user_uuid GROUP BY DATE(occurred_at) ORDER BY DATE(occurred_at)`,
		},
		{
			ID: 601, Section: "acquisition", Type: "report_table", Title: "Источники входа",
			Query: `SELECT COALESCE(NULLIF(properties->>'utm_source', ''), NULLIF(properties->>'referrer_host', ''), 'direct') AS source, COALESCE(NULLIF(properties->>'utm_medium', ''), 'none') AS medium, COUNT(DISTINCT session_id) AS sessions, COUNT(DISTINCT user_uuid) AS users FROM product_analytics_events WHERE event_name = 'session_started' GROUP BY source, medium ORDER BY sessions DESC LIMIT 20`,
		},
		{
			ID: 602, Section: "acquisition", Type: "report_table", Title: "Посадочные страницы",
			Query: `SELECT path AS landing_page, COUNT(DISTINCT session_id) AS sessions, COUNT(DISTINCT user_uuid) AS users FROM product_analytics_events WHERE event_name = 'page_view' GROUP BY path ORDER BY sessions DESC LIMIT 20`,
		},
		{
			ID: 603, Section: "acquisition", Type: "report_table", Title: "Конверсия источников в полезное действие",
			Query: `WITH sources AS (SELECT session_id, COALESCE(NULLIF(properties->>'utm_source', ''), NULLIF(properties->>'referrer_host', ''), 'direct') AS source FROM product_analytics_events WHERE event_name = 'session_started'), outcomes AS (SELECT session_id, BOOL_OR(event_name = 'login_success') AS logged_in, BOOL_OR(success = true AND category IN ('device', 'configuration', 'command', 'firmware')) AS activated FROM product_analytics_events GROUP BY session_id) SELECT source, COUNT(DISTINCT sources.session_id) AS sessions, COUNT(DISTINCT sources.session_id) FILTER (WHERE logged_in) AS logins, COUNT(DISTINCT sources.session_id) FILTER (WHERE activated) AS activated, ROUND(100.0 * COUNT(DISTINCT sources.session_id) FILTER (WHERE activated) / NULLIF(COUNT(DISTINCT sources.session_id), 0), 1) AS conversion_percent FROM sources LEFT JOIN outcomes USING (session_id) GROUP BY source ORDER BY sessions DESC LIMIT 20`,
		},
		{
			ID: 701, Section: "performance", Type: "stat", Title: "Средняя загрузка, мс",
			Query: `SELECT COALESCE(ROUND(AVG(duration_ms)), 0)::numeric AS value FROM product_analytics_events WHERE event_name = 'page_performance' AND user_uuid = user_uuid`,
		},
		{
			ID: 702, Section: "performance", Type: "report_table", Title: "Медленные операции API",
			Query: `SELECT properties->>'api_path' AS api_path, properties->>'method' AS method, ROUND(AVG(duration_ms)) AS avg_ms, MAX(duration_ms) AS max_ms, COUNT(*) AS requests FROM product_analytics_events WHERE duration_ms IS NOT NULL AND (category = 'api' OR event_name LIKE '%_success') GROUP BY api_path, method ORDER BY avg_ms DESC NULLS LAST LIMIT 20`,
		},
		{
			ID: 703, Section: "performance", Type: "report_line", Title: "Скорость и ошибки API",
			Query: `SELECT DATE(occurred_at) AS date, ROUND(AVG(duration_ms)) AS avg_ms, COUNT(*) FILTER (WHERE success = false) AS errors FROM product_analytics_events WHERE duration_ms IS NOT NULL AND user_uuid = user_uuid GROUP BY DATE(occurred_at) ORDER BY DATE(occurred_at)`,
		},
		{
			ID: 801, Section: "devices", Type: "report_table", Title: "Команды и прошивки",
			Query: `SELECT category, event_name, COUNT(*) AS operations, COUNT(*) FILTER (WHERE success = true) AS successful, COUNT(*) FILTER (WHERE success = false) AS failed, ROUND(AVG(duration_ms)) AS avg_ms FROM product_analytics_events WHERE category IN ('command', 'firmware', 'device') AND user_uuid = user_uuid GROUP BY category, event_name ORDER BY operations DESC LIMIT 25`,
		},
		{
			ID: 802, Section: "devices", Type: "report_table", Title: "Результаты выполнения команд",
			Query: `SELECT device_command_executions.status, COUNT(*) AS executions, COUNT(DISTINCT device_command_executions.device_imei) AS devices, ROUND(AVG(EXTRACT(EPOCH FROM (device_command_executions.updated_at - device_command_executions.created_at)))) AS avg_seconds FROM device_command_executions JOIN device_commands ON device_commands.id = device_command_executions.command_id JOIN devices ON devices.imei = device_command_executions.device_imei WHERE device_command_executions.created_at BETWEEN {{from}}::timestamptz AND {{to}}::timestamptz AND ({{user_uuid}} IS NULL OR device_commands.user_uuid = {{user_uuid}}) AND ({{group_id}} IS NULL OR devices.id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}})) GROUP BY device_command_executions.status ORDER BY executions DESC`,
		},
		{
			ID: 803, Section: "devices", Type: "report_table", Title: "Состояние обновления прошивок",
			Query: `SELECT syncs.firmware_update_status AS status, COUNT(*) AS devices, COUNT(*) FILTER (WHERE devices.status = true) AS online FROM syncs JOIN devices ON devices.id = syncs.device_id WHERE ({{user_uuid}} IS NULL OR devices.user_uuid = {{user_uuid}}) AND ({{group_id}} IS NULL OR devices.id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}})) GROUP BY syncs.firmware_update_status ORDER BY devices DESC`,
		},
		{
			ID: 804, Section: "devices", Type: "report_table", Title: "Подтверждение конфигураций",
			Query: `SELECT configurations.cfg_sync_status AS status, COUNT(*) AS devices, COUNT(*) FILTER (WHERE devices.status = true) AS online, COUNT(*) FILTER (WHERE configurations.cfg_sync_error IS NOT NULL) AS with_error FROM configurations JOIN devices ON devices.id = configurations.device_id WHERE ({{user_uuid}} IS NULL OR devices.user_uuid = {{user_uuid}}) AND ({{group_id}} IS NULL OR devices.id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}})) GROUP BY configurations.cfg_sync_status ORDER BY devices DESC`,
		},
	}

	titleKeys := map[int]string{
		101: "analytics-report-sessions", 102: "analytics-report-useful-actions", 103: "analytics-report-users-with-errors", 104: "analytics-report-service-activity",
		201: "analytics-report-configuration-funnel", 202: "analytics-report-device-funnel", 203: "analytics-report-time-to-value",
		301: "analytics-report-page-views", 302: "analytics-report-average-page-time", 303: "analytics-report-popular-pages", 304: "analytics-report-controls", 305: "analytics-report-search", 306: "analytics-report-role-activity",
		401: "analytics-report-client-api-errors", 402: "analytics-report-common-errors", 403: "analytics-report-validation-fields", 404: "analytics-report-error-recovery",
		501: "analytics-report-retention", 502: "analytics-report-active-users",
		601: "analytics-report-sources", 602: "analytics-report-landing-pages", 603: "analytics-report-source-conversion",
		701: "analytics-report-average-load", 702: "analytics-report-slow-api", 703: "analytics-report-api-speed-errors",
		801: "analytics-report-device-actions", 802: "analytics-report-command-results", 803: "analytics-report-firmware-status", 804: "analytics-report-configuration-confirmation",
	}
	for index := range widgets {
		widgets[index].TitleKey = titleKeys[widgets[index].ID]
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, widgets)
	return nil
}
