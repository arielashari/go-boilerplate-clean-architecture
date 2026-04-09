package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func (s *SlogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var slogLevel slog.Level
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case tracelog.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case tracelog.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case tracelog.LogLevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	s.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}

func NewPostgresConnection(cfg configs.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Name,
		cfg.Database.Postgres.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	maxLifetime, _ := time.ParseDuration(cfg.Database.Postgres.MaxConnLifetime)
	maxIdle, _ := time.ParseDuration(cfg.Database.Postgres.MaxConnIdleTime)
	healthCheckPeriod, _ := time.ParseDuration(cfg.Database.Postgres.HealthCheckPeriod)

	poolConfig.MaxConns = cfg.Database.Postgres.MaxConns
	poolConfig.MinConns = cfg.Database.Postgres.MinConns
	poolConfig.MaxConnLifetime = maxLifetime
	poolConfig.MaxConnIdleTime = maxIdle
	poolConfig.HealthCheckPeriod = healthCheckPeriod

	if cfg.App.Env != "prod" {
		poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   &SlogAdapter{logger: logger},
			LogLevel: tracelog.LogLevelDebug,
		}
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
