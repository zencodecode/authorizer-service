package middleware

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

var slidingWindowScript = redis.NewScript(`
redis.call('ZREMRANGEBYSCORE', KEYS[1], '0', ARGV[1])
local count = redis.call('ZCARD', KEYS[1])
redis.call('ZADD', KEYS[1], ARGV[2], ARGV[3])
redis.call('EXPIRE', KEYS[1], ARGV[4])
return count
`)

type RateLimiterConfig struct {
	MaxRequests             int
	Window                  time.Duration
	KeyFunc                 func(c *gin.Context) string
	Prefix                  string
	FailOpen                bool
	FallbackMaxRequests     int
	FallbackWindow          time.Duration
	RedisTimeout            time.Duration
	BreakerFailureThreshold int
	BreakerCooldown         time.Duration
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		MaxRequests:             100,
		Window:                  1 * time.Minute,
		KeyFunc:                 nil,
		Prefix:                  "ratelimit",
		FailOpen:                true,
		FallbackMaxRequests:     100,
		FallbackWindow:          1 * time.Minute,
		RedisTimeout:            150 * time.Millisecond,
		BreakerFailureThreshold: 5,
		BreakerCooldown:         10 * time.Second,
	}
}

type breakerState int

const (
	stateClosed breakerState = iota
	stateOpen
	stateHalfOpen
)

type circuitBreaker struct {
	mu        sync.Mutex
	state     breakerState
	failures  int
	threshold int
	cooldown  time.Duration
	openedAt  time.Time
}

func newCircuitBreaker(threshold int, cooldown time.Duration) *circuitBreaker {
	return &circuitBreaker{threshold: threshold, cooldown: cooldown}
}

func (b *circuitBreaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == stateOpen {
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = stateHalfOpen
			return true
		}
		return false
	}
	return true
}

func (b *circuitBreaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = stateClosed
}

func (b *circuitBreaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == stateHalfOpen || b.failures >= b.threshold {
		b.state = stateOpen
		b.openedAt = time.Now()
	}
}

type fallbackLimiter struct {
	mu      sync.Mutex
	buckets map[string]*fallbackBucket
	max     int
	window  time.Duration
}

type fallbackBucket struct {
	count       int
	windowStart time.Time
}

func newFallbackLimiter(max int, window time.Duration) *fallbackLimiter {
	return &fallbackLimiter{
		buckets: make(map[string]*fallbackBucket),
		max:     max,
		window:  window,
	}
}

func (f *fallbackLimiter) allow(key string) (allowed bool, remaining int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	b, ok := f.buckets[key]
	if !ok || now.Sub(b.windowStart) >= f.window {
		b = &fallbackBucket{windowStart: now}
		f.buckets[key] = b
	}
	b.count++

	if len(f.buckets) > 10000 {
		for k, v := range f.buckets {
			if now.Sub(v.windowStart) >= f.window {
				delete(f.buckets, k)
			}
		}
	}

	remaining = f.max - b.count
	if remaining < 0 {
		remaining = 0
	}
	return b.count <= f.max, remaining
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
	if cfg.MaxRequests <= 0 {
		cfg.MaxRequests = 100
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.FallbackMaxRequests <= 0 {
		cfg.FallbackMaxRequests = cfg.MaxRequests
	}
	if cfg.FallbackWindow <= 0 {
		cfg.FallbackWindow = cfg.Window
	}
	if cfg.RedisTimeout <= 0 {
		cfg.RedisTimeout = 150 * time.Millisecond
	}
	if cfg.BreakerFailureThreshold <= 0 {
		cfg.BreakerFailureThreshold = 5
	}
	if cfg.BreakerCooldown <= 0 {
		cfg.BreakerCooldown = 10 * time.Second
	}

	breaker := newCircuitBreaker(cfg.BreakerFailureThreshold, cfg.BreakerCooldown)
	fallback := newFallbackLimiter(cfg.FallbackMaxRequests, cfg.FallbackWindow)

	return func(c *gin.Context) {
		key := fmt.Sprintf("%s:%s", cfg.Prefix, cfg.KeyFunc(c))

		if !breaker.allow() {
			handleFallback(c, cfg, fallback, key)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.RedisTimeout)
		defer cancel()

		now := time.Now()
		windowStart := now.Add(-cfg.Window)

		member, err := uniqueMember(now)
		if err != nil {
			breaker.recordFailure()
			handleFallback(c, cfg, fallback, key)
			return
		}

		ttlSeconds := int(cfg.Window.Seconds()) + 1

		res, err := slidingWindowScript.Run(ctx, rdb, []string{key},
			strconv.FormatInt(windowStart.UnixMicro(), 10),
			strconv.FormatInt(now.UnixMicro(), 10),
			member,
			ttlSeconds,
		).Result()

		if err != nil {
			breaker.recordFailure()
			handleFallback(c, cfg, fallback, key)
			return
		}
		breaker.recordSuccess()

		count, _ := res.(int64)

		remaining := cfg.MaxRequests - int(count) - 1
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
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

func handleFallback(c *gin.Context, cfg RateLimiterConfig, fallback *fallbackLimiter, key string) {
	if !cfg.FailOpen {
		response.TooManyRequests(c, "rate limiting temporarily unavailable")
		c.Abort()
		return
	}

	allowed, remaining := fallback.allow(key)

	c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.FallbackMaxRequests))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Fallback", "true")

	if !allowed {
		c.Header("Retry-After", strconv.Itoa(int(cfg.FallbackWindow.Seconds())))
		response.TooManyRequests(c, "rate limit exceeded")
		c.Abort()
		return
	}
	c.Next()
}
func uniqueMember(now time.Time) (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d-%d", now.UnixNano(), binary.BigEndian.Uint64(b[:])), nil
}
