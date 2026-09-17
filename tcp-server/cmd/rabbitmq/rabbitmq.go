package rabbitmq

import (
	"fmt"
	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	POOL_SIZE     = 50
	PREFETCH_SIZE = 150
	WORKERS       = 100
)

func NewRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	r := &RabbitMQ{
		amqpURL: amqpURL,
	}

	var conn *amqp.Connection
	var connErr error

	/* Base config connection RabbitMQ */
	for i := 0; i < 5; i++ {
		conn, connErr = amqp.DialConfig(amqpURL, amqp.Config{
			Heartbeat: 10 * time.Second,
		})
		if connErr == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if connErr != nil {
		return nil, connErr
	}
	r.conn = conn

	/* Initial channels */
	if err := r.initChannels(); err != nil {
		return nil, err
	}

	r.mutex.Lock()
	r.isReady = true
	r.mutex.Unlock()

	/* args := amqp.Table{"x-queue-mode": "lazy"} */

	fmt.Printf("[%v] [INFO] Successfully connected to RabbitMQ\n", time.Now().Format("2006-01-02 15:04:05"))

	go r.handleReconnect()

	return r, nil
}

/* Reconnect RabbitMQ connection */
func (r *RabbitMQ) handleReconnect() {
	for {
		notifyClose := make(chan *amqp.Error, 1)

		r.mutex.Lock()
		conn := r.conn
		r.mutex.Unlock()

		if conn == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		conn.NotifyClose(notifyClose)

		err := <-notifyClose
		if err != nil {
			logger.Error("handleReconnect: RabbitMQ connection closed: %v", err)
		}

		r.mutex.Lock()
		r.isReady = false
		r.mutex.Unlock()

		logger.Error("handleReconnect: RabbitMQ disconnected. Trying to reconnect.")

		for {
			newConn, dialErr := amqp.DialConfig(r.amqpURL, amqp.Config{
				Heartbeat: 10 * time.Second,
			})
			if dialErr != nil {
				logger.Error("handleReconnect: dial failed: %v", dialErr)

				time.Sleep(5 * time.Second)
				continue
			}

			r.mutex.Lock()
			r.conn = newConn

			if initErr := r.initChannels(); initErr != nil {
				r.conn = nil
				r.mutex.Unlock()

				logger.Error("handleReconnect: initChannels failed: %v", initErr)

				_ = newConn.Close()

				time.Sleep(5 * time.Second)
				continue
			}

			if r.consumeHandler != nil {
				if startErr := r.startConsumers(); startErr != nil {
					r.conn = nil
					r.mutex.Unlock()
					logger.Error("handleReconnect: startConsumers failed: %v", startErr)
					_ = newConn.Close()
					time.Sleep(5 * time.Second)
					continue
				}
			}

			r.isReady = true
			r.mutex.Unlock()

			logger.Info("handleReconnect: RabbitMQ successfully reconnected")
			break
		}
	}
}

func (r *RabbitMQ) initChannels() error {
	/* Pool config */
	poolChannel := make(chan *amqp.Channel, POOL_SIZE)

	/* Pool initial */
	for i := 0; i < POOL_SIZE; i++ {
		ch, err := r.conn.Channel()
		if err != nil {
			return err
		}
		poolChannel <- ch
	}

	/* Consume initial */
	consumeChannel, err := r.conn.Channel()
	if err != nil {
		return err
	}

	args := amqp.Table{
		"x-message-ttl": int32(86400000), // 24 часа
	}

	qTcp, err := consumeChannel.QueueDeclare(config.QueueBrockerToTcp, true, false, false, false, args)
	if err != nil {
		return err
	}
	r.qTcp = qTcp

	qHttp, err := consumeChannel.QueueDeclare(config.QueueBrockerToHttp, true, false, false, false, args)
	if err != nil {
		return err
	}
	r.qHttp = qHttp

	r.poolChannel = poolChannel
	r.consumeChannel = consumeChannel

	return nil
}

func (r *RabbitMQ) Consume(handler func(body []byte)) error {
	r.consumeHandler = handler
	return r.startConsumers()
}

func (r *RabbitMQ) startConsumers() error {
	if err := r.consumeChannel.Qos(PREFETCH_SIZE, 0, false); err != nil {
		return err
	}

	msgs, err := r.consumeChannel.Consume(
		r.qTcp.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < WORKERS; i++ {
		go func() {
			for msg := range msgs {
				r.consumeHandler(msg.Body)
				msg.Ack(false)
			}
		}()
	}

	return nil
}

func (r *RabbitMQ) Close() {
	close(r.poolChannel)

	for ch := range r.poolChannel {
		ch.Close()
	}

	if r.consumeChannel != nil {
		_ = r.consumeChannel.Close()
	}

	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func (r *RabbitMQ) Publish(queue string, buffer []byte) error {
	r.mutex.Lock()
	ready := r.isReady
	r.mutex.Unlock()

	if !ready {
		return fmt.Errorf("Publish: RabbitMQ not ready")
	}

	ch := <-r.poolChannel
	defer func() { r.poolChannel <- ch }()

	err := ch.Publish(
		"",
		queue,
		false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/octet-stream",
			Body:         buffer,
			Timestamp:    time.Now(),
		},
	)

	/* Recreate channel */
	if err != nil {
		channel, channelErr := r.conn.Channel()
		if channelErr == nil {
			ch = channel
		}

		return err
	}

	return nil
}
