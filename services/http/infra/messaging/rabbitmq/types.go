package rabbitmq

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn           *amqp.Connection
	poolChannel    chan *amqp.Channel
	consumeChannel *amqp.Channel
	qTcp           amqp.Queue
	qHttp          amqp.Queue

	/* Handler for consume */
	consumeHandler func([]byte)

	amqpURL string
	isReady bool
	mutex   sync.Mutex
}
