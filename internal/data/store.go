package data

import "time"

// Store respresents an interface of type of data store
// which gets implemented in different types, e.g. Redis,
// InMemory, etc.
type Store interface {
	// Options returns the current options passed to Store
	Options() *Options

	// Connect connects to Store database
	Connect() error

	// Has checks if the time slot bucket 'key' exists
	Has(key string) bool

	// Get hits count based on bucket key and resource name
	Get(key string, resource string) (int, error)

	// Create bucket for specitic time slot and resource name
	Create(key string, resource string) error

	// Add new hit for particular 'resource' in particular time bucket 'key'
	Add(key string, resource string, now time.Time) (int, error)

	// Delete a specific 'key' time bucket
	Delete(key string) error
}
