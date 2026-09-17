package analytic_handler_v1

import (
	"encoding/json"
	"neomatica/neosync/infra/analytics"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	maxAnalyticsEventsBatch = 50
	maxAnalyticsProperties  = 8192
)

var analyticEventNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,95}$`)
var publicAnalyticsCategories = map[string]bool{
	"authentication": true, "error": true, "interaction": true, "navigation": true,
	"performance": true, "search": true, "session": true,
}

func (h *Handler) CreateAnalyticsEventsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	identity := middleware.GetIdentity(ctx)
	if h.Analytics == nil {
		return httperr.BadRequest("analytics unavailable")
	}

	var payload AnalyticsEventsPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if len(payload.Events) == 0 || len(payload.Events) > maxAnalyticsEventsBatch {
		return httperr.BadRequest("invalid analytics batch size")
	}

	userUUID := ""
	roleCode := ""
	if identity != nil {
		userUUID = identity.User.UserContact.UserUUID
		roleCode = identity.RoleCode
	}
	events := make([]analytics.Event, 0, len(payload.Events))
	for _, item := range payload.Events {
		if !analyticEventNamePattern.MatchString(item.EventName) || !analyticEventNamePattern.MatchString(item.Category) {
			continue
		}
		if len(item.EventID) < 8 || len(item.EventID) > 64 || len(item.SessionID) < 8 || len(item.SessionID) > 64 {
			continue
		}
		if identity == nil && !publicAnalyticsCategories[item.Category] {
			continue
		}

		properties, err := json.Marshal(item.Properties)
		if err != nil || len(properties) > maxAnalyticsProperties {
			continue
		}

		occurredAt, err := time.Parse(time.RFC3339Nano, item.OccurredAt)
		if err != nil || occurredAt.Before(time.Now().AddDate(0, -1, 0)) || occurredAt.After(time.Now().Add(5*time.Minute)) {
			occurredAt = time.Now().UTC()
		}

		deviceID := item.DeviceID
		entityType := item.EntityType
		entityID := item.EntityID
		if identity == nil {
			deviceID = nil
			entityType = ""
			entityID = ""
		}

		events = append(events, analytics.Event{
			EventID: strings.TrimSpace(item.EventID), OccurredAt: occurredAt, UserUUID: userUUID,
			RoleCode: roleCode, SessionID: strings.TrimSpace(item.SessionID),
			EventName: item.EventName, Category: item.Category, Source: "web",
			Path: item.Path, EntityType: entityType, EntityID: entityID,
			DeviceID: deviceID, Success: item.Success, DurationMs: item.DurationMs,
			ErrorCode: item.ErrorCode, AppVersion: item.AppVersion, Properties: item.Properties,
		})
	}

	accepted := h.Analytics.Enqueue(events)
	httpx.HttpResponse(w, r, http.StatusAccepted, AnalyticsEventsResponse{Accepted: accepted})
	return nil
}
