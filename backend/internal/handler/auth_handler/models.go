package auth_handler

import (
	"link-storage/pkg/types/app_errors"
	"link-storage/pkg/validators"
	"strings"
)

// loginRequest структура для входа
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (lr *loginRequest) validate() error {
	lr.Email = strings.ToLower(strings.TrimSpace(lr.Email))

	if !validators.IsEmailValid(lr.Email) {
		return app_errors.BadRequest("Неверный формат email", "loginRequest.validate")
	}

	if len(lr.Password) < 5 || len(lr.Password) > 15 {
		return app_errors.BadRequest(
			"пароль должен быть от 5 до 15 символов",
			"loginRequest.validate",
		)
	}

	return nil
}
