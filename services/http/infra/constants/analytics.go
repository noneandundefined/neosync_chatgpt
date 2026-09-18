package constants

import (
	"fmt"
	"regexp"
	"strings"
)

const SQL_Template = `
	WITH bounds AS (
		SELECT
			{{from}}::timestamptz AS from_ts,
			{{to}}::timestamptz AS to_ts
	),
	group_devices AS (
		SELECT id AS device_id
		FROM devices
		WHERE ({{group_id}} IS NULL OR id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}}))
			AND ({{model}} IS NULL OR device_model = {{model}})
	),
	current_period AS (
		{{base_query}}
			AND created_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)
			{{filters}}
	),
	previous_period AS (
		{{base_query}}
			AND created_at >= (
				(SELECT from_ts FROM bounds) - ((SELECT to_ts FROM bounds) - (SELECT from_ts FROM bounds))
			)
			AND created_at < (SELECT from_ts FROM bounds)
			{{filters}}
	)
	SELECT
		current_period.value::BIGINT AS value,
		CASE
			WHEN previous_period.value = 0 AND current_period.value = 0 THEN 0
			WHEN previous_period.value = 0 THEN 100
			ELSE ROUND(((current_period.value - previous_period.value) * 100.0 / previous_period.value), 2)
		END AS change_percent
	FROM current_period, previous_period
`

var clauseBreakPattern = regexp.MustCompile(`(?i)\s+(GROUP\s+BY|ORDER\s+BY|HAVING|LIMIT|OFFSET|UNION|INTERSECT|EXCEPT)\s+`)
var filterWherePattern = regexp.MustCompile(`(?i)\bFILTER\s*\(\s*WHERE\b`)
var topLevelWherePattern = regexp.MustCompile(`(?i)\bWHERE\b`)

// BuildAnalyticQuery wraps stat widget queries with period comparison. Use only for type=stat.
func BuildAnalyticQuery(sql string) string {
	sql = normalizeBaseQuery(sql)
	sql = injectBeforeClauses(sql, buildStatFilters(sql))

	q := SQL_Template
	q = strings.ReplaceAll(q, "{{base_query}}", sql)
	q = strings.ReplaceAll(q, "{{filters}}", "")

	return q
}

func normalizeBaseQuery(q string) string {
	q = strings.TrimSpace(q)

	if hasTopLevelWhere(q) {
		return q
	}

	return injectBeforeClauses(q, "WHERE 1=1")
}

func hasTopLevelWhere(q string) bool {
	stripped := filterWherePattern.ReplaceAllString(q, "FILTER (")
	upper := strings.ToUpper(stripped)

	fromIdx := strings.LastIndex(upper, "FROM")
	if fromIdx < 0 {
		return topLevelWherePattern.MatchString(upper)
	}

	segment := upper[fromIdx:]
	if loc := clauseBreakPattern.FindStringIndex(segment); loc != nil {
		segment = segment[:loc[0]]
	}

	return topLevelWherePattern.MatchString(segment)
}

func buildStatFilters(baseQuery string) string {
	upper := strings.ToUpper(baseQuery)
	filters := ""

	if strings.Contains(upper, "USER_UUID") {
		filters += "\n\t\tAND ({{user_uuid}} IS NULL OR user_uuid = {{user_uuid}})"
	}

	if strings.Contains(upper, "DEVICE_ID") {
		filters += "\n\t\tAND (({{group_id}} IS NULL AND {{model}} IS NULL) OR device_id IN (SELECT device_id FROM group_devices))"
	}

	return filters
}

func injectBeforeClauses(sql string, injection string) string {
	injection = strings.TrimSpace(injection)
	if injection == "" {
		return sql
	}

	loc := clauseBreakPattern.FindStringIndex(sql)
	if loc == nil {
		return strings.TrimRight(sql, " \t\n;") + "\n\t\t" + injection
	}

	return sql[:loc[0]] + "\n\t\t" + injection + sql[loc[0]:]
}

func InterpolateQuery(template string, from, to, userUUID, model string, groupID *int) string {
	query := template

	query = strings.ReplaceAll(query, "{{from}}", quoteOrNull(from))
	query = strings.ReplaceAll(query, "{{to}}", quoteOrNull(to))
	query = strings.ReplaceAll(query, "{{user_uuid}}", quoteOrNull(userUUID))
	query = strings.ReplaceAll(query, "{{group_id}}", intOrNull(groupID))
	query = strings.ReplaceAll(query, "{{model}}", quoteOrNull(model))

	return query
}

func quoteOrNull(value string) string {
	if strings.TrimSpace(value) == "" {
		return "NULL"
	}

	escaped := strings.ReplaceAll(value, "'", "''")
	return fmt.Sprintf("'%s'", escaped)
}
func intOrNull(value *int) string {
	if value == nil {
		return "NULL"
	}

	return fmt.Sprintf("%d", *value)
}

// PrepareAnalyticQuery applies widget-specific wrapping and filter interpolation.
func PrepareAnalyticQuery(queryType, query, from, to, userUUID, model string, groupID *int) string {
	q := strings.TrimSpace(query)

	switch queryType {
	case "stat":
		q = BuildAnalyticQuery(q)
	case "stat_live", "sql":
		// Ad-hoc query: execute as written (placeholders still interpolated).
	case "bar", "table", "gauge_top_reasons":
		// SQL already contains bounds/group filter CTEs.
	case "report_line", "report_table", "funnel":
		q = BuildProductAnalyticQuery(q)
	default:
		q = BuildFilteredAnalyticQuery(q)
	}

	return InterpolateQuery(q, from, to, userUUID, model, groupID)
}

func BuildProductAnalyticQuery(sql string) string {
	filter := `occurred_at BETWEEN {{from}}::timestamptz AND {{to}}::timestamptz
			AND ({{user_uuid}} IS NULL OR user_uuid = {{user_uuid}})
			AND (({{group_id}} IS NULL AND {{model}} IS NULL) OR device_id IS NULL OR device_id IN (SELECT devices.id FROM devices WHERE ({{group_id}} IS NULL OR devices.id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}})) AND ({{model}} IS NULL OR devices.device_model = {{model}})))`

	withWhere := "FROM product_analytics_events WHERE " + filter + " AND "
	withoutWhere := "FROM product_analytics_events WHERE " + filter + " "

	q := strings.ReplaceAll(sql, "FROM product_analytics_events WHERE ", withWhere)
	q = strings.ReplaceAll(q, "FROM product_analytics_events GROUP ", withoutWhere+"GROUP ")
	q = strings.ReplaceAll(q, "FROM product_analytics_events ORDER ", withoutWhere+"ORDER ")
	q = strings.ReplaceAll(q, "FROM product_analytics_events)", withoutWhere+")")

	return q
}

const SQL_FilteredTemplate = `
WITH bounds AS (
	SELECT
		{{from}}::timestamptz AS from_ts,
		{{to}}::timestamptz AS to_ts
),
group_devices AS (
	SELECT id AS device_id
	FROM devices
	WHERE ({{group_id}} IS NULL OR id IN (SELECT device_id FROM group_user_devices WHERE group_id = {{group_id}}))
		AND ({{model}} IS NULL OR device_model = {{model}})
),
filtered AS (
	{{base_query}}
)
SELECT * FROM filtered
`

func BuildFilteredAnalyticQuery(sql string) string {
	sql = normalizeBaseQuery(sql)
	if combined := buildCombinedFilters(sql); combined != "" {
		sql = injectBeforeClauses(sql, combined)
	}

	q := SQL_FilteredTemplate
	q = strings.ReplaceAll(q, "{{base_query}}", sql)

	return q
}

func buildCombinedFilters(baseQuery string) string {
	return strings.TrimSpace(buildDateFilter(baseQuery) + buildRecordFilters(baseQuery))
}

func buildDateFilter(baseQuery string) string {
	upper := strings.ToUpper(baseQuery)

	if strings.Contains(upper, "CREATED_AT") || strings.Contains(upper, "ANALYTIC_RECORDS") {
		return "\n\t\tAND created_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)"
	}

	if strings.Contains(upper, "UPDATED_AT") {
		return "\n\t\tAND updated_at BETWEEN (SELECT from_ts FROM bounds) AND (SELECT to_ts FROM bounds)"
	}

	return ""
}

func buildRecordFilters(baseQuery string) string {
	upper := strings.ToUpper(baseQuery)
	filters := ""

	if strings.Contains(upper, "USER_UUID") {
		filters += "\n\t\tAND ({{user_uuid}} IS NULL OR user_uuid = {{user_uuid}})"
	}

	if strings.Contains(upper, "DEVICE_ID") {
		filters += "\n\t\tAND (({{group_id}} IS NULL AND {{model}} IS NULL) OR device_id IN (SELECT device_id FROM group_devices))"
	}

	return filters
}
