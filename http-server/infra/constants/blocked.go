package constants

import "time"

var PenaltyDurations = []time.Duration{
	time.Minute,
	3 * time.Minute,
	5 * time.Minute,
	10 * time.Minute,
	30 * time.Minute,
	60 * time.Minute,
}
