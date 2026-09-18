package pkg

import (
	"neomatica/neosync-tcp/infra/logger"
	"runtime/debug"
)

func RecoverGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("[%s] panic recovered: %v", name, r)
				logger.Error("[%s] stacktrace: %s", name, debug.Stack())
			}
		}()

		fn()
	}()
}
