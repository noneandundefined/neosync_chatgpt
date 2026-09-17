package meta_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: описание колонок в БД */

func (h *Handler) MetaSchemaColumnsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	table := r.URL.Query().Get("table")

	var columns any
	var err error

	switch table {
	case "configurations_prof":
		columns, err = h.Store.Configurations.Get_ColumnTableConfigurationProf(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "analytics_users":
		columns, err = h.Store.Analytics.Get_ColumnTableAnalyticsUsers(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "analytic_records":
		columns, err = h.Store.Analytics.Get_ColumnTableAnalyticRecords(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "product_analytics_events":
		columns, err = h.Store.Analytics.Get_ColumnTableProductAnalyticsEvents(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "configurations":
		columns, err = h.Store.Configurations.Get_ColumnTableConfiguration(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "devices":
		columns, err = h.Store.Devices.Get_ColumnTableDevice(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "groups":
		columns, err = h.Store.Groups.Get_ColumnTableGroup(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "group_members":
		columns, err = h.Store.GroupMembers.Get_ColumnTableGroupMember(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "syncs":
		columns, err = h.Store.Syncs.Get_ColumnTableSync(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "user_cores":
		columns, err = h.Store.Users.Get_ColumnTableUserCore(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "user_roles":
		columns, err = h.Store.Users.Get_ColumnTableUserRole(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	case "user_contacts":
		columns, err = h.Store.Users.Get_ColumnTableUserContacts(ctx)
		if err != nil {
			return httperr.Db(ctx, err)
		}

	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, columns)
	return nil
}
