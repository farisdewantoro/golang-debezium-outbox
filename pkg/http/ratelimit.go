package http

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/cache"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RateLimiter struct {
	redisClient *cache.RedisClient
	config      configs.RateLimitTokenBucketConfig
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(redisClient *cache.RedisClient, config configs.RateLimitTokenBucketConfig) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		config:      config,
	}
}

// RateLimitMiddleware returns an Echo middleware function that implements rate limiting
func (rl *RateLimiter) RateLimitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// If the rate limit is not enabled, skip the middleware
			if !rl.config.Enabled {
				return next(c)
			}

			// Use IP address in the rate limit key
			key := rl.config.BaseKey + ":user_ip:" + c.RealIP()
			fmt.Printf("Rate limit key: %s\n", key)
			allowed, err := rl.redisClient.Allow(c.Request().Context(), rl.config, key)
			if err != nil {
				return err
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded. Please try again later.",
				})
			}

			return next(c)
		}
	}
}
