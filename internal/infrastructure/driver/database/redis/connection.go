package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
)

type Connection struct {
	redisClient  *redis.Client
	config       *config.Redis
	logger       *logger.Logger
	shutdownOnce sync.Once
}

func NewConnection(cfg *config.Redis, log *logger.Logger) *Connection {
	return &Connection{
		config: cfg,
		logger: log,
	}
}

func (cn *Connection) GetClient() *redis.Client {
	return cn.redisClient
}

func (cn *Connection) Ping(ctx context.Context, timeout time.Duration) (time.Duration, error) {
	if cn.redisClient == nil {
		return 0, fmt.Errorf("redis client is not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	err := cn.redisClient.Ping(pingCtx).Err()
	return time.Since(start), err
}

func (cn *Connection) Close() error {
	if cn.redisClient == nil {
		return nil
	}
	return cn.redisClient.Close()
}

func (cn *Connection) Connect(ctx context.Context) error {
	cn.logger.Info(ctx, "Starting Redis connection",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cn.config.Host, cn.config.Port),
		DB:       0,
		Protocol: 2,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	cn.redisClient = rdb

	cn.logger.Info(ctx, "Successfully connected to Redis",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	return nil
}
