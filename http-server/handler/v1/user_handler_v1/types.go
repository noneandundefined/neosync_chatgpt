package user_handler_v1

import "neomatica/neosync/infra/store/postgres/models"

type UserCreatePayload struct {
	Email            string              `json:"email" validate:"required,email"`
	Language         string              `json:"language" validate:"required,len=2"`
	Locality         string              `json:"locality" validate:"omitempty"`
	NameOrganization string              `json:"name_organization" validate:"omitempty,max=100"`
	Phone            string              `json:"phone" validate:"omitempty"`
	Password         string              `json:"password" validate:"required,min=6,max=16"`
	RoleCode         string              `json:"role_code" validate:"required,gt=0"`
	Accesses         models.AccessFields `json:"accesses" validate:"required,gt=0"`
}

type SendEmailInfoPayload struct {
	TurnstileToken string `json:"turnstile_token" validate:"required"`
	Language       string `json:"language" validate:"required,oneof=ru en es"`
}

type UpdateUserMePayload struct {
	Email            *string `json:"email,omitempty" validate:"omitempty,email"`
	Password         *string `json:"password,omitempty" validate:"omitempty,min=6,max=16"`
	Phone            *string `json:"phone,omitempty" validate:"omitempty"`
	NameOrganization *string `json:"name_organization,omitempty" validate:"omitempty"`
	Locality         *string `json:"locality,omitempty" validate:"omitempty"`
	Language         *string `json:"language,omitempty" validate:"omitempty"`
}

type UpdateCVCGPayload struct {
	CanViewChildGroups bool `json:"can_view_child_groups"`
}

type UpdateUserOnesPayload struct {
	Email            *string                `json:"email,omitempty" validate:"omitempty,email"`
	Password         *string                `json:"password,omitempty" validate:"omitempty,min=6,max=16"`
	Phone            *string                `json:"phone" validate:"omitempty"`
	NameOrganization *string                `json:"name_organization,omitempty" validate:"omitempty"`
	Locality         *string                `json:"locality,omitempty" validate:"omitempty"`
	Language         *string                `json:"language,omitempty" validate:"omitempty"`
	Accesses         models.PtrAccessFields `json:"accesses,omitempty" validate:"omitempty"`
}

type UsersTreePageResponse struct {
	Items   []models.UserTree `json:"items"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	Total   int               `json:"total"`
	HasMore bool              `json:"has_more"`
}

type UsersToOwnerPageResponse struct {
	Items   []models.UserToOwner `json:"items"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	Total   int                  `json:"total"`
	HasMore bool                 `json:"has_more"`
}

type UsersWPResponse struct {
	Items      []models.User `json:"items"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	Total      int           `json:"total"`
	TotalAll   int           `json:"total_all"`
	TotalPages int           `json:"total_pages"`
}

type DeleteUsersPayload struct {
	UUIDs []string `json:"uuids" validate:"required"`
}

type UpdateConfigurationPriority struct {
	Enabled *bool `json:"force_neosync_configuration_priority" validate:"required"`
}
