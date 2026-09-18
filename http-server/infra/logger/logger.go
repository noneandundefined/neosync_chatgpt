package logger

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	logDir   string
	logQueue = make(chan logEntry, 2000)
	timeYkb  *time.Location
	dev            bool
	packetTrace    bool
	droppedLogs    atomic.Uint64
	lastDropReport atomic.Int64
)

type logEntry struct {
	filename string
	message  string
}

func InitLogger() {
	var err error
	timeYkb, err = time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		panic("Failed to load timezone: " + err.Error())
	}

	dev = os.Getenv("GO_ENV") == "DEV"
	packetTrace = dev || strings.EqualFold(os.Getenv("LOG_PACKET_TRACE"), "true") || os.Getenv("LOG_PACKET_TRACE") == "1"

	ensureLogDir()

	/* Workers */
	go logWorker()
	go compressWorker()
}

func logWorker() {
	for entry := range logQueue {
		writeLog(entry.filename, entry.message)
	}
}

func writeLog(filename, log string) {
	ensureLogDir()
	filePath := filepath.Join(logDir, filename)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	f.WriteString(log + "\n")
}

func writeDateTime(log *strings.Builder, t time.Time) {
	log.WriteString(fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute(), t.Second()))
}

func ensureLogDir() {
	now := time.Now().In(timeYkb)
	day := fmt.Sprintf("%02d", now.Day())
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%d", now.Year())

	baseLogPath := os.Getenv("LOG_DIR")
	if baseLogPath == "" {
		baseLogPath = "./logs"
	}

	newLogDir := filepath.Join(baseLogPath, fmt.Sprintf("log_%s%s%s", day, month, year))

	if newLogDir != logDir {
		if err := os.MkdirAll(newLogDir, os.ModePerm); err != nil {
			fmt.Println("Error creating log directory:", err)
			return
		}

		logDir = newLogDir
	}
}

func logWithLevel(level, filename, format string, args ...any) {
	now := time.Now().In(timeYkb)
	var log strings.Builder

	log.WriteString("[")
	writeDateTime(&log, now)
	log.WriteString("] [")
	log.WriteString(level)
	log.WriteString("] ")
	log.WriteString(fmt.Sprintf(format, args...))

	logMessage := log.String()

	/* If DEV write in console */
	if dev {
		fmt.Println(logMessage)
	}

	select {
	case logQueue <- logEntry{filename: filename, message: logMessage}:
	default:
		reportDroppedLog()
	}
}

func reportDroppedLog() {
	droppedLogs.Add(1)

	now := time.Now().Unix()
	last := lastDropReport.Load()
	if now-last < 5 || !lastDropReport.CompareAndSwap(last, now) {
		return
	}

	dropped := droppedLogs.Swap(0)
	fmt.Printf("Log queue full: dropped %d messages in the last interval\n", dropped)
}

func Packet(format string, args ...any) {
	if !packetTrace {
		return
	}

	logWithLevel("PACKET", "packets.log", format, args...)
}

func Info(format string, args ...any) {
	logWithLevel("INFO", "info.log", format, args...)
}

func Mail(format string, args ...any) {
	logWithLevel("MAIL", "mails.log", format, args...)
}

func Warning(format string, args ...any) {
	logWithLevel("WARN", "warning.log", format, args...)
}

func Session(format string, args ...any) {
	logWithLevel("SESS", "session.log", format, args...)
}

func Error(format string, args ...any) {
	logWithLevel("ERROR", "errors.log", format, args...)
}

func compressWorker() {
	compressOldDirs()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		compressOldDirs()
	}
}

func compressOldDirs() {
	baseLogPath := filepath.Dir(logDir)

	entries, err := os.ReadDir(baseLogPath)
	if err != nil {
		fmt.Println("Error reading log directory:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(baseLogPath, entry.Name())
		if fullPath == logDir {
			continue
		}

		zipPath := fullPath + ".zip"
		if _, err := os.Stat(zipPath); err == nil {
			continue
		}

		if err := zipFolder(fullPath, zipPath); err != nil {
			fmt.Println("Error compressing folder:", fullPath, err)
		} else {
			os.RemoveAll(fullPath)
		}
	}
}

func zipFolder(src, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		fw, err := w.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, f)
		return err
	})
}
