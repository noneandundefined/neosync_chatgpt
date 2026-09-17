package analytic_handler_v1

type WidgetResponse struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Section  string `json:"section"`
	Title    string `json:"title"`
	TitleKey string `json:"title_key,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Query    string `json:"query"`
}

type QueryPayload struct {
	Type     string `json:"type" validate:"required"`
	Query    string `json:"query" validate:"required"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	UserUUID string `json:"user_uuid,omitempty"`
	GroupID  *int   `json:"group_id,omitempty"`
	Model    string `json:"model,omitempty"`
}

type AnalyticsEventPayload struct {
	EventID    string                 `json:"event_id"`
	OccurredAt string                 `json:"occurred_at"`
	SessionID  string                 `json:"session_id"`
	EventName  string                 `json:"event_name"`
	Category   string                 `json:"category"`
	Path       string                 `json:"path,omitempty"`
	EntityType string                 `json:"entity_type,omitempty"`
	EntityID   string                 `json:"entity_id,omitempty"`
	DeviceID   *uint64                `json:"device_id,omitempty"`
	Success    *bool                  `json:"success,omitempty"`
	DurationMs *int                   `json:"duration_ms,omitempty"`
	ErrorCode  string                 `json:"error_code,omitempty"`
	AppVersion string                 `json:"app_version,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type AnalyticsEventsPayload struct {
	Events []AnalyticsEventPayload `json:"events"`
}

type AnalyticsEventsResponse struct {
	Accepted int `json:"accepted"`
}
