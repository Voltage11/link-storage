package link_service

import (
	"context"
	"fmt"
	"link-storage/internal/middleware"
	"link-storage/internal/models"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"link-storage/pkg/utils/parseurl"
)

const (
	defaultTopVisitedCount = 10
)

func (s *linkService) CreateLink(ctx context.Context, linkCreate *models.LinkCreate) (*models.Link, error) {
	op := "link_service.CreateLink"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}

	link := &models.Link{
		UserID:      user.ID,
		LinkGroupID: linkCreate.LinkGroupID,
		URL:         linkCreate.URL,
		Title:       linkCreate.Title,
		Description: linkCreate.Description,
		IsFavorite:  linkCreate.IsFavorite,
	}

	// 1. Создадим запись ссылки и получим ее ID
	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, err
	}

	// 2. После создания ссылки и получения ID получим favicon, сохраним его на диск и запишем в БД
	linkUpdated, err := s.setLinkFavIconAndTitle(ctx, link.ID)
	if err != nil {
		// Даже если не удалось получить favicon, возвращаем ссылку
		return link, nil
	}

	return linkUpdated, nil
}

func (s *linkService) setLinkFavIconAndTitle(ctx context.Context, linkID int) (*models.Link, error) {
	op := "link_service.setLinkFavIconAndTitle"

	link, err := s.repo.GetLinkByID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	if link == nil {
		return nil, app_errors.NotFound("Ссылка не найдена", op)
	}

	urlInfo := parseurl.New(link.URL)
	
	// Обновляем заголовок если он пустой
	if link.Title == "" {
		link.Title = urlInfo.GetTitle()
	}

	// ИНИЦИАЛИЗИРУЕМ ПУСТУЮ СТРОКУ ДЛЯ ЛОКАЛЬНОГО ПУТИ
	localFaviconPath := ""
	
	// Пробуем скачать favicon
	favIconUrl := urlInfo.GetFaviconPath()
	if favIconUrl != "" {
		// Загрузим фактически себе на диск иконку
		localFaviconPath, err = urlInfo.DownloadFavicon(s.favIconsPath, link.UserID, link.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем
			s.logger.Warn(fmt.Sprintf("Не удалось скачать favicon для ссылки %d: %v", link.ID, err), op)
		}
	}

	// Если удалось скачать, сохраняем локальный путь, иначе оставляем пустым
	if localFaviconPath != "" {
		link.FaviconURL = localFaviconPath // ЛОКАЛЬНЫЙ ПУТЬ!
	} else {
		link.FaviconURL = "" // Или можно оставить пустым
	}

	// Обновляем запись в БД
	if err := s.repo.SetLinkFavIconAndTitle(ctx, link.ID, link.FaviconURL, link.Title); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *linkService) LinkRefreshIcon(ctx context.Context, linkID int) (*models.Link, error) {
	op := "link_service.LinkRefreshIcon"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}

	link, err := s.repo.GetLinkByID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	if link == nil {
		return nil, app_errors.NotFound("ссылка не найдена", op)
	}

	if link.UserID != user.ID {
		return nil, app_errors.NotFound("ссылка не найдена", op)
	}

	return s.setLinkFavIconAndTitle(ctx, linkID)
}

func (s *linkService) GetLinksByUserIDWithPagination(ctx context.Context, linkGroupID, page, pageSize int, name string) (*response.ListResponse[models.LinkResponse], error) {
	op := "link_service.GetLinksByUserIDWithPagination"

	user := middleware.GetCurrentUserFromContext(ctx)

	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}

	offset := pageSize * (page - 1)

	return s.repo.GetLinksByUserIDWithPagination(ctx, user.ID, linkGroupID, pageSize, offset, name)
}

func (s *linkService) LinkVisitedPlus(ctx context.Context, linkID int) error {
	op := "link_service.LinkVisitedPlus"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return app_errors.Unauthorized(op)
	}

	link, err := s.repo.GetLinkByID(ctx, linkID)
	if err != nil {
		return err
	}

	if link == nil {
		return app_errors.NotFound("ссылка не найдена", op)
	}

	if link.UserID != user.ID {
		return app_errors.NotFound("ссылка не найдена", op)
	}

	return s.repo.LinkVisitedPlus(ctx, linkID)
}

func (s *linkService) GetLinksTopVisited(ctx context.Context) ([]*models.Link, error) {
	op := "link_service.GetLinksTopVisited"
	
	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}
	return s.repo.GetLinksTopVisited(ctx, user.ID, defaultTopVisitedCount)
}