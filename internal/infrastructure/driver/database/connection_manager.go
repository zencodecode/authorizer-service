package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database/postgres"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database/redis"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
)

const healthCheckTimeout = 5 * time.Second

type ConnectionManager struct {
	Postgres *postgres.Connection
	Redis    *redis.Connection
	Cfg      config.DatabaseConfig
	Logger   *logger.Logger
}

func NewConnectionManager(ctx context.Context, cfg config.DatabaseConfig, logger *logger.Logger) (*ConnectionManager, error) {
	pg := postgres.NewConnection(&cfg.Postgres, logger)
	if err := pg.Connect(ctx); err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	rd := redis.NewConnection(&cfg.Redis, logger)
	if err := rd.Connect(ctx); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return &ConnectionManager{
		Postgres: pg,
		Redis:    rd,
		Cfg:      cfg,
		Logger:   logger,
	}, nil
}

type HealthStatus struct {
	Postgres  DatabaseHealth `json:"postgres"`
	Redis     DatabaseHealth `json:"redis"`
	Timestamp time.Time      `json:"timestamp"`
}

func (h HealthStatus) IsHealthy() bool {
	return h.Postgres.Status == "healthy" && h.Redis.Status == "healthy"
}

type DatabaseHealth struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency_ms"`
}

// CheckHealth mengecek kedua database secara paralel.
func (cm *ConnectionManager) CheckHealth(ctx context.Context) HealthStatus {
	pgCh := make(chan DatabaseHealth, 1)
	rdCh := make(chan DatabaseHealth, 1)

	go func() { pgCh <- cm.checkPostgres(ctx) }()
	go func() { rdCh <- cm.checkRedis(ctx) }()

	return HealthStatus{
		Postgres:  <-pgCh,
		Redis:     <-rdCh,
		Timestamp: time.Now(),
	}
}

func (cm *ConnectionManager) checkPostgres(ctx context.Context) DatabaseHealth {
	latency, err := cm.Postgres.Ping(ctx, healthCheckTimeout)
	if err != nil {
		cm.Logger.Warn(ctx, "postgres health check failed", "error", err.Error())
		return DatabaseHealth{Status: "unhealthy", Message: err.Error(), Latency: latency}
	}
	return DatabaseHealth{Status: "healthy", Latency: latency}
}

func (cm *ConnectionManager) checkRedis(ctx context.Context) DatabaseHealth {
	latency, err := cm.Redis.Ping(ctx, healthCheckTimeout)
	if err != nil {
		cm.Logger.Warn(ctx, "redis health check failed", "error", err.Error())
		return DatabaseHealth{Status: "unhealthy", Message: err.Error(), Latency: latency}
	}
	return DatabaseHealth{Status: "healthy", Latency: latency}
}

const shutdownTimeout = 10 * time.Second

func (cm *ConnectionManager) Shutdown(ctx context.Context) error {
	pgDone := make(chan error, 1)
	rdDone := make(chan error, 1)

	go func() { pgDone <- cm.Postgres.Close(ctx, shutdownTimeout) }()
	go func() { rdDone <- cm.Redis.Close() }()

	pgErr := <-pgDone
	rdErr := <-rdDone

	switch {
	case pgErr != nil && rdErr != nil:
		return fmt.Errorf("postgres shutdown failed: %w; redis shutdown failed: %v", pgErr, rdErr)
	case pgErr != nil:
		return fmt.Errorf("postgres shutdown failed: %w", pgErr)
	case rdErr != nil:
		return fmt.Errorf("redis shutdown failed: %w", rdErr)
	default:
		cm.Logger.Info(ctx, "all database connections closed gracefully")
		return nil
	}
}
