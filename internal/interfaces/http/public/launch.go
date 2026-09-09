package httppublic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/bootstrap"
	"github.com/zencodecode/authorizer-service/internal/interfaces/http/middleware"
	oauth "github.com/zencodecode/authorizer-service/internal/interfaces/http/public/router/v1"
	httpserver "github.com/zencodecode/authorizer-service/internal/interfaces/http/server"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

// @title   Log Service API
// @version 1.0.0
// @description OAuth2/OIDC Authorization Server API

// @contact.name  Zencode

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAccessToken
// @in header
// @name Authorization
func Launch(ctx context.Context, c *bootstrap.Container) error {
	gin.SetMode(c.Config.Interfaces.HTTPPublic.GinMode)
	r := gin.New()

	r.Use(gin.CustomRecovery(func(gc *gin.Context, recovered any) {
		c.Logger.Error(gc, "panic recovered",
			"error", fmt.Sprintf("%v", recovered),
			"path", gc.FullPath(),
		)
		response.InternalServerError(gc, "internal server error")
	}))
	r.Use(middleware.OAuthErrorHandler())
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLog(c.Logger))

	r.Use(secure.New(secure.Config{
		SSLRedirect:        c.Config.Interfaces.HTTPPublic.Environment == "production",
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		FrameDeny:          true,
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins: c.Config.Interfaces.HTTPPublic.Cors,
	}))
	// r.Use(middleware.RateLimiterMiddleware(c.ConnManager.Redis.GetClient(), middleware.DefaultRateLimiterConfig()))

	// Health & Readiness probes (outside basePath, no auth)
	r.GET("/health", func(gc *gin.Context) {
		gc.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	r.GET("/ready", func(gc *gin.Context) {
		health := c.ConnManager.CheckHealth(gc.Request.Context())
		if !health.IsHealthy() {
			gc.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "not_ready",
				"postgres": health.Postgres,
				"redis":    health.Redis,
			})
			return
		}
		gc.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"postgres": health.Postgres,
			"redis":    health.Redis,
		})
	})

	public := r.Group("")
	// if os.Getenv("ENVIRONMENT") != "production" {
	// 	docsGroup := basePath.Group("/api-docs")
	// 	docsGroup.Use(gin.BasicAuth(gin.Accounts{
	// 		"mika": "Merdeka2025!",
	// 	}))
	// 	docsGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	v1 := oauth.SetupRouter(c, public)
	v1.MountOIDC()
	v1.MountOAuth2()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.Config.Interfaces.HTTPPublic.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return httpserver.Run(ctx, srv, c.Logger, "http-private")
}
