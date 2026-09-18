package rabbitmq

import (
	"context"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/memory"
	"neomatica/neosync/types"

	"neomatica/neosync/config"

	"github.com/google/uuid"
)

func (r *RabbitMQ) SendToRabbitAndWait(ctx context.Context, rabbitmqTransit types.RabbitMQ_TransitBinary, session *memory.SessionService) (*types.RabbitMQ_TransitBinary, error) {
	requestId := uuid.New().String()

	/* rabbitmqTransit */
	rabbitmqTransit.RequestID = requestId
	rabbitmqTransit.StatCode = 100

	buffer, err := r.Encoder(rabbitmqTransit)
	if err != nil {
		logger.Error("SendToRabbitAndWait req={%s} rabbitmqReqId={%s}: Failed encoder buffer for send to queue: %s", ctx.Value("XREQID").(string), rabbitmqTransit.RequestID, err.Error())
		return nil, err
	}

	respChan := make(chan []byte, 1)
	session.Register(rabbitmqTransit.RequestID, respChan)
	defer session.Unregister(rabbitmqTransit.RequestID)

	if err := r.Publish(config.QueueBrockerToTcp, buffer); err != nil {
		logger.Error("SendToRabbitAndWait req={%s} rabbitmqReqId={%s}: Failed to publish: %s", ctx.Value("XREQID").(string), rabbitmqTransit.RequestID, err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.RabbitMQ_TimeoutRead)
	defer cancel()

	select {
	case respBuffer := <-respChan:
		resp, err := r.Decoder(respBuffer)
		if err != nil {
			logger.Error("SendToRabbitAndWait req={%s} rabbitmqReqId={%s}: Failed decoder resp: %s", ctx.Value("XREQID").(string), rabbitmqTransit.RequestID, err.Error())
			return nil, err
		}

		return &resp, nil
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ctx.Err()
		}

		return nil, ctx.Err()
	}
}

func (r *RabbitMQ) SendToRabbitAsync(ctx context.Context, rabbitmqTransit types.RabbitMQ_TransitBinary) error {
	requestId := uuid.New().String()

	rabbitmqTransit.RequestID = requestId
	rabbitmqTransit.StatCode = 100

	buffer, err := r.Encoder(rabbitmqTransit)
	if err != nil {
		logger.Error("SendToRabbitAsync req={%s} rabbitmqReqId={%s}: Failed encoder buffer for send to queue: %s", ctx.Value("XREQID").(string), rabbitmqTransit.RequestID, err.Error())
		return err
	}

	if err := r.Publish(config.QueueBrockerToTcp, buffer); err != nil {
		logger.Error("SendToRabbitAsync req={%s} rabbitmqReqId={%s}: Failed publish buffer to queue: %s", ctx.Value("XREQID").(string), rabbitmqTransit.RequestID, err.Error())
		return err
	}

	return nil
}
