package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zencodecode/authorizer-service/pkg/enval"
)

type Postgres struct {
	Host                          string
	Port                          int
	User                          string
	Password                      string
	Database                      string
	SSLMode                       string
	LogMode                       string
	MaxIdleConnection             int
	MaxOpenConnection             int
	ConnectionMaxLifetimeInSecond time.Duration
}

func LoadPostgresConfig() (Postgres, error) {
	debugMode := enval.GetAsBool("DB_DEBUG", false)
	host := enval.Get("POSTGRES_HOST", "")
	port := enval.GetInt("POSTGRES_PORT", 5432)
	user := enval.Get("POSTGRES_USER", "")
	password := enval.Get("POSTGRES_PASSWORD", "")
	database := enval.Get("POSTGRES_DATABASE", "")
	sslMode := enval.Get("POSTGRES_SSLMODE", "disable")
	maxIdleConnection := enval.GetInt("POSTGRES_MAX_IDLE_CONNECTION", 10)
	maxOpenConnection := enval.GetInt("POSTGRES_MAX_OPEN_CONNECTION", 100)
	connectionMaxLifetimeInSecond := enval.GetDuration("POSTGRES_CONNECTION_MAX_LIFETIME_IN_SECOND", 3600*time.Second)

	logMode := 0
	if debugMode == true {
		logMode = 3
	}

	pg := Postgres{
		Host:                          host,
		Port:                          port,
		User:                          user,
		Password:                      password,
		Database:                      database,
		SSLMode:                       sslMode,
		LogMode:                       strconv.Itoa(logMode),
		MaxIdleConnection:             maxIdleConnection,
		MaxOpenConnection:             maxOpenConnection,
		ConnectionMaxLifetimeInSecond: connectionMaxLifetimeInSecond,
	}

	if err := pg.Validate(); err != nil {
		return Postgres{}, fmt.Errorf("invalid postgres config: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	}

	if p.User == "" {
		return fmt.Errorf("user is required")
	}

	if p.Password == "" {
		return fmt.Errorf("password is required")
	}

	if p.Database == "" {
		return fmt.Errorf("database is required")
	}

	if p.Port < 1 || p.Port > 65535 {
		return fmt.Errorf("invalid port: must be between 1 and 65535, got %d", p.Port)
	}

	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}

	if !validSSLModes[p.SSLMode] {
		return fmt.Errorf("invalid sslmode: must be one of [disable, require, verify-ca, verify-full], got '%s'", p.SSLMode)
	}

	return nil
}
