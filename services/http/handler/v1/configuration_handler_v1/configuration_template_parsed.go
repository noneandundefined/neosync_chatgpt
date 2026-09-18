package configuration_handler_v1

import (
	"encoding/json"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/adm/admparser"
	"time"
)

func applyChangesToSectionParsed(cfgParsed []adm.FieldConfigurationParsed, changes map[string]any) []adm.FieldConfigurationParsed {
	if len(changes) == 0 {
		return cfgParsed
	}

	result := make([]adm.FieldConfigurationParsed, len(cfgParsed))
	copy(result, cfgParsed)

	for i, field := range result {
		if newVal, ok := changes[field.UID]; ok {
			result[i].Value = newVal
		}
	}

	return result
}

func buildTemplateSectionParsedFromStored(tr locale.Translator, template *models.ConfigurationTemplate, cfgData []byte, section string) (ConfigurationTemplateSectionParsed, error) {
	switch template.TypeOfSaving {
	case "modified_ones":
		var changesData ConfigurationTemplateChangesData
		if err := json.Unmarshal(cfgData, &changesData); err != nil {
			return ConfigurationTemplateSectionParsed{}, err
		}

		defaultCfg, schemaList := getDefaultConfigurationBySection(section)
		mergedCfg := applyChangesToSectionParsed(defaultCfg, changesData.Changes)

		return ConfigurationTemplateSectionParsed{
			CfgHash:      changesData.CfgHash,
			Section:      section,
			ConfigParsed: mergedCfg,
			Schema:       schemaList,
			Timestamp:    time.Now().Unix(),
		}, nil
	case "full":
		cfgParsed, err := getConfiguration(tr, cfgData)
		if err != nil {
			return ConfigurationTemplateSectionParsed{}, err
		}

		filteredCfg, schemaList := getConfigurationBySection(section, cfgParsed)

		return ConfigurationTemplateSectionParsed{
			CfgHash:      admparser.GetCfgHash(cfgData),
			Section:      section,
			ConfigParsed: filteredCfg,
			Schema:       schemaList,
			Timestamp:    time.Now().Unix(),
		}, nil
	default:
		filteredCfg, schemaList := getDefaultConfigurationBySection(section)

		return ConfigurationTemplateSectionParsed{
			CfgHash:      0,
			Section:      section,
			ConfigParsed: filteredCfg,
			Schema:       schemaList,
			Timestamp:    time.Now().Unix(),
		}, nil
	}
}
