package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type HttpPublic struct {
	Environment string
	GinMode     string
	Port        string
	BasePath    string
	Cors        []string
}

func LoadHttpPublicConfig() (HttpPublic, error) {
	env := envutil.Get("ENVIRONMENT", "production")
	corsOrigins := envutil.Get("CORS", "")
	port := envutil.Get("HTTP_PUBLIC_PORT", "8181")
	basePath := envutil.Get("BASE_PATH", "/api")

	ginMode := gin.ReleaseMode
	if env != "production" {
		ginMode = gin.DebugMode
	}

	if env == "production" && corsOrigins == "" {
		return HttpPublic{}, fmt.Errorf("CORS must be configured in production")
	}

	if corsOrigins == "" {
		corsOrigins = "*"
	}

	httpPublic := HttpPublic{
		Environment: env,
		GinMode:     ginMode,
		Port:        port,
		BasePath:    basePath,
		Cors:        strings.Split(corsOrigins, ","),
	}

	if err := httpPublic.Validate(); err != nil {
		return HttpPublic{}, fmt.Errorf("invalid http private config: %w", err)
	}

	return httpPublic, nil
}

func (h *HttpPublic) Validate() error {
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
