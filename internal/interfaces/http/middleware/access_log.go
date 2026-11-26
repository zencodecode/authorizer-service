// interfaces/http-private/middleware/access_log.go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

func AccessLog(log service.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Info(c, "request handled",
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
			"client_ip", c.ClientIP(),
		)
	}
}
