package link_handler

import (
	"link-storage/internal/service/link_service"
	"link-storage/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

type linkHandler struct {
	service link_service.LinkService
	logger  logger.AppLogger
}

func New(r *chi.Mux, service link_service.LinkService, logger logger.AppLogger) {
	h := &linkHandler{
		service: service,
		logger:  logger,
	}

	r.Route("/api/v1", func(r chi.Router) {
		// LinkGroup
		r.Use(httprate.LimitByIP(5, 1*time.Second))
		r.Post("/link-groups", h.linkGroupCreate)
		r.Put("/link-groups/{id}", h.linkGroupUpdate)
		r.Delete("/link-groups/{id}", h.linkGroupDelete)
		r.Get("/link-groups", h.linkGroupList)
		// Link
		r.Post("/links", h.linkCreate)
		r.Post("/links/refresh-icon/{id}", h.linkRefreshIcon)
		r.Post("/links/visited/{id}", h.linkVisitedPlus)
		r.Get("/links/top-visited", h.getLinkTopVisited)
		//r.Put("/links/{id}", h.linkUpdate)
		//r.Delete("/links/{id}", h.linkDelete)
		r.Get("/links", h.linkList)
	})
}
