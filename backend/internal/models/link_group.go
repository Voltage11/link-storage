package models

import (
	"link-storage/pkg/types/app_errors"
	"time"
)

type LinkGroup struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// type LinkGroupShortResponse struct {
// 	ID          int       `json:"id"`
// 	UserID      int       `json:"user_id"`
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// 	Position    int       `json:"position"`
// 	Color       string    `json:"color"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

type LinkGroupCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

func (lc *LinkGroupCreate) Validate() error {
	if len(lc.Name) < 3 || len(lc.Name) > 50 {
		return app_errors.BadRequest("Имя группы должно быть от 3 до 50 символов", "LinkGroupCreate.Validate")
	}
	return nil
}

type LinkGroupUpdate struct {
	ID          int    `json:"-;omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (lc *LinkGroupUpdate) Validate() error {
	if len(lc.Name) < 3 || len(lc.Name) > 50 {
		return app_errors.BadRequest("Имя группы должно быть от 3 до 50 символов", "LinkGroupCreate.Validate")
	}
	return nil
}
