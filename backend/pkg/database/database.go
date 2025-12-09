package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"link-storage/pkg/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Pool   *pgxpool.Pool
	logger logger.AppLogger
}

type Config struct {
	DSN           string
	MigrationPath string
	MaxConns      int32
	MinConns      int32
}

func New(cfg Config, appLogger logger.AppLogger) (*Database, error) {
	// Устанавливаем значения по умолчанию
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 25
	}
	if cfg.MinConns == 0 {
		cfg.MinConns = 5
	}

	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DSN: %w", err)
	}

	// Настройки пула соединений
	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 30 * time.Second
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка ping к бд: %w", err)
	}

	database := &Database{
		Pool:   pool,
		logger: appLogger,
	}

	// Выполняем миграции если указан путь
	if cfg.MigrationPath != "" {
		if err := database.migrate(cfg.MigrationPath); err != nil {
			return nil, fmt.Errorf("ошибка миграций: %w", err)
		}
	}

	return database, nil
}

func (d *Database) migrate(migrationPath string) error {
	op := "database.migrate"

	// Создаем отдельное соединение для миграций через database/sql
	sqlDB, err := d.createSQLDB()
	if err != nil {
		return fmt.Errorf("ошибка создания sql соединения для миграций: %w", err)
	}
	defer sqlDB.Close()

	// Проверяем соединение
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ошибка ping для миграций: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("ошибка создания драйвера миграции: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("ошибка применения миграции: %w", err)
	}

	d.logger.Info("Миграции выполнены успешно", op)
	return nil
}

// createSQLDB создает стандартное sql.DB соединение для миграций
func (d *Database) createSQLDB() (*sql.DB, error) {
	config := d.Pool.Config()

	// Используем те же параметры что и для основного пула
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.ConnConfig.User,
		string(config.ConnConfig.Password),
		config.ConnConfig.Host,
		config.ConnConfig.Port,
		config.ConnConfig.Database,
	)

	return sql.Open("pgx", dsn)
}

// Close закрывает пул соединений
func (d *Database) Close() {
	d.Pool.Close()
}

// WithTransaction выполняет операцию в транзакции
func (d *Database) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
