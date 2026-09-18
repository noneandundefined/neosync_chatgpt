package company_handler_v1

import (
	"neomatica/neosync/infra/store/postgres/models"
	"time"
)

type CreateCompanyPayload struct {
	Name                      string   `json:"name" validate:"required,max=255"`
	TTL                       string   `json:"ttl" validate:"required,oneof=1 7 30"`
	PriorityModel             string   `json:"priority_model" validate:"omitempty,max=255"`
	ExistingTaskAction        string   `json:"existing_task_action" validate:"required,oneof=skip replace overwrite"`
	LaunchMode                string   `json:"launch_mode" validate:"required,oneof=all selected"`
	IncompatibleAction        string   `json:"incompatible_action" validate:"required,oneof=exclude"`
	ConfigurationSource       string   `json:"configuration_source" validate:"omitempty,oneof=device_sources template_sources file_sources device template file"`
	ConfigurationCource       string   `json:"configuration_cource" validate:"omitempty,oneof=device_sources template_sources file_sources device template file"`
	ConfigurationDeviceCource *string  `json:"configuration_device_cource" validate:"omitempty,max=255"`
	ConfigurationTemplateID   *uint64  `json:"configuration_template_id" validate:"omitempty"`
	ConfigurationFileName     *string  `json:"configuration_file_name" validate:"omitempty,max=255"`
	ConfigurationFileBase64   *string  `json:"configuration_file_base64" validate:"omitempty"`
	IMEIs                     []string `json:"imeis" validate:"required,min=1,max=250,dive,required"`
}

type CreateCompanyResponse struct {
	ID           uint64 `json:"id"`
	TasksCreated int    `json:"tasks_created"`
	Message      string `json:"message"`
}

type CompaniesWPResponse struct {
	Items      []models.Company `json:"items"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int              `json:"total"`
	TotalAll   int              `json:"total_all"`
	TotalPages int              `json:"total_pages"`
}

type CompanyTaskStatusStats struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Cancelled int `json:"cancelled"`
	Expired   int `json:"expired"`
}

type CompanyDetailsResponse struct {
	models.CompanyWithTasks
	StatusStats CompanyTaskStatusStats `json:"status_stats"`
	ExpiresAt   time.Time              `json:"expires_at"`
	LaunchedAt  *time.Time             `json:"launched_at,omitempty"`
}
