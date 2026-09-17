package constants

import "time"

const (
	MAX_RETRIES_RABBITMQ     = 3
	MAX_RETRY_DELAY_RABBITMQ = 2 * time.Second
)
