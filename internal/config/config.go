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
		HTTPPublic  HTTPPublic
		HTTPPrivate HTTPPrivate
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

	pblCfg, err := LoadHTTPPublicConfig()
	if err != nil {
		return Config{}, err
	}

	pvtCfg, err := LoadHTTPPrivateConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Database: DatabaseConfig{Postgres: pgCfg, Redis: rdCfg},
		Auth:     AuthConfig{JWT: jwtCfg},
		Interfaces: InterfacesConfig{
			HTTPPublic:  pblCfg,
			HTTPPrivate: pvtCfg,
		},
	}, nil
}
