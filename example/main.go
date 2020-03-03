package main

import (
	"fmt"
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
)

func main() {
	resource := "foo"
	store := data.NewInMemory(&data.Options{})
	rule := ratelimiter.NewRule(10, 1, time.Second)
	limiter := ratelimiter.NewLimiter(rule, store)
	if err := limiter.Register(resource); err != nil {
		panic(err)
	}

	start := time.Now()
	fmt.Printf("limiting resource '%s' to %s\n\n", resource, rule.String())

	for i := 0; i < 25; i++ {
		fmt.Printf("hit #%-10dallowed: %-10velapsed: %f seconds\n", i+1, limiter.IsAllowed(resource), time.Now().Sub(start).Seconds())
		time.Sleep(80 * time.Millisecond)
	}

	fmt.Printf("\ntook %f seconds\n", time.Now().Sub(start).Seconds())
}
