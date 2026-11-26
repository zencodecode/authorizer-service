package config

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type RabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}

func LoadRabbitMQConfig() (RabbitMQ, error) {
	host := envutil.Get("RABBITMQ_HOST", "")
	port := envutil.GetInt("RABBITMQ_PORT", 5672)
	user := envutil.Get("RABBITMQ_USER", "")
	password := envutil.Get("RABBITMQ_PASSWORD", "")
	vhost := envutil.Get("RABBITMQ_VHOST", "/")

	rmq := RabbitMQ{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Vhost:    vhost,
	}

	if err := rmq.Validate(); err != nil {
		return RabbitMQ{}, fmt.Errorf("invalid rabbitmq config: %w", err)
	}

	return rmq, nil
}

func (r *RabbitMQ) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("host is required")
	}

	if r.Port == 0 {
		return fmt.Errorf("port is required")
	}

	if r.User == "" {
		return fmt.Errorf("user is required")
	}

	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	if r.Vhost == "" {
		return fmt.Errorf("vhost is required")
	}

	return nil
}
