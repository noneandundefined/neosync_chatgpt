package httperr

import (
	"context"
	"errors"
	"neomatica/neosync/infra/locale"
	"net"
	"net/http"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func New(msg string, code int) *ServerError {
	return &ServerError{
		Msg:  msg,
		Code: code,
	}
}

func InternalServerError(msg string) *ServerError {
	return New(msg, http.StatusInternalServerError)
}

func ServiceUnavailable(msg string) *ServerError {
	return New(msg, http.StatusServiceUnavailable)
}

func Forbidden(msg string) *ServerError {
	return New(msg, http.StatusForbidden)
}

func BadRequest(msg string) *ServerError {
	return New(msg, http.StatusBadRequest)
}

func RequestTimeout(msg string) *ServerError {
	return New(msg, http.StatusRequestTimeout)
}

func NotFound(msg string) *ServerError {
	return New(msg, http.StatusNotFound)
}

func Conflict(msg string) *ServerError {
	return New(msg, http.StatusConflict)
}

func Unauthorized(msg string) *ServerError {
	return New(msg, http.StatusUnauthorized)
}

func Db(ctx context.Context, err error) *ServerError {
	tr := ctx.Value("translator").(locale.Translator)

	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return RequestTimeout(tr.TErr(Err_ContextDeadlineExceeded.Error()))
	}

	/* Error Duplicate email */
	if errors.Is(err, Err_DuplicateEmail) {
		return Conflict(tr.TErr(Err_DuplicateEmail.Error()))
	}

	if errors.Is(err, context.Canceled) {
		return Conflict(tr.TErr(Err_ContextCanceled.Error()))
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return RequestTimeout(tr.TErr(Err_DbTimeout.Error()))
		}

		return InternalServerError(tr.TErr(Err_DbNetwork.Error()))
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": /* Error unique */
			if pqErr.Constraint == "devices_imei_key" {
				return Conflict(tr.TErr("device-imei-already-registered"))
			}

			return Conflict(tr.TErr(Err_UniqueViolation.Error()))
		case "23503": /* Not found */
			return NotFound(tr.TErr(Err_UserNotFound.Error()))
		}
	}

	switch {
	case errors.Is(err, Err_NotDeleted):
		return Conflict(tr.TErr(Err_NotDeleted.Error()))

	case errors.Is(err, Err_NotUpdated):
		return Conflict(tr.TErr(Err_NotUpdated.Error()))

	}

	return InternalServerError(tr.TErr("db-operation-failed"))
}

func Redis(ctx context.Context, err error) *ServerError {
	tr := ctx.Value("translator").(locale.Translator)

	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return NotFound(tr.TErr(Err_RedisKeyNotFound.Error()))
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return RequestTimeout(tr.TErr(Err_RedisTimeout.Error()))
		}

		return InternalServerError(tr.TErr(Err_RedisNetwork.Error()))
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return RequestTimeout(tr.TErr(Err_RedisDeadlineExceeded.Error()))
	}

	if errors.Is(err, context.Canceled) {
		return Conflict(tr.TErr(Err_RedisCanceled.Error()))
	}

	return InternalServerError(tr.TErr(Err_RedisOperationFailed.Error()))
}
