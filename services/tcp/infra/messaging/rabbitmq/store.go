package rabbitmq

import (
	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/types"
	"neomatica/neosync-tcp/util"
)

func (r *RabbitMQ) SendToRabbitAsync(rabbitmqTransit types.RabbitMQ_TransitBinary) error {
	imeiPrepare := util.PrepareImei(rabbitmqTransit.Imei)

	/* rabbitmqTransit */
	rabbitmqTransit.Imei = imeiPrepare

	buffer, err := r.Encoder(rabbitmqTransit)
	if err != nil {
		logger.Error("SendToRabbitAndWait imei={%s} rabbitmqReqId={%s}: Failed encoder buffer for send to queue: %s", rabbitmqTransit.Imei, rabbitmqTransit.RequestID, err.Error())
		return err
	}

	if err := r.Publish(config.QueueBrockerToHttp, buffer); err != nil {
		logger.Error("SendToRabbitAndWait imei={%s} rabbitmqReqId={%s}: Failed to publish: %s", rabbitmqTransit.Imei, rabbitmqTransit.RequestID, err.Error())
		return err
	}

	return nil
}
