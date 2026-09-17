package adm

import "neomatica/neosync/infra/constants"

func ParsedToPacketFields(cfg []FieldConfigurationParsed) []PacketFieldConfiguration {
	out := make([]PacketFieldConfiguration, 0, len(cfg))

	for _, field := range cfg {
		if field.Unknown {
			out = append(out, PacketFieldConfiguration{
				UID:      field.RawUID,
				RawUID:   field.RawUID,
				Size:     field.Size,
				RawValue: field.RawValue,
				Unknown:  true,
			})
			continue
		}

		var uid uint16
		found := false
		for _, s := range constants.CfgSchema {
			if s.Name == field.UID {
				uid = s.UID
				found = true
				break
			}
		}
		if !found {
			continue
		}

		out = append(out, PacketFieldConfiguration{
			UID:    uid,
			RawUID: field.RawUID,
			Size:   field.Size,
			Value:  field.Value,
		})
	}

	return out
}
