package models

import "time"

type ConfigurationTemplate struct {
	ID           uint64    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_ut"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UserUUID     string    `json:"user_uuid" db:"user_uuid"`
	Name         string    `json:"name" db:"name"`
	Model        *string   `json:"model,omitempty" db:"model"`
	TypeOfSaving string    `json:"type_of_saving" db:"type_of_saving"`
}
