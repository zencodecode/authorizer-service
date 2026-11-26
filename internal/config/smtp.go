package config

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

func LoadSMTPConfig() (SMTP, error) {
	host := envutil.Get("SMTP_HOST", "")
	port := envutil.GetInt("SMTP_PORT", 587)
	username := envutil.Get("SMTP_USERNAME", "")
	password := envutil.Get("SMTP_PASSWORD", "")
	from := envutil.Get("SMTP_FROM", "")
	fromName := envutil.Get("SMTP_FROM_NAME", "Authorizer")

	smtp := SMTP{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		FromName: fromName,
	}

	if err := smtp.Validate(); err != nil {
		return SMTP{}, fmt.Errorf("invalid smtp config: %w", err)
	}

	return smtp, nil
}

func (s *SMTP) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("host is required")
	}
	if s.From == "" {
		return fmt.Errorf("from address is required")
	}
	return nil
}
