package data

import (
	"time"
)

// Store respresents an interface of type of data store
// which gets implemented in different types, e.g. Redis,
// InMemory, etc.
type Store interface {
	Options() *Options
	Connect() error
	Has(key string) bool
	Get(key string) (time.Time, int, error)
	Put(key string) error
	Delete(key string) error
	Update(key string) error
}
