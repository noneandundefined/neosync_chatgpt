package handler

import (
	"database/sql"
	"neomatica/neosync/infra/analytics"
	"neomatica/neosync/infra/store/memory"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/infra/store/postgres/usecase"

	"neomatica/neosync/cmd/rabbitmq"
)

type BaseHandler struct {
	Db        *sql.DB
	Store     store.Storage
	UseCase   usecase.UseCase
	RMQ       *rabbitmq.RabbitMQ
	Session   *memory.SessionService
	Analytics *analytics.Collector
}
