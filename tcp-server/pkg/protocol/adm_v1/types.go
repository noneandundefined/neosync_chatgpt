package adm_v1

import (
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/infra/store/postgres/usecase"
)

type ADM_V1 struct {
	Store   store.Storage
	UseCase usecase.UseCase
}
