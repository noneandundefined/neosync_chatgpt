package types

import (
	"neomatica/neosync/infra/store/postgres/models"
	"time"
)

type RequestAuthToken struct {
	UUID      string    `json:"uuid"`
	IPAddress string    `json:"ip_address"`
	SessionID string    `json:"session_id"`
	RoleCode  string    `json:"role_code"`
	Timestamp time.Time `json:"timestamp"`
}

type AuthToken struct {
	User     models.UserAuth `json:"user"`
	RoleCode string          `json:"role_code"`
}
