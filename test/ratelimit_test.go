package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/cache"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimit(t *testing.T) {
	// Setup Redis client for testing
	redisConfig := &configs.AppConfig{
		Redis: configs.Redis{
			Address: "localhost:6379",
		},
	}

	rc, err := cache.NewRedisClient(redisConfig)
	if err != nil {
		fmt.Printf("\n⚠️  Redis is not available at %s - Skipping rate limit tests\n", redisConfig.Redis.Address)
		t.Skip("Redis is not available, skipping rate limit tests")
		return
	}
	defer rc.Close()

	tests := []struct {
		name          string
		config        configs.RateLimitTokenBucketConfig
		requests      int
		expectedAllow []bool
		sleepBetween  time.Duration
		cleanupBefore bool
	}{
		{
			name: "basic rate limiting",
			config: configs.RateLimitTokenBucketConfig{
				Capacity:  3,
				Rate:      time.Second,
				BaseKey:   "test:ratelimit:1",
				TimeoutMS: 10000,
			},
			requests:      5,
			expectedAllow: []bool{true, true, true, false, false},
			sleepBetween:  0,
			cleanupBefore: true,
		},
		{
			name: "rate limit with refill",
			config: configs.RateLimitTokenBucketConfig{
				Capacity:  2,
				Rate:      time.Second * 3, // Increase refill time to 3 seconds
				BaseKey:   "test:ratelimit:2",
				TimeoutMS: 100000,
			},
			requests:      4,
			expectedAllow: []bool{true, true, false, true},
			sleepBetween:  time.Second * 1, // Reduce sleep time - this ensures 3rd request fails
			cleanupBefore: true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up the test key before running the test
			err := rc.Del(ctx, tt.config.BaseKey)
			assert.NoError(t, err, "Failed to cleanup Redis key")

			if tt.cleanupBefore {
				rc.Close()
				rc, err = cache.NewRedisClient(redisConfig)
				assert.NoError(t, err)
			}

			for i := 0; i < tt.requests; i++ {
				allowed, err := rc.Allow(ctx, tt.config, tt.config.BaseKey)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAllow[i], allowed, "Request %d", i+1)

				if tt.sleepBetween > 0 {
					time.Sleep(tt.sleepBetween)
				}
			}
		})
	}
}
