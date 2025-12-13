package link_repository

import (
	"context"
	"fmt"

	"link-storage/internal/models"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"time"
)

func (r *linkRepository) CreateLink(ctx context.Context, link *models.Link) error {
	op := "link_repository.CreateLink"

	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now

	query := `
		INSERT INTO links (user_id, link_group_id, url, title, description, is_archived, is_favorite, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	if err := r.pool.QueryRow(ctx, query, link.UserID,
		link.LinkGroupID,
		link.URL,
		link.Title,
		link.Description,
		link.IsArchived,
		link.IsFavorite,
		link.CreatedAt,
		link.UpdatedAt).Scan(&link.ID); err != nil {
		return app_errors.HandleDBError(err, "Создание ссылки", op)
	}
	return nil
}

func (r *linkRepository) GetLinkByID(ctx context.Context, linkID int) (*models.Link, error) {
	op := "link_repository.GetLinkByID"

	query := `
		SELECT id, user_id, link_group_id, url, title, description, favicon_url, preview_image, is_archived, is_favorite,
		       click_count, last_visited, created_at, updated_at
		FROM links
		WHERE id = $1
	`

	var link models.Link

	if err := r.pool.QueryRow(ctx, query, linkID).Scan(
		&link.ID,
		&link.UserID,
		&link.LinkGroupID,
		&link.URL,
		&link.Title,
		&link.Description,
		&link.FaviconURL,
		&link.PreviewImage,
		&link.IsArchived,
		&link.IsFavorite,
		&link.ClickCount,
		&link.LastVisited,
		&link.CreatedAt,
		&link.UpdatedAt); err != nil {
		return nil, app_errors.HandleDBError(err, "Получение ссылки по ID", op)
	}

	return &link, nil
}

func (r *linkRepository) SetLinkFavIconAndTitle(ctx context.Context, linkID int, favIconPath, title string) error {
	op := "link_repository.SetLinkFavIcon"

	query := `
		UPDATE links 
			SET favicon_url = $1,
			    title = $2,
			    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.pool.Exec(ctx, query, favIconPath, title, linkID)
	if err != nil {
		return app_errors.HandleDBError(err, "Установка favicon для ссылки", op)
	}
	if result.RowsAffected() == 0 {
		return app_errors.NotFound("Ссылка не найдена", op)
	}
	return nil
}

func (r *linkRepository) GetLinksByUserIDWithPagination(ctx context.Context, userID, linkGroupID, limit, offset int, name string) (*response.ListResponse[models.LinkResponse], error) {
	op := "link_repository.GetLinksByUserIDWithPagination"

	query := `
		SELECT l.id, l.user_id, l.link_group_id, l.url, l.title, l.description, l.favicon_url, l.preview_image, l.is_archived, l.is_favorite,
			   l.click_count, l.last_visited, l.created_at, l.updated_at, g.id, g.name
		FROM links l LEFT JOIN link_groups g ON l.link_group_id = g.id
		WHERE l.user_id = $1
	`

	queryCount := `SELECT COUNT(id) FROM links WHERE user_id = $1`

	args := []any{userID}
	argsCount := []any{userID}

	if linkGroupID > 0 {
		query += fmt.Sprintf(" AND l.link_group_id = $%d", len(args)+1)
		args = append(args, linkGroupID)

		queryCount += fmt.Sprintf(" AND link_group_id = $%d", len(argsCount)+1)
		argsCount = append(argsCount, linkGroupID)
	}

	if name != "" {
		search := "%" + name + "%"
		query += fmt.Sprintf(` AND ((l.title ILIKE $%d) OR (l.url ILIKE $%d))`, len(args)+1, len(args)+2)
		args = append(args, search, search)

		queryCount += fmt.Sprintf(` AND ((title ILIKE $%d) OR (url ILIKE $%d))`, len(argsCount)+1, len(argsCount)+2)
		argsCount = append(argsCount, search, search)
	}

	query += ` ORDER BY l.title ASC`
	if limit > 0 && offset >= 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	}

	var total int

	tx, err := r.pool.Begin(ctx)

	if err != nil {
		r.logger.Error(err, op)
		return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
	}
	defer tx.Rollback(ctx)

	// Получим количество записей
	if err := tx.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
		r.logger.Error(err, op)
		return nil, app_errors.HandleDBError(err, "получение ссылок", op)
	}

	// Получим список записей
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error(err, op)
		return nil, app_errors.HandleDBError(err, "получение ссылок", op)
	}
	defer rows.Close()

	var links []*models.LinkResponse

	for rows.Next() {
		var link models.LinkResponse
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.LinkGroupID,
			&link.URL,
			&link.Title,
			&link.Description,
			&link.FaviconURL,
			&link.PreviewImage,
			&link.IsArchived,
			&link.IsFavorite,
			&link.ClickCount,
			&link.LastVisited,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.Group.ID,	
			&link.Group.Name); err != nil {
			return nil, app_errors.HandleDBError(err, "получение ссылок", op)
		}
		

		links = append(links, &link)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, app_errors.HandleDBError(err, "получение ссылок", op)
	}

	page := 1
	totalPages := 1
	if limit > 0 {
		page = offset/limit + 1
		totalPages = (total + limit - 1) / limit
	}

	return &response.ListResponse[models.LinkResponse]{
		Data:       links,
		Total:      total,
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	}, nil
}
