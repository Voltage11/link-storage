package link_repository

import (
	"context"
	"fmt"
	"link-storage/internal/models"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"time"
)

func (r *linkRepository) CreateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error {
	op := "link_repository.CreateLinkGroup"

	currentTime := time.Now()

	linkGroup.CreatedAt = currentTime
	linkGroup.UpdatedAt = currentTime

	query := `
		INSERT INTO link_groups (user_id, name, description, position, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	if err := r.pool.QueryRow(ctx, query, linkGroup.UserID, linkGroup.Name, linkGroup.Description, linkGroup.Position, linkGroup.Color, linkGroup.CreatedAt, linkGroup.UpdatedAt).Scan(&linkGroup.ID); err != nil {
		return app_errors.HandleDBError(err, "добавление группы ссылок", op)
	}

	return nil
}

func (r *linkRepository) GetLinkGroupByID(ctc context.Context, id, userID int) (*models.LinkGroup, error) {
	op := "link_repository.GetLinkGroupByID"

	query := `
		SELECT id, user_id, name, description, position, color, created_at, updated_at
		FROM link_groups
		WHERE id = $1 AND user_id = $2
	`

	var linkGroup models.LinkGroup

	if err := r.pool.QueryRow(ctc, query, id, userID).Scan(
		&linkGroup.ID,
		&linkGroup.UserID,
		&linkGroup.Name,
		&linkGroup.Description,
		&linkGroup.Position,
		&linkGroup.Color,
		&linkGroup.CreatedAt,
		&linkGroup.UpdatedAt); err != nil {
		return nil, app_errors.HandleDBError(err, "получение группы ссылок", op)
	}

	return &linkGroup, nil
}

func (r *linkRepository) HasLinkGroupWithNameByUserID(ctx context.Context, name string, userID int) (bool, error) {
	op := "link_repository.HasLinkGroupWithNameByUserID"

	query := `
		SELECT COUNT(*)
		FROM link_groups
		WHERE user_id = $1 AND
		      name ILIKE $2
	`
	var count int
	if err := r.pool.QueryRow(ctx, query, userID, name).Scan(&count); err != nil {
		return false, app_errors.HandleDBError(err, "проверка наличия группы ссылок с таким именем", op)
	}
	return count > 0, nil
}

func (r *linkRepository) UpdateLinkGroup(ctx context.Context, linkGroup *models.LinkGroup) error {
	op := "link_repository.UpdateLinkGroup"

	linkGroup.UpdatedAt = time.Now()

	query := `
		UPDATE link_groups
		SET name = $1, description = $2, color = $3, updated_at = $4
		WHERE id = $5		
	`
	result, err := r.pool.Exec(ctx, query, linkGroup.Name, linkGroup.Description, linkGroup.Color, linkGroup.UpdatedAt, linkGroup.ID)
	if err != nil {
		return app_errors.HandleDBError(err, "обновление группы ссылок", op)
	}

	if result.RowsAffected() == 0 {
		return app_errors.NotFound("группа ссылок не найдена", op)
	}

	return nil
}

func (r *linkRepository) DeleteLinkGroup(ctx context.Context, id int) error {
	op := "link_repository.DeleteLinkGroup"

	query := `
		DELETE FROM link_groups
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return app_errors.HandleDBError(err, "удаление группы ссылок", op)
	}

	if result.RowsAffected() == 0 {
		return app_errors.NotFound("группа ссылок не найдена", op)
	}

	return nil
}

func (r *linkRepository) GetLinkGroupsByUserIDWithPagination(ctx context.Context, name string, userID int, limit, offset int) (*response.ListResponse[models.LinkGroup], error) {
	op := "link_repository.GetLinkGroupsByUserIDWithPagination"

	query := `
		SELECT id, user_id, name, description, position, color, created_at, updated_at
		FROM link_groups
		WHERE user_id = $1
	`

	queryCount := `
		SELECT COUNT(*)
		FROM link_groups
		WHERE user_id = $1
	`
	args := []any{userID}
	argsCount := []any{userID}

	if name != "" {
		searchName := "%" + name + "%"
		query += fmt.Sprintf(` AND name ILIKE $%d`, len(args)+1)
		args = append(args, searchName)
		queryCount += fmt.Sprintf(` AND name ILIKE $%d`, len(argsCount)+1)
		argsCount = append(argsCount, searchName)
	}

	query += ` ORDER BY name`

	if limit > 0 && offset >= 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	}

	var total int

	// В пределах одной транзакции запросим количество записей и данные
	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
	}
	defer tx.Rollback(ctx)

	// Получим количество записей
	if err := tx.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
		return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
	}

	// Получим список записей
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
	}
	defer rows.Close()

	var linkGroups []*models.LinkGroup

	for rows.Next() {
		var linkGroup models.LinkGroup
		if err := rows.Scan(
			&linkGroup.ID,
			&linkGroup.UserID,
			&linkGroup.Name,
			&linkGroup.Description,
			&linkGroup.Position,
			&linkGroup.Color,
			&linkGroup.CreatedAt,
			&linkGroup.UpdatedAt); err != nil {
			return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
		}
		linkGroups = append(linkGroups, &linkGroup)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, app_errors.HandleDBError(err, "получение групп ссылок", op)
	}

	page := 1
	totalPages := 1
	if limit > 0 {
		page = offset/limit + 1
		totalPages = (total + limit - 1) / limit
	}

	return &response.ListResponse[models.LinkGroup]{
		Data:       linkGroups,
		Total:      total,
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	}, nil
}
