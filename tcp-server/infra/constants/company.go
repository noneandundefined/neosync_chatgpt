package constants

const (
	COMPANY_STATUS_PENDING   = "pending"
	COMPANY_STATUS_RUNNING   = "running"
	COMPANY_STATUS_COMPLETED = "completed"
	COMPANY_STATUS_FAILED    = "failed"
	COMPANY_STATUS_CANCELLED = "cancelled"

	COMPANY_TASK_STATUS_PENDING               = "pending"
	COMPANY_TASK_STATUS_QUEUED                = "queued"
	COMPANY_TASK_STATUS_SENDING               = "sending"
	COMPANY_TASK_STATUS_AWAITING_CONFIRMATION = "awaiting_confirmation"
	COMPANY_TASK_STATUS_EXPIRED               = "expired"
	COMPANY_TASK_STATUS_COMPLETED             = "completed"
	COMPANY_TASK_STATUS_FAILED                = "failed"
	COMPANY_TASK_STATUS_CANCELLED             = "cancelled"

	COMPANY_TASK_ERROR_EXPIRED    = "campaign expired"
	COMPANY_TASK_ERROR_NO_ACCOUNT = "device has no account"
)
