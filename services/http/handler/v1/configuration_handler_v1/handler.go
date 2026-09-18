package configuration_handler_v1

import (
	"neomatica/neosync/handler"
)

// var cfgParser = cfgparser.CfgParser{}

type Handler struct {
	*handler.BaseHandler
}

func NewHandler(base *handler.BaseHandler) *Handler {
	return &Handler{BaseHandler: base}
}
