package app_errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// AppError ошибка приложения
type AppError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"-"`
	Err     error  `json:"-"`
	Op      string `json:"-"`
}

func (e *AppError) Error() string {
	msg := e.Message
	if e.Op != "" {
		msg = fmt.Sprintf("[%s] %s", e.Op, msg)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New создает новую ошибку с операцией
func New(statusCode int, errType string, message string, op string, err ...error) *AppError {
	appErr := &AppError{
		Type:    errType,
		Message: message,
		Code:    statusCode,
		Op:      op,
	}

	if len(err) > 0 && err[0] != nil {
		appErr.Err = err[0]
	}

	return appErr
}

// Конструкторы с поддержкой операции

func NotFound(message, op string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, op)
}

func BadRequest(message, op string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, op)
}

func BadRequestWithError(err error, message, op string) *AppError {
	if message != "" {
		message = err.Error()
	}

	return &AppError{
		Type:    "BAD_REQUEST",
		Message: message,
		Code:    http.StatusBadRequest,
		Op:      op,
		Err:     err,
	}
}

func Conflict(message, op string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, op)
}

func Internal(err error, op string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR",
		"Внутренняя ошибка сервера", op, err)
}

func Unauthorized(op string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED",
		"Неавторизованный доступ", op)
}

func Forbidden(op string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN",
		"Доступ запрещен", op)
}

func Validation(message, op string) *AppError {
	return New(http.StatusBadRequest, "VALIDATION_ERROR", message, op)
}

// HandleDBError обрабатывает ошибки Database с указанием операции
func HandleDBError(err error, message, op string) error {
	if err == nil {
		return nil
	}

	// Если уже наша ошибка, возвращаем её
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Обработка ошибок pgx
	if errors.Is(err, pgx.ErrNoRows) {
		return NotFound(message, op)
	}

	// Обработка PostgreSQL ошибок
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return Conflict(message+" (запись уже существует)", op)
		case "23503": // foreign_key_violation
			return BadRequest("Связанная запись не найдена", op)
		case "23502": // not_null_violation
			return BadRequest("Обязательное поле не заполнено", op)
		case "23514": // check_violation
			return BadRequest("Некорректное значение поля", op)
		case "25P02": // in_failed_sql_transaction
			return BadRequest("Ошибка в транзакции", op)
		case "40001": // serialization_failure
			return BadRequest("Конфликт параллельных операций", op)
		default:
			return Internal(fmt.Errorf("%s: %w", message, err), op)
		}
	}

	// Обертываем остальные ошибки с контекстом
	return Internal(fmt.Errorf("%s: %w", message, err), op)
}

// Хелперы для проверки
func IsNotFound(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Type == "NOT_FOUND"
}

func IsConflict(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Type == "CONFLICT"
}
