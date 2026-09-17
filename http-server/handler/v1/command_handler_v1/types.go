package command_handler_v1

import "context"

type CommandPayload struct {
	SessId   string   `json:"sess_id" validate:"required"`
	SendMode string   `json:"send_mode" validate:"required,oneof=instant on_connect"`
	Command  string   `json:"command" validate:"required,min=2,max=256"`
	Imeis    []string `json:"imeis" validate:"required,min=1,max=200"`
}

type CommandWaitPayload struct {
	Command string `json:"command" validate:"required,min=2"`
}

type BleTask struct {
	ctx    context.Context
	cancel context.CancelFunc
}
