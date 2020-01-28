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
	Limit    int
	Interval int
	Unit     time.Duration
}

func (r *Rule) String() string {
	hit := "hit"
	if r.Limit > 1 {
		hit += "s"
	}
	unit := "second"
	switch r.Unit {
	case time.Second:
		unit = "second"
	case time.Minute:
		unit = "minute"
	case time.Hour:
		unit = "hour"
	}
	if r.Interval > 1 {
		unit += "s"
		return fmt.Sprintf("%d %s per %d %s", r.Limit, hit, r.Interval, unit)
	}
	return fmt.Sprintf("%d %s per %s", r.Limit, hit, unit)
}

func (r *Rule) within(start time.Time) bool {
	diff := time.Now().Sub(start)
	var inunit float64
	switch r.Unit {
	case time.Second:
		inunit = diff.Seconds()
	case time.Minute:
		inunit = diff.Minutes()
	case time.Hour:
		inunit = diff.Hours()
	}
	return inunit < 1
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
func (l *Limiter) IsAllowed(resource string) bool {
	l.Register(resource) //nolint: errcheck

	timeframe, hits, err := l.Store.Get(resource)
	if err != nil {
		return false
	}
	if l.Rule.within(timeframe) {
		if hits >= l.Rule.Limit {
			return false
		}
		l.Store.Update(resource) //nolint: errcheck
	} else {
		if err := l.Store.Delete(resource); err != nil {
			return false
		}
	}
	return true
}

// Register registers new resource to be rate-limited
func (l *Limiter) Register(resource string) error {
	if !l.Store.Has(resource) {
		return l.Store.Put(resource)
	}
	return nil
}

// NewLimiter returns new instance of Limiter
func NewLimiter(r *Rule, s data.Store) *Limiter {
	return &Limiter{
		Rule:  r,
		Store: s,
	}
}
