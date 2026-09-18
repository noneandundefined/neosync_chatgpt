package util

import (
	"fmt"
	"strings"
	"time"
)

func ConvertNumberToPrettyHexString(number int64) string {
	if number < 0 {
		number = 0xffffffff + number + 1
	}

	return strings.ToUpper(fmt.Sprintf("%X", number))
}

func ConvertUnixTimestampToPrettyFormat(unixTimestamp int64) string {
	date := time.Unix(unixTimestamp, 0).UTC()

	day := date.Day()
	month := date.Month()
	year := date.Year()

	hours := date.Hour()
	minutes := date.Minute()
	seconds := date.Second()

	return fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d", day, month, year, hours, minutes, seconds)
}
