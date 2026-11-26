package config

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type Redis struct {
	Host     string
	Port     int
	User     string
	Password string
	TLS      string
	Prefix   string
}

func LoadRedisConfig() (Redis, error) {
	host := envutil.Get("REDIS_HOST", "")
	port := envutil.GetInt("REDIS_PORT", 6379)
	user := envutil.Get("REDIS_USER", "")
	password := envutil.Get("REDIS_PASSWORD", "")
	tls := envutil.Get("REDIS_TLS", "")
	prefix := envutil.Get("REDIS_PREFIX", "")

	rd := Redis{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		TLS:      tls,
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
