package main

import (
	"context"
	"fmt"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (s *httpServer) startCronFirmwareUpdater() {
	client := &http.Client{Timeout: 30 * time.Second}

	_, err := s.cron.AddFunc("0 0 * * *", func() {
		ctx := context.Background()

		deviceFirmwareSources, err := s.store.DeviceModelSources.Get_Sources(ctx)
		if err != nil {
			return
		}

		var wg sync.WaitGroup
		for _, d := range deviceFirmwareSources {
			wg.Add(1)
			go func(d models.DeviceModelSource) {
				defer wg.Done()

				resp, err := headFirmwareWithRetry(d.FirmwareUrl, client.Head, time.Sleep)
				if err != nil {
					logger.Error("StartCronFirmwareUpdater: %s: %s", d.DeviceModel, err.Error())
					return
				}
				defer resp.Body.Close()

				etag := resp.Header.Get("ETag")
				if etag == "" {
					return
				}

				/* Convert to firmware */
				etag = strings.TrimPrefix(etag, "W/")
				etag = strings.Trim(etag, `"`)
				version, err := strconv.Atoi(etag)
				if err != nil {
					logger.Error("StartCronFirmwareUpdater: %s", err.Error())
					return
				}

				if err := s.store.Syncs.Update_SyncFirmwareVersionUpd(ctx, uint16(version), d.DeviceModel); err != nil {
					logger.Error("StartCronFirmwareUpdater: %s", err.Error())
					return
				}

				logger.Info("StartCronFirmwareUpdater: [OK] %s → version %d updated\n", d.DeviceModel, version)
			}(d)
		}
		wg.Wait()
	})

	if err != nil {
		logger.Error("StartCronFirmwareUpdater: Failed to schedule job: %s", err.Error())
	}
}

func headFirmwareWithRetry(url string, head func(string) (*http.Response, error), sleep func(time.Duration)) (*http.Response, error) {
	const maxAttempts = 3
	const retryDelay = 10 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := head(url)
		if err == nil {
			return resp, nil
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		lastErr = err
		if attempt < maxAttempts {
			sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("firmware HEAD failed after %d attempts: %w", maxAttempts, lastErr)
}

func (s *httpServer) startCronCleanupLogs() {
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
