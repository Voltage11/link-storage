package models

import (
	"link-storage/pkg/types/app_errors"
	"link-storage/pkg/validators"
	"strings"
	"time"
)

type Registration struct {
	ID             int
	Name           string
	Email          string
	PasswordHashed string
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
	IsActive       bool
	ExpiredAt      time.Time
	ActivatedAt    *time.Time
	Token          string
	VerifyCode     string
}

func (r *Registration) IsValidForActivate() bool {
	return r.ActivatedAt == nil && time.Now().Before(r.ExpiredAt) && r.IsActive
}

func (r *Registration) SendEmailConfirm() {
	// TODO реализовать отправку email для подтверждения email
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegistrationRequest) Validate() error {
	if !validators.IsEmailValid(r.Email) {
		return app_errors.BadRequest("не верный формат email", "RegistrationRequest.Validate")
	}

	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.Name == "" {
		r.Name = r.Email
	} else {
		r.Name = strings.TrimSpace(r.Name)
	}

	if len(r.Password) < 5 || len(r.Password) > 15 {
		return app_errors.BadRequest(
			"неверные учетные данные",
			"RegistrationRequest.Validate",
		)
	}

	return nil
}
