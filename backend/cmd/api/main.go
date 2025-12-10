package main

import (
	"fmt"
	"link-storage/internal/config"
	"link-storage/internal/handler/auth_handler"
	"link-storage/internal/handler/link_handler"
	"link-storage/internal/middleware"
	"link-storage/internal/repository/auth_repository"
	"link-storage/internal/repository/link_repository"
	"link-storage/internal/service/auth_service"
	"link-storage/internal/service/link_service"
	"link-storage/pkg/database"
	"link-storage/pkg/logger"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	// Config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Logger
	appLogger := logger.New(cfg.LogLevel)
	appLogger.Debug("Конфигурация загружена")

	// Database
	dbConfig := database.Config{
		DSN:           cfg.DSN(),
		MigrationPath: cfg.DB.MigrationPath,
		MaxConns:      25,
		MinConns:      5,
	}
	appDb, err := database.New(dbConfig, appLogger)
	if err != nil {
		panic("ошибка при подключении к БД")
	}
	defer appDb.Close()

	// Обработаем падение приложения
	defer func() {
		if err := recover(); err != nil {
			appLogger.Error(fmt.Errorf("PANIC: %v", err), "main")
		}
	}()

	// Repositories
	authRepo := auth_repository.New(appDb, appLogger)
	linkRepo := link_repository.New(appDb.Pool, appLogger)

	// Services
	authService := auth_service.New(authRepo, appLogger, cfg.Secret.Jwt)
	linkService := link_service.New(linkRepo, appLogger)

	// Server
	router := chi.NewRouter()

	// Middleware - ПРАВИЛЬНЫЙ ПОРЯДОК
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.SecurityHeaders)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Timeout(30 * time.Second))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(middleware.NewCORSMiddleware(cfg.Server.Cors).Handler)
	router.Use(middleware.RequireJSONContentType)
	router.Use(middleware.AuthMiddleware(authService, appLogger))

	// Handlers
	auth_handler.New(router, authService, appLogger)
	link_handler.New(router, linkService, appLogger)

	// Run
	runServer := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	appLogger.Warn(fmt.Sprintf("Запуск сервера на: %s", runServer), "main")
	log.Fatal(http.ListenAndServe(runServer, router))
}
