package config

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type Seed struct {
	AdminEmail      string
	AdminPassword   string
	AdminName       string
	AppClientSecret string
}

func LoadSeedConfig() (Seed, error) {
	seed := Seed{
		AdminEmail:      envutil.Get("SEED_ADMIN_EMAIL", ""),
		AdminPassword:   envutil.Get("SEED_ADMIN_PASSWORD", ""),
		AdminName:       envutil.Get("SEED_ADMIN_NAME", "System Admin"),
		AppClientSecret: envutil.Get("SEED_APP_CLIENT_SECRET", ""),
	}

	if err := seed.Validate(); err != nil {
		return Seed{}, fmt.Errorf("invalid seed config: %w", err)
	}

	return seed, nil
}

func (s *Seed) Validate() error {
	if s.AdminEmail == "" {
		return fmt.Errorf("SEED_ADMIN_EMAIL is required")
	}

	if s.AdminPassword == "" {
		return fmt.Errorf("SEED_ADMIN_PASSWORD is required")
	}

	if len(s.AdminPassword) < 8 {
		return fmt.Errorf("SEED_ADMIN_PASSWORD must be at least 8 characters")
	}

	if s.AppClientSecret == "" {
		return fmt.Errorf("SEED_APP_CLIENT_SECRET is required")
	}

	if len(s.AppClientSecret) < 16 {
		return fmt.Errorf("SEED_APP_CLIENT_SECRET must be at least 16 characters")
	}

	return nil
}
