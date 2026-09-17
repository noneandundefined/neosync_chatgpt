package usecase

import (
	"database/sql"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
)

type ConfigurationUseCase struct {
	Device *adm.ADMDevice
	Store  store.Storage
	Db     *sql.DB
	Cache  *memory.Cache
}
