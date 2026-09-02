package queue

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
)

type Connection struct {
	config       *config.RabbitMQ
	logger       *logger.Logger
	shutdownOnce sync.Once
}

func NewConnection(cfg *config.RabbitMQ, log *logger.Logger) *Connection {
	return &Connection{
		config: cfg,
		logger: log,
	}
}

func (cn *Connection) Close() error {
	return nil
}

func failOnError(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("failed to ping Redis: %s", msg)
	}
	return nil
}

func (cn *Connection) Connect(ctx context.Context) error {
	cn.logger.Info(ctx, "Starting RabbitMQ connection",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cn.config.User,
		cn.config.Password,
		cn.config.Host,
		cn.config.Port,
		cn.config.Vhost)

	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	cn.logger.Info(ctx, "Successfully connected to RabbitMQ",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	return nil
}
