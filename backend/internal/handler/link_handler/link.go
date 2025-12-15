package link_handler

import (
	"link-storage/internal/models"
	"link-storage/pkg/request"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"net/http"
)

func (h *linkHandler) linkCreate(w http.ResponseWriter, r *http.Request) {
	op := "linkHandler.linkCreate"

	linkCreate, err := request.ParseRequestBody[models.LinkCreate](r)
	if err != nil || linkCreate == nil {
		response.WriteError(w, app_errors.BadRequest("неверный формат запроса", op))
		return
	}

	if err := linkCreate.Validate(); err != nil {
		response.WriteError(w, err)
		return
	}

	link, err := h.service.CreateLink(r.Context(), linkCreate)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	if link == nil {
		response.WriteError(w, app_errors.New(500, "internal server error", "не удалось создать ссылку", op))
		return
	}
	response.WriteSuccess(w, link)

}

func (h *linkHandler) linkRefreshIcon(w http.ResponseWriter, r *http.Request) {
	op := "linkHandler.linkRefreshIcon"

	linkID, ok := request.GetIntFromRequest(r, "id")
	if !ok {
		response.WriteError(w, app_errors.NotFound("ссылка не найдена", op))
		return
	}

	link, err := h.service.LinkRefreshIcon(r.Context(), linkID)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	if link == nil {
		response.WriteError(w, app_errors.NotFound("ссылка не найдена", op))
		return
	}

	response.WriteSuccess(w, link)
}

func (h *linkHandler) linkList(w http.ResponseWriter, r *http.Request) {
	page, pageSize := request.GetPaginateFromRequest(r)
	name, _ := request.GetQueryValueFromRequest(r, "q")
	linkGroupID, _ := request.GetIntFromRequest(r, "link_group_id")

	linkList, err := h.service.GetLinksByUserIDWithPagination(r.Context(), linkGroupID, page, pageSize, name)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, linkList)
}

// linkVisitedPlus при каждом посещении ссылки двинем счетчик
func (h *linkHandler) linkVisitedPlus(w http.ResponseWriter, r *http.Request) {
	op := "linkHandler.linkVisitedPlus"

	linkID, ok := request.GetIntFromRequest(r, "id")
	if !ok {
		response.WriteError(w, app_errors.BadRequest("Невалидный url", op))
		return
	}

	if err := h.service.LinkVisitedPlus(r.Context(), linkID); err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, nil)
}

func (h *linkHandler) getLinkTopVisited(w http.ResponseWriter, r *http.Request) {
	links, err := h.service.GetLinksTopVisited(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}

	if links == nil {
		response.WriteSuccess(w, []models.Link{})
		return
	}

	response.WriteSuccess(w, links)
}
