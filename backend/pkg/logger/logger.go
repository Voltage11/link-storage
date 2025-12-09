package logger

import (
	"log/slog"
	"os"
)

type AppLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, op string, args ...any)
	Warn(msg string, op string, args ...any)
	Error(err error, op string, args ...any)
}

func getLoggerLevel(level string) slog.Leveler {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type appLogger struct {
	logger *slog.Logger
}

func New(levelCfg string) AppLogger {

	level := getLoggerLevel(levelCfg)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Переименовываем стандартные атрибуты
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "timestamp", Value: a.Value}
			}
			if a.Key == slog.LevelKey {
				return slog.Attr{Key: "level", Value: a.Value}
			}
			if a.Key == slog.MessageKey {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return a
		},
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &appLogger{
		logger: logger,
	}
}

func (l *appLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *appLogger) Info(msg string, op string, args ...any) {
	allArgs := make([]any, 0, len(args)+1)
	allArgs = append(allArgs, slog.String("op", op))
	allArgs = append(allArgs, args...)

	l.logger.Info(msg, allArgs...)
}

func (l *appLogger) Warn(msg string, op string, args ...any) {
	allArgs := make([]any, 0, len(args)+1)
	allArgs = append(allArgs, slog.String("op", op))
	allArgs = append(allArgs, args...)

	l.logger.Warn(msg, allArgs...)
}

func (l *appLogger) Error(err error, op string, args ...any) {
	if err == nil {
		return
	}

	allArgs := make([]any, 0, len(args)+2)
	allArgs = append(allArgs, slog.String("op", op))
	allArgs = append(allArgs, slog.String("error", err.Error()))
	allArgs = append(allArgs, args...)

	l.logger.Error("operation failed", allArgs...)
}
