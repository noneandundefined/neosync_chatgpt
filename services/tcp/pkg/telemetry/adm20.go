package telemetry

import (
	"fmt"
	"regexp"
	"strconv"
)

var admRegexp = regexp.MustCompile(
	`ADM20INFO:\s*(?P<emptinessMask>\d+);\s*(?P<cardUId>\d+);\s*(?P<tagId>\d+),(?P<tagSn>\d+),(?P<tagRssi>\d+),(?P<tagAdc>\d+),(?P<tagPeriod>\d+)`,
)

func ParseADM20Info(tl *Telemetry, line string) error {
	match := admRegexp.FindStringSubmatch(line)
	if match == nil {
		return fmt.Errorf("error parsing line")
	}

	result := &ADM20Info{}
	for i, name := range admRegexp.SubexpNames() {
		if name == "" {
			continue
		}
		val, err := strconv.Atoi(match[i])
		if err != nil {
			return fmt.Errorf("error convert to atoi format")
		}
		switch name {
		case "emptinessMask":
			result.EmptinessMask = &val
		case "cardUId":
			result.CardUId = &val
		case "tagId":
			result.TagId = &val
		case "tagSn":
			result.TagSn = &val
		case "tagRssi":
			result.TagRssi = &val
		case "tagAdc":
			result.TagAdc = &val
		case "tagPeriod":
			result.TagPeriod = &val
		}
	}

	tl.ADM20INFO = result
	return nil
}
