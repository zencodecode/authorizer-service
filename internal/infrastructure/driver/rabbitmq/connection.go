package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type Connection struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config config.RabbitMQ
	logger service.Logger
	mu     sync.Mutex
}

func NewConnection(cfg config.RabbitMQ, logger service.Logger) *Connection {
	return &Connection{
		config: cfg,
		logger: logger,
	}
}

func (cn *Connection) Connect(ctx context.Context) error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cn.config.User,
		cn.config.Password,
		cn.config.Host,
		cn.config.Port,
		cn.config.Vhost,
	)

	cn.logger.Info(ctx, "connecting to RabbitMQ",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	conn, err := amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	cn.conn = conn
	cn.ch = ch

	cn.logger.Info(ctx, "connected to RabbitMQ successfully",
		"host", cn.config.Host,
		"port", cn.config.Port,
	)

	return nil
}

// Channel returns the default channel created during Connect().
func (cn *Connection) Channel() *amqp.Channel {
	return cn.ch
}

// NewChannel creates a new independent channel from the same connection.
// Useful when publisher and consumer need separate channels.
func (cn *Connection) NewChannel() (*amqp.Channel, error) {
	if cn.conn == nil {
		return nil, fmt.Errorf("rabbitmq connection is not initialized")
	}
	return cn.conn.Channel()
}

func (cn *Connection) Close() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	var errs []error
	if cn.ch != nil {
		if err := cn.ch.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if cn.conn != nil {
		if err := cn.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing RabbitMQ: %v", errs)
	}
	return nil
}
