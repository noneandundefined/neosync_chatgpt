package constants

import "time"

const Redis_LogTTL time.Duration = 48 * time.Hour
const Redis_CommandTTL time.Duration = 4 * time.Minute
const Redis_DeviceUpdTTL time.Duration = 8 * time.Minute
const Redis_DraftTTL time.Duration = 30 * time.Minute
const Redis_ConfigurationTTL time.Duration = 3 * time.Minute
const RabbitMQ_TimeoutRead time.Duration = 18 * time.Second

const Tcp_TimeoutRead time.Duration = 20 * time.Second
const Tcp_ReadDeadline time.Duration = 10 * time.Minute
