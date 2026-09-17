package cfgmerge

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/pkg/adm"
)

var ErrDeviceConfigurationMissing = errors.New("device configuration is missing")

type TemplateChangesData struct {
	CfgHash uint32         `json:"cfg_hash"`
	Changes map[string]any `json:"changes"`
}

func ParseConfiguration(tr locale.Translator, configuration []byte) ([]adm.FieldConfigurationParsed, error) {
	admInst := adm.Adm{}

	if len(configuration) < 4 {
		return nil, nil
	}

	size := binary.LittleEndian.Uint16(configuration[0:2])
	if int(size) != len(configuration) {
		return nil, errors.New(tr.TErr("cfg-parse-error-device"))
	}

	return admInst.ParseBinaryToJson(tr, configuration)
}

func ConvertDraftToParsedWithSize(draft *redis.ConfigurationDraftInsert, cfgParsed []adm.FieldConfigurationParsed) []adm.PacketFieldConfiguration {
	fullConfig := make([]adm.PacketFieldConfiguration, 0, len(cfgParsed))

	changes := make(map[string]interface{})
	for k, v := range draft.Changes {
		changes[k] = v
	}

	for _, field := range cfgParsed {
		if field.Unknown {
			fullConfig = append(fullConfig, adm.PacketFieldConfiguration{
				UID:      field.RawUID,
				RawUID:   field.RawUID,
				Size:     field.Size,
				RawValue: field.RawValue,
				Unknown:  true,
			})
			continue
		}

		uidStr := field.UID
		var value interface{}

		if newVal, ok := changes[uidStr]; ok {
			value = newVal
		} else {
			value = field.Value
		}

		var uid uint16
		found := false
		for _, s := range constants.CfgSchema {
			if s.Name == uidStr {
				uid = s.UID
				found = true
				break
			}
		}

		if !found {
			continue
		}

		fullConfig = append(fullConfig, adm.PacketFieldConfiguration{
			UID:   uid,
			Size:  field.Size,
			Value: value,
		})
	}

	return fullConfig
}

func BuildBinaryFromDraft(tr locale.Translator, draft *redis.ConfigurationDraftInsert, cfgParsed []adm.FieldConfigurationParsed) ([]byte, error) {
	cfgConverted := ConvertDraftToParsedWithSize(draft, cfgParsed)
	admInst := adm.Adm{}

	return admInst.ParseJsonToBinary(tr, cfgConverted)
}

func IsTemplateChangesData(sourceData []byte) bool {
	var changesData TemplateChangesData
	if err := json.Unmarshal(sourceData, &changesData); err != nil {
		return false
	}

	return len(changesData.Changes) > 0
}

func BuildTaskCfgData(tr locale.Translator, configurationSource string, sourceData []byte, deviceConfiguration *models.Configuration) ([]byte, error) {
	switch configurationSource {
	case "device_sources", "file_sources":
		if len(sourceData) == 0 {
			return nil, errors.New("configuration source data is empty")
		}

		return sourceData, nil
	case "template_sources":
		if len(sourceData) == 0 {
			return nil, errors.New("configuration source data is empty")
		}

		if !IsTemplateChangesData(sourceData) {
			return sourceData, nil
		}

		if deviceConfiguration == nil || deviceConfiguration.CfgData == nil {
			return nil, ErrDeviceConfigurationMissing
		}

		var changesData TemplateChangesData
		if err := json.Unmarshal(sourceData, &changesData); err != nil {
			return nil, err
		}

		cfgParsed, err := ParseConfiguration(tr, deviceConfiguration.CfgData)
		if err != nil {
			return nil, err
		}

		if cfgParsed == nil {
			return nil, ErrDeviceConfigurationMissing
		}

		draft := &redis.ConfigurationDraftInsert{
			CfgHash: changesData.CfgHash,
			Changes: changesData.Changes,
		}

		if deviceConfiguration.CfgHash != 0 {
			draft.CfgHash = deviceConfiguration.CfgHash
		}

		return BuildBinaryFromDraft(tr, draft, cfgParsed)
	default:
		return nil, errors.New("unsupported configuration source")
	}
}
