package ratelimiter_test

import (
	"fmt"
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
)

func Example_limiter() {
	resource := "foo"
	store := data.NewInMemory(&data.Options{})
	rule := ratelimiter.NewRule(10, 1, time.Second)
	limiter := ratelimiter.NewLimiter(rule, store)
	if err := limiter.Register(resource); err != nil {
		panic(err)
	}

	fmt.Printf("limiting resource '%s' to %s\n\n", resource, rule.String())

	for i := 0; i < 25; i++ {
		fmt.Printf("hit #%-10dallowed: %v\n", i+1, limiter.IsAllowed(resource))
		time.Sleep(80 * time.Millisecond)
	}

	// Output:
	// limiting resource 'foo' to 10 hits per second
	//
	// hit #1         allowed: true
	// hit #2         allowed: true
	// hit #3         allowed: true
	// hit #4         allowed: true
	// hit #5         allowed: true
	// hit #6         allowed: true
	// hit #7         allowed: true
	// hit #8         allowed: true
	// hit #9         allowed: true
	// hit #10        allowed: true
	// hit #11        allowed: false
	// hit #12        allowed: false
	// hit #13        allowed: false
	// hit #14        allowed: true
	// hit #15        allowed: true
	// hit #16        allowed: true
	// hit #17        allowed: true
	// hit #18        allowed: true
	// hit #19        allowed: true
	// hit #20        allowed: true
	// hit #21        allowed: true
	// hit #22        allowed: true
	// hit #23        allowed: true
	// hit #24        allowed: true
	// hit #25        allowed: false
}
