package adm_v2

import (
	"database/sql"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"sync/atomic"
)

type ADM_V2 struct {
	Store                 store.Storage
	Db                    *sql.DB
	Imei                  string
	deviceModelLookupDone atomic.Bool
}

func (adm *ADM_V2) MarkDeviceModelLookupSatisfied() {
	adm.deviceModelLookupDone.Store(true)
}
