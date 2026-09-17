package constants

import (
	"strings"
	"testing"
)

func TestPrepareProductAnalyticQuery(t *testing.T) {
	groupID := 42
	query := PrepareAnalyticQuery(
		"report_table",
		"SELECT event_name, COUNT(*) FROM product_analytics_events WHERE success = false GROUP BY event_name",
		"2026-09-01T00:00:00Z",
		"2026-09-05T23:59:59Z",
		"user-uuid",
		"ADM007",
		&groupID,
	)

	for _, placeholder := range []string{"{{from}}", "{{to}}", "{{user_uuid}}", "{{model}}", "{{group_id}}"} {
		if strings.Contains(query, placeholder) {
			t.Fatalf("query contains unresolved placeholder %s", placeholder)
		}
	}

	for _, expected := range []string{"occurred_at BETWEEN", "user-uuid", "ADM007", "group_id = 42"} {
		if !strings.Contains(query, expected) {
			t.Fatalf("query does not contain filter %s", expected)
		}
	}
}

func TestPrepareStatQueryKeepsLegacyPublishedFlag(t *testing.T) {
	query := PrepareAnalyticQuery(
		"stat",
		"SELECT COUNT(*)::numeric AS value FROM analytic_records WHERE rmq_published = true AND user_uuid = user_uuid AND device_id = device_id",
		"2026-09-01T00:00:00Z",
		"2026-09-05T23:59:59Z",
		"",
		"",
		nil,
	)

	if !strings.Contains(query, "rmq_published = true") {
		t.Fatal("legacy rmq_published success condition changed")
	}
}
