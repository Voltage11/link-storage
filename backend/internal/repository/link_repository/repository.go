package link_repository

import (
	"context"
	"link-storage/internal/models"
	"link-storage/pkg/logger"
	"link-storage/pkg/response"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository interface {
	CreateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error
	GetLinkGroupByID(ctc context.Context, id, userID int) (*models.LinkGroup, error)
	HasLinkGroupWithNameByUserID(ctx context.Context, name string, userID int) (bool, error)
	UpdateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error
	DeleteLinkGroup(ctx context.Context, id int) error
	GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, userID int, limit, offset int) (*response.ListResponse[models.LinkGroup], error)
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
