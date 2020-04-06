package data

import (
	"fmt"
	"time"
)

// bucket holds list of 'hits' (timestamp) for a particular 'resource'
type bucket map[string][]time.Time

// InMemory holds the options and for with in memory
// key-value data store
type InMemory struct {
	options *Options
	cache   map[string]bucket
}

// NewInMemory returns new InMemory client
func NewInMemory(options *Options) *InMemory {
	return &InMemory{
		options: options,
		cache:   make(map[string]bucket),
	}
}

// Options returns the current options passed to InMemory
func (i *InMemory) Options() *Options {
	return i.options
}

// Connect to InMemory database
func (i *InMemory) Connect() error {
	return nil
}

// Has checks if the time slot bucket 'key' exists
func (i *InMemory) Has(key string) bool {
	if key == "" {
		return false
	}
	if _, ok := i.cache[key]; !ok {
		return false
	}
	return true
}

// Get hits count from InMemory based on bucket key and resource name
func (i *InMemory) Get(key string, resource string) (int, error) {
	if !i.Has(key) {
		return 0, fmt.Errorf("bucket '%s' not found", key)
	}
	hits, ok := i.cache[key][resource]
	if !ok {
		return 0, fmt.Errorf("resource '%s' not found", resource)
	}
	return len(hits), nil
}

// Create bucket in InMemory for specitic time slot and resource name
func (i *InMemory) Create(key string, resource string) error {
	if !i.Has(key) {
		i.cache[key] = make(bucket)
	}
	if _, ok := i.cache[key][resource]; ok {
		return fmt.Errorf("resource %s exists", key)
	}
	hits := make([]time.Time, 0)
	i.cache[key][resource] = hits
	return nil
}

// Add new hit for particular 'resource' in particular time bucket 'key'
func (i *InMemory) Add(key string, resource string, now time.Time) (int, error) {
	if !i.Has(key) {
		return 0, fmt.Errorf("bucket %s not found", key)
	}
	hits := i.cache[key][resource]
	if hits == nil {
		return 0, fmt.Errorf("resource %s not found", resource)
	}
	hits = append(hits, now)
	i.cache[key][resource] = hits
	return len(hits), nil
}

// Delete a specific 'key' time bucket
func (i *InMemory) Delete(key string) error {
	if !i.Has(key) {
		return fmt.Errorf("bucket key %s not found", key)
	}
	delete(i.cache, key)
	return nil
}
