package constants

import "time"

const RateLimiter_StateTTL time.Duration = 1 * time.Minute
const Redis_DeviceAuthTTL time.Duration = 10 * time.Minute
const Redis_CommandTTL time.Duration = 4 * time.Minute
const Redis_TelemetryTTL time.Duration = 2 * time.Minute
const Redis_TelemetryRefreshCooldown time.Duration = 15 * time.Second
const Redis_DeviceUpdTTL time.Duration = 8 * time.Minute
const Redis_DraftTTL time.Duration = 30 * time.Minute
const Redis_ConfigurationTTL time.Duration = 3 * time.Minute
const Redis_ConfigurationTmpTTL time.Duration = 30 * time.Second
const RabbitMQ_TimeoutRead time.Duration = 18 * time.Second
