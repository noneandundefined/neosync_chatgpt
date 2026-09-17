package handlers

import (
	"context"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/util"
	"net"
)

func (h *BasePackEventHandler) KeepalivePackEventHandler(device *adm.ADMDevice, received any) {
	if device == nil || device.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	imei := util.PrepareImei(*device.Imei)

	h.limiterDatabase(func() {
		if err := h.Store.Devices.Update_DeviceStatusByImei(context.Background(), imei, true); err != nil {
			logger.Error("KeepalivePackEventHandler ip={%s} imei={%s}: Failed update device status true/false", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei)
		}
	})
}
