package company_handler_v1

import (
	"neomatica/neosync/infra/store/postgres/models"
	"time"
)

func buildCompanyDetailsResponse(company *models.CompanyWithTasks) CompanyDetailsResponse {
	expiresAt := company.CreatedAt.Add(time.Duration(company.TTL) * 24 * time.Hour)

	stats := CompanyTaskStatusStats{
		Total: len(company.Tasks),
	}

	for _, task := range company.Tasks {
		switch task.Status {
		case "completed":
			stats.Completed++
		case "failed":
			stats.Failed++
		case "pending":
			stats.Pending++
		case "queued", "sending", "awaiting_confirmation":
			stats.Running++
		case "cancelled":
			stats.Cancelled++
		case "expired":
			stats.Expired++
		}
	}

	var launchedAt *time.Time
	if company.Status != "pending" {
		launchedAt = &company.UpdatedAt
	}

	return CompanyDetailsResponse{
		CompanyWithTasks: *company,
		StatusStats:      stats,
		ExpiresAt:        expiresAt,
		LaunchedAt:       launchedAt,
	}
}
