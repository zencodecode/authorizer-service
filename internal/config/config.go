package config

type (
	Config struct {
		Database   DatabaseConfig
		Auth       AuthConfig
		Interfaces InterfacesConfig
		SMTP       SMTP
	}

	DatabaseConfig struct {
		Postgres Postgres
		Redis    Redis
	}

	AuthConfig struct {
		JWT  JWT
		OIDC OIDC
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

	oidcCfg, err := LoadOIDCConfig(jwtCfg)
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

	smtpCfg, err := LoadSMTPConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Database: DatabaseConfig{Postgres: pgCfg, Redis: rdCfg},
		Auth: AuthConfig{
			JWT:  jwtCfg,
			OIDC: oidcCfg,
		},
		Interfaces: InterfacesConfig{
			HTTPPublic:  pblCfg,
			HTTPPrivate: pvtCfg,
		},
		SMTP: smtpCfg,
	}, nil
}
