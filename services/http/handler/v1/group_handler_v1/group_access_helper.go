package group_handler_v1

import (
	"context"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx/httperr"
)

type groupPermission int

const (
	groupPermOwnerOnly groupPermission = iota
	groupPermEditGroup
	groupPermManageDevices
	groupPermDelete
)

func (h *Handler) groupIsOwner(group *models.Group, userUUID string) bool {
	return group != nil && group.UserUuid == userUUID
}

func (h *Handler) groupIsMember(ctx context.Context, group *models.Group, userUUID string) (bool, error) {
	if h.groupIsOwner(group, userUUID) {
		return true, nil
	}

	m, err := h.Store.GroupMembers.Get_GroupMember(ctx, group.ID, userUUID)
	if err != nil {
		return false, err
	}

	return m != nil, nil
}

func (h *Handler) groupCanEdit(ctx context.Context, group *models.Group, userUUID string) (bool, error) {
	if h.groupIsOwner(group, userUUID) {
		return true, nil
	}

	m, err := h.Store.GroupMembers.Get_GroupMember(ctx, group.ID, userUUID)
	if err != nil {
		return false, err
	}

	if m == nil {
		return false, nil
	}

	return group.CanEditGroup, nil
}

func (h *Handler) groupCanManageDevices(ctx context.Context, group *models.Group, userUUID string) (bool, error) {
	if h.groupIsOwner(group, userUUID) {
		return true, nil
	}

	m, err := h.Store.GroupMembers.Get_GroupMember(ctx, group.ID, userUUID)
	if err != nil {
		return false, err
	}

	if m == nil {
		return false, nil
	}

	return group.CanManageDevices, nil
}

func (h *Handler) groupForbiddenFor(ctx context.Context, group *models.Group, userUUID string, perm groupPermission) error {
	tr := middleware.TranslatorFromContext(ctx)

	isMember, err := h.groupIsMember(ctx, group, userUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if !isMember {
		return httperr.Forbidden(tr.TErr("group-access-denied"))
	}

	switch perm {
	case groupPermOwnerOnly:
		return httperr.Forbidden(tr.TErr("group-owner-only"))
	case groupPermEditGroup:
		return httperr.Forbidden(tr.TErr("group-insufficient-rights-edit"))
	case groupPermManageDevices:
		return httperr.Forbidden(tr.TErr("group-insufficient-rights-manage-devices"))
	case groupPermDelete:
		return httperr.Forbidden(tr.TErr("group-insufficient-rights-delete"))
	default:
		return httperr.Forbidden(tr.TErr("group-access-denied"))
	}
}

func (h *Handler) groupShareRightsChanged(payload *UpdateGroupPayload, group *models.Group) bool {
	return payload.CanEditGroup != group.CanEditGroup ||
		payload.CanManageDevices != group.CanManageDevices ||
		payload.CanReadConfig != group.CanReadConfig ||
		payload.CanEditConfig != group.CanEditConfig ||
		payload.CanSendCommands != group.CanSendCommands
}
