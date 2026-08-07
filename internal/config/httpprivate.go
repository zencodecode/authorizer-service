package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/pkg/enval"
)

type HttpPrivate struct {
	Environment string
	GinMode     string
	Port        string
	BasePath    string
	Cors        []string
}

func LoadHttpPrivateConfig() (HttpPrivate, error) {
	env := enval.Get("ENVIRONMENT", "production")
	corsOrigins := enval.Get("CORS", "")
	port := enval.Get("HTTP_PRIVATE_PORT", "8080")
	basePath := enval.Get("BASE_PATH", "/api")

	ginMode := gin.ReleaseMode
	if env != "production" {
		ginMode = gin.DebugMode
	}

	if env == "production" && corsOrigins == "" {
		return HttpPrivate{}, fmt.Errorf("CORS must be configured in production")
	}

	if corsOrigins == "" {
		corsOrigins = "*"
	}

	httpPrivate := HttpPrivate{
		Environment: env,
		GinMode:     ginMode,
		Port:        port,
		BasePath:    basePath,
		Cors:        strings.Split(corsOrigins, ","),
	}

	if err := httpPrivate.Validate(); err != nil {
		return HttpPrivate{}, fmt.Errorf("invalid http private config: %w", err)
	}

	return httpPrivate, nil
}

func (h *HttpPrivate) Validate() error {
	if h.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if h.GinMode == "" {
		return fmt.Errorf("gin mode is required")
	}

	if h.Port == "" {
		return fmt.Errorf("port is required")
	}

	if h.BasePath == "" {
		return fmt.Errorf("base path is required")
	}

	if len(h.Cors) == 0 {
		return fmt.Errorf("cors is required")
	}

	return nil
}
