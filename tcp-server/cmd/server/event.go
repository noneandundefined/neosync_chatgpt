package main

import (
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/usecase"
	"neomatica/neosync-tcp/pkg/protocol/adm"
)

func (tcp *tcpServer) ListenEvents(device *adm.ADMDevice) {
	if device == nil {
		logger.Warning("ListenEvents: Device unavailable: cannot fetch events (device not connected)")
		return
	}

	/* Initial UseCase */
	usecase := usecase.NewUseCase(device, tcp.db, tcp.store, tcp.cache)

	device.Emitter.On("hello-pack-received", func(received any) {
		tcp.handler.HelloPackEventHandler(device, &usecase, received)
	})

	device.Emitter.On("command-pack-received", func(received any) {
		tcp.handler.CommandPackEventHandler(device, received)
	})

	device.Emitter.On("sync-pack-received", func(received any) {
		tcp.handler.SyncPackEventHandler(device, &usecase, received)
	})

	device.Emitter.On("configuration-pack-received", func(received any) {
		tcp.handler.ConfigurationPackEventHandler(device, &usecase, received)
	})

	device.Emitter.On("keepalive-pack-received", func(received any) {
		tcp.handler.KeepalivePackEventHandler(device, received)
	})

	device.Emitter.On("error", func(received any) {
		tcp.handler.ErrorPackEventHandler(device, received)
	})
}
