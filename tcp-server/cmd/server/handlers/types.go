package handlers

import (
	"neomatica/neosync-tcp/cmd/rabbitmq"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/infra/store/postgres/usecase"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/adm"
)

/* Jobs struct */
type JobHelloEvent struct {
	Device  *adm.ADMDevice
	UseCase *usecase.UseCase
	Packet  *protocol.WelcomePacket
}

type JobCommandEvent struct {
	Device *adm.ADMDevice
	Packet string
}

/* Event struct */
type BasePackEventHandler struct {
	Session *memory.SessionMemory
	RMQ     *rabbitmq.RabbitMQ
	Cache   *memory.Cache
	Store   store.Storage

	/* Queues */
	QueueHelloEvent   chan JobHelloEvent
	QueueCommandEvent chan JobCommandEvent

	/* Limiter */
	LimiterDb chan struct{}
}

func (h *BasePackEventHandler) limiterDatabase(fn func()) {
	h.LimiterDb <- struct{}{}
	defer func() { <-h.LimiterDb }()

	fn()
}
