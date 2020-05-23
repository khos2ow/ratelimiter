package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
)

func TestRuleToString(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		interval int
		unit     time.Duration
		expected string
	}{
		{
			name:     "rule to string",
			limit:    1,
			interval: 1,
			unit:     time.Second,
			expected: "1 hit per second",
		},
		{
			name:     "rule to string",
			limit:    10,
			interval: 1,
			unit:     time.Second,
			expected: "10 hits per second",
		},
		{
			name:     "rule to string",
			limit:    20,
			interval: 3,
			unit:     time.Second,
			expected: "20 hits per 3 seconds",
		},
		{
			name:     "rule to string",
			limit:    1,
			interval: 1,
			unit:     time.Minute,
			expected: "1 hit per minute",
		},
		{
			name:     "rule to string",
			limit:    50,
			interval: 1,
			unit:     time.Minute,
			expected: "50 hits per minute",
		},
		{
			name:     "rule to string",
			limit:    90,
			interval: 2,
			unit:     time.Minute,
			expected: "90 hits per 2 minutes",
		},
		{
			name:     "rule to string",
			limit:    1,
			interval: 1,
			unit:     time.Hour,
			expected: "1 hit per hour",
		},
		{
			name:     "rule to string",
			limit:    100,
			interval: 1,
			unit:     time.Hour,
			expected: "100 hits per hour",
		},
		{
			name:     "rule to string",
			limit:    200,
			interval: 2,
			unit:     time.Hour,
			expected: "200 hits per 2 hours",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			rule := ratelimiter.NewRule(tt.limit, tt.interval, tt.unit)

			assert.Equal(tt.expected, rule.String())
		})
	}
}

func TestLimiterIsAllowed(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		interval int
		unit     time.Duration
		expected string
	}{
		{
			name:     "limiter is allowed",
			limit:    1,
			interval: 1,
			unit:     time.Second,
		},
		{
			name:     "limiter is allowed",
			limit:    2,
			interval: 1,
			unit:     time.Second,
		},
		{
			name:     "limiter is allowed",
			limit:    10,
			interval: 1,
			unit:     time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			store := data.NewInMemory(&data.Options{})
			rate := ratelimiter.NewRule(tt.limit, tt.interval, tt.unit)
			limiter := ratelimiter.NewLimiter(rate, store)

			for i := 0; i < tt.limit; i++ {
				assert.True(limiter.IsAllowed("foo"))
			}
			assert.False(limiter.IsAllowed("foo"))

			time.Sleep(2 * time.Second)

			assert.True(limiter.IsAllowed("foo"))
		})
	}
}
