package models

type Role struct {
	ID      uint   `json:"id" db:"id"`
	Code    string `json:"code" db:"code"`
	NameKey string `json:"name_key" db:"name_key"`
	DescKey string `json:"desc_key" db:"desc_key"`
}
