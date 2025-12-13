package link_repository

import (
	"context"
	"link-storage/internal/models"
	"link-storage/pkg/logger"
	"link-storage/pkg/response"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository interface {
	//LinkGroup

	CreateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error
	GetLinkGroupByID(ctc context.Context, id, userID int) (*models.LinkGroup, error)
	HasLinkGroupWithNameByUserID(ctx context.Context, name string, userID int) (bool, error)
	UpdateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error
	DeleteLinkGroup(ctx context.Context, id int) error
	GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, userID int, limit, offset int) (*response.ListResponse[models.LinkGroup], error)

	// Link
	CreateLink(ctx context.Context, link *models.Link) error
	GetLinkByID(ctx context.Context, id int) (*models.Link, error)
	SetLinkFavIconAndTitle(ctx context.Context, linkID int, favIconPath, title string) error
	GetLinksByUserIDWithPagination(ctx context.Context, userID, linkGroupID, limit, offset int, name string) (*response.ListResponse[models.LinkResponse], error)
}

type linkRepository struct {
	pool   *pgxpool.Pool
	logger logger.AppLogger
}

func New(pool *pgxpool.Pool, logger logger.AppLogger) LinkRepository {
	return &linkRepository{
		pool:   pool,
		logger: logger,
	}
}
