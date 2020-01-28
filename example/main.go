package main

import (
	"fmt"
	"time"

	"github.com/khos2ow/ratelimiter"
	"github.com/khos2ow/ratelimiter/internal/data"
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
		fmt.Printf("hit #%d\t\tallowed: %v\n", i+1, limiter.IsAllowed(resource))
		time.Sleep(80 * time.Millisecond)
	}

	fmt.Printf("\ntook %f seconds\n", time.Now().Sub(start).Seconds())
}
