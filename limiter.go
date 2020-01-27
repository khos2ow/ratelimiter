package ratelimiter

import (
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
)

// Rule represents a rule for ratelimiting
// it consists of:
//   - limit    <number>    (e.g. 100)
//   - interval <number>    (e.g. 1)
//   - unit     <time unit> (e.g. second)
//
// This will create a rule or "100 requests per 1 second"
type Rule struct {
	Limit    int
	Interval int
	Unit     time.Duration
}

// NewRule creates new Rule with provided limit, interval and unit
func NewRule(l int, i int, u time.Duration) *Rule {
	return &Rule{
		Limit:    l,
		Interval: i,
		Unit:     u,
	}
}

// Limiter represents the main ratelimiting function.
// It is aware of Rule and Data Store to do the caching
// in.
type Limiter struct {
	Rule  *Rule
	Store data.Store
}

// IsAllowed checks the allowance of the requests in the
// current timeframe with previous hits in the same timeframe
// in mind.
func (l *Limiter) IsAllowed() bool {
	// TODO
	return true
}

// NewLimiter returns new instance of Limiter
func NewLimiter(r *Rule, s data.Store) *Limiter {
	return &Limiter{
		Rule:  r,
		Store: s,
	}
}
