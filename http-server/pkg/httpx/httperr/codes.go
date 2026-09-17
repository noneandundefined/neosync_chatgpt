package httperr

import "errors"

var (
	Err_DuplicateEmail          = errors.New("errors.duplicate-email")
	Err_UserNotFound            = errors.New("errors.user-not-found")
	Err_ContextDeadlineExceeded = errors.New("errors.context-deadline-exceeded")
	Err_ContextCanceled         = errors.New("errors.context-canceled")
	Err_UniqueViolation         = errors.New("errors.unique-violation")
	Err_DbTimeout               = errors.New("errors.db-timeout")
	Err_DbNetworkTemporary      = errors.New("errors.db-network-temporary")
	Err_DbNetwork               = errors.New("errors.db-network")
	Err_NotDeleted              = errors.New("errors.not-deleted")
	Err_NotUpdated              = errors.New("errors.not-updated")
)

var (
	Err_RedisKeyNotFound      = errors.New("errors.redis-key-not-found")
	Err_RedisTimeout          = errors.New("errors.redis-timeout")
	Err_RedisDeadlineExceeded = errors.New("errors.redis-deadline-exceeded")
	Err_RedisCanceled         = errors.New("errors.redis-canceled")
	Err_RedisNetwork          = errors.New("errors.redis-network")
	Err_RedisOperationFailed  = errors.New("errors.redis-operation-failed")
)
