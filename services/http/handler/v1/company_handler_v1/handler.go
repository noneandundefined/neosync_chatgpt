package company_handler_v1

import (
	"neomatica/neosync/handler"
	"neomatica/neosync/infra/worker"
)

type Handler struct {
	*handler.BaseHandler
	CompanyWorker *worker.CompanyTasksWorker
}

func NewHandler(base *handler.BaseHandler, companyWorker *worker.CompanyTasksWorker) *Handler {
	return &Handler{
		BaseHandler:   base,
		CompanyWorker: companyWorker,
	}
}
