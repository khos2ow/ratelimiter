package data

import (
	"fmt"
	"time"
)

type tracker struct {
	timeframe time.Time
	hits      int
}

// InMemory holds the options and for with in memory
// key-value data store
type InMemory struct {
	options *Options
	cache   map[string]*tracker
}

// Options returns the current options passed to InMemory
func (i InMemory) Options() *Options {
	return i.options
}

// Connect connects to InMemory database
func (i InMemory) Connect() error {
	return nil
}

// Has checks if the record with provided 'key' exists in the cache
func (i InMemory) Has(key string) bool {
	if key == "" {
		return false
	}
	_, ok := i.cache[key]
	return ok
}

// Get gets a value from InMemory based on provided key
func (i InMemory) Get(key string) (time.Time, int, error) {
	if i.Has(key) {
		t := i.cache[key]
		return t.timeframe, t.hits, nil
	}
	return time.Time{}, 0, fmt.Errorf("key '%s' not found", key)
}

// Put saves a new key/value to InMemory
func (i InMemory) Put(key string) error {
	if i.Has(key) {
		return fmt.Errorf("provided key %s already exists", key)
	}
	i.cache[key] = &tracker{
		timeframe: time.Now(),
		hits:      0,
	}
	return nil
}

// Delete deletes a key/value pari from InMemory
func (i InMemory) Delete(key string) error {
	if !i.Has(key) {
		return fmt.Errorf("provided key %s not found", key)
	}
	delete(i.cache, key)
	return nil
}

// Update updates value of a key in InMemory
func (i InMemory) Update(key string) error {
	timeframe, hits, err := i.Get(key)
	if err != nil {
		return fmt.Errorf("provided key %s not found", key)
	}
	i.cache[key] = &tracker{
		timeframe: timeframe,
		hits:      hits + 1,
	}
	return nil
}

// NewInMemory returns new InMemory client
func NewInMemory(options *Options) *InMemory {
	return &InMemory{
		options: options,
		cache:   make(map[string]*tracker),
	}
}
