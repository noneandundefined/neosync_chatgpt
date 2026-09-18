package handlers

import (
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/types"
	"net"
	"net/http"
)

func (h *BasePackEventHandler) ErrorPackEventHandler(device *adm.ADMDevice, received any) {
	if device == nil || device.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	session, exists := h.Session.GetDeviceSession(*device.Imei)
	if session == nil || !exists {
		return
	}

	requestID := session.RequestId
	if requestID == nil {
		logger.Warning("ConfigurationPackEventHandler ip={%s} imei={%s}: RequestId not found connected client", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei)
		return
	}

	if _, ok := received.(error); !ok {
		logger.Error("ErrorPackEventHandler ip={%s} imei={%s}: expected error payload, got %T", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, received)
		return
	}

	/* RabbitMQ */
	rabbitmqTransit := types.RabbitMQ_TransitBinary{
		Type:      0x00,
		RequestID: *requestID,
		Imei:      *device.Imei,
		StatCode:  http.StatusBadRequest,
		Data:      []byte(received.(error).Error()),
	}

	if err := h.RMQ.SendToRabbitAsync(rabbitmqTransit); err != nil {
		logger.Error("ErrorPackEventHandler ip={%s}: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
	}
	h.Session.RemoveRequestIdSession(*requestID)
}
