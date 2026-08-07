package config

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/pkg/enval"
)

type Redis struct {
	Host     string
	Port     int
	User     string
	Password string
	Tls      string
	Prefix   string
}

func LoadRedisConfig() (Redis, error) {
	host := enval.Get("REDIS_HOST", "")
	port := enval.GetInt("REDIS_PORT", 6379)
	user := enval.Get("REDIS_USER", "")
	password := enval.Get("REDIS_PASSWORD", "")
	tls := enval.Get("REDIS_TLS", "")
	prefix := enval.Get("REDIS_PREFIX", "")

	rd := Redis{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Tls:      tls,
		Prefix:   prefix,
	}

	if err := rd.Validate(); err != nil {
		return Redis{}, fmt.Errorf("invalid redis config: %w", err)
	}

	return rd, nil
}

func (r *Redis) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("host is required")
	}

	if r.Port == 0 {
		return fmt.Errorf("port is required")
	}

	return nil
}
