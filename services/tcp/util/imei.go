package util

import "strings"

func PrepareImei(imei string) string {
	if imei == "" {
		return ""
	}

	imei = strings.TrimSpace(imei)
	imei = strings.TrimSuffix(strings.TrimSuffix(imei, "f"), "F")

	if len(imei) == 16 && imei[0] == '0' {
		imei = imei[1:]
	}

	return imei
}
