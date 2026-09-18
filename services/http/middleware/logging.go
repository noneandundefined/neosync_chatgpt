package middleware

import (
	"bufio"
	"fmt"
	"neomatica/neosync/infra/logger"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int64
	wroteHeader bool
}

type Logger struct {
	currentLogDir string
	loc           *time.Location

	/* Files var */
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func NewLogger() *Logger {
	/* Locale */
	loc, _ := time.LoadLocation("Asia/Yekaterinburg")
	/* Logger */
	l := &Logger{loc: loc}
	/* Memory flash in file */
	go l.backgroundFlush()

	return l
}

func (l *Logger) backgroundFlush() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		l.mu.Lock()
		if l.writer != nil {
			l.writer.Flush()
		}
		l.mu.Unlock()
	}
}

func (l *Logger) openServLogFile() error {
	path := filepath.Join(l.currentLogDir, "server.log")

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = file
	l.writer = bufio.NewWriterSize(file, 64*1024)
	return nil
}

func (l *Logger) getCurrentLogDir() string {
	now := time.Now().In(l.loc)

	day := fmt.Sprintf("%02d", now.Day())
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%d", now.Year())

	dirName := fmt.Sprintf("log_%s%s%s", day, month, year)

	baseLogPath := os.Getenv("LOG_DIR")
	if baseLogPath == "" {
		baseLogPath = "./logs"
	}

	dir := filepath.Join(baseLogPath, dirName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create log directory: %s", err.Error())
	}

	return dir
}

// Обновление даты в папке
func (l *Logger) updateLogDirIfNeeded() {
	newDir := l.getCurrentLogDir()

	if newDir != l.currentLogDir {
		l.mu.Lock()
		defer l.mu.Unlock()

		if l.file != nil {
			l.writer.Flush()
			l.file.Close()
		}

		l.currentLogDir = newDir
		_ = l.openServLogFile()
	}
}

func (l *Logger) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(wrapped, r)

		if shouldSkipAccessLog(r) {
			return
		}

		entry := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" %v\n",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			getHTTPVersion(r),
			r.Method,
			r.URL.Path,
			wrapped.status,
			wrapped.size,
			r.UserAgent(),
			time.Since(start),
		)

		l.updateLogDirIfNeeded()
		l.productionLogging(entry)
	})
}

func shouldSkipAccessLog(r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return true
	}

	path := r.URL.Path
	for _, prefix := range []string{"/assets/", "/local/", "/ui/", "/favicon", "/health"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

func getHTTPVersion(r *http.Request) string {
	switch r.ProtoMajor {
	case 1:
		return "HTTP/1.1"
	case 2:
		return "HTTP/2.0"
	case 3:
		return "HTTP/3.0"
	default:
		return fmt.Sprintf("HTTP/%d.%d", r.ProtoMajor, r.ProtoMinor)
	}
}

func (l *Logger) productionLogging(log string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.writer == nil {
		_ = l.openServLogFile()
	}

	l.writer.WriteString(log)
}

func (rw *responseWriter) WriteHeader(status int) {
	if rw.wroteHeader {
		return
	}

	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

func (rw *responseWriter) Flush() {
	if fl, ok := rw.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker not supported")
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	if cn, ok := rw.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return make(chan bool)
}
