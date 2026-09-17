package config

import "time"

var DeviceRemotePort int = 12346

var ReadDeadline time.Duration = 20 * time.Second
var WriteDeadline time.Duration = 20 * time.Second
var QueueBrockerToTcp string = "brocker_command_to_tcp"
var QueueBrockerToHttp string = "brocker_command_to_http"

/* Bulk configuration delivery */
var CompanyTaskConfirmationTimeout time.Duration = 2 * time.Minute
var CompanyTaskRetryDelay time.Duration = 30 * time.Second
var CompanyTaskLeaseTimeout time.Duration = time.Minute
var CompanyTaskPollInterval time.Duration = 10 * time.Second
var CompanyTaskMaxAttempts int = 3

var CompanyTaskWorkers int = 8
