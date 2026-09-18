package configuration_handler_v1

import (
	"testing"

	"neomatica/neosync/pkg/adm"
)

func TestConfigurationHistoryDiff(t *testing.T) {
	older := map[uint16]adm.FieldConfigurationParsed{
		60000: {UID: "unknown_60000", RawUID: 60000, Value: 10},
		60001: {UID: "unknown_60001", RawUID: 60001, Value: "unchanged"},
		60002: {UID: "unknown_60002", RawUID: 60002, Value: "removed"},
	}
	newer := map[uint16]adm.FieldConfigurationParsed{
		60000: {UID: "unknown_60000", RawUID: 60000, Value: 20},
		60001: {UID: "unknown_60001", RawUID: 60001, Value: "unchanged"},
		60003: {UID: "unknown_60003", RawUID: 60003, Value: "added"},
	}

	changes := configurationHistoryDiff(older, newer)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d: %#v", len(changes), changes)
	}

	if changes[0].UID != 60000 || changes[0].OldValue != 10 || changes[0].NewValue != 20 {
		t.Fatalf("unexpected changed field: %#v", changes[0])
	}

	if changes[1].UID != 60002 || changes[1].OldValue != "removed" || changes[1].NewValue != nil {
		t.Fatalf("unexpected removed field: %#v", changes[1])
	}

	if changes[2].UID != 60003 || changes[2].OldValue != nil || changes[2].NewValue != "added" {
		t.Fatalf("unexpected added field: %#v", changes[2])
	}
}
