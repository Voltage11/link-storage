package link_service

import (
	"context"
	"link-storage/internal/middleware"
	"link-storage/internal/models"
	"link-storage/internal/repository/link_repository"
	"link-storage/pkg/logger"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
)

type LinkService interface {
	CreateLinkGroup(ctx context.Context, linkGroupCreate *models.LinkGroupCreate) (*models.LinkGroup, error)
	GetLinkGroupByID(ctc context.Context, id, userID int) (*models.LinkGroup, error)
	UpdateLinkGroup(ctx context.Context, linkGroupUpdate *models.LinkGroupUpdate) (*models.LinkGroup, error)
	DeleteLinkGroup(ctx context.Context, id int) error
	GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, page, pageSize int) (*response.ListResponse[models.LinkGroup], error)
}

type linkService struct {
	repo   link_repository.LinkRepository
	logger logger.AppLogger
}

func New(repo link_repository.LinkRepository, logger logger.AppLogger) LinkService {
	return &linkService{
		repo:   repo,
		logger: logger,
	}
}

func (s *linkService) CreateLinkGroup(ctx context.Context, linkGroupCreate *models.LinkGroupCreate) (*models.LinkGroup, error) {
	op := "link_service.CreateLinkGroup"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}

	// Перед созданием проверим, что нет группы с таким именем по пользователю
	exists, err := s.repo.HasLinkGroupWithNameByUserID(ctx, linkGroupCreate.Name, user.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, app_errors.Conflict("Группа с таким именем уже существует", op)
	}

	linkGroup := &models.LinkGroup{
		UserID:      user.ID,
		Name:        linkGroupCreate.Name,
		Description: linkGroupCreate.Description,
		Color:       linkGroupCreate.Color,
	}

	if err := s.repo.CreateLinkGroup(ctx, linkGroup); err != nil {
		return nil, err
	}

	return linkGroup, nil
}

func (s *linkService) GetLinkGroupByID(ctx context.Context, id, userID int) (*models.LinkGroup, error) {
	return s.repo.GetLinkGroupByID(ctx, id, userID)
}

func (s *linkService) UpdateLinkGroup(ctx context.Context, linkGroupUpdate *models.LinkGroupUpdate) (*models.LinkGroup, error) {
	op := "link_service.UpdateLinkGroup"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}

	existsLinkGroup, err := s.repo.GetLinkGroupByID(ctx, linkGroupUpdate.ID, user.ID)
	if err != nil {
		return nil, err
	}

	if existsLinkGroup == nil {
		return nil, app_errors.NotFound("Группа не найдена", op)
	}

	linkGroup := &models.LinkGroup{
		ID:          linkGroupUpdate.ID,
		UserID:      user.ID,
		Name:        linkGroupUpdate.Name,
		Description: linkGroupUpdate.Description,
		Color:       linkGroupUpdate.Color,
	}

	if err := s.repo.UpdateLinkGroup(ctx, linkGroup); err != nil {
		return nil, err
	}

	return linkGroup, nil
}

func (s *linkService) DeleteLinkGroup(ctx context.Context, id int) error {
	op := "link_service.DeleteLinkGroup"

	user := middleware.GetCurrentUserFromContext(ctx)
	if user == nil {
		return app_errors.Unauthorized(op)
	}

	linkGroup, err := s.repo.GetLinkGroupByID(ctx, id, user.ID)
	if err != nil {
		return err
	}

	if linkGroup == nil {
		return app_errors.NotFound("Группа не найдена", op)
	}

	if linkGroup.UserID != user.ID {
		return app_errors.NotFound("Группа не найдена", op)
	}

	return s.repo.DeleteLinkGroup(ctx, id)
}

func (s *linkService) GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, page, pageSize int) (*response.ListResponse[models.LinkGroup], error) {
	op := "link_service.GetLinkGroupsByUserIDWithPagination"

	user := middleware.GetCurrentUserFromContext(ctx)

	if user == nil {
		return nil, app_errors.Unauthorized(op)
	}
	offset := pageSize * (page - 1)
	return s.repo.GetLinkGroupsByUserIDWithPagination(ctx, name, user.ID, pageSize, offset)
}
