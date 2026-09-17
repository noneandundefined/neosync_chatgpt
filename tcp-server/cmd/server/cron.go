package main

import (
	"context"
	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/usecase"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (tcp *tcpServer) startCronStatusNotActive() {
	reconcile := func() {
		if err := tcp.store.Devices.Update_DeviceStatusNotActive(context.Background(), tcp.session.ConnectedImeis()); err != nil {
			logger.Error("StartCronStatusNotActive: Failed to update device status: %s", err.Error())
		}
	}

	reconcile()

	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for range ticker.C {
			reconcile()
		}
	}()
}

func (s *tcpServer) startCronCleanupLogs() {
	const days = 31

	_, err := s.cron.AddFunc("0 0 * * *", func() {
		logDir := os.Getenv("LOG_DIR")
		if logDir == "" {
			logger.Error("startCronCleanupLogs: LOG_DIR is empty")
			return
		}

		expiration := time.Now().AddDate(0, 0, -days)

		entries, err := os.ReadDir(logDir)
		if err != nil {
			logger.Error("startCronCleanupLogs: Failed to read log dir: %s", err.Error())
			return
		}

		for _, entry := range entries {
			name := entry.Name()
			path := filepath.Join(logDir, name)

			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.ModTime().After(expiration) {
				continue
			}

			isDailyDir, _ := filepath.Match("log_*", name)
			isDailyZip, _ := filepath.Match("log_*.zip", name)
			if !isDailyDir && !isDailyZip {
				continue
			}

			if err := os.RemoveAll(path); err != nil {
				logger.Error("startCronCleanupLogs: Failed to delete %s: %s", path, err.Error())
			}
		}
	})

	if err != nil {
		logger.Error("startCronCleanupLogs: Failed to schedule job: %s", err.Error())
	}
}

/* startCronCompanyTasks проверяет сроки даже без пакетов от устройства. */
func (tcp *tcpServer) startCronCompanyTasks() {
	for _, process := range []func(){tcp.processCompanyTaskDeadlines, tcp.processCompanyTasks} {
		go func(process func()) {
			ticker := time.NewTicker(config.CompanyTaskPollInterval)
			defer ticker.Stop()

			for range ticker.C {
				process()
			}
		}(process)
	}
}

func (tcp *tcpServer) processCompanyTaskDeadlines() {
	ctx := context.Background()

	useCase := &usecase.CompanyTaskUseCase{Store: tcp.store, Db: tcp.db}

	tasks, err := tcp.store.Companies.Get_DueCompanyTasks(ctx)
	if err != nil {
		logger.Error("CompanyTasks: Failed to get due tasks: %s", err.Error())
		return
	}

	for _, task := range tasks {
		if task.Status == constants.COMPANY_TASK_STATUS_QUEUED {
			if _, connected := tcp.session.GetDeviceSession(task.DeviceImei); !connected {
				tcp.session.MarkCompanyTaskDisconnected(task.DeviceID)
			}
		}

		if err := useCase.ProcessTaskDeadline(ctx, &task); err != nil {
			logger.Error("CompanyTasks task_id={%d}: Failed to process deadline: %s", task.ID, err.Error())
		}
	}
}

func (tcp *tcpServer) processCompanyTasks() {
	ctx := context.Background()

	imeis, err := tcp.store.Companies.Get_CompanyTaskConnectedImeis(ctx, tcp.session.ConnectedImeis())
	if err != nil {
		logger.Error("CompanyTasks: Failed to get connected devices: %s", err.Error())
		return
	}

	workers := make(chan struct{}, config.CompanyTaskWorkers)

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, imei := range imeis {
		session, exists := tcp.session.GetDeviceSession(imei)
		if !exists || session.Device == nil {
			continue
		}

		workers <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-workers }()

			session.ConfigurationMu.Lock()
			defer session.ConfigurationMu.Unlock()

			tcp.handler.LimiterDb <- struct{}{}
			defer func() { <-tcp.handler.LimiterDb }()

			d, err := tcp.store.Devices.Get_DeviceFullByImei(ctx, imei)
			if err != nil {
				logger.Error("CompanyTasks imei={%s}: Failed to get device: %s", imei, err.Error())
				return
			}

			useCase := &usecase.CompanyTaskUseCase{Device: session.Device, Store: tcp.store, Db: tcp.db}
			if _, err := useCase.TryApplyPendingTask(ctx, tcp.session, d); err != nil {
				logger.Error("CompanyTasks imei={%s}: Failed to process task: %s", imei, err.Error())
			}
		}()
	}
}
