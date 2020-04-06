package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
)

func main() {
	resource := "foo"
	store := data.NewInMemory(&data.Options{})
	rule := ratelimiter.NewRule(10, 1, time.Second)
	limiter := ratelimiter.NewLimiter(rule, store)

	start := time.Now()
	fmt.Printf("limiting resource '%s' to %s\n\n", resource, rule.String())

	for i := 0; i < 100; i++ {
		allowed, err := limiter.IsAllowed(resource)
		if err != nil {
			fmt.Printf("hit #%-10derror: %-10velapsed: %f seconds\n", i+1, err.Error(), time.Now().Sub(start).Seconds())
		} else {
			fmt.Printf("hit #%-10dallowed: %-10velapsed: %f seconds\n", i+1, allowed, time.Now().Sub(start).Seconds())
		}
		time.Sleep(time.Duration(rand.Intn(100)-rand.Intn(20)) * time.Millisecond)
	}

	fmt.Printf("\ntook %f seconds\n", time.Now().Sub(start).Seconds())
}
