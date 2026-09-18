package types

import (
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"sync"
)

type InFlightCommand struct {
	Cmd      RabbitMQ_TransitBinary
	AnswerCh chan string
}

type TCPSession struct {
	Mu              sync.Mutex
	ConfigurationMu sync.Mutex

	DeviceID  uint64
	RequestId *string
	Device    *adm.ADMDevice

	QueueCommand   []RabbitMQ_TransitBinary
	QueueTelemetry []RabbitMQ_TransitBinary

	Busy     bool
	InFlight *InFlightCommand
}

// SIZE(uint16_t) | TYPE(uint8_t) | REQID(16 bytes) | IMEI(8 bytes) | STATCODE(uint16_t) | DATA(DATA_LEN)
type RabbitMQ_TransitBinary struct {
	Type        uint8
	RequestID   string
	Imei        string
	NResp       bool
	Telemetry   bool
	StatCode    uint16
	Data        []byte
	ExecutionID uint64
}
