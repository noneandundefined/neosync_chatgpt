package models

import (
	"time"
)

type UserOwner interface {
	GetUserUUID() *string
}

type UserCore struct {
	ID                 uint64    `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	UserUUID           string    `json:"user_uuid" db:"user_uuid"`
	ParentUUID         *string   `json:"parent_uuid" validate:"omitempty" db:"parent_uuid"`
	CanViewChildGroups bool      `json:"can_view_child_groups" db:"can_view_child_groups"`
	Email              string    `json:"email" db:"email"`
	Password           string    `json:"password" db:"password"`
	RefreshToken       *string   `json:"refresh_token" db:"refresh_token"`
}

type UserContact struct {
	ID               uint64    `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	UserUUID         string    `json:"user_uuid" db:"user_uuid"`
	NameOrganization *string   `json:"name_organization" db:"name_organization"`
	Locality         *string   `json:"locality" db:"locality"`
	Phone            *string   `json:"phone" db:"phone"`
	Language         string    `json:"language" db:"language"`
}

type UserRole struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserUUID  string    `json:"user_uuid" db:"user_uuid"`
	RoleCode  string    `json:"role_code" db:"role_code"`
}

type UserAccess struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserUUID  string    `json:"user_uuid" db:"user_uuid"`
	AccessFields
}

type User struct {
	UserContact        UserContact `json:"user_contact" db:"user_contact"`
	ParentUUID         *string     `json:"parent_uuid" validate:"omitempty" db:"parent_uuid"`
	CanViewChildGroups bool        `json:"can_view_child_groups" db:"can_view_child_groups"`
	Email              string      `json:"email" db:"email"`
	ParentEmail        *string     `json:"parent_email" validate:"omitempty" db:"parent_email"`
	Role               Role        `json:"role" db:"role"`
	AccessFields
}

type UserToOwner struct {
	UserUUID string `json:"user_uuid" db:"user_uuid"`
	Email    string `json:"email" db:"email"`
	RoleCode string `json:"role_code" db:"role_code"`
}

type UserAuth struct {
	ForceNeosyncConfigurationPriority bool        `json:"force_neosync_configuration_priority" db:"force_neosync_configuration_priority"`
	UserContact                       UserContact `json:"user_contact" db:"user_contact"`
	ParentUUID                        *string     `json:"parent_uuid" validate:"omitempty" db:"parent_uuid"`
	CanViewChildGroups                bool        `json:"can_view_child_groups" db:"can_view_child_groups"`
	Email                             string      `json:"email" db:"email"`
	ParentEmail                       *string     `json:"parent_email" validate:"omitempty" db:"parent_email"`
	Password                          string      `json:"password" db:"password"`
	RefreshToken                      *string     `json:"refresh_token" db:"refresh_token"`
	Role                              Role        `json:"role" db:"role"`
	AccessFields
}

type UserLoginWithRoleCode struct {
	Login    string `json:"login" db:"login"`
	RoleCode string `json:"role_code" db:"role_code"`
}

type UserUpdate struct {
	Email            *string `json:"email" db:"email"`
	Password         *string `json:"password" db:"password"`
	Phone            *string `json:"phone" db:"phone"`
	NameOrganization *string `json:"name_organization" db:"name_organization"`
	Locality         *string `json:"locality" db:"locality"`
	Language         *string `json:"language" db:"language"`
	Accesses         PtrAccessFields
}

type UserTree struct {
	UserUUID   string  `json:"user_uuid" db:"user_uuid"`
	ParentUUID *string `json:"parent_uuid" db:"parent_uuid"`
	RoleCode   string  `json:"role_code" db:"code"`
	Email      string  `json:"email" db:"email"`
}
