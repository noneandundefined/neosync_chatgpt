package util

import (
	"neomatica/neosync/infra/constants"
	"strings"
)

func ValidDeviceModel(model string) *string {
	model = strings.TrimSpace(model)
	if model == "" {
		return nil
	}

	for _, m := range constants.DEVICE_MODELS {
		if strings.EqualFold(m, model) {
			mod := strings.ToUpper(m)
			return &mod
		}
	}

	return nil
}
