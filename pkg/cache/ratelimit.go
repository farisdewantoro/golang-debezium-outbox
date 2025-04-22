package cache

import (
	"context"
	"eventdrivensystem/configs"
	"fmt"
	"time"
)

// Allow checks if a request is allowed based on the token bucket algorithm
func (rc *RedisClient) Allow(ctx context.Context, config configs.RateLimitTokenBucketConfig, key string) (bool, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	script := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])
		local rate = tonumber(ARGV[3])
		local timeout = tonumber(ARGV[4])

		-- Get or initialize state
		local bucket = redis.call('hmget', key, 'tokens', 'last_refill')
		local tokens = bucket[1] and tonumber(bucket[1]) or capacity
		local last_refill = bucket[2] and tonumber(bucket[2]) or now

		-- Calculate refill
		local time_passed = now - last_refill
		local new_tokens = math.floor(time_passed / rate)
		
		if new_tokens > 0 then
			tokens = math.min(capacity, tokens + new_tokens) -- Add correct number of tokens
			last_refill = now -- Update last_refill to current time
		end

		local allowed = 0
		if tokens > 0 then
			tokens = tokens - 1
			allowed = 1
		end

		redis.call('hmset', key, 'tokens', tokens, 'last_refill', last_refill)
		redis.call('pexpire', key, timeout)
		return allowed
	`

	result, err := rc.client.Eval(ctx, script, []string{key},
		now,
		config.Capacity,
		config.Rate.Milliseconds(),
		config.TimeoutMS,
	).Int()

	if err != nil {
		return false, fmt.Errorf("failed to execute rate limit check: %v", err)
	}

	return result == 1, nil
}
