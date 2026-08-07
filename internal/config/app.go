package config

type (
	Config struct {
		Database   DatabaseConfig
		Auth       AuthConfig
		Interfaces InterfacesConfig
	}

	DatabaseConfig struct {
		Postgres Postgres
		Redis    Redis
	}

	AuthConfig struct {
		JWT JWT
	}

	InterfacesConfig struct {
		HttpPrivate HttpPrivate
	}
)

func Load() (Config, error) {
	pgCfg, err := LoadPostgresConfig()
	if err != nil {
		return Config{}, err
	}

	jwtCfg, err := LoadJWTConfig()
	if err != nil {
		return Config{}, err
	}

	rdCfg, err := LoadRedisConfig()
	if err != nil {
		return Config{}, err
	}

	httpCfg, err := LoadHttpPrivateConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		// GinMode:  helper.GetEnv("GIN_MODE", "release"),
		Database: DatabaseConfig{Postgres: pgCfg, Redis: rdCfg},
		Auth:     AuthConfig{JWT: jwtCfg},
		Interfaces: InterfacesConfig{
			HttpPrivate: httpCfg,
		},
	}, nil
}
