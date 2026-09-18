package adm

import (
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/adm/help"
	"strings"
	"unicode/utf8"
)

type ConfigurationProfValues struct {
	Values        map[string]any
	UnknownFields map[string]any
}

func (adm *Adm) ParseProfValues(cfg []PacketFieldConfiguration) (*ConfigurationProfValues, error) {
	values := make(map[string]any, len(cfg))
	unknownFields := make(map[string]any)

	for _, item := range cfg {
		uid := item.RawUID
		if uid == 0 {
			uid = item.UID
		}

		if item.Unknown {
			unknownFields[fmt.Sprintf("unknown_%d", uid)] = sanitizeUnknownField(item.RawValue)
			continue
		}

		name, err := constants.GetCfgName(uid)
		if err != nil {
			unknownFields[fmt.Sprintf("unknown_%d", uid)] = sanitizeForPostgresJSON(item.Value)
			continue
		}

		colName := constants.ConfigProfColumnName(name)
		schema, err := constants.GetCfgSchema(uid)
		if err != nil {
			unknownFields[colName] = sanitizeForPostgresJSON(item.Value)
			continue
		}

		values[colName] = normalizeProfValue(schema.Type, item.Value)
	}

	return &ConfigurationProfValues{
		Values:        values,
		UnknownFields: unknownFields,
	}, nil
}

func normalizeProfValue(fieldType constants.FieldType, value any) any {
	switch constants.ConfigProfSQLType(fieldType) {
	case "SMALLINT":
		return toInt64(value)
	case "INTEGER", "BIGINT":
		return toInt64(value)
	case "TEXT":
		return toText(value)
	default:
		return toJSONBValue(value)
	}
}

func toInt64(v any) any {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return toJSONBValue(v)
	}
}

func toText(v any) string {
	var s string
	switch t := v.(type) {
	case string:
		s = t
	case []byte:
		s = string(t)
	default:
		b, _ := json.Marshal(sanitizeForPostgresJSON(v))
		s = string(b)
	}
	return sanitizePostgresText(s)
}

func toJSONBValue(v any) any {
	v = sanitizeForPostgresJSON(v)
	switch t := v.(type) {
	case []byte:
		return json.RawMessage(sanitizePostgresText(string(t)))
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		return json.RawMessage(b)
	}
}

func sanitizeUnknownField(v any) any {
	switch t := v.(type) {
	case []byte:
		if len(t) == 0 {
			return ""
		}
		return help.HelpConvertToHex(t)
	default:
		return sanitizeForPostgresJSON(v)
	}
}

func sanitizePostgresText(s string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, "\x00", "")
	if utf8.ValidString(s) {
		return s
	}
	return strings.ToValidUTF8(s, "")
}

func sanitizeForPostgresJSON(v any) any {
	switch t := v.(type) {
	case nil:
		return nil
	case string:
		return sanitizePostgresText(t)
	case []byte:
		return sanitizePostgresText(string(t))
	case bool, float64, float32,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return t
	case json.RawMessage:
		var parsed any
		if err := json.Unmarshal(t, &parsed); err != nil {
			return sanitizePostgresText(string(t))
		}
		return sanitizeForPostgresJSON(parsed)
	case []string:
		out := make([]string, len(t))
		for i, s := range t {
			out[i] = sanitizePostgresText(s)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, elem := range t {
			out[i] = sanitizeForPostgresJSON(elem)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = sanitizeForPostgresJSON(val)
		}
		return out
	default:
		return sanitizeCompositeForPostgresJSON(t)
	}
}

// sanitizeCompositeForPostgresJSON normalizes slices/arrays/maps via a single JSON round-trip.
func sanitizeCompositeForPostgresJSON(v any) any {
	b, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var parsed any
	if err := json.Unmarshal(b, &parsed); err != nil {
		return v
	}
	return sanitizeForPostgresJSON(parsed)
}
