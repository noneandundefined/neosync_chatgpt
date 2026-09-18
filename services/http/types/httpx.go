package types

import (
	"compress/gzip"
	"net/http"
)

/* Структура возврата сжатого ответа от сервера */
type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}
