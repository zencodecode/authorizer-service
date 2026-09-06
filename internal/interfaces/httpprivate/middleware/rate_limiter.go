package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

type RateLimiterConfig struct {
	MaxRequests int
	Window      time.Duration
	KeyFunc     func(c *gin.Context) string
	Prefix      string
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		KeyFunc:     nil,
		Prefix:      "ratelimit",
	}
}

func RateLimiter(rdb *redis.Client, cfg RateLimiterConfig) gin.HandlerFunc {
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "ratelimit"
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key := fmt.Sprintf("%s:%s", cfg.Prefix, cfg.KeyFunc(c))
		now := time.Now()
		windowStart := now.Add(-cfg.Window)

		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixMicro(), 10))
		countCmd := pipe.ZCard(ctx, key)

		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now.UnixMicro()),
			Member: fmt.Sprintf("%d", now.UnixNano()),
		})

		pipe.Expire(ctx, key, cfg.Window+time.Second)

		_, err := pipe.Exec(ctx)
		if err != nil {
			c.Next()
			return
		}

		count := countCmd.Val()

		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, cfg.MaxRequests-int(count)-1)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now.Add(cfg.Window).Unix(), 10))

		if int(count) >= cfg.MaxRequests {
			c.Header("Retry-After", strconv.Itoa(int(cfg.Window.Seconds())))
			response.TooManyRequests(c, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
