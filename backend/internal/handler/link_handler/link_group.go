package link_handler

import (
	"link-storage/internal/models"
	"link-storage/pkg/request"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"net/http"
)

func (h *linkHandler) linkGroupCreate(w http.ResponseWriter, r *http.Request) {
	op := "link_handler.linkGroupCreate"

	linkGroupRequest, err := request.ParseRequestBody[models.LinkGroupCreate](r)

	if err != nil || linkGroupRequest == nil {
		response.WriteError(w, app_errors.BadRequest("Неверный формат запроса", op))
		return
	}

	if err := linkGroupRequest.Validate(); err != nil {
		response.WriteError(w, app_errors.BadRequestWithError(err, "", op))
		return
	}

	linkGroupCreated, err := h.service.CreateLinkGroup(r.Context(), linkGroupRequest)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, linkGroupCreated)
}

func (h *linkHandler) linkGroupUpdate(w http.ResponseWriter, r *http.Request) {
	op := "link_handler.linkGroupUpdate"

	linkGroupID, ok := request.GetIntFromRequest(r, "id")
	if !ok {
		response.WriteError(w, app_errors.BadRequest("Неверный формат запроса", op))
		return
	}

	linkGroupRequest, err := request.ParseRequestBody[models.LinkGroupUpdate](r)

	if err != nil || linkGroupRequest == nil {
		response.WriteError(w, app_errors.BadRequest("Неверный формат запроса", op))
		return
	}

	if err := linkGroupRequest.Validate(); err != nil {
		response.WriteError(w, app_errors.BadRequestWithError(err, "", op))
		return
	}

	linkGroupRequest.ID = linkGroupID

	linkGroupUpdated, err := h.service.UpdateLinkGroup(r.Context(), linkGroupRequest)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	if linkGroupUpdated == nil {
		response.WriteError(w, app_errors.NotFound("Группа не найдена", op))
		return
	}

	response.WriteSuccess(w, linkGroupUpdated)
}

func (h *linkHandler) linkGroupDelete(w http.ResponseWriter, r *http.Request) {
	op := "link_handler.linkGroupDelete"

	linkGroupID, ok := request.GetIntFromRequest(r, "id")
	if !ok {
		response.WriteError(w, app_errors.BadRequest("Неверный формат запроса", op))
		return
	}

	if err := h.service.DeleteLinkGroup(r.Context(), linkGroupID); err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, nil)
}
