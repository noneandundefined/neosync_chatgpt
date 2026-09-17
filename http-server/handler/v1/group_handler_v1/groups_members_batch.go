package group_handler_v1

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) requireGroupOwner(ctx context.Context, groupID uint64) (*models.Group, error) {
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !authToken.User.Has(constants.AccessGroupManage) {
		return nil, httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	group, err := h.Store.Groups.Get_GroupById(ctx, groupID)
	if err != nil {
		return nil, httperr.Db(ctx, err)
	}

	if group == nil {
		return nil, httperr.NotFound(tr.TErr("group-not-found"))
	}

	if !h.groupIsOwner(group, authToken.User.UserContact.UserUUID) {
		return nil, h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermOwnerOnly)
	}

	return group, nil
}

func parseGroupMembersBatch(ctx context.Context, r *http.Request) (*GroupMembersBatchPayload, uint64, error) {
	tr := middleware.TranslatorFromContext(ctx)

	groupIDStr := mux.Vars(r)["id"]
	if groupIDStr == "" {
		return nil, 0, httperr.NotFound(tr.TErr("group-id-not-found"))
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		logger.Error("parseGroupMembersBatch req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, 0, httperr.BadRequest(tr.TErr("invalid-group-id"))
	}

	var payload *GroupMembersBatchPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return nil, 0, httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return nil, 0, httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return nil, 0, httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	return payload, groupID, nil
}

func filterMemberUUIDs(group *models.Group, uuids []string) []string {
	out := make([]string, 0, len(uuids))
	seen := make(map[string]struct{}, len(uuids))

	for _, uuid := range uuids {
		if uuid == "" || uuid == group.UserUuid {
			continue
		}

		if _, ok := seen[uuid]; ok {
			continue
		}

		seen[uuid] = struct{}{}
		out = append(out, uuid)
	}

	return out
}

/* Neosync HTTPx V1 */
/* Handler: массовое предоставление доступа к группе */

func (h *Handler) GroupMembersUpsertMassiveHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	payload, groupID, err := parseGroupMembersBatch(ctx, r)
	if err != nil {
		return err
	}

	if err := smartcaptcha.VerifyRequest(tr, payload.TurnstileToken, r); err != nil {
		logger.Error("GroupMembersUpsertMassiveHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	group, err := h.requireGroupOwner(ctx, groupID)
	if err != nil {
		return err
	}

	memberUUIDs := filterMemberUUIDs(group, payload.MemberUUIDs)
	if len(memberUUIDs) == 0 {
		return httperr.BadRequest(tr.TErr("user-list-empty"))
	}

	if err := h.Store.GroupMembers.Upsert_GroupMembers(ctx, groupID, memberUUIDs); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("group-members-saved"))
	return nil
}

/* Neosync HTTPx V1 */
/* Handler: массовый отзыв доступа к группе */

func (h *Handler) GroupMembersDeleteMassiveHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	payload, groupID, err := parseGroupMembersBatch(ctx, r)
	if err != nil {
		return err
	}

	if err := smartcaptcha.VerifyRequest(tr, payload.TurnstileToken, r); err != nil {
		logger.Error("GroupMembersDeleteMassiveHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	group, err := h.requireGroupOwner(ctx, groupID)
	if err != nil {
		return err
	}

	memberUUIDs := filterMemberUUIDs(group, payload.MemberUUIDs)
	if len(memberUUIDs) == 0 {
		return httperr.BadRequest(tr.TErr("user-list-empty"))
	}

	if err := h.Store.GroupMembers.Delete_GroupMembers(ctx, groupID, memberUUIDs); err != nil {
		if err == sql.ErrNoRows {
			return httperr.NotFound(tr.TErr("group-member-not-found"))
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("group-member-removed"))
	return nil
}
