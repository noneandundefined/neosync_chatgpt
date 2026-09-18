package rabbitmq

import (
	"context"
	"neomatica/neosync/infra/logger"
	"time"
)

type AmqpOperation func(ctx context.Context) (interface{}, error)

func (r *RabbitMQ) AmqpDo(ctx context.Context, op AmqpOperation, maxRetries int, retryDelay time.Duration, logMsg string) (interface{}, error) {
	var lastErr error
	var resp interface{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, lastErr = op(ctx)
		if lastErr == nil {
			return resp, nil
		}

		if attempt < maxRetries-1 {
			logger.Warning("ARabbitMQ mqpDo: %s (retry %d/%d): %v", logMsg, attempt+1, maxRetries, lastErr)
			time.Sleep(retryDelay)
		}
	}

	return nil, lastErr
}
