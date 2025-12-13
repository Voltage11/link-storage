package models

import (
	"link-storage/pkg/types/app_errors"
	"time"
)

type Link struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	LinkGroupID  *int      `json:"link_group_id,omitempty"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	FaviconURL   string    `json:"favicon_url,omitempty"`
	PreviewImage string    `json:"preview_image,omitempty"`
	IsArchived   bool      `json:"is_archived"`
	IsFavorite   bool      `json:"is_favorite"`
	ClickCount   int       `json:"click_count"`
	LastVisited  time.Time `json:"last_visited"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LinkResponse struct {
	Link
	Group struct {
		ID          *int    `json:"id"`
		Name        *string `json:"name"`
	} `json:"link_group,omitempty"`
	// Group *LinkGroup `json:"group,omitempty"`
}

type LinkCreate struct {
	LinkGroupID *int   `json:"link_group_id,omitempty"`
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsFavorite  bool   `json:"is_favorite,omitempty"`
}

func (l *LinkCreate) Validate() error {
	if l.URL == "" {
		return app_errors.BadRequest("URL не может быть пустым", "LinkCreate.Validate")
	}
	return nil
}
