package constants

import "sort"

func ConfigProfColumnName(uidName string) string {
	if uidName == "device_id" {
		return "cfg_device_id"
	}

	return uidName
}

func ConfigProfSQLType(fieldType FieldType) string {
	switch fieldType {
	case UInt8:
		return "SMALLINT"
	case UInt16:
		return "INTEGER"
	case UInt32:
		return "BIGINT"
	case Ascii, Hex:
		return "TEXT"
	default:
		return "JSONB"
	}
}

func ConfigProfSchemas() []FieldSchema {
	schemas := make([]FieldSchema, 0, len(CfgSchema))
	for _, schema := range CfgSchema {
		schemas = append(schemas, schema)
	}

	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].UID < schemas[j].UID
	})

	return schemas
}
