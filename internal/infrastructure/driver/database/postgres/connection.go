package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	postgresClient *gorm.DB
	config         *config.Postgres
	logger         *logger.Logger
	shutdownOnce   sync.Once
}

func NewConnection(cfg *config.Postgres, log *logger.Logger) *Connection {
	return &Connection{
		config: cfg,
		logger: log,
	}
}

func (cn *Connection) GetClient() *gorm.DB {
	return cn.postgresClient
}

func (cn *Connection) retryConnect(
	ctx context.Context,
	connectFunc func(context.Context) error,
	dbType string,
) error {
	const (
		maxAttempts  = 3
		initialDelay = 1 * time.Second
		maxDelay     = 10 * time.Second
	)

	var lastErr error
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		cn.logger.Info(ctx, "Attempting database connection",
			"database", dbType,
			"attempt", attempt,
			"max_attempts", maxAttempts,
		)

		err := connectFunc(ctx)
		if err == nil {
			if attempt > 1 {
				cn.logger.Info(ctx, "Successfully connected after retries",
					"database", dbType,
					"attempts", attempt,
				)
			}
			return nil
		}

		lastErr = err
		cn.logger.Warn(ctx, "Connection attempt failed",
			"database", dbType,
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"error", err.Error(),
		)

		// Don't sleep after the last attempt
		if attempt < maxAttempts {
			cn.logger.Info(ctx, "Retrying connection",
				"database", dbType,
				"delay", delay.String(),
			)

			select {
			case <-time.After(delay):
				// Continue to next attempt
			case <-ctx.Done():
				return fmt.Errorf("connection retry cancelled: %w", ctx.Err())
			}

			// Exponential backoff: double the delay for next attempt
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return fmt.Errorf("failed to connect to %s after %d attempts: %w", dbType, maxAttempts, lastErr)
}

func (cn *Connection) Connect(ctx context.Context) error {
	cn.logger.Info(ctx, "Starting PostgreSQL connection",
		"host", cn.config.Host,
		"port", cn.config.Port,
		"database", cn.config.Database,
	)

	connectFunc := func(ctx context.Context) error {
		// Build PostgreSQL DSN for GORM
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cn.config.Host,
			cn.config.Port,
			cn.config.User,
			cn.config.Password,
			cn.config.Database,
			cn.config.SSLMode,
		)

		// Open GORM connection with auto-migration disabled
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to open GORM PostgreSQL connection: %w", err)
		}

		// Get underlying *sql.DB to configure connection pool
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB from GORM: %w", err)
		}

		// Configure connection pool
		sqlDB.SetMaxOpenConns(cn.config.MaxOpenConnection)
		sqlDB.SetMaxIdleConns(cn.config.MaxIdleConnection)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)

		// Ping to verify connection
		pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(pingCtx); err != nil {
			_ = sqlDB.Close()
			return fmt.Errorf("failed to ping PostgreSQL: %w", err)
		}

		cn.postgresClient = db
		return nil
	}

	// Use retry logic to establish connection
	if err := cn.retryConnect(ctx, connectFunc, "postgresql"); err != nil {
		cn.logger.Error(ctx, "Failed to connect to PostgreSQL",
			"error", err.Error(),
		)
		return err
	}

	cn.logger.Info(ctx, "Successfully connected to PostgreSQL",
		"host", cn.config.Host,
		"port", cn.config.Port,
		"database", cn.config.Database,
		"pool_config", map[string]any{
			"max_open_conns":     100,
			"max_idle_conns":     10,
			"conn_max_lifetime":  "1h",
			"conn_max_idle_time": "5m",
		},
	)

	return nil
}

func (cn *Connection) Close(ctx context.Context, timeout time.Duration) error {
	if cn.postgresClient == nil {
		return nil
	}
	sqlDB, err := cn.postgresClient.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql.DB: %w", err)
	}

	closeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- sqlDB.Close() }()

	select {
	case err := <-errCh:
		return err
	case <-closeCtx.Done():
		return fmt.Errorf("postgresql close timeout: %w", closeCtx.Err())
	}
}

func (cn *Connection) Ping(ctx context.Context, timeout time.Duration) (time.Duration, error) {
	if cn.postgresClient == nil {
		return 0, fmt.Errorf("postgres client is not initialized")
	}

	sqlDB, err := cn.postgresClient.DB()
	if err != nil {
		return 0, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	err = sqlDB.PingContext(pingCtx)
	return time.Since(start), err
}
