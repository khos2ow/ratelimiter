package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func Test10hpsRule(t *testing.T) {
	rule := ratelimiter.NewRule(10, 1, time.Second)

	assert.Equal(t, "10 hits per second", rule.String())
}

func Test10hp2mRule(t *testing.T) {
	rule := ratelimiter.NewRule(10, 2, time.Minute)

	assert.Equal(t, "10 hits per 2 minutes", rule.String())
}

func TestRegisteration(t *testing.T) {
	s := data.NewInMemory(&data.Options{})
	r := ratelimiter.NewRule(1, 1, time.Second)
	l := ratelimiter.NewLimiter(r, s)
	err := l.Register("foo")

	assert.Nil(t, err)
	assert.True(t, l.Store.Has("foo"))
	assert.False(t, l.Store.Has("bar"))
}

func Test1hpsLimiter(t *testing.T) {
	s := data.NewInMemory(&data.Options{})
	r := ratelimiter.NewRule(1, 1, time.Second)
	l := ratelimiter.NewLimiter(r, s)
	err := l.Register("foo")

	assert.Nil(t, err)
	assert.True(t, l.Store.Has("foo"))

	assert.True(t, l.IsAllowed("foo"))
	assert.False(t, l.IsAllowed("foo"))

	time.Sleep(2 * time.Second)

	assert.True(t, l.IsAllowed("foo"))
}

func Test2hpsLimiter(t *testing.T) {
	s := data.NewInMemory(&data.Options{})
	r := ratelimiter.NewRule(2, 1, time.Second)
	l := ratelimiter.NewLimiter(r, s)
	err := l.Register("foo")

	assert.Nil(t, err)
	assert.True(t, l.Store.Has("foo"))

	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.False(t, l.IsAllowed("foo"))

	time.Sleep(2 * time.Second)

	assert.True(t, l.IsAllowed("foo"))
}

func Test10hpsLimiter(t *testing.T) {
	s := data.NewInMemory(&data.Options{})
	r := ratelimiter.NewRule(10, 1, time.Second)
	l := ratelimiter.NewLimiter(r, s)
	err := l.Register("foo")

	assert.Nil(t, err)
	assert.True(t, l.Store.Has("foo"))

	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.True(t, l.IsAllowed("foo"))
	assert.False(t, l.IsAllowed("foo"))
}
