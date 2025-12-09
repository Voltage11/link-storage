package middleware

import (
	"context"
	"link-storage/internal/models"
	"link-storage/internal/service/auth_service"
	"link-storage/pkg/logger"
	"link-storage/pkg/response"
	"link-storage/pkg/types/app_errors"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(authService auth_service.AuthService, logger logger.AppLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			op := "middleware.AuthMiddleware"

			// OPTIONS запросы должны проходить без проверки авторизации
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r.WithContext(r.Context()))
				return
			}

			var token string

			// Проверка заголовка
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}

			if token == "" {
				response.WriteError(w, app_errors.Unauthorized(op))
				return
			}

			currentUser, err := authService.VerifyJwt(token)
			if err != nil || currentUser == nil {
				logger.Warn("Доступ с невалидным токеном", op, "token", token)
				response.WriteError(w, app_errors.Unauthorized(op))
				return
			}

			if !currentUser.IsActive {
				response.WriteError(w, app_errors.BadRequest("Пользователь не активен", op))
				return
			}

			if isAdminPath(r.URL.Path) && !currentUser.IsAdmin {
				response.WriteError(w, app_errors.Forbidden(op))
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, currentUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Вспомогательные функции остаются без изменений
func isPublicPath(path string) bool {
	if strings.HasPrefix(path, "/api/v1/auth/activate") {
		return true
	}

	publicPaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh-token",
	}
	for _, p := range publicPaths {
		if p == path {
			return true
		}
	}
	return false
}

func isAdminPath(path string) bool {
	adminPrefix := []string{"/api/v1/admin", "api/v2/admin", "api/v3/admin"}

	for _, p := range adminPrefix {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// Обновленная функция для работы с контекстом
func GetCurrentUserFromContext(ctx context.Context) *models.CurrentUser {
	if user, ok := ctx.Value(UserContextKey).(*models.CurrentUser); ok {
		return user
	}
	return nil
}

func IsCurrentUserWithID(ctx context.Context, userID int) bool {
	if user := GetCurrentUserFromContext(ctx); user != nil {
		return user.ID == userID
	}
	return false
}
