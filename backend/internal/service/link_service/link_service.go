package link_service

import (
	"context"
	"link-storage/internal/models"
	"link-storage/internal/repository/link_repository"
	"link-storage/pkg/logger"
	"link-storage/pkg/response"
)

type LinkService interface {
	// LinkGroup
	CreateLinkGroup(ctx context.Context, linkGroupCreate *models.LinkGroupCreate) (*models.LinkGroup, error)
	GetLinkGroupByID(ctc context.Context, id, userID int) (*models.LinkGroup, error)
	UpdateLinkGroup(ctx context.Context, linkGroupUpdate *models.LinkGroupUpdate) (*models.LinkGroup, error)
	DeleteLinkGroup(ctx context.Context, id int) error
	GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, page, pageSize int) (*response.ListResponse[models.LinkGroup], error)

	// Link
	CreateLink(ctx context.Context, linkCreate *models.LinkCreate) (*models.Link, error)
	LinkRefreshIcon(ctx context.Context, linkID int) (*models.Link, error)
	GetLinksByUserIDWithPagination(ctx context.Context, linkGroupID, page, pageSize int, name string) (*response.ListResponse[models.LinkResponse], error)
	LinkVisitedPlus(ctx context.Context, linkID int) error
	//GetLinkByID(ctx context.Context, id, userID int) (*models.Link, error)
	//UpdateLink(ctx context.Context, linkUpdate *models.LinkUpdate) (*models.Link, error)
	//DeleteLink(ctx context.Context, id int) error
	//GetLinksByLinkGroupIDWithPagination(ctx context.Context, linkGroupID, page, pageSize int) (*response.ListResponse[models.Link], error))
}

type linkService struct {
	repo         link_repository.LinkRepository
	logger       logger.AppLogger
	favIconsPath string
}

func New(repo link_repository.LinkRepository, logger logger.AppLogger, favIconsPath string) LinkService {
	return &linkService{
		repo:         repo,
		logger:       logger,
		favIconsPath: favIconsPath,
	}
}
