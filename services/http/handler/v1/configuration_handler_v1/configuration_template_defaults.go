package configuration_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/adm"
	"strings"
)

func getDefaultConfigurationBySection(section string) ([]adm.FieldConfigurationParsed, []constants.FieldSchema) {
	var schemaList []constants.FieldSchema
	var filteredCfg []adm.FieldConfigurationParsed

	for _, schema := range constants.CfgSchema {
		if schema.Section == "" {
			continue
		}

		sections := strings.Split(schema.Section, ";")

		for _, s := range sections {
			if !strings.EqualFold(strings.TrimSpace(s), section) {
				continue
			}

			schemaList = append(schemaList, schema)
			filteredCfg = append(filteredCfg, adm.FieldConfigurationParsed{
				UID:    schema.Name,
				RawUID: schema.UID,
				Size:   inferFieldSize(schema),
				Value:  schema.Default,
			})
			break
		}
	}

	return filteredCfg, schemaList
}

func getDefaultConfigurationAll() []adm.FieldConfigurationParsed {
	result := make([]adm.FieldConfigurationParsed, 0, len(constants.CfgSchema))

	for _, schema := range constants.CfgSchema {
		if schema.Section == "" {
			continue
		}

		result = append(result, adm.FieldConfigurationParsed{
			UID:    schema.Name,
			RawUID: schema.UID,
			Size:   inferFieldSize(schema),
			Value:  schema.Default,
		})
	}

	return result
}

func inferFieldSize(schema constants.FieldSchema) uint8 {
	switch schema.Type {
	case constants.UInt8:
		return 1
	case constants.TwoUint8:
		return 2
	case constants.UInt16:
		return 2
	case constants.TwoUint16:
		return 4
	case constants.UInt32:
		return 4
	case constants.Hex:
		if schema.Default != nil {
			if s, ok := schema.Default.(string); ok && len(s) > 0 {
				return uint8(len(s) / 2)
			}
		}

		return 6
	case constants.Ascii:
		if schema.MaxLength != nil {
			return uint8(*schema.MaxLength)
		}

		return 32
	case constants.Phones:
		if schema.MaxItems != nil {
			return uint8(*schema.MaxItems * 16)
		}

		return 64
	case constants.MacAddressArray:
		return 24
	case constants.SNSArray:
		return 32
	case constants.Uint32Array, constants.Uint64Array, constants.ByteArray:
		return 16
	default:
		return 1
	}
}
