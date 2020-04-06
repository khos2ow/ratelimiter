package ratelimiter

import (
	"fmt"
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
)

// Rule represents a rule for ratelimiting
// it consists of:
//   - limit    <number>    (e.g. 100)
//   - interval <number>    (e.g. 1)
//   - unit     <time unit> (e.g. second)
//
// This example will create a rule of "100 requests per 1 second"
type Rule struct {
	limit    int
	interval int
	unit     time.Duration
}

// NewRule creates new Rule with provided limit, interval and unit
func NewRule(limit int, interval int, unit time.Duration) *Rule {
	// TODO decide what to do when limit or interval is zero and unit is not s,m,h
	return &Rule{
		limit:    limit,
		interval: interval,
		unit:     unit,
	}
}

func (r *Rule) String() string {
	hit := "hit"
	if r.limit > 1 {
		hit += "s"
	}
	unit := "second"
	switch r.unit {
	case time.Second:
		unit = "second"
	case time.Minute:
		unit = "minute"
	case time.Hour:
		unit = "hour"
	}
	if r.interval > 1 {
		unit = fmt.Sprintf("%d %ss", r.interval, unit)
	}
	return fmt.Sprintf("%d %s per %s", r.limit, hit, unit)
}

// Limiter represents the main ratelimiting function. It
// is aware of Rule and Data Store to do the caching in.
type Limiter struct {
	rule  *Rule
	store data.Store
}

// NewLimiter returns new instance of Limiter
func NewLimiter(rule *Rule, store data.Store) *Limiter {
	return &Limiter{
		rule:  rule,
		store: store,
	}
}

// IsAllowed checks the allowance of the requests in the current
// timeframe with previous hits in the same timeframe in mind.
func (l *Limiter) IsAllowed(resource string) bool {
	bucket := time.Duration(l.rule.interval) * l.rule.unit
	key := time.Now().Truncate(bucket).String()
	if !l.store.Has(key) {
		l.store.Create(key, resource)
	}
	hits, err := l.store.Get(key, resource)
	if err != nil {
		return false
	}
	if hits <= l.rule.limit {
		if hits, err = l.store.Add(key, resource, time.Now()); err != nil {
			return false
		}
	}
	return hits <= l.rule.limit
}
